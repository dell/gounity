/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/gounity/util"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

type host struct {
	client *Client
}

var (
	HostNotFoundError          = errors.New("unable to find host")
	MultipleHostFoundError     = errors.New("Found multiple hosts with same name. Delete the duplicate entries on the array")
	MultipleHostFoundErrorCode = "0x7d13158"
	HostNotFoundErrorCode = "0x7d13005"
)

func NewHost(client *Client) *host {
	return &host{client}
}

//Find the Host by it's name. If the Host is not found, an error will be returned.
func (h *host) FindHostByName(ctx context.Context, hostName string) (*types.Host, error) {
	log := util.GetRunIdLogger(ctx)
	if len(hostName) == 0 {
		return nil, errors.New("host Name shouldn't be empty")
	}
	hResponse := &types.Host{}
	log.Info("URI", fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, api.HostAction, hostName, HostfieldsToQuery))
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, api.HostAction, hostName, HostfieldsToQuery), nil, hResponse)
	if err != nil {
		//Using the multiple host found error code(MultipleHostFoundErrorCode) for comparison
		if strings.Contains(err.Error(), MultipleHostFoundErrorCode) {
			return nil, MultipleHostFoundError
		} else if strings.Contains(err.Error(), HostNotFoundErrorCode){
			return nil, HostNotFoundError
		} else{
			return nil, err
		}
	}
	return hResponse, nil
}

//Create a new Host
func (h *host) CreateHost(ctx context.Context, hostName string, tenantId string) (*types.Host, error) {
	if len(hostName) == 0 {
		return nil, errors.New("hostname shouldn't be empty")
	}
	tenantIdStruct := types.Tenants{
		TenantId: tenantId,
	}
	hostReq := &types.HostCreateParam{
		Type:        "1", //Initiator type hardcoded as "1" for FC Initiator
		Name:        hostName,
		Description: hostName,
		OsType:      "Linux",
	}

	if tenantId != "" {
		hostReq.Tenant=&tenantIdStruct
	}

	hostResp := &types.Host{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, api.HostAction), hostReq, hostResp)
	if err != nil {
		return nil, err
	}
	return hostResp, nil
}

//Delete Host. This function is used only in unit tests
func (h *host) DeleteHost(ctx context.Context, hostName string) error {
	if len(hostName) == 0 {
		return errors.New(fmt.Sprintf("Hostname shouldn't be empty."))
	}

	hostResp := &types.Host{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceByNameUri, api.HostAction, hostName), nil, hostResp)
	if err != nil {
		return err
	}
	return nil
}

//Create Host IP Port
func (h *host) CreateHostIpPort(ctx context.Context, hostId, ip string) (*types.HostIpPort, error) {
	if len(hostId) == 0 {
		return nil, errors.New("host ID shouldn't be empty")
	}

	hostIdContent := types.HostIdContent{
		ID: hostId,
	}

	hostIpReq := &types.HostIpPortCreateParam{
		HostIdContent: &hostIdContent,
		Address:       ip,
	}

	hostIpResp := &types.HostIpPort{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, api.HostIPPortAction), hostIpReq, hostIpResp)
	if err != nil {
		return nil, err
	}
	return hostIpResp, nil
}

// FindHostIpPortById method to get host Ip port object from Unity by cli ID
func (h *host) FindHostIpPortById(ctx context.Context, hostIpID string) (*types.HostIpPort, error) {
	hostIpResp := &types.HostIpPort{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, api.HostIPPortAction, hostIpID, HostIpPortDisplayFields), nil, hostIpResp)
	if err != nil {
		return nil, err
	}
	return hostIpResp, nil
}

// ListHostInitiators lists all host initiators
func (h *host) ListHostInitiators(ctx context.Context) ([]types.HostInitiator, error) {
	listInitiatorResp := &types.ListHostInitiator{}
	hostInitiatorUri := api.UnityListHostInitiatorsUri + HostInitiatorsDisplayFields
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, hostInitiatorUri, nil, listInitiatorResp)
	if err != nil {
		return nil, err
	}
	return listInitiatorResp.HostInitiator, nil
}

//Find Host Initiator
func (h *host) FindHostInitiatorByName(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error) {
	if len(wwnOrIqn) == 0 {
		return nil, errors.New("host Initiator Name shouldn't be empty")
	}

	list, err := h.ListHostInitiators(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range list {
		if strings.ToLower(i.HostInitiatorContent.InitiatorId) == strings.ToLower(wwnOrIqn) {
			return &i, nil
		}
	}

	// @TODO The following code should work. Unity rest api having a bug querying host initiators by host initiatorId
	//hostInitiatorResp := &types.HostInitiator{}
	//err := h.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByPropertyWithFieldsUri, "hostInitiator", "initiatorId" ,wwnOrIqn, api.HostInitiatorsDisplayFields), nil, hostInitiatorResp)
	//if err != nil {
	//	log.Info("Unable to find host initiator:", wwnOrIqn)
	//	return nil, errors.New(fmt.Sprintf("Unable to find host %s", wwnOrIqn))
	//}

	return nil, errors.New("wwn or iqn not found")
}

//Find Host Initiator
func (h *host) FindHostInitiatorById(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error) {
	hostInitiatorResp := &types.HostInitiator{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, api.HostInitiatorAction, wwnOrIqn, HostInitiatorsDisplayFields), nil, hostInitiatorResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find host %s : %v", wwnOrIqn, err))
	}
	return hostInitiatorResp, nil
}

//Create Host Initiator
func (h *host) CreateHostInitiator(ctx context.Context, hostId, wwnOrIqn string, initiatorType types.InitiatorType) (*types.HostInitiator, error) {
	log := util.GetRunIdLogger(ctx)
	if len(hostId) == 0 {
		return nil, errors.New("host ID shouldn't be empty")
	}

	if len(wwnOrIqn) == 0 {
		return nil, errors.New(fmt.Sprintf("WwnOrIqn shouldn't be empty."))
	}

	hostInitiatorResp := &types.HostInitiator{}

	log.Debugf("Finding Initiator: %s", wwnOrIqn)
	initiator, err := h.FindHostInitiatorByName(ctx, wwnOrIqn)
	log.Debugf("FindHostInitiatorByName: %v Error: %v", initiator, err)
	if err != nil {
		log.Debugf("Initiator not found. Adding new Initiator: %s to host: %s \n", wwnOrIqn, hostId)
		hostIdContent := types.HostIdContent{
			ID: hostId,
		}

		hostInitiatorReq := &types.HostInitiatorCreateParam{
			HostIdContent: &hostIdContent,
			InitiatorType: initiatorType,
			InitiatorWwn:  wwnOrIqn,
		}
		err := h.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, api.HostInitiatorAction), hostInitiatorReq, hostInitiatorResp)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Create Host Initiator %s Error: %v", wwnOrIqn, err))
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == "" {
		log.Debugf("Initiator found, but parent host is not added. Updating the existing Initiator: %s to host: %s \n", wwnOrIqn, hostId)
		initiator, err = h.ModifyHostInitiator(ctx, hostId, initiator)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Modify Host Initiator %s Error: %v", wwnOrIqn, err))
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == hostId {
		log.Debugf("Initiator found and already added to existing host Initiator: %s to host: %s \n", wwnOrIqn, hostId)
	} else if initiator.HostInitiatorContent.ParentHost.ID != hostId {
		return nil, errors.New(fmt.Sprintf("Initiator found (%s), and attached to someother host (%s) instead of host: %s", wwnOrIqn, initiator.HostInitiatorContent.ParentHost.ID, hostId))
	} else {
		log.Error("Initiator unknown operation.")
	}

	return hostInitiatorResp, nil
}

//Modify Host Initiator - WILL BE DEPRECATED
func (h *host) ModifyHostInitiator(ctx context.Context, hostId string, initiator *types.HostInitiator) (*types.HostInitiator, error) {
	if initiator == nil {
		return nil, errors.New("HostInitiator shouldn't be null")
	}
	
	return h.ModifyHostInitiatorById(ctx, hostId, initiator.HostInitiatorContent.Id)
}

// ModifyHostInitiatorById
func (h *host) ModifyHostInitiatorById(ctx context.Context, hostId , initiatorId string) (*types.HostInitiator, error) {

	if hostId == "" {
		return nil, errors.New("Host ID shouldn't be null")
	}

	if initiatorId == "" {
		return nil, errors.New("Initiator ID shouldn't be null")
	}

	hostIdContent := types.HostIdContent{
		ID: hostId,
	}
	hostInitiatorReq := &types.HostInitiatorModifyParam{
		HostIdContent: &hostIdContent,
	}
	hostInitiatorResp := &types.HostInitiator{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyHostInitiators, initiatorId), hostInitiatorReq, hostInitiatorResp)
	if err != nil {
		return nil, err
	}
	return hostInitiatorResp, nil
}

//Find Host Initiator
func (h *host) FindHostInitiatorPathById(ctx context.Context, initiatorPathId string) (*types.HostInitiatorPath, error) {
	hostInitiatorPathResp := &types.HostInitiatorPath{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, api.HostInitiatorPathAction, initiatorPathId, HostInitiatorPathDisplayFields), nil, hostInitiatorPathResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find host initiator path %s : %v", initiatorPathId, err))
	}
	return hostInitiatorPathResp, nil
}

//Find FC Port
func (h *host) FindFcPortById(ctx context.Context, fcPortId string) (*types.FcPort, error) {
	fcPortResp := &types.FcPort{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, HostInitiatorPathDisplayFields, fcPortId, FcPortDisplayFields), nil, fcPortResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find Fc port %s : %v", fcPortId, err))
	}
	return fcPortResp, nil
}

func (h *host) FindTenants(ctx context.Context) (*types.Host, error) {
	tenantsResp := &types.Tenants{}
	err := h.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiGetTenantUri,api.TenantAction,TenantDisplayFields), nil, tenantsResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find tenants : %v", err))
	}
	return  tenantsResp, nil
}
