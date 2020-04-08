/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"context"
	"errors"
	"fmt"
	"github.com/dell/gounity/util"
	"net/http"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

type filesystem struct {
	client *Client
}

const (
	FsNameMaxLength = 63
)

var FilesystemNotFoundError = errors.New("unable to find filesystem")

func NewFilesystem(client *Client) *filesystem {
	return &filesystem{client}
}

//FindFilesystemByName - Find the Filesystem by it's name. If the Filesystem is not found, an error will be returned.
func (f *filesystem) FindFilesystemByName(ctx context.Context, filesystemName string) (*types.Filesystem, error) {
	log := util.GetRunIdLogger(ctx)
	if len(filesystemName) == 0 {
		return nil, errors.New("filesystem Name shouldn't be empty")
	}
	hResponse := &types.Filesystem{}
	fieldsToQuery := "id,name,description,type,sizeTotal,isReadOnly,isThinEnabled,isDataReductionEnabled,pool,nasServer"
	log.Debug("URI", fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "filesystem", filesystemName, fieldsToQuery))
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "filesystem", filesystemName, fieldsToQuery), nil, hResponse)
	if err != nil {
		log.Error("Unable to find Filesystem", err)
		return nil, FilesystemNotFoundError
	}
	return hResponse, nil
}

//CreateFilesystem - Create a new filesystem on the array
func (f *filesystem) CreateFilesystem(ctx context.Context, name, storagepool, description, nasServer string, size uint64, tieringPolicy, hostIOSize, supportedProtocol int, isThinEnabled, isDataReductionEnabled, isCacheDisabled bool) (*types.Filesystem, error) {
	log := util.GetRunIdLogger(ctx)
	if name == "" {
		return nil, errors.New("filesystem name should not be empty.")
	}

	if len(name) > FsNameMaxLength {
		return nil, errors.New(fmt.Sprintf("filesystem name %s should not exceed %d characters.", name, FsNameMaxLength))
	}

	poolApi := NewStoragePool(f.client)
	pool, err := poolApi.FindStoragePoolById(ctx, storagepool)
	if pool == nil {
		log.Errorf("Unable to get PoolID (%s) Error:%v", storagepool, err)
		return nil, errors.New(fmt.Sprintf("unable to get PoolID (%s) Error:%v", storagepool, err))
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error trying to get Storage Pool (%s) Error:%v", storagepool, err))
	}

	storagePool := types.StoragePoolID{
		PoolId: storagepool,
	}

	fileEventSettings := types.FileEventSettings{
		IsCIFSEnabled: false, //Set to false to disable CIFS
		IsNFSEnabled:  true,  //Set to true to enable NFS alone
	}

	nas := types.NasServerID{
		NasServerID: nasServer,
	}

	fsParams := types.FsParameters{
		StoragePool:            &storagePool,
		Size:                   size,
		IsThinEnabled:          isThinEnabled,
		IsDataReductionEnabled: isDataReductionEnabled,
		SupportedProtocol:      supportedProtocol,
		HostIOSize:             hostIOSize,
		IsCacheDisabled:        isCacheDisabled,
		NasServer:              &nas,
		FileEventSettings:      fileEventSettings,
	}

	if pool != nil && pool.StoragePoolContent.PoolFastVP.Status != 0 {
		log.Info("FastVP is enabled")
		fastVPParameters := types.FastVPParameters{
			TieringPolicy: tieringPolicy,
		}
		fsParams.FastVPParameters = &fastVPParameters
	} else {
		log.Info("FastVP is not enabled")
	}

	fileReqParam := types.FsCreateParam{
		Name:         name,
		Description:  description,
		FsParameters: &fsParams,
	}

	fileResp := &types.Filesystem{}
	err = f.client.executeWithRetryAuthenticate(ctx,
		http.MethodPost, fmt.Sprintf(api.UnityApiStorageResourceActionUri, "createFilesystem"), fileReqParam, fileResp)
	if err != nil {
		return nil, err
	}

	return fileResp, nil
}

//DeleteFilesystem - Delete the Filesystem by it's ID. If the Filesystem is not present on the array, an error will be returned.
func (f *filesystem) DeleteFilesystem(ctx context.Context, filesystemID string) error {
	log := util.GetRunIdLogger(ctx)
	if len(filesystemID) == 0 {
		return errors.New("filesystem ID shouldn't be empty")
	}
	filesystemResp := &types.Filesystem{}
	fieldsToQuery := "id,storageResource"

	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "filesystem", filesystemID, fieldsToQuery), nil, filesystemResp)

	if err != nil {
		log.Debugf("Unable to find filesystem %v", err)
		return FilesystemNotFoundError
	} else {
		resourceID := filesystemResp.FileContent.StorageResource.Id
		deleteErr := f.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceUri, "storageResource", resourceID), nil, nil)
		if deleteErr != nil {
			log.Errorf("Delete Filesystem Failed: %v", deleteErr)
			return errors.New(fmt.Sprintf("Delete Filesystem %s Failed. Error: %v", filesystemID, deleteErr))
		}
		log.Debugf("Delete Filesystem %s Successful", filesystemID)
		return nil
	}
}
