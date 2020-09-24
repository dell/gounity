package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var snapVolName string
var snapVolID string
var snapName string
var snapID string
var snap2Name string
var snap2ID string
var cloneVolName string
var cloneVolID string

func TestSnapshot(t *testing.T) {

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	snapVolName = "Unit-test-snap-vol-" + timeStamp
	snapName = "Unit-test-snapshot-" + timeStamp
	snap2Name = "Unit-test-snapshot2-" + timeStamp
	cloneVolName = "Unit-test-clone-vol-" + timeStamp
	ctx = context.Background()

	createSnapshotTest(t)
	findSnapshotByNameTest(t)
	findSnapshotByIDTest(t)
	listSnapshotsTest(t)
	modifySnapshotAutoDeleteParameterTest(t)
	creteLunThinCloneTest(t) //create thin clone
	deleteSnapshot(t)
}

func createSnapshotTest(t *testing.T) {

	fmt.Println("Begin - Create Snapshot Test")

	vol, err := testConf.volumeApi.CreateLun(ctx, snapVolName, testConf.poolId, "Description", 5368709120, 0, "", true, false)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(ctx, snapVolName)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}
	snapVolID = vol.VolumeContent.ResourceId

	snap, err := testConf.snapApi.CreateSnapshot(ctx, snapVolID, snapName, "Snapshot Description", "")
	fmt.Println("Create Snapshot:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Create Snapshot failed: %v", err)
	}

	snap, err = testConf.snapApi.CreateSnapshot(ctx, snapVolID, snap2Name, "Snapshot Description", "1:23:52:50")
	fmt.Println("Create Snapshot2:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Create Snapshot 2failed: %v", err)
	}

	//Negative cases
	snapVolIDTemp := ""
	_, err = testConf.snapApi.CreateSnapshot(ctx, snapVolIDTemp, snap2Name, "Snapshot Description", "")
	if err == nil {
		t.Fatalf("Create Snapshot with empty volume Id case failed: %v", err)
	}

	snapNameTemp := "snap-name-max-length-12345678901234567890123456789012345678901234567890"
	_, err = testConf.snapApi.CreateSnapshot(ctx, snapVolID, snapNameTemp, "Snapshot Description", "")
	if err == nil {
		t.Fatalf("Create Snapshot with max name characters case failed: %v", err)
	}

	_, err = testConf.snapApi.CreateSnapshot(ctx, snapVolID, snap2Name, "Snapshot Description", "1:23:99:99")
	if err == nil {
		t.Fatalf("Create Snapshot with invalid retention duration case failed: %v", err)
	}

	_, err = testConf.snapApi.CreateSnapshot(ctx, snapVolID, snap2Name, "Snapshot Description", "1:23:52:50")
	if err == nil {
		t.Fatalf("Create duplicate Snapshot case failed: %v", err)
	}

	fmt.Println("Create Snapshot Test - Successful")
}

func findSnapshotByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Snapshot by Name Test")

	snap, err := testConf.snapApi.FindSnapshotByName(ctx, snapName)
	fmt.Println("Find snapshot by Name:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}
	snapID = snap.SnapshotContent.ResourceId

	snap, err = testConf.snapApi.FindSnapshotByName(ctx, snap2Name)
	fmt.Println("Find snapshot2 by Name:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Find snapshot2 failed: %v", err)
	}
	snap2ID = snap.SnapshotContent.ResourceId

	//Negative test cases
	snapNameTemp := ""
	_, err = testConf.snapApi.FindSnapshotByName(ctx, snapNameTemp)
	if err == nil {
		t.Fatalf("Find snapshot by Name with empty name case failed: %v", err)
	}

	snapNameTemp = "dummy_snap_name_1"
	_, err = testConf.snapApi.FindSnapshotByName(ctx, snapNameTemp)
	if err == nil {
		t.Fatalf("Find snapshot by Name with empty name case failed: %v", err)
	}

	fmt.Println("Find Snapshot by Name - Successful")
}

func findSnapshotByIDTest(t *testing.T) {

	fmt.Println("Begin - Find Snapshot by Id Test")

	snap, err := testConf.snapApi.FindSnapshotById(ctx, snapID)
	fmt.Println("Find snapshot by ID:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}

	//Negative test cases
	snapIDTemp := ""
	_, err = testConf.snapApi.FindSnapshotById(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Find snapshot by Id with empty Id case failed: %v", err)
	}

	snapIDTemp = "dummy_snap_id_1"
	_, err = testConf.snapApi.FindSnapshotById(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Find snapshot by Id with empty id case failed: %v", err)
	}

	fmt.Println("Find Snapshot by Id - Successful")
}

func listSnapshotsTest(t *testing.T) {

	fmt.Println("Begin - List Snapshots Test")

	snaps, _, err := testConf.snapApi.ListSnapshots(ctx, 0, 10, snapVolID, "")
	fmt.Println("List snapshots:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots success:", len(snaps))
	} else {
		t.Fatalf("List snapshot failed: %v", err)
	}

	snaps, _, err = testConf.snapApi.ListSnapshots(ctx, 0, 10, snapVolID, snapID)
	fmt.Println("List snapshots with snap Id:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots with snap Id success:", len(snaps))
	} else {
		t.Fatalf("List snapshot with snap Id failed: %v", err)
	}

	snaps, _, err = testConf.snapApi.ListSnapshots(ctx, 6, 5, "", "")
	fmt.Println("List snapshots pagination:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots pagination success:", len(snaps))
	} else {
		t.Fatalf("List snapshot pagination failed: %v", err)
	}

	fmt.Println("List Snapshots Test - Successful")
}

func modifySnapshotAutoDeleteParameterTest(t *testing.T) {

	fmt.Println("Begin - Modify Snapshot Test")

	err := testConf.snapApi.ModifySnapshotAutoDeleteParameter(ctx, snapID)
	if err != nil {
		t.Fatalf("Modify Snapshot failed: %v", err)
	}

	//Negative test cases
	snapIDTemp := ""
	err = testConf.snapApi.ModifySnapshotAutoDeleteParameter(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Modify snapshot with empty Id case failed: %v", err)
	}

	snapIDTemp = "dummy_snap_id_1"
	err = testConf.snapApi.ModifySnapshotAutoDeleteParameter(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Modify snapshot with invalid Id case failed: %v", err)
	}

	fmt.Println("Modify Snapshot Test - Successful")
}

func creteLunThinCloneTest(t *testing.T) {

	fmt.Println("Begin - Create LUN thin clone Test")

	vol, err := testConf.volumeApi.CreteLunThinClone(ctx, cloneVolName, snapID, snapVolID)
	if err != nil {
		t.Fatalf("Create thin clone failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(ctx, cloneVolName)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}
	cloneVolID = vol.VolumeContent.ResourceId
	fmt.Println("Create LUN thin clone Test - Successful")
}

func deleteSnapshot(t *testing.T) {

	fmt.Println("Begin - Delete Snapshot Test")

	err := testConf.snapApi.DeleteSnapshot(ctx, snapID)
	if err != nil {
		t.Fatalf("Delete Snapshot failed: %v", err)
	}

	err = testConf.snapApi.DeleteSnapshot(ctx, snap2ID)
	if err != nil {
		t.Fatalf("Delete Snapshot2 failed: %v", err)
	}

	//Delete thin clone volume
	err = testConf.volumeApi.DeleteVolume(ctx, cloneVolID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(ctx, snapVolID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	//Negative test cases
	snapIDTemp := ""
	err = testConf.snapApi.DeleteSnapshot(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Delete snapshot with empty Id case failed: %v", err)
	}

	snapIDTemp = "dummy_snapshot_id_1"
	err = testConf.snapApi.DeleteSnapshot(ctx, snapIDTemp)
	if err == nil {
		t.Fatalf("Delete snapshot with invalid Id case failed: %v", err)
	}

	fmt.Println("Delete Snapshot Test - Successful")
}

/*
func deleteSnapshot(t *testing.T) {

	fmt.Println("Begin - Delete Snapshot Test")


	//Add Negative test cases
	fmt.Println("Delete Snapshot Test - Successful")
}

*/
