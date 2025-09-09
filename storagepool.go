/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/apitypes"
)

// FindStoragePoolByName - Find the volume by it's name. If the volume is not found, an error will be returned.
func (c *UnityClientImpl) FindStoragePoolByName(ctx context.Context, poolName string) (*types.StoragePool, error) {
	if len(poolName) == 0 {
		return nil, errors.New("poolName shouldn't be empty")
	}
	spResponse := &types.StoragePool{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.PoolAction, poolName, StoragePoolFields), nil, spResponse)
	if err != nil {
		return nil, fmt.Errorf("find storage pool by name failed %s err: %v", poolName, err)
	}

	return spResponse, nil
}

// FindStoragePoolByID - Find the volume by it's Id. If the volume is not found, an error will be returned.
func (c *UnityClientImpl) FindStoragePoolByID(ctx context.Context, poolID string) (*types.StoragePool, error) {
	if len(poolID) == 0 {
		return nil, errors.New("pool Id cannot be empty")
	}
	spResponse := &types.StoragePool{}

	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.PoolAction, poolID, StoragePoolFields), nil, spResponse)
	if err != nil {
		return nil, fmt.Errorf("find storage pool by ID failed %s err: %v", poolID, err)
	}

	return spResponse, nil
}
