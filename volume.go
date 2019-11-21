/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/payloads"
)

const (
	LunNameMaxLength = 63
)

var NotFoundError = errors.New("Unable to find volume")

type volume struct {
	client *Client
}

func NewVolume(client *Client) *volume {
	return &volume{client}
}

// CreateLun API create a Lun with the given arguments.
// Pre-validations: 1. Length of the Lun name should be less than 63 characters.
//                  2. Size of Lun should be in bytes.
func (v *volume) CreateLun(name, poolId, description string, size uint64, fastVPTieringPolicy int,
	hostIOLimitID string, isThinEnabled, isDataReductionEnabled bool) (*types.Volume, error) {

	if len(name) > LunNameMaxLength {
		return nil, errors.New(fmt.Sprintf("lun name %s should not exceed 63 characters.", name))
	}

	poolApi := NewStoragePool(v.client)
	pool, err := poolApi.FindStoragePoolById(poolId)
	if pool == nil {
		log.Errorf("Unable to get PoolID (%s) Error:%v", poolId, err)
		return nil, errors.New(fmt.Sprintf("unable to get PoolID (%s) Error:%v", poolId, err))
	}

	storagePool := types.StoragePoolID{
		PoolId: pool.StoragePoolContent.ID,
	}

	lunParams := types.LunParameters{
		StoragePool:            &storagePool,
		Size:                   size,
		IsThinEnabled:          isThinEnabled,
		IsDataReductionEnabled: isDataReductionEnabled,
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
		log.Info("FastVP is enabled")
		fastVPParameters := types.FastVPParameters{
			TieringPolicy: fastVPTieringPolicy,
		}
		lunParams.FastVPParameters = &fastVPParameters
	} else {
		log.Info("FastVP is not enabled")
	}

	volumeReqParam := types.LunCreateParam{
		Name:          name,
		Description:   description,
		LunParameters: &lunParams,
	}

	volumeResp := &types.Volume{}
	err = v.client.executeWithRetryAuthenticate(
		http.MethodPost, api.UnityApiCreateLunUri, volumeReqParam, volumeResp)
	if err != nil {
		return nil, err
	}
	return volumeResp, nil
}

//Find the volume by it's name. If the volume is not found, an error will be returned.
func (v *volume) FindVolumeByName(volName string) (*types.Volume, error) {
	if len(volName) == 0 {
		return nil, errors.New(fmt.Sprintf("Lun Name shouldn't be empty."))
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "lun", volName, api.LunDisplayFields), nil, volumeResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find volume by name %s", volName))
	}

	return volumeResp, nil
}

//Find the volume by it's Id. If the volume is not found, an error will be returned.
func (v *volume) FindVolumeById(volId string) (*types.Volume, error) {
	if len(volId) == 0 {
		return nil, errors.New("lun ID shouldn't be empty")
	}
	volumeResp := &types.Volume{}
	err := v.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "lun", volId, api.LunDisplayFields), nil, volumeResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find volume Id %s Error: %v", volId, err))
	}
	return volumeResp, nil
}

func (v *volume) ListVolumes(startToken int, maxEntries int) ([]types.Volume, int, error) {
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

	err := v.client.executeWithRetryAuthenticate(http.MethodGet, lunUri, nil, volumeResp)
	if err != nil {
		log.Errorf("executeWithRetryAuthenticate Error: %v", err)
	}
	return volumeResp.Volumes, nextToken, err
}

func (v *volume) DeleteVolume(volId string) error {
	if len(volId) == 0 {
		return errors.New("lun id cannot be empty")
	}
	volumeResp := &types.Volume{}

	err := v.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceUri, "lun", volId), nil, volumeResp)
	if err != nil {
		log.Info("Unable to find volume ", err)
		return NotFoundError
	} else {
		deleteErr := v.client.executeWithRetryAuthenticate(http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceUri, "storageResource", volId), nil, nil)
		if deleteErr != nil {
			log.Info("Delete Lun Failed: ", deleteErr)
			return errors.New(fmt.Sprintf("Delete Lun %s Failed. Error: %v", volId, deleteErr))
		}
		log.Infof("Delete Lun %s Successful", volId)
		return nil
	}
}

//Export volume to a host
func (v *volume) ExportVolume(volID, hostID string) error {
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
	return v.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityModifyLunUri, volID), lunModifyParam, nil)
}

//Unexport volume
func (v *volume) UnexportVolume(volID string) error {
	hostAccessArray := []types.HostAccess{}
	lunParams := types.LunHostAccessParameters{
		HostAccess: &hostAccessArray,
	}
	lunModifyParam := types.LunHostAccessModifyParam{
		LunHostAccessParameters: &lunParams,
	}
	return v.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityModifyLunUri, volID), lunModifyParam, nil)
}

// Expand volume to provided capacity
func (v *volume) ExpandVolume(volumeId string, newSize uint64) error {
	vol, err := v.FindVolumeById(volumeId)
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
	return v.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityApiModifyLunUri, volumeId), volumeReqParam, nil)
}

func (v *volume) FindHostIOLimitByName(hostIoPolicyName string) (*types.IoLimitPolicy, error) {
	if len(hostIoPolicyName) == 0 {
		return nil, errors.New("policy Name shouldn't be empty")
	}
	ioLimitPolicyResp := &types.IoLimitPolicy{}
	fieldsToQuery := "id,name,description"
	fmt.Printf(api.UnityApiGetResourceByNameWithFieldsUri, "ioLimitPolicy", hostIoPolicyName, fieldsToQuery)
	err := v.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "ioLimitPolicy", hostIoPolicyName, fieldsToQuery), nil, ioLimitPolicyResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find IO Limit Policy:%s Error: %v", hostIoPolicyName, err))
	}
	return ioLimitPolicyResp, nil
}
