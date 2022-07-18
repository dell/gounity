package gounity

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var rsName string
var rsysName string
var rsID string

func TestReplication(t *testing.T) {
	ctx = context.Background()

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	rsName = "Unit-test-cg-" + timeStamp
}

func createRS(t *testing.T) {
	fmt.Println("Begin - Create RS Test")

	replicationSessionName := "newRStest"
	srcResourceId := "res_3585"
	dstResourceId := "res_90957"
	remoteSystemName := "APM00213404194"
	maxTimeOutOfSync := "0"

	rs, err := testConf.rAPI.CreateReplicationSession(ctx, replicationSessionName, srcResourceId, dstResourceId, remoteSystemName, maxTimeOutOfSync)
	fmt.Println("Create RS :", prettyPrintJSON(rs), err)
	if err != nil {
		t.Fatalf("Create RS failed: %v", err)
	}

	fmt.Println("Create RS Test - Successful")
}

func deleteRS(t *testing.T) {
	fmt.Println("Begin - Delete RS Test")

	replicationSessionId := "81604379415_APM00213404195_0000_81604379592_APM00213404194_0000"

	err := testConf.rAPI.DeleteReplicationSession(ctx, replicationSessionId)
	if err != nil {
		t.Fatalf("Delete RS failed: %v", err)
	}

	fmt.Println("Delete RS Test - Successful")
}

func findRSBySrcResourceID(t *testing.T, resourceTestID string) {

	fmt.Println("Begin - Find RS By ResourceID Test")

	rs, err := testConf.rAPI.FindReplicationSessionBySrcResourceID(ctx, resourceTestID)
	fmt.Println("RS By ResourceID:", prettyPrintJSON(rs), err)
	if err != nil {
		t.Fatalf("Find RS By ResourceID failed: %v", err)
	}
	rsID = rs.ReplicationSessionContent.ReplicationSessionId

	//Negative cases
	emptyId := ""
	rs, err = testConf.rAPI.FindReplicationSessionBySrcResourceID(ctx, emptyId)
	if err == nil {
		t.Fatalf("Find RS By ResourceID with empty Id case failed: %v", err)
	}

	fmt.Println("Find RS By ResourceID Test - Successful")
}

func findRSByID(t *testing.T) {

	fmt.Println("Begin - Find RS By ID Test")

	rs, err := testConf.rAPI.FindReplicationSessionById(ctx, rsID)
	fmt.Println("RS By ID:", prettyPrintJSON(rs), err)
	if err != nil {
		t.Fatalf("Find RS By ID failed: %v", err)
	}

	//Negative cases
	emptyId := ""
	rs, err = testConf.rAPI.FindReplicationSessionById(ctx, emptyId)
	if err == nil {
		t.Fatalf("Find RS By ID with empty Id case failed: %v", err)
	}

	fmt.Println("Find RS By ID Test - Successful")
}

func findRSysByName(t *testing.T) {

	fmt.Println("Begin - Find RSys By Name Test")

	rs, err := testConf.rAPI.FindRemoteSystemByName(ctx, rsysName)
	fmt.Println("RSys By Name:", prettyPrintJSON(rs), err)
	if err != nil {
		t.Fatalf("Find RSys By Name failed: %v", err)
	}

	//Negative cases
	emptyName := ""
	rs, err = testConf.rAPI.FindRemoteSystemByName(ctx, emptyName)
	if err == nil {
		t.Fatalf("Find RSys By Name with empty Name case failed: %v", err)
	}

	fmt.Println("Find RSys By Name Test - Successful")
}

func findRSessionByName(t *testing.T) {

	fmt.Println("Begin - Find RSession By Name Test")

	rs, err := testConf.rAPI.FindReplicationSessionByName(ctx, rsName)
	fmt.Println("RSession By Name:", prettyPrintJSON(rs), err)
	if err != nil {
		t.Fatalf("Find RSession By Name failed: %v", err)
	}

	//Negative cases
	emptyName := ""
	rs, err = testConf.rAPI.FindReplicationSessionByName(ctx, emptyName)
	if err == nil {
		t.Fatalf("Find RSession By Name with empty Name case failed: %v", err)
	}

	fmt.Println("Find RSession By Name Test - Successful")
}
