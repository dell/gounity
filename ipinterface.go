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

	"github.com/dell/gounity/util"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

type ipinterface struct {
	client *Client
}

func NewIpInterface(client *Client) *ipinterface {
	return &ipinterface{client}
}

//ListIscsiIPInterfaces - List the IpnInterfaces configured on the array
func (f *ipinterface) ListIscsiIPInterfaces(ctx context.Context) ([]types.IPInterfaceEntries, error) {

	log := util.GetRunIdLogger(ctx)
	hResponse := &types.ListIPInterfaces{}
	log.Debugf("URI: "+api.UnityApiInstanceTypeResourcesWithFields, api.IPInterface, IscsiIPFields)
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiInstanceTypeResourcesWithFields, api.IPInterface, IscsiIPFields), nil, hResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to list Ip Interfaces %v", err))
	}
	var iscsiInterfaces []types.IPInterfaceEntries
	for _, ipInterface := range hResponse.Entries {
		IPContent := &ipInterface.IPInterfaceContent
		if IPContent != nil && ipInterface.IPInterfaceContent.Type == 2 { //2 stands for iScsi Interface in Unisphere 5.0. Verifu while qualifying higher versions
			iscsiInterfaces = append(iscsiInterfaces, ipInterface)
		}
	}
	return iscsiInterfaces, nil
}
