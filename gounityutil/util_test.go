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

package gounityutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

const MaxResourceNameLength = 63

func TestUtils(t *testing.T) {
	getLoggerTest(t)
	getRunIDLoggerTest(t)
	validateResourceNameTest(t)
	validateDurationTest(t)
	getSecuredCipherSuitesTest(t)
}

func getRunIDLoggerTest(t *testing.T) {
	fmt.Println("Begin - Get RunId Logger Test")

	// Scenario 1: rlog is of type *logrus.Entry and has data
	log := GetLogger()
	ctx := context.Background()
	entry := log.WithField("runid", "1111")
	ctx = context.WithValue(ctx, UnityLog, entry)

	// This should cover the line: entry = rlog.(*logrus.Entry)
	logEntry := GetRunIDLogger(ctx)
	logEntry.Info("Hi This is log test1")

	if len(logEntry.Data) == 0 {
		t.Error("Expected logEntry data to have fields, but it was empty")
	}

	// Scenario 2: rlog is of type *logrus.Entry but has no data
	emptyEntry := &logrus.Entry{}
	ctx = context.WithValue(ctx, UnityLog, emptyEntry)

	// This should skip the line and use log.WithContext(ctx)
	logEntry = GetRunIDLogger(ctx)
	logEntry.Info("Hi This is log test2")

	if len(logEntry.Data) != 0 {
		t.Error("Expected logEntry data to be empty, but it had fields")
	}

	// Additional Scenario: rlog is not of type *logrus.Entry
	ctx = context.WithValue(ctx, UnityLog, "invalid type")

	// This should go directly to log.WithContext(ctx)
	logEntry = GetRunIDLogger(ctx)
	logEntry.Info("Hi This is log test3")

	fmt.Println("Get RunId Logger Test Successful")
}

func getLoggerTest(_ *testing.T) {
	fmt.Println("Begin - Get Logger Test")
	// debug flag needs to be true to hit a test condition, reset it after
	formerDebug := Debug
	Debug = true
	_ = GetLogger()
	Debug = formerDebug
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

func getSecuredCipherSuitesTest(t *testing.T) {
	fmt.Println("Begin - Get Secured Cipher Suites Test")

	suites := GetSecuredCipherSuites()
	if len(suites) == 0 {
		t.Fatalf("No secured cipher suites found")
	}

	// Check if all returned suites are valid TLS cipher suites
	for _, suite := range suites {
		found := false
		for _, v := range tls.CipherSuites() {
			if suite == v.ID {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Invalid cipher suite ID found: %d", suite)
		}
	}

	fmt.Println("Get Secured Cipher Suites Test Successful")
}

func TestChangeLogLevel(t *testing.T) {
	// Ensure singletonLog is initialized before any tests
	GetLogger()

	t.Run("Debug level case", func(t *testing.T) {
		ChangeLogLevel("debug")
		if singletonLog.Level != logrus.DebugLevel {
			t.Errorf("expected DebugLevel, got %v", singletonLog.Level)
		}
	})

	t.Run("Warn level case", func(t *testing.T) {
		ChangeLogLevel("warn")
		if singletonLog.Level != logrus.WarnLevel {
			t.Errorf("expected WarnLevel, got %v", singletonLog.Level)
		}
	})

	t.Run("Warning level case", func(t *testing.T) {
		ChangeLogLevel("warning")
		if singletonLog.Level != logrus.WarnLevel {
			t.Errorf("expected WarnLevel, got %v", singletonLog.Level)
		}
	})

	t.Run("Error level case", func(t *testing.T) {
		ChangeLogLevel("error")
		if singletonLog.Level != logrus.ErrorLevel {
			t.Errorf("expected ErrorLevel, got %v", singletonLog.Level)
		}
	})

	t.Run("Info level case", func(t *testing.T) {
		ChangeLogLevel("info")
		if singletonLog.Level != logrus.InfoLevel {
			t.Errorf("expected InfoLevel, got %v", singletonLog.Level)
		}
	})

	t.Run("Default level case with unknown input", func(t *testing.T) {
		ChangeLogLevel("unknown")
		if singletonLog.Level != logrus.InfoLevel {
			t.Errorf("expected InfoLevel, got %v", singletonLog.Level)
		}
	})
}
