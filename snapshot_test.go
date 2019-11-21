package gounity

import (
	"fmt"
	"github.com/dell/gounity/payloads"
	"testing"
	"time"
)

func TestCreateSanapshot(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")

	var vol *payloads.Volume
	var err error
	vol, err = testConf.volumeApi.CreateLun(volName, testConf.poolId, "Description", 5368709120, 0, "", true, false)
	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	snap, err := testConf.snapApi.CreateSnapshot(vol.VolumeContent.ResourceId, "testsnap-1", "description", "1:23:52:50", "true")
	fmt.Println("Create Snapshot:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Create Snapshot failed: %v", err)
	}

	snap, err = testConf.snapApi.FindSnapshotById(snap.SnapshotContent.ResourceId)
	fmt.Println("Find snapshot:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}

	snap, err = testConf.snapApi.FindSnapshotByName(snap.SnapshotContent.Name)
	fmt.Println("Find snapshot:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}

	snaps, _, err := testConf.snapApi.ListSnapshots(0, 10, vol.VolumeContent.ResourceId, "")
	fmt.Println("List snapshots:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots success:", len(snaps))
	} else {
		t.Fatalf("List snapshot failed: %v", err)
	}

	err = testConf.snapApi.DeleteSnapshot(snap.SnapshotContent.ResourceId)
	fmt.Println("Delete Snapshot:", err)
	if err != nil {
		t.Fatalf("Delete Snapshot failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}
}
