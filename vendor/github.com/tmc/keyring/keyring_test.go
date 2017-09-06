package keyring

import (
	"fmt"
	"testing"
)

func TestBasicSetGet(t *testing.T) {
	var (
		pw  string
		err error
	)
	pw, err = Get("keyring-test", "jack")
	if err != nil {
		// ok on initial invokation
		fmt.Println("Get() error:", err)
	}
	err = Set("keyring-test", "jack", "test")
	if err != nil {
		t.Error("Set() error:", err)
	}
	pw, err = Get("keyring-test", "jack")
	if err != nil {
		t.Error("Get() error:", err)
	}

	if pw != "test" {
		fmt.Errorf("expected 'test', got '%s'", pw)
	}
}
