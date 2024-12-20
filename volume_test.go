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

	"github.com/dell/gounity/mocks"
	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	volName         = "unit-test-vol"
	cloneVolumeName = "unit-test-clone-vol"
	volID           = "unity-volume-id"
	cloneVolumeID   string
	hostIOLimitID   string
	anyArgs         = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
	ctx             = context.Background()
)

// func TestVolume(t *testing.T) {
// 	now := time.Now()
// 	timeStamp := now.Format("20060102150405")
// 	volName = "Unit-test-vol-" + timeStamp
// 	cloneVolumeName = "Unit-test-clone-vol-" + timeStamp
// 	ctx = context.Background()

// 	findHostIOLimitByNameTest(t)
// 	createLunTest(t)
// 	findVolumeByNameTest(t)
// 	findVolumeByIDTest(t)
// 	listVolumesTest(t)
// 	exportVolumeTest(t)
// 	unexportVolumeTest(t)
// 	expandVolumeTest(t)
// 	createCloneFromVolumeTest(t)
// 	modifyVolumeExportTest(t)
// 	deleteVolumeTest(t)
// 	getMaxVolumeSizeTest(t)
// 	// creteLunThinCloneTest(t) - Will be added to snapshot_test
// }

func TestFindHostIOLimitByName(t *testing.T) {
	fmt.Println("Begin - Find Host IO Limit by Name Test")

	// Mock the client.DoWithHeaders to return nil
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()

	// Call the FindHostIOLimitByName function
	hostIOLimit, err := testConf.volumeAPI.FindHostIOLimitByName(ctx, testConf.hostIOLimitName)
	fmt.Println("hostIOLimit:", prettyPrintJSON(hostIOLimit), "Error:", err)
	assert.NotNil(t, hostIOLimit.IoLimitPolicyContent)

	// Negative cases

	// Mock the client.DoWithHeaders to return an error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("not found")).Once()

	// Call the FindHostIOLimitByName function with a dummy name
	_, err = testConf.volumeAPI.FindHostIOLimitByName(ctx, "dummy_hostio_1")
	if err == nil {
		t.Fatalf("Find Host IO Limit negative case failed: %v", err)
	}

	// Call the FindHostIOLimitByName function with an empty name
	_, err = testConf.volumeAPI.FindHostIOLimitByName(ctx, "")
	if err == nil {
		t.Fatalf("Find Host IO Limit with empty name case failed: %v", err)
	}

	fmt.Println("Find Host IO Limit by Name Test - Successful")
}

func TestCreateLun(t *testing.T) {
	fmt.Println("Begin - Create LUN Test")

	// Mock FindStoragePoolByID to return nil
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock isFeatureLicensed to return expected response
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.LicenseInfo")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.LicenseInfo)
			*resp = types.LicenseInfo{LicenseInfoContent: types.LicenseInfoContent{IsInstalled: true, IsValid: true}}
		}).Twice()
	// Mock create request to return nil
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()

	_, err := testConf.volumeAPI.CreateLun(ctx, volName, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err != nil {
		t.Fatalf("Create LUN failed: %v", err)
	}

	// Negative cases
	volNameTemp := ""
	_, err = testConf.volumeAPI.CreateLun(ctx, volNameTemp, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with empty name case failed: %v", err)
	}

	volNameTemp = "vol-name-max-length-12345678901234567890123456789012345678901234567890"
	_, err = testConf.volumeAPI.CreateLun(ctx, volNameTemp, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN exceeding max name length case failed: %v", err)
	}

	// Mock FindStoragePoolByID to return error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("storage pool not found")).Once()
	poolIDTemp := "dummy_pool_1"
	_, err = testConf.volumeAPI.CreateLun(ctx, volName, poolIDTemp, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with invalid pool name case failed: %v", err)
	}

	// Mock FindStoragePoolByID to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock isFeatureLicensed to return expected response
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.LicenseInfo")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.LicenseInfo)
			*resp = types.LicenseInfo{LicenseInfoContent: types.LicenseInfoContent{IsInstalled: true, IsValid: true}}
		}).Twice()
	// Mock create volume to return error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume already exists")).Once()
	_, err = testConf.volumeAPI.CreateLun(ctx, volName, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with same name case failed: %v", err)
	}

	fmt.Println("Create LUN Test - Successful")
}

func TestFindVolumeByName(t *testing.T) {
	fmt.Println("Begin - Find Volume By Name Test")

	// Mock FindVolumeByName to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	vol, err := testConf.volumeAPI.FindVolumeByName(ctx, volName)
	fmt.Println("Find volume by Name:", prettyPrintJSON(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Name failed: %v", err)
	}
	assert.NotNil(t, vol.VolumeContent.ResourceID)

	// Negative cases
	volNameTemp := ""
	_, err = testConf.volumeAPI.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with empty name case failed: %v", err)
	}

	// Mock FindVolumeByName to return error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume not found")).Once()
	volNameTemp = "dummy_volume_1"
	_, err = testConf.volumeAPI.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with invalid name case failed: %v", err)
	}

	fmt.Println("Find Volume by Name Test - Successful")
}

func TestFindVolumeByID(t *testing.T) {
	fmt.Println("Begin - Find Volume By Name Test")

	// Mock FindVolumeByID to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	vol, err := testConf.volumeAPI.FindVolumeByID(ctx, volID)
	fmt.Println("Find volume by Name:", prettyPrintJSON(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Id failed: %v", err)
	}

	// Negative cases
	volIDTemp := ""
	_, err = testConf.volumeAPI.FindVolumeByID(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with empty Id case failed: %v", err)
	}

	// Mock FindVolumeByID to return error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume not found")).Once()
	volIDTemp = "dummy_vol_sv_1"
	_, err = testConf.volumeAPI.FindVolumeByID(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with invalid Id case failed: %v", err)
	}
	fmt.Println("Find Volume by Id Test - Successful")
}

func TestListVolumes(t *testing.T) {
	fmt.Println("Begin - List Volumes Test")

	// Mock ListVolumes to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*types.ListVolumes)
		resp.Volumes = make([]types.Volume, 10)
	}).Once()
	vols, _, err := testConf.volumeAPI.ListVolumes(ctx, 11, 10)
	fmt.Println("List volumes count: ", len(vols))
	if len(vols) <= 10 {
		fmt.Println("List volume success")
	} else {
		t.Fatalf("List volumes failed: %v", err)
	}

	fmt.Println("List Volume Test - Successful")
}

func TestExportVolume(t *testing.T) {
	fmt.Println("Begin - Export Volume Test")

	// Mock FindHostByName to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	host, err := testConf.hostAPI.FindHostByName(ctx, testConf.nodeHostName)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	// Mock ExportVolume to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.volumeAPI.ExportVolume(ctx, volID, host.HostContent.ID)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}

	// Negative case for Delete Volume
	// Mock executeWithRetryAuthenticate to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock FindVolumeByID to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock DeleteVolume to return error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("failed to delete exported volume")).Once()
	err = testConf.volumeAPI.DeleteVolume(ctx, volID)
	if err == nil {
		t.Fatalf("Delete volume on exported volume case failed: %v", err)
	}

	fmt.Println("Export Volume Test - Successful")
}

func TestUnexportVolume(t *testing.T) {
	fmt.Println("Begin - Unexport Volume Test")

	// Mock UnexportVolume to return no error
	testConf.volumeAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err := testConf.volumeAPI.UnexportVolume(ctx, volID)
	if err != nil {
		t.Fatalf("UnExportVolume failed: %v", err)
	}
	fmt.Println("Unexport Volume Test - Successful")
}

func expandVolumeTest(t *testing.T) {
	fmt.Println("Begin - Expand Volume Test")

	err := testConf.volumeAPI.ExpandVolume(ctx, volID, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume failed: %v", err)
	}

	err = testConf.volumeAPI.ExpandVolume(ctx, volID, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume with same size failed: %v", err)
	}

	// Negative cases
	volIDTemp := "dummy_vol_sv_1"
	err = testConf.volumeAPI.ExpandVolume(ctx, volIDTemp, 5368709120)
	if err == nil {
		t.Fatalf("Expand volume with invalid Id case failed: %v", err)
	}

	err = testConf.volumeAPI.ExpandVolume(ctx, volID, 4368709120)
	if err == nil {
		t.Fatalf("Expand volume with smaller size case failed: %v", err)
	}

	fmt.Println("Expand Volume Test - Successful")
}

func createCloneFromVolumeTest(t *testing.T) {
	fmt.Println("Begin - Create clone from Volume Test")

	_, err := testConf.volumeAPI.CreateCloneFromVolume(ctx, cloneVolumeName, volID)
	if err != nil {
		t.Fatalf("Clone volume failed: %v", err)
	}

	vol, err := testConf.volumeAPI.FindVolumeByName(ctx, cloneVolumeName)
	fmt.Println("Find volume by Name:", prettyPrintJSON(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Name failed: %v", err)
	}

	cloneVolumeID = vol.VolumeContent.ResourceID

	// Negative Test Case
	// Creating clone with same name
	_, err = testConf.volumeAPI.CreateCloneFromVolume(ctx, cloneVolumeName, volID)
	if err == nil {
		t.Fatalf("Clone volume with a same existing volume name test case failed: %v", err)
	}

	// Creating clone with invalid volume ID
	volIDTemp := "dummy-vol-1"
	_, err = testConf.volumeAPI.CreateCloneFromVolume(ctx, cloneVolumeName, volIDTemp)
	if err == nil {
		t.Fatalf("Clone volume with invalid volume ID test case failed: %v", err)
	}

	fmt.Println("Create clone from Volume Test - Successful")
}

func modifyVolumeExportTest(t *testing.T) {
	fmt.Println("Begin - Modify Volume Export Test")

	hostIDList := []string{}
	for _, hostName := range testConf.hostList {
		host, err := testConf.hostAPI.FindHostByName(ctx, hostName)
		if err != nil {
			t.Fatalf("Find host by name %s failed. Error: %v", hostName, err)
		}
		hostIDList = append(hostIDList, host.HostContent.ID)
	}

	err := testConf.volumeAPI.ModifyVolumeExport(ctx, volID, hostIDList)
	if err != nil {
		t.Fatalf("Modify Volume Export failed: %v", err)
	}

	// Modify Volume name
	volName = volName + "_renamed"
	err = testConf.volumeAPI.RenameVolume(ctx, volName, volID)
	if err != nil {
		t.Fatalf("Rename existing volume failed. Error: %v", err)
	}

	// Negative Test case

	volIDTemp := "dummy_vol_1"
	err = testConf.volumeAPI.RenameVolume(ctx, volName, volIDTemp)
	if err == nil {
		t.Fatalf("Rename existing volume failed. Error: %v", err)
	}

	// Unexport volume from host
	err = testConf.volumeAPI.UnexportVolume(ctx, volID)
	if err != nil {
		t.Fatalf("Unexport volume failed. Error: %v", err)
	}
	fmt.Println("Modify Volume Export Test Successful")
}

func deleteVolumeTest(t *testing.T) {
	fmt.Println("Begin - Delete Volume Test")

	// Deletion of volume, Volume won't get deleted as clone exists
	err := testConf.volumeAPI.DeleteVolume(ctx, volID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	// Deletion of clone and volume
	err = testConf.volumeAPI.DeleteVolume(ctx, cloneVolumeID)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}

	// Negative cases
	volIDTemp := ""
	err = testConf.volumeAPI.DeleteVolume(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Delete volume with empty Id case failed: %v", err)
	}

	volIDTemp = "dummy_vol_sv_1"
	err = testConf.volumeAPI.DeleteVolume(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Delete volume with invalid Id case failed: %v", err)
	}

	fmt.Println("Delete Volume Test - Successful")
}

func getMaxVolumeSizeTest(t *testing.T) {
	fmt.Println("Begin - Get Max Volume Size")

	// Positive case
	systemLimitID := "Limit_MaxLUNSize"
	_, err := testConf.volumeAPI.GetMaxVolumeSize(ctx, systemLimitID)
	if err != nil {
		t.Fatalf("Get maximum volume size failed: %v", err)
	}

	// Negative cases
	systemLimitID = ""
	_, err = testConf.volumeAPI.GetMaxVolumeSize(ctx, systemLimitID)
	if err == nil {
		t.Fatalf("Get maximum volume size with empty systemLimitID case failed: %v", err)
	}

	systemLimitID = "dummy_name"
	_, err = testConf.volumeAPI.GetMaxVolumeSize(ctx, systemLimitID)
	if err == nil {
		t.Fatalf("Get maximum volume size with invalid systemLimitID case failed: %v", err)
	}
	fmt.Println("Get Max Volume Size - Successful")
}
