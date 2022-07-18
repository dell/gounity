package gounity

import (
	"context"
	"errors"
	"fmt"
	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
	"net/http"
	"net/url"
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
	log.Debugf("Remote system name: %s Id: %s", remoteSystemName, remoteSystemNameResp.RemoteSystemContent.Id)
	return remoteSystemNameResp, nil
}

func (r *Replication) FindRemoteSystemById(ctx context.Context, remoteSystemId string) (*types.RemoteSystem, error) {
	log := util.GetRunIDLogger(ctx)
	remoteSystemId, err := util.ValidateResourceName(remoteSystemId, api.MaxResourceNameLength)
	if err != nil {
		return nil, err
	}
	remoteSystem := &types.RemoteSystem{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.RemoteSystemAction, remoteSystemId, RemoteSystemFields), nil, remoteSystem)
	if err != nil {
		return nil, fmt.Errorf("unable to find Remote system %v. Error: %v", remoteSystemId, err)
	}
	log.Debugf("Remote system Id: %s", remoteSystemId)
	return remoteSystem, nil
}

func (r *Replication) CreateReplicationSession(ctx context.Context, replicationSessionName, srcResourceId, dstResourceId, remoteSystemName string, maxTimeOutOfSync string) (*types.ReplicationSession, error) {
	var createRS types.CreateReplicationSessionParam
	if len(srcResourceId) == 0 {
		return nil, errors.New("storage Resource ID cannot be empty")
	}
	var err error
	createRS.Name, err = util.ValidateResourceName(replicationSessionName, api.MaxResourceNameLength)
	if err != nil {
		return nil, fmt.Errorf("invalid replication session name. Error:%v", err)
	}
	remoteSystem, err := r.FindRemoteSystemByName(ctx, remoteSystemName)
	if err != nil {
		return nil, fmt.Errorf("can't find remote system %v. Error:%v", remoteSystem, err)
	}

	createRS.RemoteSystem = &remoteSystem.RemoteSystemContent
	createRS.MaxTimeOutOfSync = maxTimeOutOfSync
	createRS.SrcResourceId = srcResourceId
	createRS.DstResourceId = dstResourceId
	createRS.OverwriteDestination = true
	rsResp := &types.ReplicationSession{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.ReplicationSessionAction), createRS, rsResp)
	if err != nil {
		return nil, err
	}
	return rsResp, nil
}

//DeleteConsistencyGroup - Delete ReplicationSession by ID.
func (r *Replication) DeleteReplicationSession(ctx context.Context, rID string) error {
	if len(rID) == 0 {
		return errors.New("ReplicationSession Id cannot be empty")
	}

	_, err := r.FindReplicationSessionById(ctx, rID)

	if err != nil {
		return err
	}

	deleteErr := r.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.ReplicationSessionAction, rID), nil, nil)

	if deleteErr != nil {
		return fmt.Errorf("delete ReplicationSession %s Failed. Error: %v", rID, deleteErr)
	}

	return nil
}

func (r *Replication) FindReplicationSessionBySrcResourceID(ctx context.Context, srcResourceId string) (*types.ReplicationSession, error) {
	if len(srcResourceId) == 0 {
		return nil, fmt.Errorf("SrcResourceId shouldn't be empty")
	}
	filter := fmt.Sprintf("srcResourceId eq %s", "\""+srcResourceId+"\"")
	queryURI := fmt.Sprintf(api.UnityInstancesFilterWithFields, api.ReplicationSessionAction, ReplicationSessionFields, url.QueryEscape(filter))
	listReplSession := &types.ListReplicationSession{}
	err := r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, listReplSession)
	if err != nil {
		return nil, err
	}
	if len(listReplSession.ReplicationSessions) == 0 {
		return nil, nil
	}
	rs := &listReplSession.ReplicationSessions[0]
	return rs, nil
}

func (r *Replication) FindReplicationSessionByName(ctx context.Context, rsName string) (*types.ReplicationSession, error) {
	log := util.GetRunIDLogger(ctx)
	if rsName == "" {
		return nil, errors.New("Replication session name cannot be empty")
	}
	rsResp := &types.ReplicationSession{}
	err := r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.ReplicationSessionAction, rsName, ReplicationSessionFields), nil, rsResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find Replication Session name %s Error: %v", rsName, err)
	}
	log.Debugf("Replcation Session name: %s Id: %s", rsResp.ReplicationSessionContent.Name, rsResp.ReplicationSessionContent.ReplicationSessionId)
	return rsResp, nil
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

func (r *Replication) DeleteReplicationSessionById(ctx context.Context, sessionId string) error {
	log := util.GetRunIDLogger(ctx)
	if len(sessionId) == 0 {
		return errors.New("Replication session Id cannot be empty")
	}

	_, err := r.FindReplicationSessionById(ctx, sessionId)
	if err != nil {
		return err
	}

	deleteErr := r.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.ReplicationSessionAction, sessionId), nil, nil)
	if deleteErr != nil {
		return fmt.Errorf("delete replication session %s Failed. Error: %v", sessionId, deleteErr)
	}
	log.Debugf("Delete Replication session %s Successful", sessionId)

	return nil
}
