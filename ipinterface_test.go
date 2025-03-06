/*
 Copyright © 2020-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

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
	"testing"

	mocksapi "github.com/dell/gounity/mocks/api"
	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListIscsiIPInterfaces(t *testing.T) {
	assert := assert.New(t)

	// Initial Setup
	t.Log("Begin - List IP Interfaces Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock ListIscsiIPInterfaces to return example data
	expectedIPInterfaces := &types.ListIPInterfaces{
		Entries: []types.IPInterfaceEntries{
			{IPInterfaceContent: types.IPInterfaceContent{Type: 2, IPAddress: "192.168.1.100"}},
			{IPInterfaceContent: types.IPInterfaceContent{Type: 2, IPAddress: "192.168.1.101"}},
		},
	}

	mockClient := testConf.client.(*UnityClientImpl).api.(*mocksapi.Client)
	mockClient.On("DoWithHeaders", mock.Anything, "GET", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.ListIPInterfaces")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.ListIPInterfaces)
			*resp = *expectedIPInterfaces
		}).Once()

	// Call the method
	ipInterfaces, err := testConf.client.ListIscsiIPInterfaces(ctx)

	// Verify the results for the main case
	assert.NoError(err, "List IP Interfaces should not return an error")
	assert.Len(ipInterfaces, 2, "Expected 2 IP interfaces")
	for _, ipInterface := range ipInterfaces {
		t.Logf("IP Address of interface: %s", ipInterface.IPInterfaceContent.IPAddress)
	}

	// Negative Cases

	// Case: API returns an error
	mockClient.On("DoWithHeaders", mock.Anything, "GET", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.ListIPInterfaces")).Return(
		fmt.Errorf("API call error"),
	).Once()
	_, err = testConf.client.ListIscsiIPInterfaces(ctx)
	assert.Error(err, "Expected error when API call returns an error")
	t.Log("Negative case: API call error - successful")

	// Case: No iSCSI interfaces found
	mockClient.On("DoWithHeaders", mock.Anything, "GET", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.ListIPInterfaces")).Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(5).(*types.ListIPInterfaces)
			*resp = types.ListIPInterfaces{
				Entries: []types.IPInterfaceEntries{
					{IPInterfaceContent: types.IPInterfaceContent{Type: 1, IPAddress: "192.168.1.102"}}, // Not an iSCSI interface
				},
			}
		}).Once()
	ipInterfaces, err = testConf.client.ListIscsiIPInterfaces(ctx)
	assert.NoError(err, "List IP Interfaces with no iSCSI interfaces should not return an error")
	assert.Len(ipInterfaces, 0, "Expected 0 iSCSI IP interfaces")
	t.Log("Negative case: No iSCSI interfaces - successful")

	// Mock network error
	mockClient.On("DoWithHeaders", mock.Anything, "GET", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*types.ListIPInterfaces")).Return(
		fmt.Errorf("network error"),
	).Once()
	_, err = testConf.client.ListIscsiIPInterfaces(ctx)
	assert.Error(err, "Expected network error")
	t.Log("Negative case: Network error successfully validated")

	t.Log("List IP Interfaces Test - Successful")
}
