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
	"fmt"
	"net/http"

	util "github.com/dell/gounity/gounityutil"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/apitypes"
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
