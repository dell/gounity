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

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

// Storagepool structure
type Storagepool struct {
	client *Client
}

// NewStoragePool returns storagepool
func NewStoragePool(client *Client) *Storagepool {
	return &Storagepool{client}
}

// FindStoragePoolByName - Find the volume by it's name. If the volume is not found, an error will be returned.
func (sp *Storagepool) FindStoragePoolByName(ctx context.Context, poolName string) (*types.StoragePool, error) {
	if len(poolName) == 0 {
		return nil, errors.New("poolName shouldn't be empty")
	}
	spResponse := &types.StoragePool{}
	err := sp.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.PoolAction, poolName, StoragePoolFields), nil, spResponse)
	if err != nil {
		return nil, fmt.Errorf("find storage pool by name failed %s err: %v", poolName, err)
	}

	return spResponse, nil
}

// FindStoragePoolByID - Find the volume by it's Id. If the volume is not found, an error will be returned.
func (sp *Storagepool) FindStoragePoolByID(ctx context.Context, poolID string) (*types.StoragePool, error) {
	if len(poolID) == 0 {
		return nil, errors.New("pool Id cannot be empty")
	}
	spResponse := &types.StoragePool{}

	err := sp.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.PoolAction, poolID, StoragePoolFields), nil, spResponse)
	if err != nil {
		return nil, fmt.Errorf("find storage pool by ID failed %s err: %v", poolID, err)
	}

	return spResponse, nil
}
