package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
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

func (r *Replication) CreateReplicationSession(ctx context.Context, replicationSessionName, srcResourceId, dstResourceId, remoteSystemName string, maxTimeOutOfSync int32) (*types.ReplicationSession, error) {
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
	createRS.MaxTimeOutOfSync = maxTimeOutOfSync
	createRS.SrcResourceId = srcResourceId
	createRS.DstResourceId = dstResourceId
	rsResp := &types.ReplicationSession{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.ReplicationSessionAction), createRS, rsResp)
	if err != nil {
		return nil, err
	}
	return rsResp, nil
}

func (r *Replication) FindReplicationSessionIdBySrcResourceID(ctx context.Context, srcResourceId string) (*types.ReplicationSession, error) {
	filter := fmt.Sprintf("srcResourceId eq %s", srcResourceId)
	queryURI := fmt.Sprintf(api.UnityInstancesFilter, api.ReplicationSessionAction, url.QueryEscape(filter))
	rsIdResult := &types.ReplicationSession{}
	err := r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, rsIdResult)
	if err != nil {
		return nil, err
	}
	return rsIdResult, nil
}

func (r *Replication) FindReplicationSessionById(ctx context.Context, rsId string) (*types.ReplicationSession, error) {
	log := util.GetRunIDLogger(ctx)
	if rsId == "" {
		return nil, errors.New("Replication session ID cannot be empty")
	}
	rsResp := &types.ReplicationSession{}
	err := r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.ReplicationSessionAction, rsId, ReplicationSessionFields), nil, rsResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find Replication Session id %s Error: %v", rsId, err)
	}
	log.Debugf("Replcation Session name: %s Id: %s", rsResp.ReplicationSessionContent.Name, rsResp.ReplicationSessionContent.ReplicationSessionId)
	return rsResp, nil
}
