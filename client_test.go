package uptimerobot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}
func TestNew(t *testing.T) {
	u := New("dummy")
	if u == nil {
		t.Error("New() => nil, want client object")
	}
	if u.apiKey != "dummy" {
		t.Error("New() did not set API key on client")
	}
}

func fakeAccountDetailsHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"stat": "ok",
			"account": {
				"email": "test@domain.com",
				"monitor_limit": 50,
				"monitor_interval": 1,
				"up_monitors": 1,
				"down_monitors": 0,
				"paused_monitors": 2
			}
		      }`)),
	}, nil
}

func fakeGetMonitorsHandler(req *http.Request) (*http.Response, error) {
	data, err := os.Open("testdata/getMonitors.json")
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       data,
	}, nil
}

func fakeGetMonitorsBySearchHandler(req *http.Request) (*http.Response, error) {
	var f string
	if req.FormValue("search") != "" {
		f = "testdata/getMonitorsBySearch.json"
	} else {
		f = "testdata/getMonitors.json"
	}
	data, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       data,
	}, nil
}

func fakeNewMonitorHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"stat": "ok",
			"monitor": {
				"id": 777810874,
				"status": 1
			}
		      }`)),
	}, nil
}

func TestGetAccountDetails(t *testing.T) {
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeAccountDetailsHandler,
	}
	c.http = &mockClient
	a, err := c.GetAccountDetails()
	if err != nil {
		t.Error(err)
	}
	wantEmail := "test@domain.com"
	if a.Email != wantEmail {
		t.Errorf("GetAccountDetails() => email %q, want %q", a.Email, wantEmail)
	}
}

func TestGetMonitors(t *testing.T) {
	want := []string{"Google", "My Web Page"}
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeGetMonitorsHandler,
	}
	c.http = &mockClient
	monitors, err := c.GetMonitors()
	if err != nil {
		t.Error(err)
	}
	for i, m := range monitors {
		if m.FriendlyName != want[i] {
			t.Errorf("GetMonitors[%d] => %q, want %q", i, m.FriendlyName, want[i])
		}
	}
}

func TestGetMonitorsBySearch(t *testing.T) {
	want := "My Web Page"
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeGetMonitorsBySearchHandler,
	}
	c.http = &mockClient
	monitors, err := c.GetMonitorsBySearch(want)
	if err != nil {
		t.Error(err)
	}
	got := monitors[0].FriendlyName
	if got != want {
		t.Errorf("GetMonitorBySearch(%q) => %q", want, got)
	}
}

func TestNewMonitor(t *testing.T) {
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeNewMonitorHandler,
	}
	c.http = &mockClient
	want := Monitor{
		FriendlyName: "My test monitor",
		URL:          "http://example.com",
	}
	got, err := c.NewMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != 777810874 {
		t.Errorf("NewMonitor() => ID %d, want 777810874", got.ID)
	}
}
