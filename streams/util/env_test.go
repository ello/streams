package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ello/ello-go/streams/util"
)

func ExampleGetEnvWithDefault() {
	//Key doesn't exist
	myval := util.GetEnvWithDefault("mykey", "This is the default")

	fmt.Printf("%v", myval) // Output: This is the default
}

func TestGetEnvWithDefault(t *testing.T) {
	key := "AbC_123"
	key2 := "AAAZZZ___"
	val := "zzddzz023adfg12345"
	_ = os.Setenv(key, val)

	result := util.GetEnvWithDefault(key, "default")
	if result != val {
		t.Error("Default value should not be returned")
	}

	result = util.GetEnvWithDefault(key2, "default")
	if result != "default" {
		t.Error("Default value should be returned")
	}
}

func TestGetEnvIntWithDefault(t *testing.T) {
	key := "AbC_123"
	key2 := "AAAZZZ___"
	val := "1"
	_ = os.Setenv(key, val)

	result := util.GetEnvIntWithDefault(key, 10)
	if result != 1 {
		t.Error("Default value should not be returned")
	}

	result = util.GetEnvIntWithDefault(key2, 10)
	if result != 10 {
		t.Error("Default value should be returned")
	}
}

func TestIsEnvPresent(t *testing.T) {
	key := "AbC_123"
	key2 := "AAAZZZ___"
	val := "zzddzz023adfg12345"
	_ = os.Setenv(key, val)

	result := util.IsEnvPresent(key)
	if !result {
		t.Error("Key is present, result should be true")
	}

	result = util.IsEnvPresent(key2)
	if result {
		t.Error("Key is not present, result should be false")
	}
}
