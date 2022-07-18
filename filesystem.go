/*
Copyright (c) 2019 Dell Corporation
All Rights Reserved
*/

package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dell/gounity/util"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

//Filesystem structure
type Filesystem struct {
	client *Client
}

//FsNameMaxLength provides the allowed max length for filesystem name
const (
	FsNameMaxLength = 63
)

//AccessType type is string
type AccessType string

//AccessType constants
const (
	ReadOnlyAccessType      = AccessType("READ_ONLY")
	ReadWriteAccessType     = AccessType("READ_WRITE")
	ReadOnlyRootAccessType  = AccessType("READ_ONLY_ROOT")
	ReadWriteRootAccessType = AccessType("READ_WRITE_ROOT")
)

//NFSShareDefaultAccess is string
type NFSShareDefaultAccess string

//NFSShareDefaultAccess constants
const (
	NoneDefaultAccess          = NFSShareDefaultAccess("0")
	ReadOnlyDefaultAccess      = NFSShareDefaultAccess("1")
	ReadWriteDefaultAccess     = NFSShareDefaultAccess("2")
	ReadOnlyRootDefaultAccess  = NFSShareDefaultAccess("3")
	ReadWriteRootDefaultAccess = NFSShareDefaultAccess("4")
)

//ErrorFilesystemNotFound stores error for filesystem not found
var ErrorFilesystemNotFound = errors.New("Unable to find filesystem")

//FilesystemNotFoundErrorCode stores error code for filesystem not found
var FilesystemNotFoundErrorCode = "0x7d13005"

//AttachedSnapshotsErrorCode stores error code for attached snapshots
var AttachedSnapshotsErrorCode = "0x6000c17"

//MarkFilesystemForDeletion stores filesystem for deletion mark
var MarkFilesystemForDeletion = "csi-marked-filesystem-for-deletion(do not remove this from description)"

//NewFilesystem function returns filesystem
func NewFilesystem(client *Client) *Filesystem {
	return &Filesystem{client}
}

//FindFilesystemByName - Find the Filesystem by it's name. If the Filesystem is not found, an error will be returned.
func (f *Filesystem) FindFilesystemByName(ctx context.Context, filesystemName string) (*types.Filesystem, error) {
	if len(filesystemName) == 0 {
		return nil, errors.New("Filesystem Name shouldn't be empty")
	}
	fileSystemResp := &types.Filesystem{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.FileSystemAction, filesystemName, FileSystemDisplayFields), nil, fileSystemResp)
	if err != nil {
		if strings.Contains(err.Error(), FilesystemNotFoundErrorCode) {
			return nil, ErrorFilesystemNotFound
		}
		return nil, err
	}
	return fileSystemResp, nil
}

//FindFilesystemByID - Find the Filesystem by it's Id. If the Filesystem is not found, an error will be returned.
func (f *Filesystem) FindFilesystemByID(ctx context.Context, filesystemID string) (*types.Filesystem, error) {
	log := util.GetRunIDLogger(ctx)
	if len(filesystemID) == 0 {
		return nil, errors.New("Filesystem Id shouldn't be empty")
	}
	fileSystemResp := &types.Filesystem{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.FileSystemAction, filesystemID, FileSystemDisplayFields), nil, fileSystemResp)
	if err != nil {
		log.Debugf("Unable to find filesystem Id %s Error: %v", filesystemID, err)
		if strings.Contains(err.Error(), FilesystemNotFoundErrorCode) {
			return nil, ErrorFilesystemNotFound
		}
		return nil, err
	}
	return fileSystemResp, nil
}

//GetFilesystemIDFromResID - Returns the filesystem ID for the filesystem
func (f *Filesystem) GetFilesystemIDFromResID(ctx context.Context, filesystemResID string) (string, error) {
	if filesystemResID == "" {
		return "", errors.New("Filesystem Resource Id shouldn't be empty")
	}

	fileSystemResp := &types.StorageResourceParameters{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.StorageResourceAction, filesystemResID, StorageResourceDisplayFields), nil, fileSystemResp)
	if err != nil {
		return "", fmt.Errorf("get filesystem Id for %s failed with error: %v", filesystemResID, err)
	}
	return fileSystemResp.StorageResourceContent.Filesystem.ID, nil
}

//CreateFilesystem - Create a new filesystem on the array
func (f *Filesystem) CreateFilesystem(ctx context.Context, name, storagepool, description, nasServer string, size uint64, tieringPolicy, hostIOSize, supportedProtocol int, isThinEnabled, isDataReductionEnabled bool, isReplicationDestination bool) (*types.Filesystem, error) {
	log := util.GetRunIDLogger(ctx)
	if name == "" {
		return nil, errors.New("filesystem name should not be empty")
	}

	if len(name) > FsNameMaxLength {
		return nil, fmt.Errorf("filesystem name %s should not exceed %d characters", name, FsNameMaxLength)
	}

	poolAPI := NewStoragePool(f.client)
	pool, err := poolAPI.FindStoragePoolByID(ctx, storagepool)

	if err != nil {
		return nil, fmt.Errorf("unable to get PoolID (%s) Error:%v", storagepool, err)
	}

	storagePool := types.StoragePoolID{
		PoolID: storagepool,
	}

	fileEventSettings := types.FileEventSettings{
		IsCIFSEnabled: false, //Set to false to disable CIFS
		IsNFSEnabled:  true,  //Set to true to enable NFS alone
	}

	nas := types.NasServerID{
		NasServerID: nasServer,
	}

	fsParams := types.FsParameters{
		StoragePool:       &storagePool,
		Size:              size,
		SupportedProtocol: supportedProtocol,
		HostIOSize:        hostIOSize,
		NasServer:         &nas,
		FileEventSettings: fileEventSettings,
	}

	replParameters := types.ReplicationParameters{
		IsReplicationDestination: isReplicationDestination,
	}

	volAPI := NewVolume(f.client)
	thinProvisioningLicenseInfoResp, err := volAPI.isFeatureLicensed(ctx, ThinProvisioning)
	if err != nil {
		return nil, fmt.Errorf("unable to get license info for feature: %s", ThinProvisioning)
	}

	dataReductionLicenseInfoResp, err := volAPI.isFeatureLicensed(ctx, DataReduction)
	if err != nil {
		return nil, fmt.Errorf("unable to get license info for feature: %s", DataReduction)
	}

	if thinProvisioningLicenseInfoResp.LicenseInfoContent.IsInstalled && thinProvisioningLicenseInfoResp.LicenseInfoContent.IsValid {
		fsParams.IsThinEnabled = strconv.FormatBool(isThinEnabled)
	} else if isThinEnabled == true {
		return nil, fmt.Errorf("thin provisioning is not supported on array and hence cannot create Filesystem")
	}

	if dataReductionLicenseInfoResp.LicenseInfoContent.IsInstalled && dataReductionLicenseInfoResp.LicenseInfoContent.IsValid {
		fsParams.IsDataReductionEnabled = strconv.FormatBool(isDataReductionEnabled)
	} else if isDataReductionEnabled == true {
		return nil, fmt.Errorf("data reduction is not supported on array and hence cannot create Filesystem")
	}

	if pool != nil && pool.StoragePoolContent.PoolFastVP.Status != 0 {
		log.Debug("FastVP is enabled")
		fastVPParameters := types.FastVPParameters{
			TieringPolicy: tieringPolicy,
		}
		fsParams.FastVPParameters = &fastVPParameters
	} else {
		log.Debug("FastVP is not enabled")
		if tieringPolicy != 0 {
			return nil, fmt.Errorf("fastVP is not enabled and requested tiering policy is: %d ", tieringPolicy)
		}
	}

	fileReqParam := types.FsCreateParam{
		Name:                  name,
		Description:           description,
		FsParameters:          &fsParams,
		ReplicationParameters: &replParameters,
	}

	fileResp := &types.Filesystem{}
	err = f.client.executeWithRetryAuthenticate(ctx,
		http.MethodPost, fmt.Sprintf(api.UnityAPIStorageResourceActionURI, api.CreateFSAction), fileReqParam, fileResp)
	if err != nil {
		return nil, err
	}

	return fileResp, nil
}

//DeleteFilesystem delete by its ID. If the Filesystem is not present on the array, an error will be returned.
func (f *Filesystem) DeleteFilesystem(ctx context.Context, filesystemID string) error {
	log := util.GetRunIDLogger(ctx)
	if len(filesystemID) == 0 {
		return errors.New("Filesystem Id cannot be empty")
	}

	filesystemResp, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return err
	}
	resourceID := filesystemResp.FileContent.StorageResource.ID
	deleteErr := f.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.StorageResourceAction, resourceID), nil, nil)
	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), AttachedSnapshotsErrorCode) {
			err := f.updateDescription(ctx, filesystemID, MarkFilesystemForDeletion)
			if err != nil {
				return fmt.Errorf("mark filesystem %s for deletion failed. Error: %v", filesystemID, err)
			}
			return nil
		}
		return fmt.Errorf("delete Filesystem %s Failed. Error: %v", filesystemID, deleteErr)
	}
	log.Debugf("Delete Filesystem %s Successful", filesystemID)
	return nil
}

//Update description of filesystem
func (f *Filesystem) updateDescription(ctx context.Context, filesystemID, description string) error {
	if len(filesystemID) == 0 {
		return errors.New("Filesystem Id cannot be empty")
	}
	filesystemResp, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return err
	}
	resourceID := filesystemResp.FileContent.StorageResource.ID

	filesystemModifyParam := types.FsModifyParameters{
		Description: description,
	}
	err = f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, resourceID), filesystemModifyParam, nil)
	if err != nil {
		return fmt.Errorf("update filesystem: %s description failed with error: %v", resourceID, err)
	}
	return nil
}

//Updates IsReplicationDestination parameter of filesystem
func (f *Filesystem) UpdateReplicationDestinationParameter(ctx context.Context, resourceID string, isReplicationDestination bool) error {
	log := util.GetRunIDLogger(ctx)
	log.Debugf("Updating Filesystem %s, isReplicationDestination parameter is %v", resourceID, isReplicationDestination)
	if len(resourceID) == 0 {
		return errors.New("Filesystem Id cannot be empty")
	}

	filesystemModifyParam := types.FsModifyParameters{
		ReplicationParameters: types.ReplicationParameters{
			isReplicationDestination,
		},
	}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, resourceID), filesystemModifyParam, nil)
	if err != nil {
		return fmt.Errorf("update filesystem: %s isReplicationDestination failed with error: %v", resourceID, err)
	}
	return nil
}

//CreateNFSShare - Create NFS Share for a File system
func (f *Filesystem) CreateNFSShare(ctx context.Context, name, path, filesystemID string, nfsShareDefaultAccess NFSShareDefaultAccess) (*types.Filesystem, error) {
	if len(filesystemID) == 0 {
		return nil, errors.New("Filesystem Id cannot be empty")
	}

	filesystemResp, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return nil, err
	}
	resourceID := filesystemResp.FileContent.StorageResource.ID

	nfsShareParam := types.NFSShareParameters{
		DefaultAccess: string(nfsShareDefaultAccess),
	}

	nfsShareCreateReqParam := types.NFSShareCreateParam{
		Name:               name,
		Path:               path,
		NFSShareParameters: &nfsShareParam,
	}

	nfsShares := []types.NFSShareCreateParam{nfsShareCreateReqParam}
	filesystemModifyParam := types.FsModifyParameters{
		NFSShares: &nfsShares,
	}

	err = f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, resourceID), filesystemModifyParam, nil)
	if err != nil {
		return nil, fmt.Errorf("create NFS Share failed. Error: %v", err)
	}

	filesystemResp, err = f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return nil, ErrorFilesystemNotFound
	}
	return filesystemResp, nil
}

//CreateNFSShareFromSnapshot - Create NFS Share for a File system Snapshot
func (f *Filesystem) CreateNFSShareFromSnapshot(ctx context.Context, name, path, snapshotID string, nfsShareDefaultAccess NFSShareDefaultAccess) (*types.NFSShare, error) {
	if len(snapshotID) == 0 {
		return nil, errors.New("Snapshot Id cannot be empty")
	}

	snapshotContent := types.SnapshotIDContent{
		ID: snapshotID,
	}

	nfsShareCreateReq := types.NFSShareCreateFromSnapParam{
		Name:          name,
		Path:          path,
		DefaultAccess: string(nfsShareDefaultAccess),
		Snapshot:      snapshotContent,
	}

	nfsShareResp := &types.NFSShare{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.NfsShareAction), nfsShareCreateReq, nfsShareResp)
	if err != nil {
		return nil, fmt.Errorf("create NFS Share: %s failed. Error: %v", name, err)
	}

	return nfsShareResp, nil
}

//FindNFSShareByName - Find the NFS Share by it's name. If the NFS Share is not found, an error will be returned.
func (f *Filesystem) FindNFSShareByName(ctx context.Context, nfsSharename string) (*types.NFSShare, error) {
	if len(nfsSharename) == 0 {
		return nil, errors.New("NFS Share Name shouldn't be empty")
	}
	nfsShareResp := &types.NFSShare{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.NfsShareAction, nfsSharename, NFSShareDisplayfields), nil, nfsShareResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find NFS Share. Error: %v", err)
	}
	return nfsShareResp, nil
}

//FindNFSShareByID - Find the NFS Share by it's Id. If the NFS Share is not found, an error will be returned.
func (f *Filesystem) FindNFSShareByID(ctx context.Context, nfsShareID string) (*types.NFSShare, error) {
	if len(nfsShareID) == 0 {
		return nil, errors.New("NFS Share Id shouldn't be empty")
	}
	nfsShareResp := &types.NFSShare{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.NfsShareAction, nfsShareID, NFSShareDisplayfields), nil, nfsShareResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find NFS Share: %s. Error: %v", nfsShareID, err)
	}
	return nfsShareResp, nil
}

//ModifyNFSShareHostAccess - Modify the host access on NFS Share
func (f *Filesystem) ModifyNFSShareHostAccess(ctx context.Context, filesystemID, nfsShareID string, hostIDs []string, accessType AccessType) error {
	log := util.GetRunIDLogger(ctx)
	if len(filesystemID) == 0 {
		return errors.New("Filesystem Id cannot be empty")
	}

	filesystemResp, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return ErrorFilesystemNotFound
	}
	resourceID := filesystemResp.FileContent.StorageResource.ID

	hostsIdsContent := []types.HostIDContent{}
	for _, hostID := range hostIDs {
		hostIDContent := types.HostIDContent{
			ID: hostID,
		}
		hostsIdsContent = append(hostsIdsContent, hostIDContent)
	}

	nfsShareParameters := types.NFSShareParameters{}
	if accessType == ReadOnlyAccessType {
		nfsShareParameters.ReadOnlyHosts = &hostsIdsContent
	} else if accessType == ReadWriteAccessType {
		nfsShareParameters.ReadWriteHosts = &hostsIdsContent
	} else if accessType == ReadOnlyRootAccessType {
		nfsShareParameters.ReadOnlyRootAccessHosts = &hostsIdsContent
	} else if accessType == ReadWriteRootAccessType {
		nfsShareParameters.RootAccessHosts = &hostsIdsContent
	}

	nfsShare := types.StorageResourceParam{
		ID: nfsShareID,
	}

	nfsShareModifyContent := types.NFSShareModifyContent{
		NFSShare:           &nfsShare,
		NFSShareParameters: &nfsShareParameters,
	}
	nfsSharesModifyContent := []types.NFSShareModifyContent{nfsShareModifyContent}

	nfsShareModifyReq := types.NFSShareModify{
		NFSSharesModifyContent: &nfsSharesModifyContent,
	}

	err = f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, resourceID), nfsShareModifyReq, nil)
	if err != nil {
		return fmt.Errorf("modify NFS Share failed. Error: %v", err)
	}
	log.Debugf("Modify NFS share: %s successful. Added host with access %s", nfsShareID, accessType)
	return nil
}

//ModifyNFSShareCreatedFromSnapshotHostAccess - Modify the host access on NFS Share
func (f *Filesystem) ModifyNFSShareCreatedFromSnapshotHostAccess(ctx context.Context, nfsShareID string, hostIDs []string, accessType AccessType) error {
	if nfsShareID == "" {
		return errors.New("NFS Share Id cannot be empty")
	}

	hostsIdsContent := []types.HostIDContent{}
	for _, hostID := range hostIDs {
		hostIDContent := types.HostIDContent{
			ID: hostID,
		}
		hostsIdsContent = append(hostsIdsContent, hostIDContent)
	}

	nfsShareModifyReq := types.NFSShareCreateFromSnapModify{}

	if accessType == ReadOnlyAccessType {
		nfsShareModifyReq.ReadOnlyHosts = &hostsIdsContent
	} else if accessType == ReadWriteAccessType {
		nfsShareModifyReq.ReadWriteHosts = &hostsIdsContent
	} else if accessType == ReadOnlyRootAccessType {
		nfsShareModifyReq.ReadOnlyRootAccessHosts = &hostsIdsContent
	} else if accessType == ReadWriteRootAccessType {
		nfsShareModifyReq.RootAccessHosts = &hostsIdsContent
	}

	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyNFSShareURI, api.NfsShareAction, nfsShareID), nfsShareModifyReq, nil)
	if err != nil {
		return fmt.Errorf("modify NFS Share %s failed. Error: %v", nfsShareID, err)
	}
	return nil
}

//DeleteNFSShare by its ID. If the NFSShare is not present on the array, an error will be returned.
func (f *Filesystem) DeleteNFSShare(ctx context.Context, filesystemID, nfsShareID string) error {
	log := util.GetRunIDLogger(ctx)

	if len(filesystemID) == 0 {
		return errors.New("Filesystem Id cannot be empty")
	}
	filesystemResp, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return ErrorFilesystemNotFound
	}
	resourceID := filesystemResp.FileContent.StorageResource.ID

	if len(nfsShareID) == 0 {
		return errors.New("NFS Share Id cannot be empty")
	}
	_, err = f.FindNFSShareByID(ctx, nfsShareID)
	if err != nil {
		return fmt.Errorf("unable to find NFS Share. Error: %v", err)
	}

	nfsShare := types.StorageResourceParam{
		ID: nfsShareID,
	}

	nfsShareDeleteContent := types.NFSShareModifyContent{
		NFSShare: &nfsShare,
	}
	nfsSharesDeleteContent := []types.NFSShareModifyContent{nfsShareDeleteContent}

	nfsShareDeleteReq := types.NFSShareDelete{
		NFSSharesDeleteContent: &nfsSharesDeleteContent,
	}

	deleteErr := f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, resourceID), nfsShareDeleteReq, nil)
	if deleteErr != nil {
		return fmt.Errorf("delete NFS Share: %s Failed. Error: %v", nfsShareID, deleteErr)
	}
	log.Infof("Delete NFS Share: %s Successful", nfsShareID)
	return nil
}

//DeleteNFSShareCreatedFromSnapshot by its ID. If the NFSShare is not present on the array, an error will be returned.
func (f *Filesystem) DeleteNFSShareCreatedFromSnapshot(ctx context.Context, nfsShareID string) error {
	if len(nfsShareID) == 0 {
		return errors.New("NFS Share Id cannot be empty")
	}

	_, err := f.FindNFSShareByID(ctx, nfsShareID)
	if err != nil {
		return fmt.Errorf("unable to find NFS Share %s. Error: %v", nfsShareID, err)
	}

	err = f.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.NfsShareAction, nfsShareID), nil, nil)
	if err != nil {
		return fmt.Errorf("delete NFS Share: %s Failed. Error: %v", nfsShareID, err)
	}
	return nil
}

//FindNASServerByID - Find the NAS Server by it's Id. If the NAS Server is not found, an error will be returned.
func (f *Filesystem) FindNASServerByID(ctx context.Context, nasServerID string) (*types.NASServer, error) {
	if len(nasServerID) == 0 {
		return nil, errors.New("NAS Server Id shouldn't be empty")
	}
	nasServerResp := &types.NASServer{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.NasServerAction, nasServerID, NasServerDisplayfields), nil, nasServerResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find NAS Server: %s. Error: %v", nasServerID, err)
	}
	return nasServerResp, nil
}

//ExpandFilesystem Filesystem Expand volume to provided capacity
func (f *Filesystem) ExpandFilesystem(ctx context.Context, filesystemID string, newSize uint64) error {
	log := util.GetRunIDLogger(ctx)
	filesystem, err := f.FindFilesystemByID(ctx, filesystemID)
	if err != nil {
		return fmt.Errorf("unable to find filesystem Id %s. Error: %v", filesystemID, err)
	}
	if filesystem.FileContent.SizeTotal == newSize {
		log.Infof("New Volume size (%d) is same as existing Volume size (%d). Ignoring expand volume operation.", newSize, filesystem.FileContent.SizeTotal)
		return nil
	} else if filesystem.FileContent.SizeTotal > newSize {
		return fmt.Errorf("requested new capacity smaller than existing capacity")
	}
	fsExpandParams := types.FsExpandParameters{
		Size: newSize,
	}
	fsExpandReqParam := types.FsExpandModifyParam{
		FsParameters: &fsExpandParams,
	}
	return f.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyFilesystemURI, filesystem.FileContent.StorageResource.ID), fsExpandReqParam, nil)
}

func (f *Filesystem) FindFileSystemGroupByPrefix(ctx context.Context, prefix string) (*types.ListFileSystem, error) {
	log := util.GetRunIDLogger(ctx)
	if len(prefix) == 0 {
		return nil, fmt.Errorf("Filesystem prefix cannot be empty")
	}

	filter := fmt.Sprintf("name lk %s", "\""+prefix+"%\"")
	queryURI := fmt.Sprintf(api.UnityInstancesFilterWithFields, api.StorageResourceAction, StorageResourceDisplayFields, url.QueryEscape(filter))
	listFileSystems := &types.ListFileSystem{}
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, listFileSystems)
	if err != nil {
		return nil, err
	}
	if len(listFileSystems.Filesystems) == 0 {
		log.Info("List of File Systems is empty")
		return nil, nil
	}
	return listFileSystems, nil
}
