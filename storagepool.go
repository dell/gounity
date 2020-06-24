/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"errors"
	"fmt"
	"net/http"

	"context"
	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

type storagepool struct {
	client *Client
}

func NewStoragePool(client *Client) *storagepool {
	return &storagepool{client}
}

//Find the volume by it's name. If the volume is not found, an error will be returned.
func (sp *storagepool) FindStoragePoolByName(ctx context.Context, poolName string) (*types.StoragePool, error) {
	if len(poolName) == 0 {
		return nil, errors.New("poolName shouldn't be empty")
	}
	spResponse := &types.StoragePool{}

	fieldsToQuery := "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasCompressionEnabledLuns,hasCompressionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"

	err := sp.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "pool", poolName, fieldsToQuery), nil, spResponse)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to find pool %s err: %v", poolName, err))
		return nil, err
	}

	return spResponse, nil
}

//Find the volume by it's Id. If the volume is not found, an error will be returned.
func (sp *storagepool) FindStoragePoolById(ctx context.Context, poolId string) (*types.StoragePool, error) {
	if len(poolId) == 0 {
		return nil, errors.New("pool Id cannot be empty")
	}
	spResponse := &types.StoragePool{}

	fieldsToQuery := "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasCompressionEnabledLuns,hasCompressionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"
	err := sp.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "pool", poolId, fieldsToQuery), nil, spResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to find the PoolID %s", poolId))
	}

	return spResponse, nil
}
