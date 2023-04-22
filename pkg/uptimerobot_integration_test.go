//go:build integration
// +build integration

package uptimerobot

import (
	"log"
	"os"
	"testing"
	"time"

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

func TestIntegration(t *testing.T) {
	t.Parallel()
	ID, err := client.CreateMonitor(exampleMonitor("create_test"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetMonitor(ID)
	if !cmp.Equal(ID, got.ID) {
		t.Error(cmp.Diff(ID, got.ID))
	}
	time.Sleep(10 * time.Second) // avoid rate limit
	if err = client.DeleteMonitor(ID); err != nil {
		t.Error(err)
	}
	time.Sleep(10 * time.Second) // avoid rate limit
	_, err = client.GetAccountDetails()
	if err != nil {
		t.Error(err)
	}
}
