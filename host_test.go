package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
)

var hostName string
var hostID string
var hostIPPortID string
var iqnInitiatorID string
var wwnInitiatorPathID string
var fcPortID string
var iqnInitiator *types.HostInitiator

func TestHost(t *testing.T) {
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	hostName = "Unit-test-host-" + timeStamp
	ctx = context.Background()

	createHostTest(t)
	findHostByNameTest(t)
	createHostIPPortTest(t)
	findHostIPPortByIDTest(t)
	createHostInitiatorTest(t)
	listHostInitiatorsTest(t)
	findHostInitiatorByNameTest(t)
	findHostInitiatorByIDTest(t)
	modifyHostInitiatorTest(t)
	modifyHostInitiatorByIdTest(t)
	findHostInitiatorPathByIDTest(t)
	findFcPortByIDTest(t)
	deleteHostTest(t)
}

func createHostTest(t *testing.T) {

	fmt.Println("Begin - Create Host Test")

	host, err := testConf.hostApi.CreateHost(ctx, hostName)
	if err != nil {
		t.Fatalf("Create Host failed: %v", err)
	}
	hostID = host.HostContent.ID

	//Negative test cases
	hostNameTemp := ""
	_, err = testConf.hostApi.CreateHost(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Create Host with empty hostName - Negative case failed")
	}

	fmt.Println("Create Host Test Successful")
}

func findHostByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Host by name Test")

	_, err := testConf.hostApi.FindHostByName(ctx, hostName)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	//Negative test cases
	hostNameTemp := ""
	_, err = testConf.hostApi.FindHostByName(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Find Host with empty hostName - Negative case failed")
	}

	hostNameTemp = "dummy-host-1"
	_, err = testConf.hostApi.FindHostByName(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Find Host with invalid hostName - Negative case failed")
	}

	fmt.Println("Find Host by name Successful")
}

func createHostIPPortTest(t *testing.T) {

	fmt.Println("Begin - Create Host IP Port Test")

	hostIPPort, err := testConf.hostApi.CreateHostIpPort(ctx, hostID, testConf.nodeHostIp)
	if err != nil {
		t.Fatalf("CreateHostIpPort failed: %v", err)
	}

	hostIPPortID = hostIPPort.HostIpContent.ID
	//Negative test cases
	hostIDTemp := ""
	_, err = testConf.hostApi.CreateHostIpPort(ctx, hostIDTemp, testConf.nodeHostIp)
	if err == nil {
		t.Fatalf("Create Host IP Port with empty hostID - Negative case failed")
	}

	hostIDTemp = "Host_dummy_1"
	_, err = testConf.hostApi.CreateHostIpPort(ctx, hostIDTemp, testConf.nodeHostIp)
	if err == nil {
		t.Fatalf("Create Host IP Port with invalid hostID - Negative case failed")
	}

	fmt.Println("Create Host IP Port Test Successful")
}

func findHostIPPortByIDTest(t *testing.T) {

	fmt.Println("Begin - Find Host IP Port Test")

	_, err := testConf.hostApi.FindHostIpPortById(ctx, hostIPPortID)
	if err != nil {
		t.Fatalf("Find Host IP Port failed: %v", err)
	}

	//Negative test cases
	hostIPPortIDTemp := "dummy-ip-port-id-1"
	_, err = testConf.hostApi.FindHostIpPortById(ctx, hostIPPortIDTemp)
	if err == nil {
		t.Fatalf(" Find Host IP Port with invalid hostID - Negative case failed")
	}

	fmt.Println("Find Host IP Port Test Successful")
}

func createHostInitiatorTest(t *testing.T) {

	fmt.Println("Begin - Create Host Initiator Test")

	fmt.Println("WWNs: ", testConf.wwns)
	for _, wwn := range testConf.wwns {
		fmt.Printf("Adding new Initiator: %s to host: %s \n", hostName, wwn)
		initiator, err := testConf.hostApi.CreateHostInitiator(ctx, hostID, wwn, api.FCInitiatorType)
		fmt.Println("CreateHostInitiator:", initiator, err)
		if err != nil {
			t.Fatalf("CreateHostInitiator %s Error: %v", wwn, err)
		}
	}

	//Negative case
	hostIDTemp := "host_dummy_1"
	_, err := testConf.hostApi.CreateHostInitiator(ctx, hostIDTemp, testConf.iqn, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator Idempotency with invalid hostID - Negative case failed")
	}

	//Add Iqn
	initiator, err := testConf.hostApi.CreateHostInitiator(ctx, hostID, testConf.iqn, api.ISCSCIInitiatorType)
	fmt.Println("CreateHostInitiator:", initiator, err)
	if err != nil {
		t.Fatalf("CreateHostInitiator %s Error: %v", testConf.iqn, err)
	}
	iqnInitiatorID = initiator.HostInitiatorContent.Id

	//Test idempotency for parent host check
	initiator, err = testConf.hostApi.CreateHostInitiator(ctx, hostID, testConf.iqn, api.ISCSCIInitiatorType)
	if err != nil {
		t.Fatalf("CreateHostInitiator %s Error: %v", testConf.iqn, err)
	}

	//Negative test cases
	hostIDTemp = ""
	iqnTemp := ""
	_, err = testConf.hostApi.CreateHostInitiator(ctx, hostIDTemp, testConf.iqn, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator with empty hostID - Negative case failed")
	}

	_, err = testConf.hostApi.CreateHostInitiator(ctx, hostID, iqnTemp, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator with empty iqn - Negative case failed")
	}

	//Test idempotency for parent host check
	hostIDTemp = "host_dummy_1"
	_, err = testConf.hostApi.CreateHostInitiator(ctx, hostIDTemp, testConf.iqn, api.ISCSCIInitiatorType)
	if err == nil {
		t.Fatalf("Create Host Initiator Idempotency with invalid hostID - Negative case failed")
	}

	//@TODO: Cheack and add positive case to modify parent host
	fmt.Println("Create Host Initiator Test Successful")
}

func listHostInitiatorsTest(t *testing.T) {

	fmt.Println("Begin - List Host Initiators Test")
	list, err := testConf.hostApi.ListHostInitiators(ctx)
	fmt.Println("List Host initiators", list, err)
	if err != nil {
		t.Fatalf("ListHostInitiators error: %v", err)
	}

	fmt.Println("List Host Initiators Test Successful")

}

func findHostInitiatorByNameTest(t *testing.T) {

	fmt.Println("Begin - Find Host Initiator by Name Test")

	initiator, err := testConf.hostApi.FindHostInitiatorByName(ctx, testConf.iqn)
	fmt.Println("FindHostInitiatorByName:", initiator, err)
	if err != nil {
		t.Fatalf("FindHostInitiatorByName %s Error: %v", testConf.iqn, err)
	}
	iqnInitiator = initiator

	//Check if call for wwn is required

	//Negative test cases
	iqnTemp := ""
	_, err = testConf.hostApi.FindHostInitiatorByName(ctx, iqnTemp)
	if err == nil {
		t.Fatalf("Find Host Initiator with empty iqn - Negative case failed")
	}

	fmt.Println("Find Host Initiator by Name Test Successful")
}

func findHostInitiatorByIDTest(t *testing.T) {

	fmt.Println("Begin - Find Host Initiator by Id Test")

	//parameterize this
	fcHostName := "lglal016"

	host, err := testConf.hostApi.FindHostByName(ctx, fcHostName)
	if err != nil {
		t.Fatalf("Find Host failed: %v", err)
	}

	for _, fcInitiator := range host.HostContent.FcInitiators {
		initiatorID := fcInitiator.Id
		initiator, err := testConf.hostApi.FindHostInitiatorById(ctx, initiatorID)
		fmt.Println("FindHostInitiatorById:", initiator, err)
		if err != nil {
			t.Fatalf("FindHostInitiatorById %s Error: %v", initiatorID, err)
		}

		if len(initiator.HostInitiatorContent.Paths) > 0 {
			wwnInitiatorPathID = initiator.HostInitiatorContent.Paths[0].Id
			break
		}
	}

	//Negative test cases
	initiatorIDTemp := "dummy-ip-port-id-1"
	_, err = testConf.hostApi.FindHostInitiatorById(ctx, initiatorIDTemp)
	if err == nil {
		t.Fatalf(" Find Host IP Port with invalid initiator ID - Negative case failed")
	}
	fmt.Println("Find Host Initiator by Id Test Successful")
}

func modifyHostInitiatorTest(t *testing.T) {

	fmt.Println("Begin - Modify Host Initiator Test")

	initiator, err := testConf.hostApi.ModifyHostInitiator(ctx, hostID, iqnInitiator)
	fmt.Println("ModifyHostInitiator:", initiator, err)
	if err != nil {
		t.Fatalf("ModifyHostInitiator %s Error: %v", iqnInitiatorID, err)
	}

	_, err = testConf.hostApi.ModifyHostInitiator(ctx, hostID, nil)
	if err == nil {
		t.Fatalf("Modify Host initiator with nil initiator - Negative case failed")
	}

	hostIDTemp := "host_dummy_1"
	_, err = testConf.hostApi.ModifyHostInitiator(ctx, hostIDTemp, iqnInitiator)
	if err == nil {
		t.Fatalf("Modify Host initiator with invalid initiator - Negative case failed")
	}

	fmt.Println("Modify Host Initiator Test Successful")
}

func modifyHostInitiatorByIdTest(t *testing.T) {

	fmt.Println("Begin - Modify Host Initiator By ID Test")

	initiator, err := testConf.hostApi.ModifyHostInitiatorById(ctx, hostID)
	fmt.Println("ModifyHostInitiator:", initiator, err)
	if err != nil {
		t.Fatalf("ModifyHostInitiator %s Error: %v", iqnInitiatorID, err)
	}

	_, err = testConf.hostApi.ModifyHostInitiatorById(ctx, "")
	if err == nil {
		t.Fatalf("Modify Host initiator with nil initiator - Negative case failed")
	}

	hostIDTemp := "host_dummy_1"
	_, err = testConf.hostApi.ModifyHostInitiatorById(ctx, hostIDTemp)
	if err == nil {
		t.Fatalf("Modify Host initiator with invalid initiator - Negative case failed")
	}

	fmt.Println("Modify Host Initiator By ID Test Successful")
}

func findHostInitiatorPathByIDTest(t *testing.T) {

	fmt.Println("Begin - Find Initiator Path Test")

	////initiatorPathID := iqnInitiator.HostInitiatorContent.Paths[0].Id
	hostInitiatorPath, err := testConf.hostApi.FindHostInitiatorPathById(ctx, wwnInitiatorPathID)
	if err != nil {
		//Change to log if required for vm execution
		t.Fatalf("Find Host Initiator Path failed: %v", err)
	}
	fcPortID = hostInitiatorPath.HostInitiatorPathContent.FcPortID.Id

	//Negative test cases
	initiatorPathIDTemp := "Host_initiator_path_dummy_1"
	_, err = testConf.hostApi.FindHostInitiatorPathById(ctx, initiatorPathIDTemp)
	if err == nil {
		t.Fatalf("Find Host Initiator path with invalid Id - Negative case failed")
	}

	fmt.Println("Find Initiator Path Test Successful")
}

func findFcPortByIDTest(t *testing.T) {

	fmt.Println("Begin - Find FC Port Test")

	_, err := testConf.hostApi.FindFcPortById(ctx, fcPortID)
	if err != nil {
		//Change to log if required for vm execution
		t.Fatalf("Find FC Port failed: %v", err)
	}

	//Negative test cases
	fcPortIDTemp := "Fc_Port_dummy_1"
	_, err = testConf.hostApi.FindFcPortById(ctx, fcPortIDTemp)
	if err == nil {
		t.Fatalf("Find FC Port with invalid Id - Negative case failed")
	}

	fmt.Println("Find FC Port Test Successful")
}

func deleteHostTest(t *testing.T) {

	fmt.Println("Begin - Delete Host Test")

	err := testConf.hostApi.DeleteHost(ctx, hostName)
	if err != nil {
		t.Fatalf("Delete Host failed: %v", err)
	}

	hostNameTemp := ""
	err = testConf.hostApi.DeleteHost(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Delete Host with empty hostName - Negative case failed")
	}

	hostNameTemp = "dummy-host-1"
	err = testConf.hostApi.DeleteHost(ctx, hostNameTemp)
	if err == nil {
		t.Fatalf("Delete Host with invalid hostName - Negative case failed")
	}

	fmt.Println("Delete Host Test Successful")
}
