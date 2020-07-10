package util

import (
	"context"
	"github.com/dell/gounity/api"
	"testing"
)

func TestGetRunIdLogger(t *testing.T) {
	log := GetLogger()
	ctx := context.Background()
	entry := log.WithField("runid", "1111")
	ctx = context.WithValue(ctx, UnityLog, entry)

	logEntry := GetRunIdLogger(ctx)
	logEntry.Info("Hi This is log test1")

	entry = entry.WithField("arrayid", "arr0000")
	ctx = context.WithValue(ctx, UnityLog, entry)
	logEntry = GetRunIdLogger(ctx)
	logEntry.Info("Hi This is log test2")
}

//To validate the ValidateResourceName function
func TestValidateResourceName(t *testing.T) {
	_, err := ValidateResourceName("", api.MaxResourceNameLength)
	if err != NameEmptyError {
		t.Fatalf("%v", err)
	}

	_, err = ValidateResourceName(" ", api.MaxResourceNameLength)
	if err != NameEmptyError {
		t.Fatalf("%v", err)
	}
	_, err = ValidateResourceName("SomeResource123having space", api.MaxResourceNameLength)
	if err != InvalidCharacters {
		t.Fatalf("%v", err)
	}
	_, err = ValidateResourceName("MoreThan40Charactersaaaaaaaaaaaaaa100000000000000000000000000000000000000000000", api.MaxResourceNameLength)
	if err != NameTooLongError {
		t.Fatalf("%v", err)
	}

	_, err = ValidateResourceName("Valid_Name-9:1", api.MaxResourceNameLength)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

//To validate the ValidateResourceName function
func TestValidateDuration(t *testing.T) {
	_, err := ValidateDuration("")
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = ValidateDuration("1:2:3:")
	if err != nil {
		t.Logf("%v", err)
	} else {
		t.Fatalf("%v", err)
	}
	_, err = ValidateDuration("1:2::4")
	if err != nil {
		t.Logf("%v", err)
	} else {
		t.Fatalf("%v", err)
	}
	_, err = ValidateDuration("1:2:3:a")
	if err != nil {
		t.Logf("%v", err)
	} else {
		t.Fatalf("%v", err)
	}

	_, err = ValidateDuration("1:2:3:4")
	if err != nil {
		t.Fatalf("%v", err)
	}

	_, err = ValidateDuration("-1:2:3:4")
	if err != nil {
		t.Logf("%v", err)
	} else {
		t.Fatalf("%v", err)
	}

	_, err = ValidateDuration("1:200:3:4")
	if err != nil {
		t.Logf("%v", err)
	} else {
		t.Fatalf("%v", err)
	}

	sec, _ := ValidateDuration("0:7:1:40")
	if sec != uint64(25300) {
		t.Fatal("invalid time")
	}
}