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

type storagepool struct {
	client *Client
}

func NewStoragePool(client *Client) *storagepool {
	return &storagepool{client}
}

//Find the volume by it's name. If the volume is not found, an error will be returned.
func (sp *storagepool) FindStoragePoolByName(poolName string) (*types.StoragePool, error) {
	if len(poolName) == 0 {
		log.Error("pool name cannot be empty")
		return nil, errors.New("poolName shouldn't be empty")
	}
	spResponse := &types.StoragePool{}

	fieldsToQuery := "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasCompressionEnabledLuns,hasCompressionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"

	err := sp.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "pool", poolName, fieldsToQuery), nil, spResponse)
	if err != nil {
		log.Error("Unable to find pool")
		err = errors.New(fmt.Sprintf("Unable to find pool %s err: %v", poolName, err))
		return nil, err
	}

	return spResponse, nil
}

//Find the volume by it's Id. If the volume is not found, an error will be returned.
func (sp *storagepool) FindStoragePoolById(poolId string) (*types.StoragePool, error) {
	if len(poolId) == 0 {
		log.Error("pool Id cannot be empty")
		return nil, errors.New("pool Id cannot be empty")
	}
	spResponse := &types.StoragePool{}

	fieldsToQuery := "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasCompressionEnabledLuns,hasCompressionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"

	log.Info("URI", fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "pool", poolId, fieldsToQuery))
	err := sp.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "pool", poolId, fieldsToQuery), nil, spResponse)
	if err != nil {
		log.Errorf("Unable to find PoolID. Error: %v ", err)
		return nil, errors.New(fmt.Sprintf("unable to find the PoolID %s", poolId))
	}

	return spResponse, nil
}
