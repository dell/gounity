package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

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

	snapshotResp := &types.Snapshot{}
	err = s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, "snap"), createSnapshot, snapshotResp)
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

	deleteErr := s.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceUri, "snap", snapshotId), nil, nil)
	if deleteErr != nil {
		log.Infof("Delete Snapshot Id-%s Failed: %v ", snapshotId, deleteErr)
		return deleteErr
	}
	log.Infof("Delete Snapshot ID-%s Successful", snapshotId)
	return nil
}

// ListSnapshot lists all snapshots based on Snapshot ID or source-volume-id
// Returns a chunk of data on a single page, as specified by the maxEntries and page (startToken) parameters.
func (s *snapshot) ListSnapshots(ctx context.Context, startToken int, maxEntries int, sourceVolumeId, snapshotId string) ([]types.Snapshot, int, error) {
	snapResp := &types.ListSnapshot{}

	if snapshotId != "" {
		snapshotUri := fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "snap", snapshotId, api.SnapshotDisplayFields)

		snapshotResp := &types.Snapshot{}
		err := s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, snapshotUri, nil, snapshotResp)
		if err != nil {
			return nil, 0, nil
		}
		return []types.Snapshot{*snapshotResp}, 0, nil
	} else {
		nextToken := startToken + 1
		snapshotUri := fmt.Sprintf(api.UnityApiInstanceTypeResourcesWithFields, "snap", api.SnapshotDisplayFields)
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
	err = s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "snap", snapshotName, api.SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find Snapshot Name %s Error: %v", snapshotName, err))
	}
	log.Infof("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
	return snapshotResp, nil
}

// To find snapshot using snapshot-id
func (s *snapshot) FindSnapshotById(ctx context.Context, snapshotId string) (*types.Snapshot, error) {
	log := util.GetRunIdLogger(ctx)
	if snapshotId == "" {
		return nil, errors.New("snapshot ID cannot be empty")
	}
	snapshotResp := &types.Snapshot{}
	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "snap", snapshotId, api.SnapshotDisplayFields), nil, snapshotResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find Snapshot id %s Error: %v", snapshotId, err))
	}
	log.Infof("Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
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

	err := s.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifySnapshotUri, "snap", snapshotId), modifySnapshot, snapshotResp)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to modify Snapshot %s Error: %v", snapshotId, err))
	}
	log.Infof("Changed AutoDelete to false for Snapshot name: %s Id: %s", snapshotResp.SnapshotContent.Name, snapshotResp.SnapshotContent.ResourceId)
	return nil
}
