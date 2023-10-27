/*
 Copyright © 2019 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package util

import (
	"context"
	"fmt"
	"testing"
)

var ctx context.Context

const MaxResourceNameLength = 63

func TestUtils(t *testing.T) {
	getRunIDLoggerTest(t)
	getLoggetTest(t)
	validateResourceNameTest(t)
	validateDurationTest(t)
}

func getRunIDLoggerTest(t *testing.T) {
	fmt.Println("Begin - Get RunId Logger Test")

	log := GetLogger()
	ctx := context.Background()
	entry := log.WithField("runid", "1111")
	ctx = context.WithValue(ctx, UnityLogStruct{UnityLog}, entry)

	logEntry := GetRunIDLogger(ctx)
	logEntry.Info("Hi This is log test1")

	entry = entry.WithField("arrayid", "arr0000")
	ctx = context.WithValue(ctx, UnityLogStruct{UnityLog}, entry)
	logEntry = GetRunIDLogger(ctx)
	logEntry.Info("Hi This is log test2")

	fmt.Println("Get RunId Logger Test Successful")
}

func getLoggetTest(t *testing.T) {
	fmt.Println("Begin - Get Logger Test")

	_ = GetLogger()

	fmt.Println("Get Logger Test Successful")
}

func validateResourceNameTest(t *testing.T) {
	fmt.Println("Begin - Validate Resource Name Test")

	_, err := ValidateResourceName("", MaxResourceNameLength)
	if err != ErrorNameEmpty {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)
	_, err = ValidateResourceName(" ", MaxResourceNameLength)
	if err != ErrorNameEmpty {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)
	_, err = ValidateResourceName("SomeResource123having space", MaxResourceNameLength)
	if err != ErrorInvalidCharacters {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)
	_, err = ValidateResourceName("MoreThan40Charactersaaaaaaaaaaaaaa100000000000000000000000000000000000000000000", MaxResourceNameLength)
	if err != ErrorNameTooLong {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)
	_, err = ValidateResourceName("Valid_Name-9:@*&^%$1", MaxResourceNameLength)
	if err != ErrorInvalidCharacters {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)
	_, err = ValidateResourceName("Valid_Name-9:1", MaxResourceNameLength)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("Validate Resource Name Error: ", err)

	fmt.Println("Validate Resource Name Test Successful")
}

func validateDurationTest(t *testing.T) {
	fmt.Println("Begin - Validate Duration Test")

	_, err := ValidateDuration("")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:2")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1d:23:52:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:23h:52:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:23:52m:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:23:52:50s")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("-1:23:52:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:28:52:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	_, err = ValidateDuration("1:23:70:50")
	if err == nil {
		t.Fatalf("ValidateDuration Negative test failed: %v", err)
	}
	fmt.Println("Error: ", err)
	fmt.Println("Validate Duration Test Successful")
}
