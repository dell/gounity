/*
 Copyright © 2020-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

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

	"github.com/dell/gounity/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	fsName            string
	fsID              string
	nfsShareName      string
	nfsShareID        string
	storageResourceID string
	snapshotID        string
	nfsShareIDBySnap  string
)

const (
	NFSShareLocalPath  = "/"
	NFSShareNamePrefix = "csishare-"
)

func TestFindNasServer(t *testing.T) {
	fmt.Println("Begin - Find Nas Server Test")
	ctx := context.Background()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.fileAPI.FindNASServerByID(ctx, testConf.nasServer)
	if err != nil {
		t.Fatalf("Find filesystem by name failed: %v", err)
	}

	// Test case :  GET using invalid ID
	nasServer := "nas_dummy_1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err = testConf.fileAPI.FindNASServerByID(ctx, nasServer)
	if err == nil {
		t.Fatal("Find Nas Server - Negative case failed")
	}

	// Test case :  GET using empty ID
	nasServer = ""
	_, err = testConf.fileAPI.FindNASServerByID(ctx, nasServer)
	assert.Equal(t, errors.New("NAS Server Id shouldn't be empty"), err)
	fmt.Println("Find Nas Server Test Successful")
}

func TestCreateFilesystem(t *testing.T) {
	fmt.Println("Begin - Create Filesystem Test")
	ctx := context.Background()
	fsName = ""
	_, err := testConf.fileAPI.CreateFilesystem(ctx, fsName, testConf.poolID, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	assert.Equal(t, errors.New("filesystem name should not be empty"), err)

	// Negative cases
	fsNameTemp := "dummy-fs-1234567890123456789012345678901234567890123456789012345678"
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsNameTemp, testConf.poolID, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	assert.Equal(t, errors.New("filesystem name dummy-fs-1234567890123456789012345678901234567890123456789012345678 should not exceed 63 characters"), err)

	poolIDTemp := "dummy_pool_1"
	fsName = "xfs"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsName, poolIDTemp, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	assert.Equal(t, errors.New("thin provisioning is not supported on array and hence cannot create Filesystem"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsName, poolIDTemp, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, false, true)
	assert.Equal(t, errors.New("data reduction is not supported on array and hence cannot create Filesystem"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsName, poolIDTemp, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, false, false)
	assert.Equal(t, nil, err)

	fmt.Println("Create Filesystem test successful")
}

func TestFindFilesystem(t *testing.T) {
	fmt.Println("Begin - Find Filesystem Test")
	ctx := context.Background()
	fsName = ""
	_, err := testConf.fileAPI.FindFilesystemByName(ctx, fsName)
	assert.Equal(t, errors.New("Filesystem Name shouldn't be empty"), err)

	_, err = testConf.fileAPI.FindFilesystemByID(ctx, "")
	assert.Equal(t, errors.New("Filesystem Id shouldn't be empty"), err)

	_, err = testConf.fileAPI.GetFilesystemIDFromResID(ctx, "")
	assert.Equal(t, errors.New("Filesystem Resource Id shouldn't be empty"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.fileAPI.GetFilesystemIDFromResID(ctx, "ID")
	assert.Equal(t, nil, err)

	// Test case :  GET using invalid fsName/ID
	fsNameTemp := "dummy-fs-1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.fileAPI.FindFilesystemByName(ctx, fsNameTemp)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.fileAPI.FindFilesystemByID(ctx, fsNameTemp)
	assert.Equal(t, nil, err)

	fmt.Println("Find Filesystem test successful")
}

func TestCreateNfsShare(t *testing.T) {
	fmt.Println("Begin - Create NFS Share Test")
	ctx := context.Background()

	_, err := testConf.fileAPI.CreateNFSShare(ctx, nfsShareName, NFSShareLocalPath, fsID, NoneDefaultAccess)
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	// Test case : Create NFS share using snapshot
	_, err = testConf.snapAPI.CreateSnapshot(ctx, storageResourceID, "snapName", "Snapshot Description", "")
	assert.Equal(t, errors.New("storage Resource ID cannot be empty"), err)

	snapshotID = ""
	_, err = testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, snapshotID, NoneDefaultAccess)
	assert.Equal(t, errors.New("Snapshot Id cannot be empty"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	_, err = testConf.fileAPI.CreateNFSShare(ctx, nfsShareName, NFSShareLocalPath, "fsID", NoneDefaultAccess)
	if err != nil {
		t.Fatalf("Create NFS Share Negative scenario failed: %v", err)
	}

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, "snapshotID", NoneDefaultAccess)
	if err != nil {
		t.Fatalf("Create NFS Share from snapshot negative case failed: %v", err)
	}

	fmt.Println("Create NFS Share Test Successful")
}

func TestFindNfsShare(t *testing.T) {
	fmt.Println("Begin - Find NFS Share Test")
	ctx := context.Background()
	_, err := testConf.fileAPI.FindNFSShareByName(ctx, nfsShareName)
	assert.Equal(t, errors.New("NFS Share Name shouldn't be empty"), err)

	_, err = testConf.fileAPI.FindNFSShareByID(ctx, nfsShareID)
	assert.Equal(t, errors.New("NFS Share Id shouldn't be empty"), err)

	// Test case :  GET using invalid shareName/ID
	nfsShareNameTemp := "dummy-fs-1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	_, err = testConf.fileAPI.FindNFSShareByName(ctx, nfsShareNameTemp)
	assert.Equal(t, nil, err)

	fmt.Println("Find NFS Share Test Successful")
}

func TestModifyNfsShare(t *testing.T) {
	fmt.Println("Begin - Modify NFS Share Test")
	ctx := context.Background()
	testConf.hostAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.hostAPI.FindHostByName(ctx, testConf.nodeHostName)
	if err != nil {
		t.Fatalf("Find host failed: %v", err)
	}

	var hostIDList []string
	hostIDList = append(hostIDList, "host.HostContent.ID")

	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsID, nfsShareID, hostIDList, ReadOnlyAccessType)
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, "nfsShareIDBySnap", []string{"host1", "host2"}, ReadOnlyAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, "nfsShareIDBySnap", []string{"host1", "host2"}, ReadWriteRootAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, "nfsShareIDBySnap", []string{"host1", "host2"}, ReadOnlyRootAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, "nfsShareIDBySnap", []string{"host1", "host2"}, ReadWriteAccessType)
	assert.Equal(t, nil, err)

	fsIDTemp := "dummy-fs-1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadWriteRootAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadOnlyRootAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadWriteAccessType)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadOnlyAccessType)
	assert.Equal(t, nil, err)

	fsIDTemp = ""
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadWriteRootAccessType)
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	nfsShareIDBySnapTemp := "dummy-nsf-share-1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnapTemp, hostIDList, ReadOnlyAccessType)
	assert.Equal(t, nil, err)

	fmt.Println("Modify NFS Share Test Successful")
}

func TestDescription(t *testing.T) {
	fmt.Println("Begin - Update Description of Filesystem Test")
	ctx := context.Background()
	// Positive scenario is covered under DeleteFilesystemTest()
	// Negative test case

	filesystemIDTemp := ""
	err := testConf.fileAPI.updateDescription(ctx, filesystemIDTemp, "Description of filesystem")
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	filesystemIDTemp = "dummy_fs_1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.updateDescription(ctx, filesystemIDTemp, "Description of filesystem")
	assert.Equal(t, nil, err)
}

func TestDeleteNfsShare(t *testing.T) {
	fmt.Println("Begin - Delete NFS Share Test")
	ctx := context.Background()
	// Test case :  Delete using invalid shareID and fsID
	nfsShareIDTemp := "dummy-fs-1"
	fsIDTemp := "dummy-fs-1"

	err := testConf.fileAPI.DeleteNFSShare(ctx, fsID, nfsShareIDTemp)
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.DeleteNFSShare(ctx, fsIDTemp, nfsShareIDTemp)
	assert.Equal(t, nil, err)

	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.DeleteNFSShareCreatedFromSnapshot(ctx, nfsShareIDTemp)
	assert.Equal(t, nil, err)

	// Test case :  Delete using empty shareID and fsID

	nfsShareIDTemp = ""

	err = testConf.fileAPI.DeleteNFSShare(ctx, fsID, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share with empty nfs share ID failed")
	}

	fsIDTemp = ""
	err = testConf.fileAPI.DeleteNFSShare(ctx, fsIDTemp, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share with empty fsID failed")
	}

	err = testConf.fileAPI.DeleteNFSShareCreatedFromSnapshot(ctx, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share created by snapshot with empty nfs share ID failed")
	}

	//@TODO: Check and Add negative test cases

	fmt.Println("Delete NFS Share Test Successful")
}

func TestExpandFilesystem(t *testing.T) {
	fmt.Println("Begin - Expand Filesystem Test")
	ctx := context.Background()
	err := testConf.fileAPI.ExpandFilesystem(ctx, fsID, 7516192768)
	assert.Equal(t, errors.New("unable to find filesystem Id . Error: Filesystem Id shouldn't be empty"), err)

	// Negative cases
	fsIDTemp := "dummy_fs_sv_1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.ExpandFilesystem(ctx, fsIDTemp, 7368709120)
	assert.Equal(t, nil, err)

	err = testConf.fileAPI.ExpandFilesystem(ctx, fsID, 4368709120)
	if err == nil {
		t.Fatalf("Expand filesystem with smaller size case failed: %v", err)
	}

	fmt.Println("Expand Filesystem Test Successful")
}

func TestDeleteFilesystem(t *testing.T) {
	fmt.Println("Begin - Delete Filesystem Test")
	ctx := context.Background()
	// Clear existing expectations
	testConf.volumeAPI.client.api.(*mocks.Client).ExpectedCalls = nil

	// Test case: Delete using empty fsName/ID
	fsIDTemp := ""
	err := testConf.fileAPI.DeleteFilesystem(ctx, fsIDTemp)
	assert.Equal(t, errors.New("Filesystem Id cannot be empty"), err)

	fsIDTemp = "dummy-fs-1"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Twice()
	err = testConf.fileAPI.DeleteFilesystem(ctx, fsIDTemp)
	assert.Equal(t, nil, err)

	fsIDTemp = "fsID"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("error")).Once()
	err = testConf.fileAPI.DeleteFilesystem(ctx, "fsID")
	assert.ErrorContainsf(t, err, "Error", "delete Filesystem %s Failed.", fsIDTemp)

	fsIDTemp = "fsID"
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(errors.New(AttachedSnapshotsErrorCode)).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	testConf.fileAPI.client.api.(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(errors.New(AttachedSnapshotsErrorCode)).Once()
	err = testConf.fileAPI.DeleteFilesystem(ctx, "fsID")
	assert.ErrorContainsf(t, err, "Error", "mark filesystem %s for deletion failed.", fsIDTemp)

	fmt.Println("Delete Filesystem Test Successful")
}
