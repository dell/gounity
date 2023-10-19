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
	"strconv"
	"strings"
	"time"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

// LicenseType is string
type LicenseType string

// Constants
const (
	LunNameMaxLength             = 63
	SnapForClone                 = "csi-snapforclone-"
	ThinProvisioning LicenseType = "THIN_PROVISIONING"
	DataReduction    LicenseType = "DATA_REDUCTION"
)

// DependentClonesErrorCode stores error code of dependent clones
var DependentClonesErrorCode = "0x6701673"

// ErrorDependentClones stores dependent clones error message
var ErrorDependentClones = errors.New("the specified volume cannot be deleted because it has one or more dependent thin clones")

// VolumeNotFoundErrorCode stores Volume not found error code
var VolumeNotFoundErrorCode = "0x7d13005"

// ErrorVolumeNotFound stores Volume not found error
var ErrorVolumeNotFound = errors.New("Unable to find volume")

// ErrorCreateSnapshotFailed stores Create snapshot failed error message
var ErrorCreateSnapshotFailed = errors.New("create Snapshot Failed")

// ErrorCloningFailed stores Cloning failed error message
var ErrorCloningFailed = errors.New("volume Cloning Failed")

// MarkVolumeForDeletion stores mark of volume deletion
var MarkVolumeForDeletion = "csi-marked-vol-for-deletion"

// Volume structure
type Volume struct {
	client *Client
}

// NewVolume function returns volume
func NewVolume(client *Client) *Volume {
	return &Volume{client}
}

// CreateLun API create a Lun with the given arguments.
// Pre-validations: 1. Length of the Lun name should be less than 63 characters.
//  2. Size of Lun should be in bytes.
func (v *Volume) CreateLun(ctx context.Context, name, poolID, description string, size uint64, fastVPTieringPolicy int,
	hostIOLimitID string, isThinEnabled, isDataReductionEnabled bool) (*types.Volume, error) {
	log := util.GetRunIDLogger(ctx)

	if name == "" {
		return nil, errors.New("lun name should not be empty")
	}

	if len(name) > LunNameMaxLength {
		return nil, fmt.Errorf("lun name %s should not exceed 63 characters", name)
	}

	poolAPI := NewStoragePool(v.client)
	pool, err := poolAPI.FindStoragePoolByID(ctx, poolID)

	if err != nil {
		return nil, fmt.Errorf("unable to get PoolID (%s) Error:%v", poolID, err)
	}

	storagePool := types.StoragePoolID{
		PoolID: pool.StoragePoolContent.ID,
	}

	lunParams := types.LunParameters{
		StoragePool: &storagePool,
		Size:        size,
	}

	thinProvisioningLicenseInfoResp, err := v.isFeatureLicensed(ctx, ThinProvisioning)
	if err != nil {
		return nil, fmt.Errorf("unable to get license info for feature: %s", ThinProvisioning)
	}

	dataReductionLicenseInfoResp, err := v.isFeatureLicensed(ctx, DataReduction)
	if err != nil {
		return nil, fmt.Errorf("unable to get license info for feature: %s", DataReduction)
	}

	if thinProvisioningLicenseInfoResp.LicenseInfoContent.IsInstalled && thinProvisioningLicenseInfoResp.LicenseInfoContent.IsValid {
		lunParams.IsThinEnabled = strconv.FormatBool(isThinEnabled)
	} else if isThinEnabled == true {
		return nil, fmt.Errorf("thin Provisioning is not supported on array and hence cannot create Volume")
	}

	if dataReductionLicenseInfoResp.LicenseInfoContent.IsInstalled && dataReductionLicenseInfoResp.LicenseInfoContent.IsValid {
		lunParams.IsDataReductionEnabled = strconv.FormatBool(isDataReductionEnabled)
	} else if isDataReductionEnabled == true {
		return nil, fmt.Errorf("data Reduction is not supported on array and hence cannot create Volume")
	}

	if hostIOLimitID != "" {
		ioLimitPolicyParam := types.IoLimitPolicyParam{
			ID: hostIOLimitID,
		}
		ioLimitParameters := types.HostIoLimitParameters{
			IoLimitPolicyParam: &ioLimitPolicyParam,
		}

		lunParams.IoLimitParameters = &ioLimitParameters
	}

	if pool != nil && pool.StoragePoolContent.PoolFastVP.Status != 0 {
		log.Debug("FastVP is enabled")
		fastVPParameters := types.FastVPParameters{
			TieringPolicy: fastVPTieringPolicy,
		}
		lunParams.FastVPParameters = &fastVPParameters
	} else {
		log.Debug("FastVP is not enabled")
		if fastVPTieringPolicy != 0 {
			return nil, fmt.Errorf("fastVP is not enabled and requested tiering policy is: %d ", fastVPTieringPolicy)
		}
	}

	volumeReqParam := types.LunCreateParam{
		Name:          name,
		Description:   description,
		LunParameters: &lunParams,
	}

	volumeResp := &types.Volume{}
	err = v.client.executeWithRetryAuthenticate(ctx,
		http.MethodPost, fmt.Sprintf(api.UnityAPIStorageResourceActionURI, api.CreateLunAction), volumeReqParam, volumeResp)
	if err != nil {
		return nil, err
	}
	return volumeResp, nil
}

// FindVolumeByName - Find the volume by it's name. If the volume is not found, an error will be returned.
func (v *Volume) FindVolumeByName(ctx context.Context, volName string) (*types.Volume, error) {
	if len(volName) == 0 {
		return nil, fmt.Errorf("lun Name shouldn't be empty")
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.LunAction, volName, LunDisplayFields), nil, volumeResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find volume by name %s", volName)
	}

	return volumeResp, nil
}

// FindVolumeByID - Find the volume by it's Id. If the volume is not found, an error will be returned.
func (v *Volume) FindVolumeByID(ctx context.Context, volID string) (*types.Volume, error) {
	log := util.GetRunIDLogger(ctx)
	if len(volID) == 0 {
		return nil, errors.New("lun ID shouldn't be empty")
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.LunAction, volID, LunDisplayFields), nil, volumeResp)
	if err != nil {
		if strings.Contains(err.Error(), VolumeNotFoundErrorCode) {
			log.Debugf("Unable to find volume Id %s Error: %v", volID, err)
			return nil, ErrorVolumeNotFound
		}
		return nil, err
	}
	return volumeResp, nil
}

// ListVolumes - list volumes
func (v *Volume) ListVolumes(ctx context.Context, startToken int, maxEntries int) ([]types.Volume, int, error) {
	log := util.GetRunIDLogger(ctx)
	volumeResp := &types.ListVolumes{}
	nextToken := startToken + 1
	lunURI := fmt.Sprintf(api.UnityAPIInstanceTypeResourcesWithFields, api.LunAction, LunDisplayFields)

	if maxEntries != 0 {
		lunURI = fmt.Sprintf(lunURI+"&per_page=%d", maxEntries)

		//startToken should exists only when maxEntries are present
		if startToken != 0 {
			lunURI = fmt.Sprintf(lunURI+"&page=%d", startToken)
		}
	}

	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, lunURI, nil, volumeResp)
	if err != nil {
		log.Errorf("executeWithRetryAuthenticate Error: %v", err)
	}
	return volumeResp.Volumes, nextToken, err
}

// DeleteVolume - Delete Volume by its ID. If the Volume is not present on the array, an error will be returned.
func (v *Volume) DeleteVolume(ctx context.Context, volumeID string) error {
	log := util.GetRunIDLogger(ctx)
	if len(volumeID) == 0 {
		return errors.New("Volume Id cannot be empty")
	}
	volumeResp := &types.Volume{}

	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceURI, api.StorageResourceAction, volumeID), nil, volumeResp)

	volResp, err := v.FindVolumeByID(ctx, volumeID)

	if err != nil {
		return err
	}
	sourceVolID := ""
	if volResp.VolumeContent.IsThinClone {
		//Check if parent volume is marked for deletion
		sourceVolID = volResp.VolumeContent.ParentVolume.ID
	}

	deleteErr := v.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.StorageResourceAction, volumeID), nil, nil)

	deleteSourceVol := false
	if sourceVolID != "" {
		sourceVolResp, err := v.FindVolumeByID(ctx, sourceVolID)
		if err != nil && err != ErrorVolumeNotFound {
			return fmt.Errorf("find Source Volume %s Failed. Error: %v", sourceVolID, err)
		}
		if strings.Contains(sourceVolResp.VolumeContent.Name, MarkVolumeForDeletion) {
			deleteSourceVol = true
		}
	}
	if deleteSourceVol {
		deleteSourceErr := v.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.StorageResourceAction, volResp.VolumeContent.ParentVolume.ID), nil, nil)
		if deleteSourceErr != nil {
			log.Warnf("Deletion of source volume: %s marked for deletion failed with error: %v", volResp.VolumeContent.ParentVolume.ID, deleteSourceErr)
		} else {
			log.Debugf("Deletion of source volume: %s marked for deletion successful", volResp.VolumeContent.ParentVolume.ID)
		}
	}
	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), DependentClonesErrorCode) {
			newName := MarkVolumeForDeletion + strconv.FormatInt(time.Now().Unix(), 10)
			err := v.RenameVolume(ctx, newName, volumeID)
			if err != nil {
				//Unable to mark volume for deletion
				log.Warnf("Unable to mark volume %s with dependent clones for deletion", volumeID)
			} else {
				log.Debugf("Volume %s has dependent clones and marked for deletion.", volumeID)
			}
			return nil
		}
		return fmt.Errorf("delete Volume %s Failed. Error: %v", volumeID, deleteErr)
	}
	log.Debugf("Delete Storage Resource %s Successful", volumeID)
	return nil
}

// ExportVolume - Export volume to a host
func (v *Volume) ExportVolume(ctx context.Context, volID, hostID string) error {
	hostIDContent := types.HostIDContent{
		ID: hostID,
	}

	hostAccess := types.HostAccess{
		HostIDContent: &hostIDContent,
		AccessMask:    "1", //Hardcoded as 1 so that the host can have access to production LUNs only.
	}
	hostAccessArray := []types.HostAccess{hostAccess}
	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunURI, volID), lunModifyParam, nil)
}

// ModifyVolumeExport - Export volume to multiple hosts / Modify the host access list on a given Volume
func (v *Volume) ModifyVolumeExport(ctx context.Context, volID string, hostIDList []string) error {

	hostAccessArray := []types.HostAccess{}
	for _, hostID := range hostIDList {
		hostIDContent := types.HostIDContent{
			ID: hostID,
		}

		hostAccess := types.HostAccess{
			HostIDContent: &hostIDContent,
			AccessMask:    "1", //Hardcoded as 1 so that the host can have access to production LUNs only.
		}
		hostAccessArray = append(hostAccessArray, hostAccess)
	}

	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunURI, volID), lunModifyParam, nil)
}

// UnexportVolume - Unexport volume
func (v *Volume) UnexportVolume(ctx context.Context, volID string) error {
	hostAccessArray := []types.HostAccess{}
	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunURI, volID), lunModifyParam, nil)
}

// ExpandVolume - Expand volume to provided capacity
func (v *Volume) ExpandVolume(ctx context.Context, volumeID string, newSize uint64) error {
	log := util.GetRunIDLogger(ctx)
	vol, err := v.FindVolumeByID(ctx, volumeID)
	if err != nil {
		return fmt.Errorf("unable to find volume Id %s Error: %v", volumeID, err)
	}
	if vol.VolumeContent.SizeTotal == newSize {
		log.Infof("New Volume size (%d) is same as existing Volume size(%d). Ignoring expand volume operation.", newSize, vol.VolumeContent.SizeTotal)
		return nil
	} else if vol.VolumeContent.SizeTotal > newSize {
		return fmt.Errorf("requested new capacity smaller than existing capacity")
	}
	lunParams := types.LunExpandParameters{
		Size: newSize,
	}
	volumeReqParam := types.LunExpandModifyParam{
		LunParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIModifyLunURI, volumeID), volumeReqParam, nil)
}

// FindHostIOLimitByName - Find Host IO limit
func (v *Volume) FindHostIOLimitByName(ctx context.Context, hostIoPolicyName string) (*types.IoLimitPolicy, error) {
	if len(hostIoPolicyName) == 0 {
		return nil, errors.New("policy Name shouldn't be empty")
	}
	ioLimitPolicyResp := &types.IoLimitPolicy{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.IOLimitPolicy, hostIoPolicyName, HostIOLimitFields), nil, ioLimitPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find IO Limit Policy:%s Error: %v", hostIoPolicyName, err)
	}
	return ioLimitPolicyResp, nil
}

// CreteLunThinClone - Create a lun thin clone
func (v *Volume) CreteLunThinClone(ctx context.Context, name, snapID, volID string) (*types.Volume, error) {
	snapIDContent := types.SnapshotIDContent{
		ID: snapID,
	}
	createLunThinCloneParam := types.CreateLunThinCloneParam{
		SnapIDContent: &snapIDContent,
		Name:          name,
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPICreateLunThinCloneURI, volID), createLunThinCloneParam, volumeResp)
	return volumeResp, err
}

// isFeatureLicensed - Get License information
func (v *Volume) isFeatureLicensed(ctx context.Context, featureName LicenseType) (*types.LicenseInfo, error) {
	licenseInfoResp := &types.LicenseInfo{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.LicenseAction, featureName, LicenseInfoDisplayFields), nil, licenseInfoResp)
	if err != nil {
		return nil, fmt.Errorf("unable to get license info for feature: %s", featureName)
	}
	return licenseInfoResp, nil
}

// CreateCloneFromVolume - Volume cloning
func (v *Volume) CreateCloneFromVolume(ctx context.Context, name, volID string) (*types.Volume, error) {
	log := util.GetRunIDLogger(ctx)
	snapAPI := NewSnapshot(v.client)
	//Create snapshot for cloning
	snapName := SnapForClone + strconv.FormatInt(time.Now().Unix(), 10)
	snapResp, err := snapAPI.CreateSnapshot(ctx, volID, snapName, "", "")
	if err != nil {
		return nil, ErrorCreateSnapshotFailed
	}
	//Clone Volume
	cloned := true
	volResp, err := v.CreteLunThinClone(ctx, name, snapResp.SnapshotContent.ResourceID, volID)
	if err != nil {
		cloned = false
	}
	//Delete Snapshot
	err = snapAPI.DeleteSnapshot(ctx, snapResp.SnapshotContent.ResourceID)
	if err != nil {
		//If delete snapshot created to clone volume failed then error is only logged not returned
		log.Warnf("Unable to Delete Snapshot: %s created to clone Volume: %s", snapName, volID)
	}
	if !cloned {
		return nil, ErrorCloningFailed
	}
	return volResp, nil
}

// RenameVolume - Rename Volume
func (v *Volume) RenameVolume(ctx context.Context, newName, volID string) error {
	lunParams := types.LunParameters{
		Name: newName,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunURI, volID), lunParams, nil)
}

// GetMaxVolumeSize - Returns the max size of a volume supported by the array
func (v *Volume) GetMaxVolumeSize(ctx context.Context, systemLimitID string) (*types.MaxVolumSizeInfo, error) {
	volumeResp := &types.MaxVolumSizeInfo{}
	if len(systemLimitID) == 0 {
		return nil, errors.New("system limit ID shouldn't be empty")
	}
	lunURI := fmt.Sprintf(api.UnityAPIGetMaxVolumeSize, systemLimitID, api.Limit, api.Unit)
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, lunURI, nil, volumeResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find system limit by ID %s", systemLimitID)
	}

	return volumeResp, nil
}
