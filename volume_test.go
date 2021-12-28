package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var volName string
var cloneVolumeName string
var volID string
var cloneVolumeID string
var hostIOLimitID string

func TestVolume(t *testing.T) {
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName = "Unit-test-vol-" + timeStamp
	cloneVolumeName = "Unit-test-clone-vol-" + timeStamp
	ctx = context.Background()

	findHostIOLimitByNameTest(t)
	createLunTest(t)
	findVolumeByNameTest(t)
	findVolumeByIdTest(t)
	listVolumesTest(t)
	exportVolumeTest(t)
	unexportVolumeTest(t)
	expandVolumeTest(t)
	createCloneFromVolumeTest(t)
	deleteVolumeTest(t)
	//creteLunThinCloneTest(t) - Will be added to snapshot_test

}

func findHostIOLimitByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Host IO Limit by Name Test")

	if testConf.hostIOLimitName != "" {
		hostIOLimit, err := testConf.volumeApi.FindHostIOLimitByName(ctx, testConf.hostIOLimitName)
		fmt.Println("hostIOLimit:", prettyPrintJson(hostIOLimit), "Error:", err)
		hostIOLimitID = hostIOLimit.IoLimitPolicyContent.Id

		//Negative case
		hostIOTemp := "dummy_hostio_1"
		_, err = testConf.volumeApi.FindHostIOLimitByName(ctx, hostIOTemp)
		if err == nil {
			t.Fatalf("Find Host IO Limit negative case failed: %v", err)
		}

		hostIOTemp = ""
		_, err = testConf.volumeApi.FindHostIOLimitByName(ctx, hostIOTemp)
		if err == nil {
			t.Fatalf("Find Host IO Limit with empty name case failed: %v", err)
		}

		fmt.Println("Find Host IO Limit by Name Test - Successful")
	} else {
		fmt.Println("Skipping Host IO Limit by Name Test - Parameter not configured")
	}
}

func createLunTest(t *testing.T) {

	fmt.Println("Begin - Create LUN Test")

	_, err := testConf.volumeApi.CreateLun(ctx, volName, testConf.poolId, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err != nil {
		t.Fatalf("Create LUN failed: %v", err)
	}

	//Negative cases
	volNameTemp := ""
	_, err = testConf.volumeApi.CreateLun(ctx, volNameTemp, testConf.poolId, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with empty name case failed: %v", err)
	}

	volNameTemp = "vol-name-max-length-12345678901234567890123456789012345678901234567890"
	_, err = testConf.volumeApi.CreateLun(ctx, volNameTemp, testConf.poolId, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN exceeding max name length case failed: %v", err)
	}

	poolIDTemp := "dummy_pool_1"
	_, err = testConf.volumeApi.CreateLun(ctx, volName, poolIDTemp, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with invalid pool name case failed: %v", err)
	}

	_, err = testConf.volumeApi.CreateLun(ctx, volName, testConf.poolId, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with same name case failed: %v", err)
	}

	fmt.Println("Create LUN Test - Successful")
}

func findVolumeByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Volume By Name Test")

	vol, err := testConf.volumeApi.FindVolumeByName(ctx, volName)
	fmt.Println("Find volume by Name:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Name failed: %v", err)
	}
	volID = vol.VolumeContent.ResourceId

	//Negative cases
	volNameTemp := ""
	_, err = testConf.volumeApi.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with empty name case failed: %v", err)
	}

	volNameTemp = "dummy_volume_1"
	_, err = testConf.volumeApi.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with invalid name case failed: %v", err)
	}

	fmt.Println("Find Volume by Name Test - Successful")
}

func findVolumeByIdTest(t *testing.T) {

	fmt.Println("Begin - Find Volume By Name Test")

	vol, err := testConf.volumeApi.FindVolumeById(ctx, volID)
	fmt.Println("Find volume by Name:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Id failed: %v", err)
	}

	//Negative cases
	volIDTemp := ""
	_, err = testConf.volumeApi.FindVolumeById(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with empty Id case failed: %v", err)
	}

	volIDTemp = "dummy_vol_sv_1"
	_, err = testConf.volumeApi.FindVolumeById(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with invalid Id case failed: %v", err)
	}
	fmt.Println("Find Volume by Id Test - Successful")
}

func listVolumesTest(t *testing.T) {

	fmt.Println("Begin - List Volumes Test")

	vols, _, err := testConf.volumeApi.ListVolumes(ctx, 11, 10)
	fmt.Println("List volumes count: ", len(vols))
	if len(vols) <= 10 {
		fmt.Println("List volume success")
	} else {
		t.Fatalf("List volumes failed: %v", err)
	}

	fmt.Println("List Volume Test - Successful")
}

func exportVolumeTest(t *testing.T) {

	fmt.Println("Begin - Export Volume Test")

	host, err := testConf.hostApi.FindHostByName(ctx, testConf.nodeHostName)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	err = testConf.volumeApi.ExportVolume(ctx, volID, host.HostContent.ID)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}

	//Negative case for Delete Volume
	err = testConf.volumeApi.DeleteVolume(ctx, volID)
	if err == nil {
		t.Fatalf("Delete volume on exported volume case failed: %v", err)
	}

	fmt.Println("Export Volume Test - Successful")
}

func unexportVolumeTest(t *testing.T) {

	fmt.Println("Begin - Unexport Volume Test")

	err := testConf.volumeApi.UnexportVolume(ctx, volID)
	if err != nil {
		t.Fatalf("UnExportVolume failed: %v", err)
	}
	fmt.Println("Unexport Volume Test - Successful")
}

func expandVolumeTest(t *testing.T) {

	fmt.Println("Begin - Expand Volume Test")

	err := testConf.volumeApi.ExpandVolume(ctx, volID, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume failed: %v", err)
	}

	err = testConf.volumeApi.ExpandVolume(ctx, volID, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume with same size failed: %v", err)
	}

	//Negative cases
	volIDTemp := "dummy_vol_sv_1"
	err = testConf.volumeApi.ExpandVolume(ctx, volIDTemp, 5368709120)
	if err == nil {
		t.Fatalf("Expand volume with invalid Id case failed: %v", err)
	}

	err = testConf.volumeApi.ExpandVolume(ctx, volID, 4368709120)
	if err == nil {
		t.Fatalf("Expand volume with smaller size case failed: %v", err)
	}

	fmt.Println("Expand Volume Test - Successful")
}

func createCloneFromVolumeTest(t *testing.T) {

	fmt.Println("Begin - Create clone from Volume Test")

	_, err := testConf.volumeApi.CreateCloneFromVolume(ctx, cloneVolumeName, volID)
	if err != nil {
		t.Fatalf("Clone volume failed: %v", err)
	}

	vol, err := testConf.volumeApi.FindVolumeByName(ctx, cloneVolumeName)
	fmt.Println("Find volume by Name:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Name failed: %v", err)
	}

	cloneVolumeID = vol.VolumeContent.ResourceId

	fmt.Println("Create clone from Volume Test - Successful")
}

func deleteVolumeTest(t *testing.T) {

	fmt.Println("Begin - Delete Volume Test")

	err := testConf.volumeApi.DeleteVolume(ctx, cloneVolumeID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(ctx, volID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	//Negative cases
	volIDTemp := ""
	err = testConf.volumeApi.DeleteVolume(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Delete volume with empty Id case failed: %v", err)
	}

	volIDTemp = "dummy_vol_sv_1"
	err = testConf.volumeApi.DeleteVolume(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Delete volume with invalid Id case failed: %v", err)
	}

	fmt.Println("Delete Volume Test - Successful")
}
