package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dell/gounity/types"
)

var cgCreateParams *types.ConsistencyGroupCreate
var cgName string
var cgID string

func TestCG(t *testing.T) {
	ctx = context.Background()

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	cgName = "Unit-test-cg-" + timeStamp
}

func findCGByNameTest(t *testing.T) {

	fmt.Println("Begin - Find CG By Name Test")

	cg, err := testConf.cgAPI.GetConsistencyGroupByName(ctx, cgName)
	fmt.Println("Find CG by Name:", prettyPrintJSON(cg), err)
	if err != nil {
		t.Fatalf("Find CG by Name failed: %v", err)
	}
	cgID = cg.ConsistencyGroupContent.ResourceID

	//Negative cases
	emptyName := ""
	_, err = testConf.cgAPI.GetConsistencyGroupByName(ctx, emptyName)
	if err == nil {
		t.Fatalf("Find CG by Name with empty name case failed: %v", err)
	}

	cgNameTemp := "dummy_CG_1"
	_, err = testConf.cgAPI.GetConsistencyGroupByName(ctx, cgNameTemp)
	if err == nil {
		t.Fatalf("Find CG by Name with invalid name case failed: %v", err)
	}

	fmt.Println("Find CG by Name Test - Successful")
}

func getCGByIDTest(t *testing.T) {

	fmt.Println("Begin - Get CG By ID Test")

	cg, err := testConf.cgAPI.GetConsistencyGroup(ctx, cgID)
	fmt.Println("Find CG by Name:", prettyPrintJSON(cg), err)
	if err != nil {
		t.Fatalf("Find CG by Id failed: %v", err)
	}

	//Negative cases
	emptyID := ""
	_, err = testConf.cgAPI.GetConsistencyGroup(ctx, emptyID)
	if err == nil {
		t.Fatalf("Find ConsistencyGroup by Id with empty Id case failed: %v", err)
	}

	cgIDTemp := "dummy_cg_sv_1"
	_, err = testConf.cgAPI.GetConsistencyGroup(ctx, cgIDTemp)
	if err == nil {
		t.Fatalf("Find CG by Id with invalid Id case failed: %v", err)
	}
	fmt.Println("Find CG by Id Test - Successful")
}

func createConsistencyGroupTest(t *testing.T) {

	fmt.Println("Begin - Create CG Test")

	_, err := testConf.cgAPI.CreateConsistencyGroup(ctx, cgCreateParams)
	if err != nil {
		t.Fatalf("Create cg failed: %v", err)
	}
	cg, err := testConf.cgAPI.GetConsistencyGroupByName(ctx, cgName)
	fmt.Println("Created CG:", prettyPrintJSON(cg), err)
	// fmt.Println(" CG id:", cg.ConsistencyGroupContent.ResourceID)
	if err != nil {
		t.Fatalf("Cannot fined CG by name: %v", err)
	}

	//Negative cases
	// empty name
	cgCreateParams.Name = ""
	_, err = testConf.cgAPI.CreateConsistencyGroup(ctx, cgCreateParams)
	if err == nil {
		t.Fatalf("Create cg with empty name case failed: %v", err)
	}

	// too long name
	cgCreateParams.Name = "cg-name-max-length-12345678901234567890123456789012345678901234567890-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	_, err = testConf.cgAPI.CreateConsistencyGroup(ctx, cgCreateParams)
	if err == nil {
		t.Fatalf("Create cg exceeding max name length case failed: %v", err)
	}

	fmt.Println("Create cg Test - Successful")
}

func deleteCGTest(t *testing.T) {

	fmt.Println("Begin - Delete CG Test")

	//Negative cases
	cgIDTemp := ""
	err := testConf.cgAPI.DeleteConsistencyGroup(ctx, cgIDTemp)
	if err == nil {
		t.Fatalf("Delete ConsistencyGroup with empty Id case failed: %v", err)
	}

	cgIDTemp = "dummy_cg_1"
	err = testConf.cgAPI.DeleteConsistencyGroup(ctx, cgIDTemp)
	if err == nil {
		t.Fatalf("Delete ConsistencyGroup with invalid Id case failed: %v", err)
	}

	err = testConf.cgAPI.DeleteConsistencyGroup(ctx, cgID)
	if err != nil {
		t.Fatalf("Delete ConsistencyGroup failed: %v", err)
	}

	fmt.Println("Delete ConsistencyGroup Test - Successful")

}
