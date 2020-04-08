/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package gounity

import (
	"context"
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
	fieldsToQuery := "id,ipAddress,type"
	log.Debugf("URI: "+api.UnityApiInstanceTypeResourcesWithFields, "ipInterface", fieldsToQuery)
	err := f.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityApiInstanceTypeResourcesWithFields, "ipInterface", fieldsToQuery), nil, hResponse)
	if err != nil {
		log.Error("Unable to list Ip Interfaces", err)
		return nil, err
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
