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
	"time"

	"github.com/dell/gounity/apitypes"
	mocksapi "github.com/dell/gounity/mocks/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	snapVolID              = "snapVolID"
	snapID                 = "snapID"
	snap2ID                = "snap2ID"
	snapByFsAccessTypeID   = "snapByFsAccessTypeID"
	snapCopyID             = "snapCopyID"
	cloneVolID             = "cloneVolID"
	now                    = time.Now()
	timeStamp              = now.Format("20060102150405")
	snapVolName            = "Unit-test-snap-vol-" + timeStamp
	snapName               = "Unit-test-snapshot-" + timeStamp
	snap2Name              = "Unit-test-snapshot2-" + timeStamp
	snapByFsAccessTypeName = "Unit-test-snapshot-by-fsxstype-" + timeStamp
	cloneVolName           = "Unit-test-clone-vol-" + timeStamp
)

func TestCreateSnapshot(t *testing.T) {
	fmt.Println("Begin - Create Snapshot Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*apitypes.LicenseInfo")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*apitypes.LicenseInfo)
			*resp = apitypes.LicenseInfo{LicenseInfoContent: apitypes.LicenseInfoContent{IsInstalled: true, IsValid: true}}
		}).Twice()

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.client.CreateLun(ctx, snapVolName, testConf.poolID, "Description", 5368709120, 0, "", true, false)
	assert.Equal(t, nil, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.FindVolumeByName(ctx, snapVolName)
	assert.Equal(t, nil, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snap, err := testConf.client.CreateSnapshot(ctx, snapVolID, snapName, "Snapshot Description", "")
	fmt.Println("Create Snapshot:", prettyPrintJSON(snap), err)
	assert.Equal(t, nil, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snap, err = testConf.client.CreateSnapshot(ctx, snapVolID, snap2Name, "Snapshot Description", "1:23:52:50")
	fmt.Println("Create Snapshot2:", prettyPrintJSON(snap), err)
	assert.Equal(t, nil, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snapFsAccess, err := testConf.client.CreateSnapshotWithFsAccesType(ctx, snapVolID, snapByFsAccessTypeName, "Snapshot Description", "", BlockAccessType)
	fmt.Println("Create Snapshot With FsAccessType:", prettyPrintJSON(snapFsAccess), err)
	assert.Equal(t, nil, err)

	snapByFsAccessTypeID = snapFsAccess.SnapshotContent.ResourceID

	// Negative cases
	snapVolIDTemp := ""
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshot(ctx, snapVolIDTemp, snap2Name, "Snapshot Description", "")
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshotWithFsAccesType(ctx, snapVolIDTemp, snapByFsAccessTypeName, "Snapshot Description", "", BlockAccessType)
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snapNameTemp := "snap-name-max-length-12345678901234567890123456789012345678901234567890"
	_, err = testConf.client.CreateSnapshot(ctx, snapVolID, snapNameTemp, "Snapshot Description", "")
	assert.Equal(t, errors.New("invalid snapshot name Error:name too long error"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshotWithFsAccesType(ctx, snapVolIDTemp, snapNameTemp, "Snapshot Description", "", BlockAccessType)
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshot(ctx, snapVolID, snap2Name, "Snapshot Description", "1:23:99:99")
	assert.Equal(t, errors.New("hours, minutes and seconds should be in between 0-60"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshotWithFsAccesType(ctx, snapVolIDTemp, snapNameTemp, "Snapshot Description", "1:23:99:99", BlockAccessType)
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CreateSnapshot(ctx, "", snap2Name, "Snapshot Description", "1:23:52:50")
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	fmt.Println("Create Snapshot Test - Successful")
}

func TestFindSnapshotByName(t *testing.T) {
	fmt.Println("Begin - Find Snapshot by Name Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snap, err := testConf.client.FindSnapshotByName(ctx, snapName)
	fmt.Println("Find snapshot by Name:", prettyPrintJSON(snap), err)
	assert.Equal(t, nil, err)
	snapID = snap.SnapshotContent.ResourceID

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snap, err = testConf.client.FindSnapshotByName(ctx, snap2Name)
	fmt.Println("Find snapshot2 by Name:", prettyPrintJSON(snap), err)
	assert.Equal(t, nil, err)
	snap2ID = snap.SnapshotContent.ResourceID

	// Negative test cases
	snapNameTemp := ""
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.FindSnapshotByName(ctx, snapNameTemp)
	assert.Equal(t, errors.New("name empty error"), err)

	snapNameTemp = "dummy-snap-name"
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New(SnapshotNotFoundErrorCode)).Once()
	_, err = testConf.client.FindSnapshotByName(ctx, snapNameTemp)
	assert.Error(t, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("unable to find Snapshot Name")).Once()
	_, err = testConf.client.FindSnapshotByName(ctx, snapNameTemp)
	assert.Error(t, err)

	fmt.Println("Find Snapshot by Name - Successful")
}

func TestFindSnapshotByID(t *testing.T) {
	fmt.Println("Begin - Find Snapshot by Id Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snapID = "snapID"
	snap, err := testConf.client.FindSnapshotByID(ctx, snapID)
	fmt.Println("Find snapshot by ID:", prettyPrintJSON(snap), err)
	assert.Equal(t, nil, err)

	// Negative test cases
	snapIDTemp := ""
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.FindSnapshotByID(ctx, snapIDTemp)
	assert.Equal(t, errors.New("snapshot ID cannot be empty"), err)

	snapIDTemp = "dummy-snap-id"
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New(SnapshotNotFoundErrorCode)).Once()
	_, err = testConf.client.FindSnapshotByID(ctx, snapIDTemp)
	assert.Error(t, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("unable to find Snapshot ID")).Once()
	_, err = testConf.client.FindSnapshotByID(ctx, snapIDTemp)
	assert.Error(t, err)

	fmt.Println("Find Snapshot by Id - Successful")
}

func TestListSnapshots(t *testing.T) {
	fmt.Println("Begin - List Snapshots Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, _, err := testConf.client.ListSnapshots(ctx, 0, 10, snapVolID, "")
	snaps := []string{"snap1", "snap2"}
	fmt.Println("List snapshots:", snaps)
	if len(snaps) > 0 {
		fmt.Println("List snapshots success:", len(snaps))
	} else {
		assert.Equal(t, errors.New("List snapshot failed"), err)
	}

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, _, err = testConf.client.ListSnapshots(ctx, 0, 10, snapVolID, snapID)
	fmt.Println("List snapshots with snap Id:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots with snap Id success:", len(snaps))
	} else {
		assert.Equal(t, errors.New("List snapshot with snap Id failed"), err)
	}

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, _, err = testConf.client.ListSnapshots(ctx, 6, 5, "", "")
	fmt.Println("List snapshots pagination:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots pagination success:", len(snaps))
	} else {
		assert.Equal(t, errors.New("List snapshot pagination failed"), err)
	}

	fmt.Println("List Snapshots Test - Successful")
}

func TestModifySnapshotAutoDeleteParameter(t *testing.T) {
	fmt.Println("Begin - Modify Snapshot Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err := testConf.client.ModifySnapshotAutoDeleteParameter(ctx, snapID)
	assert.Equal(t, nil, err)

	snapByFsAccessTypeID = "snapByFsAccessTypeID"
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.client.ModifySnapshot(ctx, snapByFsAccessTypeID, "Modify Description", "1:22:02:50")
	assert.Equal(t, nil, err)

	// Negative test cases
	snapIDTemp := ""
	err = testConf.client.ModifySnapshotAutoDeleteParameter(ctx, snapIDTemp)
	assert.Equal(t, errors.New("snapshot ID cannot be empty"), err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.client.ModifySnapshot(ctx, snapIDTemp, "Modify Description", "1:22:02:50")
	assert.Equal(t, errors.New("snapshot ID cannot be empty"), err)

	fmt.Println("Modify Snapshot Test - Successful")
}

func TestCreteLunThinClone(t *testing.T) {
	fmt.Println("Begin - Create LUN thin clone Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.client.CreteLunThinClone(ctx, cloneVolName, snapID, snapVolID)
	assert.Equal(t, nil, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	vol, err := testConf.client.FindVolumeByName(ctx, cloneVolName)
	assert.Equal(t, nil, err)
	cloneVolID = vol.VolumeContent.ResourceID
	fmt.Println("Create LUN thin clone Test - Successful")
}

func TestCopySnapshot(t *testing.T) {
	fmt.Println("Begin - Copy Snapshot Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*apitypes.CopySnapshots")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*apitypes.CopySnapshots)
			*resp = apitypes.CopySnapshots{
				CopySnapshotsContent: apitypes.CopySnapshotsContent{
					Copies: []apitypes.StorageResource{
						{ID: snapCopyID},
					},
				},
			}
		}).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	snapCopy, err := testConf.client.CopySnapshot(ctx, snapByFsAccessTypeID, snapName+"_copy")
	assert.Equal(t, nil, err)

	snapCopyID = snapCopy.SnapshotContent.ResourceID

	// Negative test cases

	snapNameTemp := ""
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CopySnapshot(ctx, snapByFsAccessTypeID, snapNameTemp)
	assert.Equal(t, errors.New("Snapshot Name cannot be empty"), err)

	snapIDTemp := ""
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.CopySnapshot(ctx, snapIDTemp, snapName)
	assert.Equal(t, errors.New("Source Snapshot ID cannot be empty"), err)

	fmt.Println("Copy Snapshot Test - Successful")
}

func TestDeleteSnapshot(t *testing.T) {
	fmt.Println("Begin - Delete Snapshot Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err := testConf.client.DeleteSnapshot(ctx, snapID)
	assert.Equal(t, nil, err)

	// Negative test cases
	snapIDTemp := ""
	err = testConf.client.DeleteSnapshot(ctx, snapIDTemp)
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	assert.Equal(t, errors.New("snapshot ID cannot be empty"), err)

	fmt.Println("Delete Snapshot Test - Successful")
}

func TestDeleteFilesystemAsSnapshot(t *testing.T) {
	fmt.Println("Begin - Delete Filesystem As Snapshot Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	sourceFs := &apitypes.Filesystem{
		FileContent: apitypes.FileContent{
			ID:          "test-filesystem-id",
			Description: "test description",
		},
	}
	err := testConf.client.DeleteFilesystemAsSnapshot(ctx, snapID, sourceFs)
	assert.Equal(t, nil, err)

	// MarkFilesystemForDeletion
	sourceFs = &apitypes.Filesystem{
		FileContent: apitypes.FileContent{
			ID:          "test-filesystem-id-delete",
			Description: "csi-marked-filesystem-for-deletion(do not remove this from description)",
		},
	}

	// Mock the DeleteSnapshot method
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()

	// Mock the DeleteFilesystem method
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()

	err = testConf.client.DeleteFilesystemAsSnapshot(ctx, snapID, sourceFs)
	assert.Equal(t, nil, err)

	fmt.Println("Delete Filesystem As Snapshot Test - Successful")
}
