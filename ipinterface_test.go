/*
 Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.

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
