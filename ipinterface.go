/*
Copyright (c) 2019 Dell Corporation
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

// ListIscsiIPInterfaces - List the IpnInterfaces configured on the array
func (c *UnityClientImpl) ListIscsiIPInterfaces(ctx context.Context) ([]types.IPInterfaceEntries, error) {
	log := util.GetRunIDLogger(ctx)
	hResponse := &types.ListIPInterfaces{}
	log.Debugf("URI: "+api.UnityAPIInstanceTypeResourcesWithFields, api.IPInterface, IscsiIPFields)
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIInstanceTypeResourcesWithFields, api.IPInterface, IscsiIPFields), nil, hResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to list Ip Interfaces %v", err)
	}
	var iscsiInterfaces []types.IPInterfaceEntries
	for _, ipInterface := range hResponse.Entries {
		IPContent := &ipInterface.IPInterfaceContent                      // #nosec G601
		if IPContent != nil && ipInterface.IPInterfaceContent.Type == 2 { // 2 stands for iScsi Interface in Unisphere 5.0. Verify while qualifying higher versions
			iscsiInterfaces = append(iscsiInterfaces, ipInterface)
		}
	}
	return iscsiInterfaces, nil
}
