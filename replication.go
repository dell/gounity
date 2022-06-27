package gounity

import (
	"context"
	"errors"
	"fmt"
	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
	"net/http"
)

type Replication struct {
	client *Client
}

func NewReplicationSession(client *Client) *Replication {
	return &Replication{client}
}

func (r *Replication) FindRemoteSystemByName(ctx context.Context, remoteSystemName string) (*types.RemoteSystem, error) {
	log := util.GetRunIDLogger(ctx)
	remoteSystemName, err := util.ValidateResourceName(remoteSystemName, api.MaxResourceNameLength)
	if err != nil {
		return nil, err
	}
	remoteSystemNameResp := &types.RemoteSystem{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.RemoteSystemAction, remoteSystemName, RemoteSystemFields), nil, remoteSystemNameResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find Remote system %v. Error: %v", remoteSystemName, err)
	}
	log.Debugf("Remote system name: %s Id: %s", remoteSystemName, remoteSystemNameResp.RemoteSystemContent.RemoteSystemId)
	return remoteSystemNameResp, nil
}

func (r *Replication) CreateReplicationSession(ctx context.Context, replicationSessionName, srcResourceId, remoteStoragePool, FsName, remoteSystemName string, maxTimeOutOfSync int32) (*types.ReplicationSession, error) {
	var createRS types.CreateReplicationSessionParam
	if len(srcResourceId) == 0 {
		return nil, errors.New("storage Resource ID cannot be empty")
	}
	var err error
	createRS.Name, err = util.ValidateResourceName(replicationSessionName, api.MaxResourceNameLength)
	if err != nil {
		return nil, fmt.Errorf("invalid replication session name. Error:%v", err)
	}
	remoteSystemId, err := r.FindRemoteSystemByName(ctx, remoteSystemName)
	if err != nil {
		return nil, fmt.Errorf("can't find remote system %v. Error:%v", remoteSystemId, err)
	}
	remoteSystem := types.RemoteSystemContent{
		RemoteSystemId: remoteSystemId.RemoteSystemContent.RemoteSystemId,
	}
	createRS.RemoteSystemId = &remoteSystem
	createRS.DstResourceConfig.StoragePool.PoolID = remoteStoragePool
	createRS.DstResourceConfig.Name = FsName
	createRS.MaxTimeOutOfSync = maxTimeOutOfSync
	rsResp := &types.ReplicationSession{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.ReplicationSessionAction), createRS, rsResp)
	if err != nil {
		return nil, err
	}
	return rsResp, nil
}
