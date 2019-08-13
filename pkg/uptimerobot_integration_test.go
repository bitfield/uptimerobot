// +build integration

package uptimerobot

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var client Client

func init() {
	key := os.Getenv("UPTIMEROBOT_API_KEY")
	if key == "" {
		log.Fatal("'UPTIMEROBOT_API_KEY' must be set for integration tests")
	}
	client = New(key)
	debug := os.Getenv("UPTIMEROBOT_DEBUG")
	if debug != "" {
		client.Debug = os.Stdout
	}
}

func exampleMonitor(name string) Monitor {
	return Monitor{
		FriendlyName: name,
		URL:          "http://example.com/" + name,
		Type:         TypeHTTP,
		SubType:      SubTypeHTTP,
		Port:         80,
	}
}

func TestCreateGetIntegration(t *testing.T) {
	t.Parallel()
	mon := exampleMonitor("create_test")
	ID, err := client.CreateMonitor(mon)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteMonitor(ID)
	got, err := client.GetMonitor(ID)
	if !cmp.Equal(ID, got.ID) {
		t.Error(cmp.Diff(ID, got.ID))
	}
}

func TestDeleteIntegration(t *testing.T) {
	t.Parallel()
	toDelete := exampleMonitor("delete_test")
	ID, err := client.CreateMonitor(toDelete)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.DeleteMonitor(ID); err != nil {
		t.Error(err)
	}
	_, err = client.GetMonitor(ID)
	if err == nil {
		t.Error("want error getting deleted check, but got nil")
	}
}

func TestAccountDetailsIntegration(t *testing.T) {
	t.Parallel()
	_, err := client.GetAccountDetails()
	if err != nil {
		t.Error(err)
	}
}
