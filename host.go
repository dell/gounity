/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/payloads"
)

type host struct {
	client *Client
}

var HostNotFoundError = errors.New("unable to find host")

func NewHost(client *Client) *host {
	return &host{client}
}

//Find the Host by it's name. If the Host is not found, an error will be returned.
func (h *host) FindHostByName(hostName string) (*types.Host, error) {
	if len(hostName) == 0 {
		return nil, errors.New("host Name shouldn't be empty")
	}
	hResponse := &types.Host{}
	fieldsToQuery := "id,name,description,fcHostInitiators"
	log.Info("URI", fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "host", hostName, fieldsToQuery))
	err := h.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByNameWithFieldsUri, "host", hostName, fieldsToQuery), nil, hResponse)
	if err != nil {
		log.Error("Unable to find host", err)
		return nil, HostNotFoundError
	}
	return hResponse, nil
}

//Create a new Host
func (h *host) CreateHost(hostName string) (*types.Host, error) {
	if len(hostName) == 0 {
		return nil, errors.New("hostname shouldn't be empty")
	}

	hostReq := &types.HostCreateParam{
		Type:        "1", //Initiator type hardcoded as "1" for FC Initiator
		Name:        hostName,
		Description: hostName,
		OsType:      "Linux",
	}

	hostResp := &types.Host{}
	err := h.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, "host"), hostReq, hostResp)
	if err != nil {
		return nil, err
	}
	return hostResp, nil
}

//Delete Host. This function is used only in unit tests
func (h *host) DeleteHost(hostName string) error {
	if len(hostName) == 0 {
		return errors.New(fmt.Sprintf("Hostname shouldn't be empty."))
	}

	hostResp := &types.Host{}
	err := h.client.executeWithRetryAuthenticate(http.MethodDelete, fmt.Sprintf(api.UnityApiGetResourceByNameUri, "host", hostName), nil, hostResp)
	if err != nil {
		return err
	}
	return nil
}

//Create Host IP Port
func (h *host) CreateHostIpPort(hostId, ip string) (*types.HostIpPort, error) {
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
	err := h.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, "hostIPPort"), hostIpReq, hostIpResp)
	if err != nil {
		return nil, err
	}
	return hostIpResp, nil
}

// ListHostInitiators lists all host initiators
func (h *host) ListHostInitiators() ([]types.HostInitiator, error) {
	listInitiatorResp := &types.ListHostInitiator{}
	hostInitiatorUri := "/api/types/hostInitiator/instances?fields=" + api.HostInitiatorsDisplayFields
	err := h.client.executeWithRetryAuthenticate(http.MethodGet, hostInitiatorUri, nil, listInitiatorResp)
	if err != nil {
		return nil, err
	}
	return listInitiatorResp.HostInitiator, nil
}

//Find Host Initiator
func (h *host) FindHostInitiatorByName(wwnOrIqn string) (*types.HostInitiator, error) {
	if len(wwnOrIqn) == 0 {
		return nil, errors.New("host Initiator Name shouldn't be empty")
	}

	list, err := h.ListHostInitiators()
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
func (h *host) FindHostInitiatorById(wwnOrIqn string) (*types.HostInitiator, error) {
	hostInitiatorResp := &types.HostInitiator{}
	err := h.client.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceWithFieldsUri, "hostInitiator", wwnOrIqn, api.HostInitiatorsDisplayFields), nil, hostInitiatorResp)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to find host %s", wwnOrIqn))
		return nil, errors.New(fmt.Sprintf("Unable to find host %s : %v", wwnOrIqn, err))
	}
	return hostInitiatorResp, nil
}

//Create Host Initiator
func (h *host) CreateHostInitiator(hostId, wwnOrIqn string) (*types.HostInitiator, error) {
	if len(hostId) == 0 {
		return nil, errors.New("host ID shouldn't be empty")
	}

	if len(wwnOrIqn) == 0 {
		return nil, errors.New(fmt.Sprintf("WwnOrIqn shouldn't be empty."))
	}

	hostInitiatorResp := &types.HostInitiator{}

	log.Infof("Finding Initiator: %s", wwnOrIqn)
	initiator, err := h.FindHostInitiatorByName(wwnOrIqn)
	log.Infof("FindHostInitiatorByName: %v Error: %v", initiator, err)
	if err != nil {
		log.Infof("Initiator not found. Adding new Initiator: %s to host: %s \n", wwnOrIqn, hostId)
		hostIdContent := types.HostIdContent{
			ID: hostId,
		}

		hostInitiatorReq := &types.HostInitiatorCreateParam{
			HostIdContent: &hostIdContent,
			InitiatorType: "1",
			InitiatorWwn:  wwnOrIqn,
		}
		err := h.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf(api.UnityApiInstanceTypeResources, "hostInitiator"), hostInitiatorReq, hostInitiatorResp)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Create Host Initiator %s Error: %v", wwnOrIqn, err))
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == "" {
		log.Infof("Initiator found, but parent host is not added. Updating the existing Initiator: %s to host: %s \n", wwnOrIqn, hostId)
		initiator, err = h.ModifyHostInitiator(hostId, initiator)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Modify Host Initiator %s Error: %v", wwnOrIqn, err))
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == hostId {
		log.Infof("Initiator found and already added to existing host Initiator: %s to host: %s \n", wwnOrIqn, hostId)
	} else if initiator.HostInitiatorContent.ParentHost.ID != hostId {
		return nil, errors.New(fmt.Sprintf("Initiator found (%s), and attached to someother host (%s) instead of host: %s", wwnOrIqn, initiator.HostInitiatorContent.ParentHost.ID, hostId))
	} else {
		log.Error("Initiator unknown operation.")
	}

	return hostInitiatorResp, nil
}

//Modify Host Initiator
func (h *host) ModifyHostInitiator(hostId string, initiator *types.HostInitiator) (*types.HostInitiator, error) {
	if initiator == nil {
		return nil, errors.New("HostInitiator shouldn't be null")
	}

	hostIdContent := types.HostIdContent{
		ID: hostId,
	}
	hostInitiatorReq := &types.HostInitiatorModifyParam{
		HostIdContent: &hostIdContent,
	}

	initiator.HostInitiatorContent.ParentHost.ID = hostId
	hostInitiatorResp := &types.HostInitiator{}
	err := h.client.executeWithRetryAuthenticate(http.MethodPost, fmt.Sprintf("/api/instances/hostInitiator/%s/action/modify", initiator.HostInitiatorContent.Id), hostInitiatorReq, hostInitiatorResp)
	if err != nil {
		return nil, err
	}
	return hostInitiatorResp, nil
}
