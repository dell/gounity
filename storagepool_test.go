package gounity

import (
	"fmt"
	"testing"
)

func TestGetStoragePool(t *testing.T) {
	pool, err := testConf.poolApi.FindStoragePoolById(testConf.poolId)

	fmt.Println("GetPool:", pool, err)
	if err != nil {
		t.Fatalf("Get Pool failed: %v", err)
	}

	pool, err = testConf.poolApi.FindStoragePoolByName(pool.StoragePoolContent.Name)

	fmt.Println("GetPool by name:", pool, err)
	if err != nil {
		t.Fatalf("Get Pool by name failed: %v", err)
	}
	//produce 404 status code
	poolId := "afasd89798asdfasfa089798" //poolid should not exists in unity
	pool, err = testConf.poolApi.FindStoragePoolById(poolId)

	fmt.Println("GetPool:", pool, err)
	if err != nil {
		t.Logf("Expected failure: %v", err)
	} else {
		t.Fatal("Pool should not be found")
	}
}
