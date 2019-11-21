package gounity

import (
	"fmt"
	"github.com/dell/gounity/payloads"
	"testing"
	"time"
)

func TestCreateHost(t *testing.T) {
	var host *payloads.Host
	var err error
	err = testConf.hostApi.DeleteHost(testConf.nodeHostName)
	fmt.Println("Delete Host error (Ignore this error):", err)

	host, err = testConf.hostApi.CreateHost(testConf.nodeHostName)
	fmt.Println("Create Host:", host, err)
	if err != nil {
		//t.Fatalf("Create Host failed: %v", err)
	}

	host, err = testConf.hostApi.FindHostByName(testConf.nodeHostName)
	fmt.Println("Find Host:", host, err)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	//Create Host Ip Port
	hostIpPort, err := testConf.hostApi.CreateHostIpPort(host.HostContent.ID, testConf.nodeHostIp)
	fmt.Println("CreateHostIpPort:", hostIpPort, err)
	if err != nil {
		t.Fatalf("CreateHostIpPort failed: %v", err)
	}
	fmt.Println("List Host initiators")
	list, err := testConf.hostApi.ListHostInitiators()
	fmt.Println("List Host initiators", list, err)
	if err != nil {
		t.Fatalf("ListHostInitiators error: %v", err)
	}

	fmt.Println("WWN or Iqns: ", testConf.wwnOrIqns)
	for _, wwnOrIqn := range testConf.wwnOrIqns {
		fmt.Printf("Adding new Initiator: %s to host: %s \n", testConf.nodeHostName, wwnOrIqn)
		initiator, err := testConf.hostApi.CreateHostInitiator(host.HostContent.ID, wwnOrIqn)
		fmt.Println("CreateHostInitiator:", initiator, err)
		if err != nil {
			t.Fatalf("CreateHostInitiator %s Error: %v", wwnOrIqn, err)
		}
	}

	fmt.Println("----------Initiator attached to some other host test -----------")
	now := time.Now()
	nodeHostName := "TestHost-" + now.Format("20060102150405")
	host, err = testConf.hostApi.CreateHost(nodeHostName)
	fmt.Println("Create Host:", host, err)
	if err != nil {
		t.Fatalf("Create Host failed: %v", err)
	}
	fmt.Println("WWN or Iqns: ", testConf.wwnOrIqns)
	for _, wwnOrIqn := range testConf.wwnOrIqns {
		fmt.Printf("Adding new Initiator: %s to host: %s \n", nodeHostName, wwnOrIqn)
		initiator, err := testConf.hostApi.CreateHostInitiator(host.HostContent.ID, wwnOrIqn)
		fmt.Println("CreateHostInitiator:", initiator, err)
		if err == nil {
			t.Fatalf("CreateHostInitiator %s Error: %v", wwnOrIqn, err)
		}
	}
	err = testConf.hostApi.DeleteHost(nodeHostName)
	fmt.Println("Delete Host:", err)
	if err != nil {
		t.Fatalf("Delete Host failed: %v", err)
	}

	err = testConf.hostApi.DeleteHost(testConf.nodeHostName)
	fmt.Println("Delete Host:", err)
	if err != nil {
		t.Fatalf("Delete Host failed: %v", err)
	}
}
