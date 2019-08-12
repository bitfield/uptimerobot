package uptimerobot

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
		if r.Method != "POST" {
			t.Errorf("want POST request, got %q", r.Method)
		}
		wantURL := "/v2/newMonitor"
		if r.URL.EscapedPath() != wantURL {
			t.Errorf("want %q, got %q", wantURL, r.URL.EscapedPath())
		}
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
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
		FriendlyName: "My test monitor",
		URL:          "http://example.com",
		Type:         MonitorType("HTTP"),
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
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getAccountDetails.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	a, err := client.GetAccountDetails()
	if err != nil {
		t.Error(err)
	}
	wantEmail := "test@domain.com"
	if a.Email != wantEmail {
		t.Errorf("GetAccountDetails() => email %q, want %q", a.Email, wantEmail)
	}
}

func TestGetAlertContacts(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getAlertContacts.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := []string{"John Doe", "My Twitter"}
	contacts, err := client.GetAlertContacts()
	if err != nil {
		t.Error(err)
	}
	for i, m := range contacts {
		if m.FriendlyName != want[i] {
			t.Errorf("GetAlertContacts[%d] => %q, want %q", i, m.FriendlyName, want[i])
		}
	}
}

func TestGetMonitorByID(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getMonitorByID.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	var want int64 = 777749809
	got, err := client.GetMonitorByID(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != want {
		t.Errorf("GetMonitor() => ID %d, want %d", got.ID, want)
	}
}

func TestGetMonitors(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getMonitors.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := []string{
		"Google",
		"My Web Page",
		"My FTP Server",
		"PortTest",
	}
	monitors, err := client.GetMonitors()
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
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getMonitorsBySearch.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.http = ts.Client()
	client.URL = ts.URL
	want := "My Web Page"
	monitors, err := client.GetMonitorsBySearch(want)
	if err != nil {
		t.Error(err)
	}
	got := monitors[0].FriendlyName
	if got != want {
		t.Errorf("GetMonitorBySearch(%q) => %q", want, got)
	}
}

func TestPauseMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/pauseMonitor.json")
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
		ID: 677810870,
	}
	got, err := client.PauseMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != want.ID {
		t.Errorf("PauseMonitor() => ID %d, want %d", got.ID, want.ID)
	}
}

func TestStartMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/startMonitor.json")
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
		ID: 677810870,
	}
	got, err := client.StartMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != want.ID {
		t.Errorf("StartMonitor() => ID %d, want %d", got.ID, want.ID)
	}
}

func TestEnsureMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/getMonitorsBySearch.json")
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
		FriendlyName: "My Web Page",
		URL:          "http://mywebpage.com",
		Type:         MonitorType("HTTP"),
	}
	got, err := client.EnsureMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != 777712827 {
		t.Errorf("EnsureMonitor() => ID %d, want 777712827", got.ID)
	}
}

func TestDeleteMonitor(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		data, err := os.Open("testdata/deleteMonitor.json")
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
		ID: 777810874,
	}
	got, err := client.DeleteMonitor(want)
	if err != nil {
		t.Error(err)
	}
	if got.ID != want.ID {
		t.Errorf("NewMonitor() => ID %d, want %d", got.ID, want.ID)
	}
}

func TestBuildAlertContacts(t *testing.T) {
	t.Parallel()
	contacts := []string{"2353888", "0132759"}
	want := "2353888_0_0-0132759_0_0"
	got := buildAlertContactList(contacts)
	if got != want {
		t.Errorf("buildAlertContacts() => %q, want %q", got, want)
	}
}

func TestRender(t *testing.T) {
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
	if got != want {
		t.Errorf("render(%q) = %q, want %q", input, got, want)
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
