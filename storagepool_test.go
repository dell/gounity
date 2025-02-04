/*
 Copyright © 2019-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package gounity

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mocksapi "github.com/dell/gounity/mocks/api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var storagePoolName string

func TestFindStoragePoolByID(t *testing.T) {
	fmt.Println("Begin - Find Storage Pool by ID Test")
	testConf.client.getAPI().(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Positive case
	pool, err := testConf.client.FindStoragePoolByID(ctx, testConf.poolID)
	fmt.Println("Find Storage Pool by ID:", prettyPrintJSON(pool), err)
	if err != nil {
		t.Fatalf("Find Storage Pool by ID failed: %v", err)
	}
	storagePoolName = pool.StoragePoolContent.Name

	// Negative cases
	// Case 1: Empty ID
	emptyID := ""
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("invalid ID")).Once()
	_, err = testConf.client.FindStoragePoolByID(ctx, emptyID)
	if err == nil {
		t.Fatalf("Find Storage Pool by ID with empty ID case - failed: %v", err)
	}

	// Case 2: Invalid ID
	invalidID := "dummy_pool_id_1"
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("invalid ID")).Once()
	_, err = testConf.client.FindStoragePoolByID(ctx, invalidID)
	if err == nil {
		t.Fatalf("Find Storage Pool by ID with invalid ID case - failed: %v", err)
	}

	fmt.Println("Find Storage Pool by ID Test - Successful")
}

func TestFindStoragePoolByNameTest(t *testing.T) {
	testConf.client.getAPI().(*mocksapi.Client).ExpectedCalls = nil
	assert := require.New(t)
	fmt.Println("Begin - Find Storage Pool by Name Test")
	ctx := context.Background()

	// Mock setup for valid pool name
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/pool/name:valid_pool_name?fields=id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasDataReductionEnabledLuns,hasDataReductionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Positive Case
	storagePoolName := "valid_pool_name" // Ensure this is set to a valid name
	pool, err := testConf.client.FindStoragePoolByName(ctx, storagePoolName)
	fmt.Println("Find volume by Name:", prettyPrintJSON(pool), err)
	assert.NoError(err, "Find Pool by Name failed")
	assert.NotNil(pool, "Pool should not be nil")

	// Mock setup for empty pool name
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/pool/name:?fields=id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasDataReductionEnabledLuns,hasDataReductionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Negative Case: Empty pool name
	storagePoolNameTemp := ""
	pool, err = testConf.client.FindStoragePoolByName(ctx, storagePoolNameTemp)
	assert.Error(err, "Expected error for empty pool name")
	assert.Nil(pool, "Pool should be nil for empty name")

	// Mock setup for invalid pool name
	testConf.client.getAPI().(*mocksapi.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/pool/name:dummy_pool_name_1?fields=id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasDataReductionEnabledLuns,hasDataReductionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("pool not found")).Once()

	// Negative Case: Invalid pool name
	storagePoolNameTemp = "dummy_pool_name_1"
	pool, err = testConf.client.FindStoragePoolByName(ctx, storagePoolNameTemp)
	assert.Error(err, "Expected error for invalid pool name")
	assert.Nil(pool, "Pool should be nil for invalid name")

	fmt.Println("Find Storage Pool by Name Test - Successful")
}
