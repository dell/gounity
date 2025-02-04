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
	"testing"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/mocks"
	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	hostName           string
	hostID             string
	hostIPPortID       string
	iqnInitiatorID     string
	wwnInitiatorPathID string
	fcPortID           string
	iqnInitiator       *types.HostInitiator
)

func TestCreateHost(t *testing.T) {
	assert := require.New(t)
	fmt.Println("Begin - Create Host Test")
	testConf.client.getAPI().(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock setup for valid host creation
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Positive Case
	hostName := "valid_host_name" // Ensure this is set to a valid name
	host, err := testConf.client.CreateHost(ctx, hostName, testConf.tenant)
	if err != nil {
		t.Fatalf("Create Host failed: %v", err)
	}
	hostID = host.HostContent.ID
	assert.NoError(err, "Create Host failed")
	assert.NotNil(host, "Host should not be nil")

	// Negative Cases
	// Mock setup for empty host name
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("hostname shouldn't be empty")).Once()

	hostNameTemp := ""
	_, err = testConf.client.CreateHost(ctx, hostNameTemp, testConf.tenant)
	assert.Error(err, "Expected error for empty host name")
	assert.EqualError(err, "hostname shouldn't be empty", "Unexpected error message")

	// Mock setup for invalid tenant ID
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("hostname shouldn't be empty")).Once()

	tenantIDTemp := "tenant_invalid_1"
	_, err = testConf.client.CreateHost(ctx, hostName, tenantIDTemp)
	assert.Error(err, "Expected error for invalid tenant ID")
	assert.EqualError(err, "hostname shouldn't be empty", "Unexpected error message")

	fmt.Println("Create Host Test Successful")
}

func TestFindHostByName(t *testing.T) {
	fmt.Println("Begin - Find Host by name Test")
	ctx := context.Background()

	// Negative test cases
	hostNameTemp := ""
	_, err := testConf.client.FindHostByName(ctx, hostNameTemp)
	assert.Equal(t, errors.New("host Name shouldn't be empty"), err)

	hostNameTemp = "dummy-host-1"
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("error")).Once()
	_, err = testConf.client.FindHostByName(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Find Host with invalid hostName - Negative case failed")
	}

	fmt.Println("Find Host by name Successful")
}

func TestCreateHostIPPort(t *testing.T) {
	assert := require.New(t)
	fmt.Println("Begin - Create Host IP Port Test")
	testConf.client.getAPI().(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock setup for valid host IP port creation
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Positive Case
	hostID := "valid_host_id" // Ensure this is set to a valid ID
	hostIPPort, err := testConf.client.CreateHostIPPort(ctx, hostID, testConf.nodeHostIP)
	if err != nil {
		t.Fatalf("CreateHostIPPort failed: %v", err)
	}
	hostIPPortID = hostIPPort.HostIPContent.ID
	assert.NoError(err, "CreateHostIPPort failed")
	assert.NotNil(hostIPPort, "Host IP Port should not be nil")

	// Negative Cases
	// Mock setup for empty host ID
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("host ID shouldn't be empty")).Once()

	hostIDTemp := ""
	_, err = testConf.client.CreateHostIPPort(ctx, hostIDTemp, testConf.nodeHostIP)
	assert.Error(err, "Expected error for empty host ID")
	assert.EqualError(err, "host ID shouldn't be empty", "Unexpected error message")

	// Mock setup for invalid host ID
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("host ID shouldn't be empty")).Once()

	hostIDTemp = "Host_dummy_1"
	_, err = testConf.client.CreateHostIPPort(ctx, hostIDTemp, testConf.nodeHostIP)
	assert.Error(err, "Expected error for invalid host ID")
	assert.EqualError(err, "host ID shouldn't be empty", "Unexpected error message")

	fmt.Println("Create Host IP Port Test Successful")
}

func TestFindHostIPPortByID(t *testing.T) {
	assert := require.New(t)
	fmt.Println("Begin - Find Host IP Port Test")

	testConf.client.getAPI().(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock setup for valid host IP port retrieval
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/hostIPPort/"+hostIPPortID+"?fields=id,address", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	_, err := testConf.client.FindHostIPPortByID(ctx, hostIPPortID)
	if err != nil {
		t.Fatalf("Find Host IP Port failed: %v", err)
	}

	// Mock setup for invalid host IP port ID
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/hostIPPort/dummy-ip-port-id-1?fields=id,address", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("host IP port not found")).Once()

	// Negative test case: Invalid host IP port ID
	hostIPPortIDTemp := "dummy-ip-port-id-1"
	_, err = testConf.client.FindHostIPPortByID(ctx, hostIPPortIDTemp)
	assert.Error(err, "Expected error for invalid host IP port ID")
	assert.EqualError(err, "host IP port not found", "Unexpected error message")

	fmt.Println("Find Host IP Port Test Successful")
}

func TestCreateHostInitiator(t *testing.T) {
	fmt.Println("Begin - Create Host Initiator Test")

	testConf.client.getAPI().(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Initialize hostID and WWNs
	hostID := "valid_host_id"                            // Replace with a valid host ID
	testConf.wwns = []string{"valid_wwn1", "valid_wwn2"} // Replace with actual valid WWNs

	// Mock setup for valid host IP port retrieval
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/instances/hostIPPort/"+hostIPPortID+"?fields=id,address", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Mock setup for host initiator retrieval
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/types/hostInitiator/instances?fields=id,health,type,initiatorId,isIgnored,parentHost,paths", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(4)

	// Mock setup for host initiator creation
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "POST", "/api/types/hostInitiator/instances", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(len(testConf.wwns) + 1)

	// Mock setup for invalid hostID
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "POST", "/api/types/hostInitiator/instances", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("invalid hostID")).Once()

	if hostID == "" {
		t.Fatalf("hostID should not be empty")
	}

	fmt.Println("WWNs: ", testConf.wwns)
	for _, wwn := range testConf.wwns {
		if wwn == "" {
			t.Fatalf("wwn should not be empty")
		}
		fmt.Printf("Adding new Initiator: %s to host: %s \n", hostName, wwn)
		initiator, err := testConf.client.CreateHostInitiator(ctx, hostID, wwn, api.FCInitiatorType)
		fmt.Println("CreateHostInitiator:", initiator, err)
		if err != nil {
			t.Fatalf("CreateHostInitiator %s Error: %v", wwn, err)
		}
	}

	// Negative case
	hostIDTemp := "host_dummy_1"

	// Add Iqn
	initiator, err := testConf.client.CreateHostInitiator(ctx, hostID, testConf.iqn, api.ISCSCIInitiatorType)
	fmt.Println("CreateHostInitiator:", initiator, err)
	if err != nil {
		t.Fatalf("CreateHostInitiator %s Error: %v", testConf.iqn, err)
	}
	iqnInitiatorID = initiator.HostInitiatorContent.ID

	// Negative test cases
	hostIDTemp = ""
	iqnTemp := ""
	_, err = testConf.client.CreateHostInitiator(ctx, hostIDTemp, testConf.iqn, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator with empty hostID - Negative case failed")
	}

	_, err = testConf.client.CreateHostInitiator(ctx, hostID, iqnTemp, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator with empty iqn - Negative case failed")
	}

	// Test idempotency for parent host check
	hostIDTemp = "host_dummy_1"
	_, err = testConf.client.CreateHostInitiator(ctx, hostIDTemp, testConf.iqn, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator Idempotency with invalid hostID - Negative case failed")
	}

	//@TODO: Check and add positive case to modify parent host
	fmt.Println("Create Host Initiator Test Successful")
}

func TestListHostInitiatorsTest(t *testing.T) {
	fmt.Println("Begin - List Host Initiators Test")
	testConf.client.getAPI().(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()

	// Mock setup for listing host initiators
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", mock.Anything, "GET", "/api/types/hostInitiator/instances?fields=id,health,type,initiatorId,isIgnored,parentHost,paths", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	list, err := testConf.client.ListHostInitiators(ctx)
	fmt.Println("List Host initiators", list, err)
	if err != nil {
		t.Fatalf("ListHostInitiators error: %v", err)
	}

	fmt.Println("List Host Initiators Test Successful")
}

func TestModifyHostInitiator(t *testing.T) {
	fmt.Println("Begin - Modify Host Initiator Test")
	ctx := context.Background()
	_, err := testConf.client.ModifyHostInitiator(ctx, hostID, iqnInitiator)
	assert.Equal(t, errors.New("HostInitiator shouldn't be null"), err)

	hostInitiatorContent := types.HostInitiatorContent{
		ID: "id",
	}
	hostInitiator := types.HostInitiator{
		HostInitiatorContent: hostInitiatorContent,
	}
	_, err = testConf.client.ModifyHostInitiator(ctx, hostID, &hostInitiator)
	if err == nil {
		t.Fatalf("Modify Host initiator with nil initiator - Negative case failed")
	}

	hostIDTemp := "host_dummy_1"
	_, err = testConf.client.ModifyHostInitiator(ctx, hostIDTemp, iqnInitiator)
	if err == nil {
		t.Fatalf("Modify Host initiator with invalid initiator - Negative case failed")
	}

	fmt.Println("Modify Host Initiator Test Successful")
}

func TestModifyHostInitiatorByID(t *testing.T) {
	fmt.Println("Begin - Modify Host Initiator By ID Test")
	ctx := context.Background()

	_, err := testConf.client.ModifyHostInitiatorByID(ctx, "", "")
	if err == nil {
		t.Fatalf("Modify Host initiator with nil initiator - Negative case failed")
	}
	_, err = testConf.client.ModifyHostInitiatorByID(ctx, "hostId", "")
	assert.Equal(t, errors.New("Initiator ID shouldn't be null"), err)

	hostIDTemp := "host_dummy_1"
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.ModifyHostInitiatorByID(ctx, hostIDTemp, "Initiator-123")
	if err != nil {
		t.Fatalf("Modify Host initiator with invalid initiator - Negative case failed")
	}

	fmt.Println("Modify Host Initiator By ID Test Successful")
}

func TestFindHostInitiatorPathByID(t *testing.T) {
	fmt.Println("Begin - Find Initiator Path Test")

	ctx := context.Background()
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	hostInitiatorPath, err := testConf.client.FindHostInitiatorPathByID(ctx, wwnInitiatorPathID)
	if err != nil {
		// Change to log if required for vm execution
		t.Fatalf("Find Host Initiator Path failed: %v", err)
	}
	fcPortID = hostInitiatorPath.HostInitiatorPathContent.FcPortID.ID

	// Negative test cases
	initiatorPathIDTemp := "Host_initiator_path_dummy_1"
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.FindHostInitiatorPathByID(ctx, initiatorPathIDTemp)
	if err != nil {
		t.Fatalf("Find Host Initiator path with invalid Id - Negative case failed")
	}

	fmt.Println("Find Initiator Path Test Successful")
}

func TestFindFcPortByID(t *testing.T) {
	fmt.Println("Begin - Find FC Port Test")
	ctx := context.Background()
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.client.FindFcPortByID(ctx, fcPortID)
	if err != nil {
		// Change to log if required for vm execution
		t.Fatalf("Find FC Port failed: %v", err)
	}

	// Negative test cases
	fcPortIDTemp := "Fc_Port_dummy_1"
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err = testConf.client.FindFcPortByID(ctx, fcPortIDTemp)
	if err != nil {
		t.Fatalf("Find FC Port with invalid Id - Negative case failed")
	}

	fmt.Println("Find FC Port Test Successful")
}

func TestFindTenants(t *testing.T) {
	fmt.Println("Begin - Find Tenants Test")
	ctx := context.Background()
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.client.FindTenants(ctx)
	if err != nil {
		t.Fatalf("Find Tenants failed: %v", err)
	}

	fmt.Println("Find Tenants Test Successful")
}

func TestDeleteHost(t *testing.T) {
	fmt.Println("Begin - Delete Host Test")

	ctx := context.Background()
	hostNameTemp := ""
	err := testConf.client.DeleteHost(ctx, hostNameTemp)
	assert.Equal(t, errors.New("hostname shouldn't be empty"), err)

	hostNameTemp = "dummy-host-1"
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	err = testConf.client.DeleteHost(ctx, hostNameTemp)
	if err != nil {
		t.Fatalf("Delete Host with invalid hostName - Negative case failed")
	}

	fmt.Println("Delete Host Test Successful")
}

func TestFindHostInitiatorByID(t *testing.T) {
	fmt.Println("Begin - Find HostInitiator By ID")
	ctx := context.Background()
	testConf.client.getAPI().(*mocks.Client).On("DoWithHeaders", anyArgs...).Return(nil).Once()
	_, err := testConf.client.FindHostInitiatorByID(ctx, "")
	if err != nil {
		t.Fatalf("Find initiator by empty id - Negative case failed")
	}
}
