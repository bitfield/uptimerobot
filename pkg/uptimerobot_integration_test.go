// +build integration

package uptimerobot

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getAPIKey(t *testing.T) string {
	key := os.Getenv("UPTIMEROBOT_API_KEY")
	if key == "" {
		t.Fatal("'UPTIMEROBOT_API_KEY' must be set for integration tests")
	}
	return key
}

func exampleMonitor(name string) Monitor {
	return Monitor{
		FriendlyName: name,
		URL:          "http://example.com",
		Type:         MonitorType("HTTP"),
		SubType:      MonitorSubType("HTTP (80)"),
		KeywordType:  0.0,
		Port:         80.0,
	}
}

func TestCreateGetIntegration(t *testing.T) {
	t.Parallel()
	client := New(getAPIKey(t))
	want := exampleMonitor("create_test")
	// client.Debug = os.Stdout
	result, err := client.NewMonitor(want)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteMonitor(result)
	got, err := client.GetMonitorByID(result.ID)
	want.ID = result.ID
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestDeleteIntegration(t *testing.T) {
	t.Parallel()
	client := New(getAPIKey(t))
	toDelete := exampleMonitor("delete_test")
	// client.Debug = os.Stdout
	result, err := client.NewMonitor(toDelete)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.DeleteMonitor(result)
	if err != nil {
		t.Error(err)
	}
	_, err = client.GetMonitorByID(result.ID)
	if err == nil {
		t.Error("want error getting deleted check, but got nil")
	}
}
