package gounity

import (
	"context"
	"fmt"
	"testing"
)

func TestListIPInterfaces(t *testing.T) {
	ctx := context.Background()

	ipInterfaces, err := testConf.ipinterfaceAPI.ListIscsiIPInterfaces(ctx)

	if err != nil {
		t.Fatalf("List Ip Interfaces failed: %v", err)
	}

	for _, ipInterface := range ipInterfaces {
		fmt.Println("Ip Address of interface: ", ipInterface.IPInterfaceContent.IPAddress)
	}
	fmt.Println("List Ip Interfaces success")
}
