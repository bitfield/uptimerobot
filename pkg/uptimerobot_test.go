package uptimerobot

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshalMonitor(t *testing.T) {
	t.Parallel()
	m := Monitor{
		ID:            777749809,
		FriendlyName:  "Google",
		URL:           "http://www.google.com",
		Type:          MonitorType("HTTP"),
		Port:          80,
		AlertContacts: []string{"3", "5", "7"},
	}
	got, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Error(err)
	}
	want, err := ioutil.ReadFile("testdata/marshal.json")
	if err != nil {
		t.Fatal(err)
	}
	// Convert the actual data and the expected data to maps, for ease of
	// comparison
	wantMap := map[string]interface{}{}
	err = json.Unmarshal(want, &wantMap)
	if err != nil {
		t.Fatal(err)
	}
	gotMap := map[string]interface{}{}
	err = json.Unmarshal(got, &gotMap)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantMap, gotMap) {
		t.Error(cmp.Diff(wantMap, gotMap))
	}
}

func TestUnmarshalMonitor(t *testing.T) {
	t.Parallel()
	want := Monitor{
		ID:           777749809,
		FriendlyName: "Google",
		URL:          "http://www.google.com",
		Type:         MonitorType("HTTP"),
		Port:         80,
	}
	data, err := ioutil.ReadFile("testdata/unmarshal.json")
	if err != nil {
		t.Fatal(err)
	}
	got := Monitor{}
	err = got.UnmarshalJSON(data)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestNewMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("want POST request, got %q", r.Method)
		}
		wantURL := "/v2/newMonitor"
		if r.URL.EscapedPath() != wantURL {
			t.Errorf("want %q, got %q", wantURL, r.URL.EscapedPath())
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		want, err := ioutil.ReadFile("testdata/requestNewMonitor.json")
		if err != nil {
			t.Fatal(err)
		}
		// Convert the received body and the expected body to maps, for
		// ease of comparison
		wantMap := map[string]interface{}{}
		err = json.Unmarshal(want, &wantMap)
		if err != nil {
			t.Fatal(err)
		}
		bodyMap := map[string]interface{}{}
		err = json.Unmarshal(body, &bodyMap)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(wantMap, bodyMap) {
			t.Error(cmp.Diff(wantMap, bodyMap))
		}
		w.WriteHeader(http.StatusOK)
		data, err := os.Open("testdata/newMonitor.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		FriendlyName:  "My test monitor",
		URL:           "http://example.com",
		Type:          MonitorType("HTTP"),
		Port:          80,
		AlertContacts: []string{"3", "5", "7"},
	}
	got, err := client.NewMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != 777810874 {
		t.Errorf("NewMonitor() => ID %d, want 777810874", got.ID)
	}
}

func TestGetAccountDetails(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getAccountDetails.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	got, err := client.GetAccountDetails()
	if err != nil {
		t.Error(err)
	}
	want := Account{
		Email:           "test@domain.com",
		MonitorLimit:    50,
		MonitorInterval: 1,
		UpMonitors:      1,
		DownMonitors:    0,
		PausedMonitors:  2,
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetAlertContacts(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getAlertContacts.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := []AlertContact{
		{
			ID:           "0993765",
			FriendlyName: "John Doe",
			Type:         2,
			Status:       1,
			Value:        "johndoe@gmail.com",
		},
		{
			ID:           "2403924",
			FriendlyName: "My Twitter",
			Type:         3,
			Status:       0,
			Value:        "sampleTwitterAccount",
		},
	}
	got, err := client.GetAlertContacts()
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMonitorByID(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getMonitorByID.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID:           777749809,
		FriendlyName: "Google",
		URL:          "http://www.google.com",
		Type:         MonitorType("HTTP"),
		Port:         80,
	}
	got, err := client.GetMonitorByID(want.ID)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMonitors(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getMonitors.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := []Monitor{
		{
			ID:           777749809,
			FriendlyName: "Google",
			URL:          "http://www.google.com",
			Type:         MonitorType("HTTP"),
			Port:         80,
		},
		{
			ID:           777712827,
			FriendlyName: "My Web Page",
			URL:          "http://mywebpage.com/",
			Type:         MonitorType("HTTP"),
		},
		{
			ID:           777559666,
			FriendlyName: "My FTP Server",
			URL:          "ftp.mywebpage.com",
			Type:         MonitorType("port"),
			SubType:      MonitorSubType("FTP (21)"),
			Port:         21,
		},
		{
			ID:           781397847,
			FriendlyName: "PortTest",
			URL:          "mywebpage.com",
			Type:         MonitorType("port"),
			SubType:      MonitorSubType("Custom Port"),
			Port:         8000,
		},
	}
	got, err := client.GetMonitors()
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMonitorsBySearch(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getMonitorsBySearch.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := []Monitor{
		{
			ID:           777712827,
			FriendlyName: "My Web Page",
			URL:          "http://mywebpage.com/",
			Type:         MonitorType("HTTP"),
		},
	}
	got, err := client.GetMonitorsBySearch("My Web Page")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPauseMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/pauseMonitor.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID: 677810870,
	}
	got, err := client.PauseMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestStartMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/startMonitor.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID: 677810870,
	}
	got, err := client.StartMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestEnsureMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getMonitorsBySearch.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID:           777712827,
		FriendlyName: "My Web Page",
		URL:          "http://mywebpage.com/",
		Type:         MonitorType("HTTP"),
	}
	got, err := client.EnsureMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestDeleteMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/deleteMonitor.json")
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID:   777810874,
		Type: MonitorType("HTTP"),
	}
	got, err := client.DeleteMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRenderMonitor(t *testing.T) {
	t.Parallel()
	input := Monitor{
		ID:            777749809,
		FriendlyName:  "Google",
		URL:           "http://www.google.com",
		Type:          MonitorType("HTTP"),
		Port:          80,
		AlertContacts: []string{"3", "5", "7"},
	}
	wantBytes, err := ioutil.ReadFile("testdata/monitor_template.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(wantBytes)
	got := render(monitorTemplate, input)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRenderAccount(t *testing.T) {
	t.Parallel()
	input := Account{
		Email:           "j.random@example.com",
		MonitorLimit:    300,
		MonitorInterval: 1,
		UpMonitors:      208,
		DownMonitors:    2,
		PausedMonitors:  0,
	}
	wantBytes, err := ioutil.ReadFile("testdata/account_template.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(wantBytes)
	got := render(accountTemplate, input)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFriendlyType(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// cannedResponseServer returns a test TLS server which
func cannedResponseServer(t *testing.T, path string) *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
}
