package gounity

import (
	"context"
	"fmt"
	"github.com/dell/gounity/types"
	"testing"
	"time"
)

func TestCreateSanapshot(t *testing.T) {
	now := time.Now()
	volName := "test-" + now.Format("20060102150405")
	ctx := context.Background()

	var vol *types.Volume
	var err error
	vol, err = testConf.volumeApi.CreateLun(ctx, volName, testConf.poolId, "Description", 5368709120, 0, "", true, false)
	fmt.Println("Create volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Create volume failed: %v", err)
	}

	vol, err = testConf.volumeApi.FindVolumeByName(ctx, volName)
	fmt.Println("Find volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find volume failed: %v", err)
	}

	snap, err := testConf.snapApi.CreateSnapshot(ctx, vol.VolumeContent.ResourceId, "testsnap-1", "description", "1:23:52:50")
	fmt.Println("Create Snapshot:", prettyPrintJson(snap), err)
	if err != nil {
		t.Fatalf("Create Snapshot failed: %v", err)
	}

	snap, err = testConf.snapApi.FindSnapshotById(ctx, snap.SnapshotContent.ResourceId)
	fmt.Println("Find snapshot:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}

	snap, err = testConf.snapApi.FindSnapshotByName(ctx, snap.SnapshotContent.Name)
	fmt.Println("Find snapshot:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Find snapshot failed: %v", err)
	}

	snaps, _, err := testConf.snapApi.ListSnapshots(ctx, 0, 10, vol.VolumeContent.ResourceId, "")
	fmt.Println("List snapshots:", len(snaps))
	if len(snaps) > 0 {
		fmt.Println("List snapshots success:", len(snaps))
	} else {
		t.Fatalf("List snapshot failed: %v", err)
	}

	err = testConf.snapApi.DeleteSnapshot(ctx, snap.SnapshotContent.ResourceId)
	fmt.Println("Delete Snapshot:", err)
	if err != nil {
		t.Fatalf("Delete Snapshot failed: %v", err)
	}

	err = testConf.volumeApi.DeleteVolume(ctx, vol.VolumeContent.ResourceId)
	fmt.Println("Delete volume:", prettyPrintJson(vol), err)
	if err != nil {
		t.Fatalf("Delete volume failed: %v", err)
	}
}
