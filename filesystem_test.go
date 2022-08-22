package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var fsName string
var fsID string
var nfsShareName string
var nfsShareID string
var storageResourceID string
var snapshotID string
var nfsShareIDBySnap string
var ctx context.Context

const (
	NFSShareLocalPath  = "/"
	NFSShareNamePrefix = "csishare-"
)

func TestFilesystem(t *testing.T) {

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	fsName = "Unit-test-fs-" + timeStamp
	snapName = "Unit-test-snapshot-" + timeStamp
	ctx = context.Background()

	findNasServerTest(t)
	createFilesystemTest(t)
	findFilesystemTest(t)
	createNfsShareTest(t)
	findNfsShareTest(t)
	modifyNfsShareTest(t)
	updateDescriptionTest(t)
	deleteNfsShareTest(t)
	expandFilesystemTest(t)
	deleteFilesystemTest(t)
}

func findNasServerTest(t *testing.T) {

	fmt.Println("Begin - Find Nas Server Test")

	_, err := testConf.fileAPI.FindNASServerByID(ctx, testConf.nasServer)

	if err != nil {
		t.Fatalf("Find filesystem by name failed: %v", err)
	}

	//Test case :  GET using invalid ID
	nasServer := "nas_dummy_1"

	_, err = testConf.fileAPI.FindNASServerByID(ctx, nasServer)
	if err == nil {
		t.Fatal("Find Nas Server - Negative case failed")
	}

	//Test case :  GET using empty ID
	nasServer = ""

	_, err = testConf.fileAPI.FindNASServerByID(ctx, nasServer)
	if err == nil {
		t.Fatal("Find NAS server using empty ID - Negative case failed")
	}

	fmt.Println("Find Nas Server Test Successful")
}

func createFilesystemTest(t *testing.T) {

	fmt.Println("Begin - Create Filesystem Test")

	_, err := testConf.fileAPI.CreateFilesystem(ctx, fsName, testConf.poolID, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	if err != nil {
		t.Fatalf("Create filesystem failed: %v", err)
	}

	//Negative cases

	fsNameTemp := ""
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsNameTemp, testConf.poolID, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	if err == nil {
		t.Fatal("Create filesystem with empty name - Negative case failed")
	}

	fsNameTemp = "dummy-fs-1234567890123456789012345678901234567890123456789012345678"
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsNameTemp, testConf.poolID, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	if err == nil {
		t.Fatal("Create filesystem with fs name more than 63 characters - Negative case failed")
	}

	poolIDTemp := "dummy_pool_1"
	_, err = testConf.fileAPI.CreateFilesystem(ctx, fsName, poolIDTemp, "Unit test resource", testConf.nasServer, 5368709120, 0, 8192, 0, true, false)
	if err == nil {
		t.Fatal("Create filesystem with invalid storage pool - Negative case failed")
	}

	fmt.Println("Create Filesystem test successful")

}

func findFilesystemTest(t *testing.T) {

	fmt.Println("Begin - Find Filesystem Test")

	filesystem, err := testConf.fileAPI.FindFilesystemByName(ctx, fsName)
	if err != nil {
		t.Fatalf("Find filesystem by name failed: %v", err)
	}

	filesystem, err = testConf.fileAPI.FindFilesystemByID(ctx, filesystem.FileContent.ID)
	if err != nil {
		t.Fatalf("Find filesystem by Id failed: %v", err)
	}

	fsID, err = testConf.fileAPI.GetFilesystemIDFromResID(ctx, filesystem.FileContent.StorageResource.ID)
	if err != nil {
		t.Fatalf("Find filesystem by Resource Id failed: %v", err)
	}

	nfsShareName = NFSShareNamePrefix + filesystem.FileContent.Name
	storageResourceID = filesystem.FileContent.StorageResource.ID

	fmt.Println("Filesystem ID: " + fsID)

	//Test case :  GET using invalid fsName/ID
	fsNameTemp := "dummy-fs-1"

	filesystem, err = testConf.fileAPI.FindFilesystemByName(ctx, fsNameTemp)
	if err == nil {
		t.Fatal("Find filesystem by name - Negative case failed")
	}

	filesystem, err = testConf.fileAPI.FindFilesystemByID(ctx, fsNameTemp)
	if err == nil {
		t.Fatal("Find filesystem by Id - Negative case failed")
	}

	_, err = testConf.fileAPI.GetFilesystemIDFromResID(ctx, fsNameTemp)
	if err == nil {
		t.Fatalf("Find filesystem by Resource Id - Negative case failed")
	}

	//Test case :  GET using empty fsName/ID
	fsNameTemp = ""

	filesystem, err = testConf.fileAPI.FindFilesystemByName(ctx, fsNameTemp)
	if err == nil {
		t.Fatal("Find filesystem by name using empty fsName - Negative case failed")
	}

	filesystem, err = testConf.fileAPI.FindFilesystemByID(ctx, fsNameTemp)
	if err == nil {
		t.Fatal("Find filesystem by Id using empty fsID - Negative case failed")
	}

	_, err = testConf.fileAPI.GetFilesystemIDFromResID(ctx, fsNameTemp)
	if err == nil {
		t.Fatalf("Find filesystem by Resource Id failed: %v", err)
	}

	fmt.Println("Find Filesystem test successul")
}

func createNfsShareTest(t *testing.T) {

	fmt.Println("Begin - Create NFS Share Test")

	_, err := testConf.fileAPI.CreateNFSShare(ctx, nfsShareName, NFSShareLocalPath, fsID, NoneDefaultAccess)
	if err != nil {
		t.Fatalf("Create NFS Share failed: %v", err)
	}

	//Test case : Create NFS share using snapshot
	snapshot, err := testConf.snapAPI.CreateSnapshot(ctx, storageResourceID, snapName, "Snapshot Description", "")
	if err != nil {
		t.Fatalf("Create snapshot of filesystem failed: %v", err)
	}

	snapshotID = snapshot.SnapshotContent.ResourceID

	nfsShareBySnap, err := testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, snapshotID, NoneDefaultAccess)
	if err != nil {
		t.Fatalf("Create NFS Share from snapshot failed: %v", err)
	}

	nfsShareIDBySnap = nfsShareBySnap.NFSShareContent.ID

	//Test case :  Create using invalid fsID
	fsIDTemp := "dummy-fs-1"
	_, err = testConf.fileAPI.CreateNFSShare(ctx, nfsShareName, NFSShareLocalPath, fsIDTemp, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share with invalid fsID - Negative case failed")
	}

	fsIDTemp = ""
	_, err = testConf.fileAPI.CreateNFSShare(ctx, nfsShareName, NFSShareLocalPath, fsIDTemp, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share with empty fsID - Negative case failed")
	}

	nfsShareNameTemp := ""
	_, err = testConf.fileAPI.CreateNFSShare(ctx, nfsShareNameTemp, NFSShareLocalPath, fsID, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share with empty share name - Negative case failed")
	}

	snapshotIDTemp := ""
	_, err = testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, snapshotIDTemp, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share from snapshot with empty snapshot Id case failed: %v", err)
	}

	snapshotIDTemp = "dummy_snap_1"
	_, err = testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, snapshotIDTemp, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share from snapshot with invalid snapshot Id case failed: %v", err)
	}

	_, err = testConf.fileAPI.CreateNFSShareFromSnapshot(ctx, nfsShareName+"_by_snap", NFSShareLocalPath, snapshotID, NoneDefaultAccess)
	if err == nil {
		t.Fatalf("Create NFS Share from snapshot with an existing nfs share name case failed: %v", err)
	}

	fmt.Println("Create NFS Share Test Successful")

}

func findNfsShareTest(t *testing.T) {

	fmt.Println("Begin - Find NFS Share Test")

	nfsShare, err := testConf.fileAPI.FindNFSShareByName(ctx, nfsShareName)
	if err != nil {
		t.Fatalf("Find NFS Share by name failed: %v", err)
	}

	nfsShareID = nfsShare.NFSShareContent.ID

	_, err = testConf.fileAPI.FindNFSShareByID(ctx, nfsShareID)
	if err != nil {
		t.Fatalf("Find NFS Share by ID failed: %v", err)
	}

	//Test case :  GET using invalid shareName/ID
	nfsShareNameTemp := "dummy-fs-1"

	_, err = testConf.fileAPI.FindNFSShareByName(ctx, nfsShareNameTemp)
	if err == nil {
		t.Fatal("Find NFS Share by name - Negative case failed")
	}

	_, err = testConf.fileAPI.FindNFSShareByID(ctx, nfsShareNameTemp)
	if err == nil {
		t.Fatal("Find NFS Share by Id - Negative case failed")
	}

	//Test case :  GET using empty fsName/ID
	nfsShareNameTemp = ""

	_, err = testConf.fileAPI.FindNFSShareByName(ctx, nfsShareNameTemp)
	if err == nil {
		t.Fatal("Find NFS Share by name using empty share Name - Negative case failed")
	}

	_, err = testConf.fileAPI.FindNFSShareByID(ctx, nfsShareNameTemp)
	if err == nil {
		t.Fatal("Find filesystem by Id using empty share ID - Negative case failed")
	}

	fmt.Println("Find NFS Share Test Successful")

}

func modifyNfsShareTest(t *testing.T) {

	fmt.Println("Begin - Modify NFS Share Test")

	host, err := testConf.hostAPI.FindHostByName(ctx, testConf.nodeHostName)
	if err != nil {
		t.Fatalf("Find host failed: %v", err)
	}

	var hostIDList []string
	hostIDList = append(hostIDList, host.HostContent.ID)

	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsID, nfsShareID, hostIDList, ReadOnlyAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share by name failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsID, nfsShareID, hostIDList, ReadWriteAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share by name failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsID, nfsShareID, hostIDList, ReadOnlyRootAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share by name failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsID, nfsShareID, hostIDList, ReadWriteRootAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share by name failed: %v", err)
	}

	//Test cases : Modify NFS share created from snapshot

	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnap, hostIDList, ReadWriteRootAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnap, hostIDList, ReadOnlyRootAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnap, hostIDList, ReadWriteAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnap, hostIDList, ReadOnlyAccessType)
	if err != nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	fsIDTemp := "dummy-fs-1"
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadWriteRootAccessType)
	if err == nil {
		t.Fatalf("Modify NFS Share with invalid fs ID - Negative case Failed")
	}

	fsIDTemp = ""
	err = testConf.fileAPI.ModifyNFSShareHostAccess(ctx, fsIDTemp, nfsShareID, hostIDList, ReadWriteRootAccessType)
	if err == nil {
		t.Fatalf("Modify NFS Share with empty fs ID - Negative case Failed")
	}

	nfsShareIDBySnapTemp := ""
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnapTemp, hostIDList, ReadOnlyAccessType)
	if err == nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	nfsShareIDBySnapTemp = "dummy-nsf-share-1"
	err = testConf.fileAPI.ModifyNFSShareCreatedFromSnapshotHostAccess(ctx, nfsShareIDBySnapTemp, hostIDList, ReadOnlyAccessType)
	if err == nil {
		t.Fatalf("Modify NFS Share created by snapshot failed: %v", err)
	}

	fmt.Println("Modify NFS Share Test Successful")

}

func updateDescriptionTest(t *testing.T) {
	fmt.Println("Begin - Update Description of Filesystem Test")

	//Positive scenario is covered under DeleteFilesystemTest()
	//Negative test case

	filesystemIDTemp := ""
	err := testConf.fileAPI.updateDescription(ctx, filesystemIDTemp, "Description of filesystem")
	if err == nil {
		t.Fatalf("Update filesystem description failed: %v", err)
	}

	filesystemIDTemp = "dummy_fs_1"
	err = testConf.fileAPI.updateDescription(ctx, filesystemIDTemp, "Description of filesystem")
	if err == nil {
		t.Fatalf("Update filesystem description failed: %v", err)
	}

}

func deleteNfsShareTest(t *testing.T) {

	fmt.Println("Begin - Delete NFS Share Test")

	err := testConf.fileAPI.DeleteNFSShare(ctx, fsID, nfsShareID)
	if err != nil {
		t.Fatalf("Delete NFS Share failed: %v", err)
	}

	err = testConf.fileAPI.DeleteNFSShareCreatedFromSnapshot(ctx, nfsShareIDBySnap)
	if err != nil {
		t.Fatalf("Delete NFS Share created by Snapshot failed: %v", err)
	}

	//Test case :  Delete using invalid shareID and fsID
	nfsShareIDTemp := "dummy-fs-1"
	fsIDTemp := "dummy-fs-1"

	err = testConf.fileAPI.DeleteNFSShare(ctx, fsID, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share with invalid nfs share ID failed")
	}

	err = testConf.fileAPI.DeleteNFSShare(ctx, fsIDTemp, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share with invalid fs ID failed")
	}

	err = testConf.fileAPI.DeleteNFSShareCreatedFromSnapshot(ctx, nfsShareIDTemp)
	if err == nil {
		t.Fatalf("Delete NFS Share created by snapshot with invalid nfs share ID failed")
	}

	//Test case :  Delete using empty shareID and fsID

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

func expandFilesystemTest(t *testing.T) {

	fmt.Println("Begin - Expand Filesystem Test")

	err := testConf.fileAPI.ExpandFilesystem(ctx, fsID, 7516192768)
	if err != nil {
		t.Fatalf("Expand filesystem failed: %v", err)
	}

	err = testConf.fileAPI.ExpandFilesystem(ctx, fsID, 7516192768)
	if err != nil {
		t.Fatalf("Expand filesystem with same size failed: %v", err)
	}

	//Negative cases
	fsIDTemp := "dummy_fs_sv_1"
	err = testConf.fileAPI.ExpandFilesystem(ctx, fsIDTemp, 7368709120)
	if err == nil {
		t.Fatalf("Expand filesystem with invalid Id case failed: %v", err)
	}

	err = testConf.fileAPI.ExpandFilesystem(ctx, fsID, 4368709120)
	if err == nil {
		t.Fatalf("Expand filesystem with smaller size case failed: %v", err)
	}

	fmt.Println("Expand Filesystem Test Successful")
}

func deleteFilesystemTest(t *testing.T) {

	fmt.Println("Begin - Delete Filesystem Test")

	//Test case : Delete fail if snapshot exists

	err := testConf.fileAPI.DeleteFilesystem(ctx, fsID)
	if err != nil {
		t.Fatalf("Delete filesystem failed: %v", err)
	}

	filesystem, err := testConf.fileAPI.FindFilesystemByID(ctx, fsID)
	if err != nil {
		t.Fatalf("Find filesystem by resource Id failed: %v", err)
	}
	fmt.Println("filesystem values ", filesystem)

	err = testConf.snapAPI.DeleteFilesystemAsSnapshot(ctx, snapshotID, filesystem)
	if err != nil {
		t.Fatalf("Delete snapshot failed: %v", err)
	}

	//@TODO: Add negative cases after export - before unexport

	//Test case :  Delete using invalid fsName/ID
	fsIDTemp := "dummy-fs-1"
	err = testConf.fileAPI.DeleteFilesystem(ctx, fsIDTemp)
	if err == nil {
		t.Fatal("Delete filesystem - invaid fsID failed")
	}

	//Test case: Delete using empty fsName/ID
	fsIDTemp = ""
	err = testConf.fileAPI.DeleteFilesystem(ctx, fsIDTemp)
	if err == nil {
		t.Fatal("Delete filesystem - empty fsID failed")
	}

	fmt.Println("Delete Filesystem Test Successful")
}
