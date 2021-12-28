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

	pool, err := testConf.poolApi.FindStoragePoolById(ctx, testConf.poolId)
	fmt.Println("Find volume by Id:", prettyPrintJson(pool), err)
	if err != nil {
		t.Fatalf("Find Pool by Id failed: %v", err)
	}
	storagePoolName = pool.StoragePoolContent.Name

	//Negative cases
	storagePoolIdTemp := ""
	pool, err = testConf.poolApi.FindStoragePoolById(ctx, storagePoolIdTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with empty Id case - failed: %v", err)
	}

	storagePoolIdTemp = "dumy_pool_id_1"
	pool, err = testConf.poolApi.FindStoragePoolById(ctx, storagePoolIdTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with invalid Id case - failed: %v", err)
	}

	fmt.Println("Find Storage Pool by Id Test - Successful")
}

func findStoragePoolByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Storage Pool by Name Test")

	pool, err := testConf.poolApi.FindStoragePoolByName(ctx, storagePoolName)
	fmt.Println("Find volume by Name:", prettyPrintJson(pool), err)
	if err != nil {
		t.Fatalf("Find Pool by Name failed: %v", err)
	}

	//Negative Cases
	storagePoolNameTemp := ""
	pool, err = testConf.poolApi.FindStoragePoolByName(ctx, storagePoolNameTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with empty Name case - failed: %v", err)
	}

	storagePoolNameTemp = "dummy_pool_name_1"
	pool, err = testConf.poolApi.FindStoragePoolByName(ctx, storagePoolNameTemp)
	if err == nil {
		t.Fatalf("Find Pool by Id with invalid Name case - failed: %v", err)
	}

	fmt.Println("Find Storage Pool by Name Test - Successful")
}
