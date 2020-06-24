package gounity

import (
	"context"
	"fmt"
	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"testing"
	"time"
)

func TestCreateHost(t *testing.T) {
	ctx := context.Background()
	var host *types.Host
	var err error
	err = testConf.hostApi.DeleteHost(ctx, testConf.nodeHostName)
	fmt.Println("Delete Host error (Ignore this error):", err)

	host, err = testConf.hostApi.FindHostByName(ctx, testConf.nodeHostName)
	fmt.Println("Find Host:", host, err)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	if err == HostNotFoundError {
		host, err = testConf.hostApi.CreateHost(ctx, testConf.nodeHostName)
		fmt.Println("Create Host:", host, err)
		if err != nil {
			//t.Fatalf("Create Host failed: %v", err)
		}
	}

	if host != nil {
		//Create Host Ip Port
		hostIpPort, err := testConf.hostApi.CreateHostIpPort(ctx, host.HostContent.ID, testConf.nodeHostIp)
		fmt.Println("CreateHostIpPort:", hostIpPort, err)
		if err != nil {
			t.Fatalf("CreateHostIpPort failed: %v", err)
		}
		fmt.Println("List Host initiators")
		list, err := testConf.hostApi.ListHostInitiators(ctx)
		fmt.Println("List Host initiators", list, err)
		if err != nil {
			t.Fatalf("ListHostInitiators error: %v", err)
		}

		fmt.Println("WWN or Iqns: ", testConf.wwnOrIqns)
		for _, wwnOrIqn := range testConf.wwnOrIqns {
			fmt.Printf("Adding new Initiator: %s to host: %s \n", testConf.nodeHostName, wwnOrIqn)
			initiator, err := testConf.hostApi.CreateHostInitiator(ctx, host.HostContent.ID, wwnOrIqn, api.FCInitiatorType)
			fmt.Println("CreateHostInitiator:", initiator, err)
			if err != nil {
				t.Fatalf("CreateHostInitiator %s Error: %v", wwnOrIqn, err)
			}
		}

		fmt.Println("----------Initiator attached to some other host test -----------")
		now := time.Now()
		nodeHostName := "TestHost-" + now.Format("20060102150405")
		host, err = testConf.hostApi.CreateHost(ctx, nodeHostName)
		fmt.Println("Create Host:", host, err)
		if err != nil {
			t.Fatalf("Create Host failed: %v", err)
		}
		fmt.Println("WWN or Iqns: ", testConf.wwnOrIqns)
		for _, wwnOrIqn := range testConf.wwnOrIqns {
			fmt.Printf("Adding new Initiator: %s to host: %s \n", nodeHostName, wwnOrIqn)
			initiator, err := testConf.hostApi.CreateHostInitiator(ctx, host.HostContent.ID, wwnOrIqn, api.FCInitiatorType)
			fmt.Println("CreateHostInitiator:", initiator, err)
			if err == nil {
				t.Fatalf("CreateHostInitiator %s Error: %v", wwnOrIqn, err)
			}
		}
		err = testConf.hostApi.DeleteHost(ctx, nodeHostName)
		fmt.Println("Delete Host:", err)
		if err != nil {
			t.Fatalf("Delete Host failed: %v", err)
		}

		err = testConf.hostApi.DeleteHost(ctx, testConf.nodeHostName)
		fmt.Println("Delete Host:", err)
		if err != nil {
			t.Fatalf("Delete Host failed: %v", err)
		}
	}
}
