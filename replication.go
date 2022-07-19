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

//Replication structure
type Replication struct {
	client *Client
}

//NewReplicationSession structure
func NewReplicationSession(client *Client) *Replication {
	return &Replication{client}
}

//FindRemoteSystemByName finds remote system by name
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
	log.Debugf("Remote system name: %s Id: %s", remoteSystemName, remoteSystemNameResp.RemoteSystemContent.ID)
	return remoteSystemNameResp, nil
}

//FindRemoteSystemByID finds remote system by id
func (r *Replication) FindRemoteSystemByID(ctx context.Context, remoteSystemID string) (*types.RemoteSystem, error) {
	log := util.GetRunIDLogger(ctx)
	remoteSystemID, err := util.ValidateResourceName(remoteSystemID, api.MaxResourceNameLength)
	if err != nil {
		return nil, err
	}
	remoteSystem := &types.RemoteSystem{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.RemoteSystemAction, remoteSystemID, RemoteSystemFields), nil, remoteSystem)
	if err != nil {
		return nil, fmt.Errorf("unable to find Remote system %v. Error: %v", remoteSystemID, err)
	}
	log.Debugf("Remote system Id: %s", remoteSystemID)
	return remoteSystem, nil
}

//CreateReplicationSession creates replication session
func (r *Replication) CreateReplicationSession(ctx context.Context, replicationSessionName, srcResourceID, dstResourceID, remoteSystemName string, maxTimeOutOfSync string) (*types.ReplicationSession, error) {
	var createRS types.CreateReplicationSessionParam
	if len(srcResourceID) == 0 {
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
	createRS.SrcResourceID = srcResourceID
	createRS.DstResourceID = dstResourceID
	createRS.OverwriteDestination = true
	rsResp := &types.ReplicationSession{}
	err = r.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.ReplicationSessionAction), createRS, rsResp)
	if err != nil {
		return nil, err
	}
	return rsResp, nil
}

//DeleteConsistencyGroup - Delete ReplicationSession by ID.
func (r *Replication) DeleteConsistencyGroup(ctx context.Context, rID string) error {
	if len(rID) == 0 {
		return errors.New("ReplicationSession Id cannot be empty")
	}

	_, err := r.FindReplicationSessionByID(ctx, rID)

	if err != nil {
		return err
	}

	deleteErr := r.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.ReplicationSessionAction, rID), nil, nil)

	if deleteErr != nil {
		return fmt.Errorf("delete ReplicationSession %s Failed. Error: %v", rID, deleteErr)
	}

	return nil
}

//FindReplicationSessionBySrcResourceID finds replication session by storage resource id
func (r *Replication) FindReplicationSessionBySrcResourceID(ctx context.Context, srcResourceID string) (*types.ReplicationSession, error) {
	if len(srcResourceID) == 0 {
		return nil, fmt.Errorf("SrcResourceID shouldn't be empty")
	}
	filter := fmt.Sprintf("srcResourceID eq %s", "\""+srcResourceID+"\"")
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

//FindReplicationSessionByName finds replication session by session name
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
	log.Debugf("Replcation Session name: %s Id: %s", rsResp.ReplicationSessionContent.Name, rsResp.ReplicationSessionContent.ReplicationSessionID)
	return rsResp, nil
}

//FindReplicationSessionByID finds replication session by session id
func (r *Replication) FindReplicationSessionByID(ctx context.Context, rsID string) (*types.ReplicationSession, error) {
	log := util.GetRunIDLogger(ctx)
	if rsID == "" {
		return nil, errors.New("Replication session ID cannot be empty")
	}
	rsResp := &types.ReplicationSession{}
	err := r.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.ReplicationSessionAction, rsID, ReplicationSessionFields), nil, rsResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find Replication Session id %s Error: %v", rsID, err)
	}
	log.Debugf("Replcation Session name: %s Id: %s", rsResp.ReplicationSessionContent.Name, rsResp.ReplicationSessionContent.ReplicationSessionID)
	return rsResp, nil
}

//DeleteReplicationSessionByID deletes replication session by session id
func (r *Replication) DeleteReplicationSessionByID(ctx context.Context, sessionID string) error {
	log := util.GetRunIDLogger(ctx)
	if len(sessionID) == 0 {
		return errors.New("Replication session Id cannot be empty")
	}

	_, err := r.FindReplicationSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	deleteErr := r.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.ReplicationSessionAction, sessionID), nil, nil)
	if deleteErr != nil {
		return fmt.Errorf("delete replication session %s Failed. Error: %v", sessionID, deleteErr)
	}
	log.Debugf("Delete Replication session %s Successful", sessionID)

	return nil
}
