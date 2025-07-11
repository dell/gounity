/*
 Copyright © 2019-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

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

	util "github.com/dell/gounity/gounityutil"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/apitypes"
)

// Host not found error variables
var (
	ErrorHostNotFound          = errors.New("unable to find host")
	ErrorMultipleHostFound     = errors.New("Found multiple hosts with same name. Delete the duplicate entries on the array")
	MultipleHostFoundErrorCode = "0x7d13158"
	HostNotFoundErrorCode      = "0x7d13005"
)

// FindHostByName Finds the Host by it's name. If the Host is not found, an error will be returned.
func (c *UnityClientImpl) FindHostByName(ctx context.Context, hostName string) (*types.Host, error) {
	log := util.GetRunIDLogger(ctx)
	if len(hostName) == 0 {
		return nil, errors.New("host Name shouldn't be empty")
	}
	hResponse := &types.Host{}
	log.Info("URI", fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.HostAction, hostName, HostfieldsToQuery))
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.HostAction, hostName, HostfieldsToQuery), nil, hResponse)
	if err != nil {
		// Using the multiple host found error code(MultipleHostFoundErrorCode) for comparison
		if strings.Contains(err.Error(), MultipleHostFoundErrorCode) {
			return nil, ErrorMultipleHostFound
		} else if strings.Contains(err.Error(), HostNotFoundErrorCode) {
			return nil, ErrorHostNotFound
		}
		return nil, err
	}
	return hResponse, nil
}

// CreateHost Create a new Host
func (c *UnityClientImpl) CreateHost(ctx context.Context, hostName string, tenantID string) (*types.Host, error) {
	if len(hostName) == 0 {
		return nil, errors.New("hostname shouldn't be empty")
	}
	tenantIDStruct := types.Tenants{
		TenantID: tenantID,
	}
	hostReq := &types.HostCreateParam{
		Type:        "1", // Initiator type hardcoded as "1" for FC Initiator
		Name:        hostName,
		Description: hostName,
		OsType:      "Linux",
	}

	if tenantID != "" {
		hostReq.Tenant = &tenantIDStruct
	}

	hostResp := &types.Host{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.HostAction), hostReq, hostResp)
	if err != nil {
		return nil, err
	}
	return hostResp, nil
}

// DeleteHost function is used only in unit tests
func (c *UnityClientImpl) DeleteHost(ctx context.Context, hostName string) error {
	if len(hostName) == 0 {
		return fmt.Errorf("hostname shouldn't be empty")
	}

	hostResp := &types.Host{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceByNameURI, api.HostAction, hostName), nil, hostResp)
	if err != nil {
		return err
	}
	return nil
}

// CreateHostIPPort - Create Host IP Port
func (c *UnityClientImpl) CreateHostIPPort(ctx context.Context, hostID, ip string) (*types.HostIPPort, error) {
	if len(hostID) == 0 {
		return nil, errors.New("host ID shouldn't be empty")
	}

	hostIDContent := types.HostIDContent{
		ID: hostID,
	}

	hostIPReq := &types.HostIPPortCreateParam{
		HostIDContent: &hostIDContent,
		Address:       ip,
	}

	hostIPResp := &types.HostIPPort{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.HostIPPortAction), hostIPReq, hostIPResp)
	if err != nil {
		return nil, err
	}
	return hostIPResp, nil
}

// FindHostIPPortByID method to get host Ip port object from Unity by cli ID
func (c *UnityClientImpl) FindHostIPPortByID(ctx context.Context, hostIPID string) (*types.HostIPPort, error) {
	hostIPResp := &types.HostIPPort{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.HostIPPortAction, hostIPID, HostIPPortDisplayFields), nil, hostIPResp)
	if err != nil {
		return nil, err
	}
	return hostIPResp, nil
}

// ListHostInitiators lists all host initiators
func (c *UnityClientImpl) ListHostInitiators(ctx context.Context) ([]types.HostInitiator, error) {
	listInitiatorResp := &types.ListHostInitiator{}
	hostInitiatorURI := api.UnityListHostInitiatorsURI + HostInitiatorsDisplayFields
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, hostInitiatorURI, nil, listInitiatorResp)
	if err != nil {
		return nil, err
	}
	return listInitiatorResp.HostInitiator, nil
}

// FindHostInitiatorByName - Find Host Initiator by name
func (c *UnityClientImpl) FindHostInitiatorByName(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error) {
	if len(wwnOrIqn) == 0 {
		return nil, errors.New("host Initiator Name shouldn't be empty")
	}

	list, err := c.ListHostInitiators(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range list {
		if strings.ToLower(i.HostInitiatorContent.InitiatorID) == strings.ToLower(wwnOrIqn) {
			return &i, nil
		}
	}

	// @TODO The following code should work. Unity rest api having a bug querying host initiators by host initiatorID
	//hostInitiatorResp := &types.HostInitiator{}
	//err := c.executeWithRetryAuthenticate(http.MethodGet, fmt.Sprintf(api.UnityApiGetResourceByPropertyWithFieldsUri, "hostInitiator", "initiatorID" ,wwnOrIqn, api.HostInitiatorsDisplayFields), nil, hostInitiatorResp)
	//if err != nil {
	//	log.Info("Unable to find host initiator:", wwnOrIqn)
	//	return nil, errors.New(fmt.Sprintf("Unable to find host %s", wwnOrIqn))
	//}

	return nil, errors.New("wwn or iqn not found")
}

// FindHostInitiatorByID - Find Host Initiator
func (c *UnityClientImpl) FindHostInitiatorByID(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error) {
	hostInitiatorResp := &types.HostInitiator{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.HostInitiatorAction, wwnOrIqn, HostInitiatorsDisplayFields), nil, hostInitiatorResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find host %s : %v", wwnOrIqn, err)
	}
	return hostInitiatorResp, nil
}

// CreateHostInitiator - Create Host Initiator
func (c *UnityClientImpl) CreateHostInitiator(ctx context.Context, hostID, wwnOrIqn string, initiatorType types.InitiatorType) (*types.HostInitiator, error) {
	log := util.GetRunIDLogger(ctx)
	if len(hostID) == 0 {
		return nil, errors.New("host ID shouldn't be empty")
	}

	if len(wwnOrIqn) == 0 {
		return nil, fmt.Errorf("wwnOrIqn shouldn't be empty")
	}

	hostInitiatorResp := &types.HostInitiator{}

	log.Debugf("Finding Initiator: %s", wwnOrIqn)
	initiator, err := c.FindHostInitiatorByName(ctx, wwnOrIqn)
	log.Debugf("FindHostInitiatorByName: %v Error: %v", initiator, err)
	if err != nil {
		log.Debugf("Initiator not found. Adding new Initiator: %s to host: %s \n", wwnOrIqn, hostID)
		hostIDContent := types.HostIDContent{
			ID: hostID,
		}

		hostInitiatorReq := &types.HostInitiatorCreateParam{
			HostIDContent: &hostIDContent,
			InitiatorType: initiatorType,
			InitiatorWwn:  wwnOrIqn,
		}
		err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.HostInitiatorAction), hostInitiatorReq, hostInitiatorResp)
		if err != nil {
			return nil, fmt.Errorf("create Host Initiator %s Error: %v", wwnOrIqn, err)
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == "" {
		log.Debugf("Initiator found, but parent host is not added. Updating the existing Initiator: %s to host: %s \n", wwnOrIqn, hostID)
		initiator, err = c.ModifyHostInitiator(ctx, hostID, initiator)
		if err != nil {
			return nil, fmt.Errorf("modify Host Initiator %s Error: %v", wwnOrIqn, err)
		}
	} else if initiator.HostInitiatorContent.ParentHost.ID == hostID {
		log.Debugf("Initiator found and already added to existing host Initiator: %s to host: %s \n", wwnOrIqn, hostID)
	} else if initiator.HostInitiatorContent.ParentHost.ID != hostID {
		return nil, fmt.Errorf("initiator found (%s), and attached to someother host (%s) instead of host: %s", wwnOrIqn, initiator.HostInitiatorContent.ParentHost.ID, hostID)
	} else {
		log.Error("Initiator unknown operation.")
	}

	return hostInitiatorResp, nil
}

// ModifyHostInitiator - WILL BE DEPRECATED
func (c *UnityClientImpl) ModifyHostInitiator(ctx context.Context, hostID string, initiator *types.HostInitiator) (*types.HostInitiator, error) {
	if initiator == nil {
		return nil, errors.New("HostInitiator shouldn't be null")
	}

	return c.ModifyHostInitiatorByID(ctx, hostID, initiator.HostInitiatorContent.ID)
}

// ModifyHostInitiatorByID function modifies host initiator by ID
func (c *UnityClientImpl) ModifyHostInitiatorByID(ctx context.Context, hostID, initiatorID string) (*types.HostInitiator, error) {
	if hostID == "" {
		return nil, errors.New("Host ID shouldn't be null")
	}

	if initiatorID == "" {
		return nil, errors.New("Initiator ID shouldn't be null")
	}

	hostIDContent := types.HostIDContent{
		ID: hostID,
	}
	hostInitiatorReq := &types.HostInitiatorModifyParam{
		HostIDContent: &hostIDContent,
	}
	hostInitiatorResp := &types.HostInitiator{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyHostInitiators, initiatorID), hostInitiatorReq, hostInitiatorResp)
	if err != nil {
		return nil, err
	}
	return hostInitiatorResp, nil
}

// FindHostInitiatorPathByID Finds Host Initiator
func (c *UnityClientImpl) FindHostInitiatorPathByID(ctx context.Context, initiatorPathID string) (*types.HostInitiatorPath, error) {
	hostInitiatorPathResp := &types.HostInitiatorPath{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.HostInitiatorPathAction, initiatorPathID, HostInitiatorPathDisplayFields), nil, hostInitiatorPathResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find host initiator path %s : %v", initiatorPathID, err)
	}
	return hostInitiatorPathResp, nil
}

// FindFcPortByID Finds FC Port
func (c *UnityClientImpl) FindFcPortByID(ctx context.Context, fcPortID string) (*types.FcPort, error) {
	fcPortResp := &types.FcPort{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, HostInitiatorPathDisplayFields, fcPortID, FcPortDisplayFields), nil, fcPortResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find Fc port %s : %v", fcPortID, err)
	}
	return fcPortResp, nil
}

// FindTenants finds tenants
func (c *UnityClientImpl) FindTenants(ctx context.Context) (*types.TenantInfo, error) {
	tenantsResp := &types.TenantInfo{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetTenantURI, api.TenantAction, TenantDisplayFields), nil, tenantsResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find tenants : %v", err)
	}
	return tenantsResp, nil
}
