/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"context"
	"errors"
	"fmt"
	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
	"net/http"
	"strconv"
)

type LicenseType string

const (
	LunNameMaxLength             = 63
	ThinProvisioning LicenseType = "THIN_PROVISIONING"
	DataReduction    LicenseType = "DATA_REDUCTION"
)

var VolumeNotFoundError = errors.New("Unable to find volume")

type volume struct {
	client *Client
}

func NewVolume(client *Client) *volume {
	return &volume{client}
}

// CreateLun API create a Lun with the given arguments.
// Pre-validations: 1. Length of the Lun name should be less than 63 characters.
//                  2. Size of Lun should be in bytes.
func (v *volume) CreateLun(ctx context.Context, name, poolId, description string, size uint64, fastVPTieringPolicy int,
	hostIOLimitID string, isThinEnabled, isDataReductionEnabled bool) (*types.Volume, error) {
	log := util.GetRunIdLogger(ctx)

	if name == "" {
		return nil, errors.New("lun name should not be empty.")
	}

	if len(name) > LunNameMaxLength {
		return nil, errors.New(fmt.Sprintf("lun name %s should not exceed 63 characters.", name))
	}

	poolApi := NewStoragePool(v.client)
	pool, err := poolApi.FindStoragePoolById(ctx, poolId)
	if pool == nil {
		return nil, errors.New(fmt.Sprintf("unable to get PoolID (%s) Error:%v", poolId, err))
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error trying to get Storage Pool (%s) Error:%v", poolId, err))
	}

	storagePool := types.StoragePoolID{
		PoolId: pool.StoragePoolContent.ID,
	}

	lunParams := types.LunParameters{
		StoragePool: &storagePool,
		Size:        size,
	}

	thinProvisioningLicenseInfoResp, err := v.isFeatureLicensed(ctx, ThinProvisioning)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get license info for feature: %s", ThinProvisioning))
	}

	dataReductionLicenseInfoResp, err := v.isFeatureLicensed(ctx, DataReduction)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get license info for feature: %s", DataReduction))
	}

	if thinProvisioningLicenseInfoResp.LicenseInfoContent.IsInstalled && thinProvisioningLicenseInfoResp.LicenseInfoContent.IsValid {
		lunParams.IsThinEnabled = strconv.FormatBool(isThinEnabled)
	} else if isThinEnabled == true {
		return nil, errors.New(fmt.Sprintf("Thin Provisioning is not supported on array and hence cannot create Volume."))
	}

	if dataReductionLicenseInfoResp.LicenseInfoContent.IsInstalled && dataReductionLicenseInfoResp.LicenseInfoContent.IsValid {
		lunParams.IsDataReductionEnabled = strconv.FormatBool(isDataReductionEnabled)
	} else if isDataReductionEnabled == true {
		return nil, errors.New(fmt.Sprintf("Data Reduction is not supported on array and hence cannot create Volume."))
	}

	if hostIOLimitID != "" {
		ioLimitPolicyParam := types.IoLimitPolicyParam{
			Id: hostIOLimitID,
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
	}

	volumeReqParam := types.LunCreateParam{
		Name:          name,
		Description:   description,
		LunParameters: &lunParams,
	}

	volumeResp := &types.Volume{}
	err = v.client.executeWithRetryAuthenticate(ctx,
		http.MethodPost, fmt.Sprintf(api.UnityApiStorageResourceActionUri, "createLun"), volumeReqParam, volumeResp)
	if err != nil {
		return nil, err
	}
	return volumeResp, nil
}

//Find the volume by it's name. If the volume is not found, an error will be returned.
func (v *volume) FindVolumeByName(ctx context.Context, volName string) (*types.Volume, error) {
	if len(volName) == 0 {
		return nil, errors.New(fmt.Sprintf("Lun Name shouldn't be empty."))
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "lun", volName, api.LunDisplayFields), nil, volumeResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find volume by name %s", volName))
	}

	return volumeResp, nil
}

//Find the volume by it's Id. If the volume is not found, an error will be returned.
func (v *volume) FindVolumeById(ctx context.Context, volId string) (*types.Volume, error) {
	log := util.GetRunIdLogger(ctx)
	if len(volId) == 0 {
		return nil, errors.New("lun ID shouldn't be empty")
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "lun", volId, api.LunDisplayFields), nil, volumeResp)
	if err != nil {
		log.Debugf("Unable to find volume Id %s Error: %v", volId, err)
		return nil, VolumeNotFoundError
	}
	return volumeResp, nil
}

func (v *volume) ListVolumes(ctx context.Context, startToken int, maxEntries int) ([]types.Volume, int, error) {
	log := util.GetRunIdLogger(ctx)
	volumeResp := &types.ListVolumes{}
	nextToken := startToken + 1
	lunUri := fmt.Sprintf(api.UnityApiInstanceTypeResourcesWithFields, "lun", api.LunDisplayFields)

	if maxEntries != 0 {
		lunUri = fmt.Sprintf(lunUri+"&per_page=%d", maxEntries)

		//startToken should exists only when maxEntries are present
		if startToken != 0 {
			lunUri = fmt.Sprintf(lunUri+"&page=%d", startToken)
		}
	}

	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, lunUri, nil, volumeResp)
	if err != nil {
		log.Errorf("executeWithRetryAuthenticate Error: %v", err)
	}
	return volumeResp.Volumes, nextToken, err
}

//Delete Volume by its ID. If the Volume is not present on the array, an error will be returned.
func (v *volume) DeleteVolume(ctx context.Context, volumeId string) error {
	log := util.GetRunIdLogger(ctx)
	if len(volumeId) == 0 {
		return errors.New("Volume Id cannot be empty")
	}
	volumeResp := &types.Volume{}

	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceUri, "storageResource", volumeId), nil, volumeResp)
	if err != nil {
		return VolumeNotFoundError
	} else {
		deleteErr := v.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceUri, "storageResource", volumeId), nil, nil)
		if deleteErr != nil {
			return errors.New(fmt.Sprintf("Delete Volume %s Failed. Error: %v", volumeId, deleteErr))
		}
		log.Debug("Delete Storage Resource %s Successful", volumeId)
		return nil
	}
}

//Export volume to a host
func (v *volume) ExportVolume(ctx context.Context, volID, hostID string) error {
	hostIdContent := types.HostIdContent{
		ID: hostID,
	}

	hostAccess := types.HostAccess{
		HostIdContent: &hostIdContent,
		AccessMask:    "1", //Hardcoded as 1 so that the host can have access to production LUNs only.
	}
	hostAccessArray := []types.HostAccess{hostAccess}
	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunUri, volID), lunModifyParam, nil)
}

//Unexport volume
func (v *volume) UnexportVolume(ctx context.Context, volID string) error {
	hostAccessArray := []types.HostAccess{}
	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyLunUri, volID), lunModifyParam, nil)
}

// Expand volume to provided capacity
func (v *volume) ExpandVolume(ctx context.Context, volumeId string, newSize uint64) error {
	log := util.GetRunIdLogger(ctx)
	vol, err := v.FindVolumeById(ctx, volumeId)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to find volume Id %s Error: %v", volumeId, err))
	}
	if vol.VolumeContent.SizeTotal == newSize {
		log.Infof("New Volume size (%d) is same as existing Volume size(%d). Ignoring expand volume operation.", newSize, vol.VolumeContent.SizeTotal)
		return nil
	} else if vol.VolumeContent.SizeTotal > newSize {
		return errors.New(fmt.Sprintf("requested new capacity smaller than existing capacity"))
	}
	lunParams := types.LunExpandParameters{
		Size: newSize,
	}
	volumeReqParam := types.LunExpandModifyParam{
		LunParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiModifyLunUri, volumeId), volumeReqParam, nil)
}

func (v *volume) FindHostIOLimitByName(ctx context.Context, hostIoPolicyName string) (*types.IoLimitPolicy, error) {
	if len(hostIoPolicyName) == 0 {
		return nil, errors.New("policy Name shouldn't be empty")
	}
	ioLimitPolicyResp := &types.IoLimitPolicy{}
	fieldsToQuery := "id,name,description"
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "ioLimitPolicy", hostIoPolicyName, fieldsToQuery), nil, ioLimitPolicyResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find IO Limit Policy:%s Error: %v", hostIoPolicyName, err))
	}
	return ioLimitPolicyResp, nil
}

// Create a lun thin clone
func (v *volume) CreteLunThinClone(ctx context.Context, name, snapId, volId string) (*types.Volume, error) {
	snapIdContent := types.SnapshotIdContent{
		Id: snapId,
	}
	createLunThinCloneParam := types.CreateLunThinCloneParam{
		SnapIdContent: &snapIdContent,
		Name:          name,
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiCreateLunThinCloneUri, volId), createLunThinCloneParam, volumeResp)
	return volumeResp, err
}

//Get License information
func (v *volume) isFeatureLicensed(ctx context.Context, featureName LicenseType) (*types.LicenseInfo, error) {
	licenseInfoResp := &types.LicenseInfo{}
	err := v.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "license", featureName, api.LicenseInfoDisplayFields), nil, licenseInfoResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get license info for feature: %s", featureName))
	}
	return licenseInfoResp, nil
}
