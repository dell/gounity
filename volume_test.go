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
	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	volName         = "unit-test-vol"
	cloneVolumeName = "unit-test-clone-vol"
	volID           = "unity-volume-id"
	hostIOLimitID   string
	anyArgs         = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
)

func TestFindHostIOLimitByName(t *testing.T) {
	fmt.Println("Begin - Find Host IO Limit by Name Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock the client.DoWithHeaders to return nil
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()

	// Call the FindHostIOLimitByName function
	hostIOLimit, err := testConf.client.FindHostIOLimitByName(ctx, testConf.hostIOLimitName)
	fmt.Println("hostIOLimit:", prettyPrintJSON(hostIOLimit), "Error:", err)
	assert.NotNil(t, hostIOLimit.IoLimitPolicyContent)

	// Negative cases

	// Mock the client.DoWithHeaders to return an error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("not found")).Once()

	// Call the FindHostIOLimitByName function with a dummy name
	_, err = testConf.client.FindHostIOLimitByName(ctx, "dummy_hostio_1")
	if err == nil {
		t.Fatalf("Find Host IO Limit negative case failed: %v", err)
	}

	// Call the FindHostIOLimitByName function with an empty name
	_, err = testConf.client.FindHostIOLimitByName(ctx, "")
	if err == nil {
		t.Fatalf("Find Host IO Limit with empty name case failed: %v", err)
	}

	fmt.Println("Find Host IO Limit by Name Test - Successful")
}

func TestCreateLun(t *testing.T) {
	fmt.Println("Begin - Create LUN Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil

	ctx := context.Background()

	// Mock FindStoragePoolByID to return nil
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock isFeatureLicensed to return expected response
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.LicenseInfo")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.LicenseInfo)
			*resp = types.LicenseInfo{LicenseInfoContent: types.LicenseInfoContent{IsInstalled: true, IsValid: true}}
		}).Twice()
	// Mock create request to return nil
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()

	_, err := testConf.client.CreateLun(ctx, volName, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err != nil {
		t.Fatalf("Create LUN failed: %v", err)
	}

	// Negative cases
	volNameTemp := ""
	_, err = testConf.client.CreateLun(ctx, volNameTemp, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with empty name case failed: %v", err)
	}

	volNameTemp = "vol-name-max-length-12345678901234567890123456789012345678901234567890"
	_, err = testConf.client.CreateLun(ctx, volNameTemp, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN exceeding max name length case failed: %v", err)
	}

	// Mock FindStoragePoolByID to return error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("storage pool not found")).Once()
	poolIDTemp := "dummy_pool_1"
	_, err = testConf.client.CreateLun(ctx, volName, poolIDTemp, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with invalid pool name case failed: %v", err)
	}

	// Mock FindStoragePoolByID to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock isFeatureLicensed to return expected response
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.LicenseInfo")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.LicenseInfo)
			*resp = types.LicenseInfo{LicenseInfoContent: types.LicenseInfoContent{IsInstalled: true, IsValid: true}}
		}).Twice()
	// Mock create volume to return error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume already exists")).Once()
	_, err = testConf.client.CreateLun(ctx, volName, testConf.poolID, "Description", 2368709120, 0, hostIOLimitID, true, false)
	if err == nil {
		t.Fatalf("Create LUN with same name case failed: %v", err)
	}

	fmt.Println("Create LUN Test - Successful")
}

func TestFindVolumeByName(t *testing.T) {
	fmt.Println("Begin - Find Volume By Name Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	// Mock FindVolumeByName to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	vol, err := testConf.client.FindVolumeByName(ctx, volName)
	fmt.Println("Find volume by Name:", prettyPrintJSON(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Name failed: %v", err)
	}
	assert.NotNil(t, vol.VolumeContent.ResourceID)

	// Negative cases
	volNameTemp := ""
	_, err = testConf.client.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with empty name case failed: %v", err)
	}

	// Mock FindVolumeByName to return error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume not found")).Once()
	volNameTemp = "dummy_volume_1"
	_, err = testConf.client.FindVolumeByName(ctx, volNameTemp)
	if err == nil {
		t.Fatalf("Find volume by Name with invalid name case failed: %v", err)
	}

	fmt.Println("Find Volume by Name Test - Successful")
}

func TestFindVolumeByID(t *testing.T) {
	fmt.Println("Begin - Find Volume By Name Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock FindVolumeByID to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	vol, err := testConf.client.FindVolumeByID(ctx, volID)
	fmt.Println("Find volume by Name:", prettyPrintJSON(vol), err)
	if err != nil {
		t.Fatalf("Find volume by Id failed: %v", err)
	}

	// Negative cases
	volIDTemp := ""
	_, err = testConf.client.FindVolumeByID(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with empty Id case failed: %v", err)
	}

	// Mock FindVolumeByID to return error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume not found")).Once()
	volIDTemp = "dummy_vol_sv_1"
	_, err = testConf.client.FindVolumeByID(ctx, volIDTemp)
	if err == nil {
		t.Fatalf("Find volume by Id with invalid Id case failed: %v", err)
	}
	fmt.Println("Find Volume by Id Test - Successful")
}

func TestListVolumes(t *testing.T) {
	fmt.Println("Begin - List Volumes Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	// Mock ListVolumes to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*types.ListVolumes)
		resp.Volumes = make([]types.Volume, 10)
	}).Once()
	vols, _, err := testConf.client.ListVolumes(ctx, 11, 10)
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
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock FindHostByName to return a valid host object
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("FindHostByName", ctx, testConf.nodeHostName).Return(&types.Host{
		HostContent: types.HostContent{
			ID: "valid_host_id",
		},
	}, nil)

	// Find the host
	host, err := testConf.client.FindHostByName(ctx, testConf.nodeHostName)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	// Mock ExportVolume to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	err = testConf.client.ExportVolume(ctx, volID, host.HostContent.ID)
	if err != nil {
		t.Fatalf("ExportVolume failed: %v", err)
	}

	// Negative case for Delete Volume
	// Mock executeWithRetryAuthenticate to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// Mock FindVolumeByID to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// Mock DeleteVolume to return a specific error type
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.Error{
		ErrorContent: types.ErrorContent{
			Message: []types.ErrorMessage{
				{EnUS: "failed to delete exported volume"},
			},
			HTTPStatusCode: 500,
			ErrorCode:      1234,
		},
	}).Once()
	err = testConf.client.DeleteVolume(ctx, volID)
	if err == nil {
		t.Fatalf("Delete volume on exported volume case failed: %v", err)
	}

	fmt.Println("Export Volume Test - Successful")
}

func TestUnexportVolume(t *testing.T) {
	fmt.Println("Begin - Unexport Volume Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock UnexportVolume to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := testConf.client.UnexportVolume(ctx, volID)
	if err != nil {
		t.Fatalf("UnexportVolume failed: %v", err)
	}

	fmt.Println("Unexport Volume Test - Successful")
}

func TestExpandVolumeTest(t *testing.T) {
	fmt.Println("Begin - Expand Volume Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("FindVolumeByID", ctx, "volID").Return(&types.VolumeContent{SizeTotal: 5368709120}, nil).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("FindVolumeByID", ctx, "dummy_vol_sv_1").Return(nil, errors.New("unable to find volume Id")).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("executeWithRetryAuthenticate", ctx, "POST", "api/modify/lun/volID", mock.Anything, nil).Return(nil).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	err := testConf.client.ExpandVolume(ctx, "volID", 5368709120)
	if err != nil {
		t.Fatalf("Expand volume failed: %v", err)
	}

	err = testConf.client.ExpandVolume(ctx, volID, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume with same size failed: %v", err)
	}

	// Negative cases
	volIDTemp := "dummy_vol_sv_1"
	err = testConf.client.ExpandVolume(ctx, volIDTemp, 5368709120)
	if err != nil {
		t.Fatalf("Expand volume with invalid Id case failed: %v", err)
	}

	err = testConf.client.ExpandVolume(ctx, volID, 4368709120)
	if err != nil {
		t.Fatalf("Expand volume with smaller size case failed: %v", err)
	}

	fmt.Println("Expand Volume Test - Successful")
}

func TestCreateCloneFromVolume(t *testing.T) {
	fmt.Println("Begin - Create clone from Volume Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock responses for CreateSnapshot, CreteLunThinClone
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
	_, err := testConf.client.CreateCloneFromVolume(ctx, cloneVolumeName, volID)
	assert.Nil(t, err)

	// Negative Test Case: Create Snapshot Failed
	_, err = testConf.client.CreateCloneFromVolume(ctx, cloneVolumeName, "")
	assert.ErrorIs(t, err, ErrorCreateSnapshotFailed)

	// Negative Test Case: Creating clone with same name
	// Mock responses for CreateSnapshot
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// Mock CreteLunThinCloneto return error volume with a same exists
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("volume with a same exists")).Once()
	// Mock responses for DeleteSnapshot
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	_, err = testConf.client.CreateCloneFromVolume(ctx, cloneVolumeName, volID)
	assert.ErrorIs(t, err, ErrorCloningFailed)

	fmt.Println("Create clone from Volume Test - Successful")
}

func TestModifyVolumeExportTest(t *testing.T) {
	fmt.Println("Begin - Modify Volume Export Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Clear existing expectations
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil

	// Mock the DoWithHeaders method to handle multiple calls
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// Mock the RenameVolume method to return an error for non-existent volume ID
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, "POST", "/api/instances/storageResource/dummy_vol_1/action/modifyLun", mock.Anything, mock.Anything).Return(fmt.Errorf("volume not found")).Once()

	// Create a list of host IDs
	hostIDList := []string{}
	for _, hostName := range testConf.hostList {
		host, err := testConf.client.FindHostByName(ctx, hostName)
		if err != nil {
			t.Fatalf("Find host by name %s failed. Error: %v", hostName, err)
		}
		hostIDList = append(hostIDList, host.HostContent.ID)
	}

	// Modify the volume export
	if err := testConf.client.ModifyVolumeExport(ctx, volID, hostIDList); err != nil {
		t.Fatalf("Modify Volume Export failed: %v", err)
	}

	// Rename the volume
	volName += "_renamed"
	if err := testConf.client.RenameVolume(ctx, volName, volID); err != nil {
		t.Fatalf("Rename existing volume failed. Error: %v", err)
	}

	// Negative test case: Attempt to rename a non-existent volume
	volIDTemp := "dummy_vol_1"
	if err := testConf.client.RenameVolume(ctx, volName, volIDTemp); err != nil {
		t.Fatalf("Expected error when renaming non-existent volume, got none")
	}

	// Unexport the volume from the host
	if err := testConf.client.UnexportVolume(ctx, volID); err != nil {
		t.Fatalf("Unexport volume failed. Error: %v", err)
	}

	fmt.Println("Modify Volume Export Test Successful")
}

func TestDeleteVolumeTest(t *testing.T) {
	fmt.Println("Begin - Delete Volume Test")
	ctx := context.Background()

	// Clear existing expectations
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil

	err := testConf.client.DeleteVolume(ctx, "")
	assert.ErrorContains(t, err, "Volume Id cannot be empty")

	// Mock the executeWithRetryAuthenticate, FindVolumeByID, executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Times(3)
	err = testConf.client.DeleteVolume(ctx, volID)
	assert.Nil(t, err)

	// Mock the executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock the FindVolumeByID method to return voume not found error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("volume not found")).Once()
	err = testConf.client.DeleteVolume(ctx, volID)
	assert.Errorf(t, err, "volume not found")

	// Mock the executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	// Mock the FindVolumeByID method to return voume not found error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New(DependentClonesErrorCode)).Once()
	// Mock the RenameVolume method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.client.DeleteVolume(ctx, volID)
	assert.Nil(t, err)

	// Mock the executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock the FindVolumeByID method to return expected response
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.Volume)
			*resp = types.Volume{
				VolumeContent: types.VolumeContent{
					IsThinClone:  true,
					ParentVolume: types.StorageResource{ID: volID},
				},
			}
		}).Once()
	// Mock the executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	// Mock the sourceVolID FindVolumeByID method to return expected response
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.Volume)
			*resp = types.Volume{
				VolumeContent: types.VolumeContent{
					Name: MarkVolumeForDeletion,
				},
			}
		}).Once()
	// Mock the deleteSourceVol executeWithRetryAuthenticate method to return no error
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.client.DeleteVolume(ctx, volID)
	assert.Nil(t, err)

	fmt.Println("Delete Volume Test - Successful")
}

func TestGetMaxVolumeSizeTest(t *testing.T) {
	fmt.Println("Begin - Get Max Volume Size")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock the DoWithHeaders method to handle multiple calls
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// Mock the GetMaxVolumeSize method to return an error for invalid systemLimitID
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/types/systemLimit/instances/dummy_name", mock.Anything, mock.Anything).Return(fmt.Errorf("invalid systemLimitID")).Once()

	// Positive case: Get maximum volume size with a valid system limit ID
	systemLimitID := "Limit_MaxLUNSize"
	if _, err := testConf.client.GetMaxVolumeSize(ctx, systemLimitID); err != nil {
		t.Fatalf("Get maximum volume size failed: %v", err)
	}

	// Negative case: Attempt to get maximum volume size with an empty system limit ID
	systemLimitID = ""
	if _, err := testConf.client.GetMaxVolumeSize(ctx, systemLimitID); err == nil {
		t.Fatalf("Expected error when getting maximum volume size with empty systemLimitID, got none")
	}

	// Negative case: Attempt to get maximum volume size with an invalid system limit ID
	systemLimitID = "dummy_name"
	if _, err := testConf.client.GetMaxVolumeSize(ctx, systemLimitID); err != nil {
		t.Fatalf("Expected error when getting maximum volume size with invalid systemLimitID, got none")
	}

	fmt.Println("Get Max Volume Size - Successful")
}
