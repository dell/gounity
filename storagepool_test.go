/*
 Copyright © 2019 Dell Inc. or its subsidiaries. All Rights Reserved.

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
	"fmt"
	"testing"
)

var storagePoolName string

func TestStoragePool(t *testing.T) {
	ctx = context.Background()

	findStoragePoolByIDTest(t)
	findStoragePoolByNameTest(t)
}

func findStoragePoolByIDTest(t *testing.T) {
	fmt.Println("Begin - Find Storage Pool by Id Test")

	pool, err := testConf.poolAPI.FindStoragePoolByID(ctx, testConf.poolID)
	fmt.Println("Find volume by Id:", prettyPrintJSON(pool), err)
	if err != nil {
		t.Fatalf("Find Pool by Id failed: %v", err)
	}
	storagePoolName = pool.StoragePoolContent.Name

	// Negative cases
	storagePoolIDTemp := ""
	pool, err = testConf.poolAPI.FindStoragePoolByID(ctx, storagePoolIDTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with empty Id case - failed: %v", err)
	}

	storagePoolIDTemp = "dumy_pool_id_1"
	pool, err = testConf.poolAPI.FindStoragePoolByID(ctx, storagePoolIDTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with invalid Id case - failed: %v", err)
	}

	fmt.Println("Find Storage Pool by Id Test - Successful")
}

func findStoragePoolByNameTest(t *testing.T) {
	fmt.Println("Begin - Find Storage Pool by Name Test")

	pool, err := testConf.poolAPI.FindStoragePoolByName(ctx, storagePoolName)
	fmt.Println("Find volume by Name:", prettyPrintJSON(pool), err)
	if err != nil {
		t.Fatalf("Find Pool by Name failed: %v", err)
	}

	// Negative Cases
	storagePoolNameTemp := ""
	pool, err = testConf.poolAPI.FindStoragePoolByName(ctx, storagePoolNameTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with empty Name case - failed: %v", err)
	}

	storagePoolNameTemp = "dummy_pool_name_1"
	pool, err = testConf.poolAPI.FindStoragePoolByName(ctx, storagePoolNameTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with invalid Name case - failed: %v", err)
	}

	fmt.Println("Find Storage Pool by Name Test - Successful")
}
