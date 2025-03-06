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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

// FilesystemAccessType is integer
type FilesystemAccessType int

// FilesystemAccessType constants
const (
	BlockAccessType      FilesystemAccessType = 0 // Parameter not applicable for block
	CheckpointAccessType FilesystemAccessType = 1 // Checkpoint access to enable access through a .ckpt folder in the file system.
	ProtocolAccessType   FilesystemAccessType = 2 // Protocol access to enable access through a file share.
)

// SnapshotNotFoundErrorCode stores snapshot not found error code
var SnapshotNotFoundErrorCode = "0x7d13005"

// ErrorSnapshotNotFound stores Snapshot not found error
var ErrorSnapshotNotFound = errors.New("Unable to find filesystem")

// CreateSnapshot creates a snapshot of a volume
//
// Parameters:
//
// - `storageResourceID` : the array to check
// - `name` : the value to search for
//
// Returns:
// - *types.Snapshot
// - an error if create snapshot fails
func (c *UnityClientImpl) CreateSnapshot(ctx context.Context, storageResourceID, snapshotName, description, retentionDuration string) (*types.Snapshot, error) {
	return c.CreateSnapshotWithFsAccesType(ctx, storageResourceID, snapshotName, description, retentionDuration, BlockAccessType)
}

// CreateSnapshotWithFsAccesType - Creates snashot with FsAccess type
func (c *UnityClientImpl) CreateSnapshotWithFsAccesType(ctx context.Context, storageResourceID, snapshotName, _, retentionDuration string, filesystemAccessType FilesystemAccessType) (*types.Snapshot, error) {
	var createSnapshot types.CreateSnapshotParam
	if len(storageResourceID) == 0 {
		return nil, errors.New("storage Resource ID cannot be empty")
	}
	var err error
	createSnapshot.Name, err = util.ValidateResourceName(snapshotName, api.MaxResourceNameLength)
	if err != nil {
		return nil, fmt.Errorf("invalid snapshot name Error:%v", err)
	}

	if retentionDuration != "" {
		seconds, err := util.ValidateDuration(retentionDuration)
		if err != nil {
			return nil, err
		}

		if seconds != 0 {
			createSnapshot.RetentionDuration = seconds
		}
	}
	storageResource := types.StorageResourceParam{
		ID: storageResourceID,
	}
	createSnapshot.StorageResource = &storageResource
	createSnapshot.FilesystemAccessType = int(filesystemAccessType)

	snapshotResp := &types.Snapshot{}
	err = c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.SnapAction), createSnapshot, snapshotResp)
	if err != nil {
		return nil, err
	}
	return snapshotResp, nil
}

// DeleteFilesystemAsSnapshot - Delete Snapshots acting as filesystem on array
func (c *UnityClientImpl) DeleteFilesystemAsSnapshot(ctx context.Context, snapshotID string, sourceFs *types.Filesystem) error {
	log := util.GetRunIDLogger(ctx)
	deleteSourceFs := false
	if strings.Contains(sourceFs.FileContent.Description, MarkFilesystemForDeletion) {
		deleteSourceFs = true
	}
	err := c.DeleteSnapshot(ctx, snapshotID)
	if err != nil {
		return err
	}
	if deleteSourceFs {
		// Try deleting the marked filesystem for deletion
		err = c.DeleteFilesystem(ctx, sourceFs.FileContent.ID)
		if err != nil {
			log.Warnf("Deletion of source filesystem: %s marked for deletion failed with error: %v", sourceFs.FileContent.ID, err)
		}
	}
	return nil
}

// DeleteSnapshot deletes a snapshot based on Snapshot ID
//
// Parameters:
//
// - `snapshotID` : User need to provide snapshot CLI Id.
//
// Returns:
// - an error if delete snapshot fails
func (c *UnityClientImpl) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	log := util.GetRunIDLogger(ctx)
	if snapshotID == "" {
		return errors.New("snapshot ID cannot be empty")
	}

	deleteErr := c.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.SnapAction, snapshotID), nil, nil)
	if deleteErr != nil {
		return fmt.Errorf("delete Snapshot Id-%s Failed: %v ", snapshotID, deleteErr)
	}
	log.Debugf("Delete Snapshot ID-%s Successful", snapshotID)
	return nil
}

// ListSnapshots lists all snapshots based on Snapshot ID or source-volume-id
// Returns a chunk of data on a single page, as specified by the maxEntries and page (startToken) parameters.
func (c *UnityClientImpl) ListSnapshots(ctx context.Context, startToken int, maxEntries int, sourceVolumeID, snapshotID string) ([]types.Snapshot, int, error) {
	snapResp := &types.ListSnapshot{}

	if snapshotID != "" {
		snapshotURI := fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.SnapAction, snapshotID, SnapshotDisplayFields)
		snapshotResp := &types.Snapshot{}
		err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, snapshotURI, nil, snapshotResp)
		if err != nil {
			return nil, 0, nil
		}
		return []types.Snapshot{*snapshotResp}, 0, nil
	}
	nextToken := startToken + 1
	snapshotURI := fmt.Sprintf(api.UnityAPIInstanceTypeResourcesWithFields, api.SnapAction, SnapshotDisplayFields)
	// Pagination will apply only for list all snapshots. If user provides snapshotID or sourceVolumeID then pagination will not apply
	if sourceVolumeID == "" {
		if maxEntries != 0 {
			snapshotURI = fmt.Sprintf(snapshotURI+"&per_page=%d", maxEntries)

			// startToken should exists only when maxEntries are present
			if startToken != 0 {
				snapshotURI = fmt.Sprintf(snapshotURI+"&page=%d", startToken)
			}
		}
	}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, snapshotURI, nil, snapResp)
	if err != nil {
		return nil, 0, err
	}

	var snapshots []types.Snapshot
	if sourceVolumeID != "" {
		for _, snapshot := range snapResp.Snapshots {
			if snapshot.SnapshotContent.StorageResource.ID == sourceVolumeID {
				snapshots = append(snapshots, snapshot)
			}
		}
		return snapshots, 0, nil
	}

	return snapResp.Snapshots, nextToken, nil
}

// FindSnapshotByName - To find snapshot using snapshot-name
func (c *UnityClientImpl) FindSnapshotByName(ctx context.Context, snapshotName string) (*types.Snapshot, error) {
	log := util.GetRunIDLogger(ctx)
	snapshotName, err := util.ValidateResourceName(snapshotName, api.MaxResourceNameLength)
	if err != nil {
		return nil, err
	}
	snapshotResp := &types.Snapshot{}
	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.SnapAction, snapshotName, SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		if strings.Contains(err.Error(), SnapshotNotFoundErrorCode) {
			return nil, ErrorSnapshotNotFound
		}
		return nil, fmt.Errorf("unable to find Snapshot Name %s Error: %v", snapshotName, err)
	}
	log.Debugf("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceID)
	return snapshotResp, nil
}

// FindSnapshotByID - To find snapshot using snapshot-id
func (c *UnityClientImpl) FindSnapshotByID(ctx context.Context, snapshotID string) (*types.Snapshot, error) {
	log := util.GetRunIDLogger(ctx)
	if snapshotID == "" {
		return nil, errors.New("snapshot ID cannot be empty")
	}
	snapshotResp := &types.Snapshot{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.SnapAction, snapshotID, SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		if strings.Contains(err.Error(), SnapshotNotFoundErrorCode) {
			return nil, ErrorSnapshotNotFound
		}
		return nil, fmt.Errorf("unable to find Snapshot id %s Error: %v", snapshotID, err)
	}
	log.Debugf("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceID)
	return snapshotResp, nil
}

// ModifySnapshotAutoDeleteParameter - Modify Snapshot (currently used to disable auto-delete parameter)
func (c *UnityClientImpl) ModifySnapshotAutoDeleteParameter(ctx context.Context, snapshotID string) error {
	log := util.GetRunIDLogger(ctx)
	if snapshotID == "" {
		return errors.New("snapshot ID cannot be empty")
	}

	modifySnapshot := types.CreateSnapshotParam{
		IsAutoDelete: false,
	}
	snapshotResp := &types.Snapshot{}

	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifySnapshotURI, api.SnapAction, snapshotID), modifySnapshot, snapshotResp)
	if err != nil {
		return fmt.Errorf("unable to modify Snapshot %s Error: %v", snapshotID, err)
	}
	log.Debugf("Changed AutoDelete to false for Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceID)
	return nil
}

// CopySnapshot - Creates a copy of the source snapshot which can be used for NFS export, and returns the ID of the copy snapshot
func (c *UnityClientImpl) CopySnapshot(ctx context.Context, sourceSnapshotID, name string) (*types.Snapshot, error) {
	if name == "" {
		return nil, errors.New("Snapshot Name cannot be empty")
	}

	if sourceSnapshotID == "" {
		return nil, errors.New("Source Snapshot ID cannot be empty")
	}

	copySnapshotReq := types.CopySnapshot{
		Name:  name,
		Child: true,
	}

	snapsResp := &types.CopySnapshots{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityCopySnapshotURI, api.SnapAction, sourceSnapshotID), copySnapshotReq, snapsResp)
	if err != nil {
		return nil, fmt.Errorf("unable to Copy Snapshot %s. Error: %v", sourceSnapshotID, err)
	}

	snapResp, err := c.FindSnapshotByID(ctx, snapsResp.CopySnapshotsContent.Copies[0].ID)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

// ModifySnapshot - Modify Snapshot's description and retention duration parameters
func (c *UnityClientImpl) ModifySnapshot(ctx context.Context, snapshotID, description, retentionDuration string) error {
	if snapshotID == "" {
		return errors.New("snapshot ID cannot be empty")
	}

	modifySnapshot := types.CreateSnapshotParam{
		Description: description,
	}
	if retentionDuration != "" {
		seconds, err := util.ValidateDuration(retentionDuration)
		if err != nil {
			return err
		}

		if seconds != 0 {
			modifySnapshot.RetentionDuration = seconds
		}
	}
	snapshotResp := &types.Snapshot{}

	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifySnapshotURI, api.SnapAction, snapshotID), modifySnapshot, snapshotResp)
	if err != nil {
		return fmt.Errorf("unable to modify Snapshot %s Error: %v", snapshotID, err)
	}
	return nil
}
