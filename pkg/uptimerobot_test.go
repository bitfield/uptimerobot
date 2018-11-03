package uptimerobot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

func badAccountDetailsHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"stat": "false",
			"error": {"message": "Somebody set up us the bomb"}}`)),
	}, nil
}

func TestAPIErrorResponse(t *testing.T) {
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: badAccountDetailsHandler,
	}
	c.http = &mockClient
	_, err := c.GetAccountDetails()
	if err == nil {
		t.Error("API call with error response returned non-nil error")
	}
}

func TestDebugFlag(t *testing.T) {
	c := New("dummy")
	out := &bytes.Buffer{}
	c.Debug = out
	mockClient := MockHTTPClient{
		DoFunc: fakeAccountDetailsHandler,
	}
	c.http = &mockClient
	_, err := c.GetAccountDetails()
	if err != nil {
		t.Error("GetAccountDetails() returned non-nil in debug mode")
	}
	want := "POST /v2/getAccountDetails HTTP/1.1"
	if !strings.Contains(out.String(), want) {
		t.Errorf("GetAccountDetails() debugged %v, want %q ...", out.String(), want)
	}
}

func fakeGetAlertContactsHandler(req *http.Request) (*http.Response, error) {
	data, err := os.Open("testdata/getAlertContacts.json")
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       data,
	}, nil
}
func TestGetAlertContacts(t *testing.T) {
	want := []string{"John Doe", "My Twitter"}
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeGetAlertContactsHandler,
	}
	c.http = &mockClient
	contacts, err := c.GetAlertContacts()
	if err != nil {
		t.Error(err)
	}
	for i, m := range contacts {
		if m.FriendlyName != want[i] {
			t.Errorf("GetAlertContacts[%d] => %q, want %q", i, m.FriendlyName, want[i])
		}
	}
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

func TestGetMonitors(t *testing.T) {
	want := []string{
		"Google",
		"My Web Page",
		"My FTP Server",
		"PortTest",
	}
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

func fakeNewMonitorHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"stat": "ok",
			"monitor": {
				"id": 777810874,
				"status": 1,
				"type": 1
			}
		      }`)),
	}, nil
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
		Type:         MonitorType("HTTP"),
	}
	got, err := c.NewMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != 777810874 {
		t.Errorf("NewMonitor() => ID %d, want 777810874", got.ID)
	}
}

func TestEnsureMonitor(t *testing.T) {
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeGetMonitorsBySearchHandler,
	}
	c.http = &mockClient
	want := Monitor{
		FriendlyName: "My Web Page",
		URL:          "http://mywebpage.com",
		Type:         MonitorType("HTTP"),
	}
	got, err := c.EnsureMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != 777712827 {
		t.Errorf("EnsureMonitor() => ID %d, want 777712827", got.ID)
	}
}

func TestDeleteMonitor(t *testing.T) {
	c := New("dummy")
	mockClient := MockHTTPClient{
		DoFunc: fakeNewMonitorHandler,
	}
	c.http = &mockClient
	want := Monitor{
		ID: 777810874,
	}
	got, err := c.DeleteMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != want.ID {
		t.Errorf("NewMonitor() => ID %d, want %d", got.ID, want.ID)
	}
}

func TestBuildAlertContacts(t *testing.T) {
	contacts := []string{"2353888", "0132759"}
	want := "2353888_0_0-0132759_0_0"
	got := buildAlertContactList(contacts)
	if got != want {
		t.Errorf("buildAlertContacts() => %q, want %q", got, want)
	}
}

func TestRender(t *testing.T) {
	input := Account{
		Email:           "j.random@example.com",
		MonitorLimit:    300,
		MonitorInterval: 1,
		UpMonitors:      208,
		DownMonitors:    2,
		PausedMonitors:  0,
	}
	want := `Email: j.random@example.com
Monitor limit: 300
Monitor interval: 1
Up monitors: 208
Down monitors: 2
Paused monitors: 0`
	got := render(accountTemplate, input)
	if got != want {
		t.Errorf("render(%q) = %q, want %q", input, got, want)
	}
}

func TestFriendlyType(t *testing.T) {
	m := Monitor{
		Type: 1,
	}
	want := "HTTP"
	got := m.FriendlyType()
	if got != want {
		t.Errorf("FriendlyType(1) = %q, want %q", got, want)
	}
}

func TestFriendlySubType(t *testing.T) {
	mHTTPS := Monitor{
		Type: 4,
		// SubType is interface{}, so numeric JSON values will be parsed
		// as float64. Therefore, our test data must be float64.
		SubType: 2.0,
	}
	wantHTTPS := "HTTPS (443)"
	gotHTTPS := mHTTPS.FriendlySubType()
	if gotHTTPS != wantHTTPS {
		t.Errorf("FriendlySubType(HTTPS) = %q, want %q", gotHTTPS, wantHTTPS)
	}
	mCustom := Monitor{
		Type:    4,
		SubType: 99.0,
		Port:    8080,
	}
	wantCustom := "Custom Port (8080)"
	gotCustom := mCustom.FriendlySubType()
	if gotCustom != wantCustom {
		t.Errorf("FriendlySubType(Custom8080) = %q, want %q", gotCustom, wantCustom)
	}
}

func TestFriendlyKeywordType(t *testing.T) {
	m := Monitor{
		Type: 2,
		// KeywordType is interface{}, so numeric JSON values will be parsed
		// as float64. Therefore, our test data must be float64.
		KeywordType: 1.0,
	}
	want := "exists"
	got := m.FriendlyKeywordType()
	if got != want {
		t.Errorf("FriendlyKeywordType() = %q, want %q", got, want)
	}
}

func TestStringMethods(t *testing.T) {
	m := Monitor{Type: 1}
	if m.String() == "" {
		t.Error("m.String() = empty, want non-empty string")
	}
	a := Account{}
	if a.String() == "" {
		t.Error("a.String() = empty, want non-empty string")
	}
	c := AlertContact{}
	if c.String() == "" {
		t.Error("c.String() = empty, want non-empty string")
	}
}
