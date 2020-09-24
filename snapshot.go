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

type FilesystemAccessType int

const (
	BlockAccessType      FilesystemAccessType = 0 //Parameter not applicable for block
	CheckpointAccessType FilesystemAccessType = 1 //Checkpoint access to enable access through a .ckpt folder in the file system.
	ProtocolAccessType   FilesystemAccessType = 2 //Protocol access to enable access through a file share.
)

var SnapshotNotFoundErrorCode = "0x7d13005"
var SnapshotNotFoundError = errors.New("Unable to find filesystem")

type snapshot struct {
	client *Client
}

func NewSnapshot(client *Client) *snapshot {
	return &snapshot{client}
}

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

func (s *snapshot) CreateSnapshot(ctx context.Context, storageResourceID, snapshotName, description, retentionDuration string) (*types.Snapshot, error) {
	return s.CreateSnapshotWithFsAccesType(ctx, storageResourceID, snapshotName, description, retentionDuration, BlockAccessType)
}

func (s *snapshot) CreateSnapshotWithFsAccesType(ctx context.Context, storageResourceID, snapshotName, description, retentionDuration string, filesystemAccessType FilesystemAccessType) (*types.Snapshot, error) {
	var createSnapshot types.CreateSnapshotParam
	if len(storageResourceID) == 0 {
		return nil, errors.New("storage Resource ID cannot be empty")
	}
	var err error
	createSnapshot.Name, err = util.ValidateResourceName(snapshotName, api.MaxResourceNameLength)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid snapshot name Error:%v", err))
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
	err = s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, api.SnapAction), createSnapshot, snapshotResp)
	if err != nil {
		return nil, err
	}
	return snapshotResp, nil
}

// DeleteSnapshot deletes a snapshot based on Snapshot ID
//
// Parameters:
//
// - `snapshotId` : User need to provide snapshot CLI Id.
//
// Returns:
// - an error if delete snapshot fails
func (s *snapshot) DeleteSnapshot(ctx context.Context, snapshotId string) error {
	log := util.GetRunIdLogger(ctx)
	if snapshotId == "" {
		return errors.New("snapshot ID cannot be empty")
	}

	deleteErr := s.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceUri, api.SnapAction, snapshotId), nil, nil)
	if deleteErr != nil {
		return errors.New(fmt.Sprintf("Delete Snapshot Id-%s Failed: %v ", snapshotId, deleteErr))
	}
	log.Debugf("Delete Snapshot ID-%s Successful", snapshotId)
	return nil
}

// ListSnapshot lists all snapshots based on Snapshot ID or source-volume-id
// Returns a chunk of data on a single page, as specified by the maxEntries and page (startToken) parameters.
func (s *snapshot) ListSnapshots(ctx context.Context, startToken int, maxEntries int, sourceVolumeId, snapshotId string) ([]types.Snapshot, int, error) {
	snapResp := &types.ListSnapshot{}

	if snapshotId != "" {
		snapshotUri := fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, api.SnapAction, snapshotId, SnapshotDisplayFields)
		snapshotResp := &types.Snapshot{}
		err := s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, snapshotUri, nil, snapshotResp)
		if err != nil {
			return nil, 0, nil
		}
		return []types.Snapshot{*snapshotResp}, 0, nil
	} else {
		nextToken := startToken + 1
		snapshotUri := fmt.Sprintf(api.UnityApiInstanceTypeResourcesWithFields, api.SnapAction, SnapshotDisplayFields)
		//Pagination will apply only for list all snapshots. If user provides snapshotId or sourceVolumeId then pagination will not apply
		if sourceVolumeId == "" {
			if maxEntries != 0 {
				snapshotUri = fmt.Sprintf(snapshotUri+"&per_page=%d", maxEntries)

				//startToken should exists only when maxEntries are present
				if startToken != 0 {
					snapshotUri = fmt.Sprintf(snapshotUri+"&page=%d", startToken)
				}
			}
		}
		err := s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, snapshotUri, nil, snapResp)
		if err != nil {
			return nil, 0, err
		}

		var snapshots []types.Snapshot
		if sourceVolumeId != "" {
			for _, snapshot := range snapResp.Snapshots {
				if snapshot.SnapshotContent.StorageResource.Id == sourceVolumeId {
					snapshots = append(snapshots, snapshot)
				}
			}
			return snapshots, 0, nil
		}

		return snapResp.Snapshots, nextToken, nil
	}
}

// To find snapshot using snapshot-name
func (s *snapshot) FindSnapshotByName(ctx context.Context, snapshotName string) (*types.Snapshot, error) {
	log := util.GetRunIdLogger(ctx)
	snapshotName, err := util.ValidateResourceName(snapshotName, api.MaxResourceNameLength)
	if err != nil {
		return nil, err
	}
	snapshotResp := &types.Snapshot{}
	err = s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, api.SnapAction, snapshotName, SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		if strings.Contains(err.Error(), SnapshotNotFoundErrorCode) {
			return nil, SnapshotNotFoundError
		}
		return nil, errors.New(fmt.Sprintf("Unable to find Snapshot Name %s Error: %v", snapshotName, err))
	}
	log.Debugf("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
	return snapshotResp, nil
}

// To find snapshot using snapshot-id
func (s *snapshot) FindSnapshotById(ctx context.Context, snapshotId string) (*types.Snapshot, error) {
	log := util.GetRunIdLogger(ctx)
	if snapshotId == "" {
		return nil, errors.New("snapshot ID cannot be empty")
	}
	snapshotResp := &types.Snapshot{}
	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, api.SnapAction, snapshotId, SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		if strings.Contains(err.Error(), SnapshotNotFoundErrorCode) {
			return nil, SnapshotNotFoundError
		}
		return nil, errors.New(fmt.Sprintf("Unable to find Snapshot id %s Error: %v", snapshotId, err))
	}
	log.Debugf("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
	return snapshotResp, nil
}

// Modify Snapshot (currently used to disable auto-delete parameter)
func (s *snapshot) ModifySnapshotAutoDeleteParameter(ctx context.Context, snapshotId string) error {
	log := util.GetRunIdLogger(ctx)
	if snapshotId == "" {
		return errors.New("snapshot ID cannot be empty")
	}

	modifySnapshot := types.CreateSnapshotParam{
		IsAutoDelete: false,
	}
	snapshotResp := &types.Snapshot{}

	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifySnapshotUri, api.SnapAction, snapshotId), modifySnapshot, snapshotResp)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to modify Snapshot %s Error: %v", snapshotId, err))
	}
	log.Debugf("Changed AutoDelete to false for Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
	return nil
}

//CopySnapshot - Creates a copy of the source snapshot which can be used for NFS export, and returns the ID of the copy snapshot
func (s *snapshot) CopySnapshot(ctx context.Context, sourceSnapshotId, name string) (*types.Snapshot, error) {
	if name == "" {
		return nil, errors.New("Snapshot Name cannot be empty")
	}

	if sourceSnapshotId == "" {
		return nil, errors.New("Source Snapshot ID cannot be empty")
	}

	copySnapshotReq := types.CopySnapshot{
		Name:  name,
		Child: true,
	}

	snapsResp := &types.CopySnapshots{}
	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityCopySnapshotUri, api.SnapAction, sourceSnapshotId), copySnapshotReq, snapsResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to Copy Snapshot %s. Error: %v", sourceSnapshotId, err))
	}

	snapResp, err := s.FindSnapshotById(ctx, snapsResp.CopySnapshotsContent.Copies[0].Id)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

//Modify Snapshot's description and retention duration parameters
func (s *snapshot) ModifySnapshot(ctx context.Context, snapshotId, description, retentionDuration string) error {
	if snapshotId == "" {
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

	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifySnapshotUri, api.SnapAction, snapshotId), modifySnapshot, snapshotResp)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to modify Snapshot %s Error: %v", snapshotId, err))
	}
	return nil
}
