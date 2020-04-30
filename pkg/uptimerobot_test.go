package uptimerobot

import (
	"encoding/json"
	"fmt"
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
		Type:          TypeHTTP,
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
		Type:         TypeHTTP,
		Port:         80,
		Status:       StatusUnknown,
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

func TestCreate(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		input        Monitor
		requestFile  string
		responseFile string
	}{
		{
			name: "Simple HTTP",
			input: Monitor{
				FriendlyName:  "My test monitor",
				URL:           "http://example.com",
				Type:          TypeHTTP,
				AlertContacts: []string{"3", "5", "7"},
			},
			requestFile:  "testdata/requestNewHttpMonitor.json",
			responseFile: "testdata/newMonitor.json",
		},
		{
			name: "Simple HTTPS",
			input: Monitor{
				FriendlyName:  "My HTTPS test monitor",
				URL:           "https://example.com",
				Type:          TypeHTTP,
				AlertContacts: []string{"3", "5", "7"},
			},
			requestFile:  "testdata/requestNewHttpsMonitor.json",
			responseFile: "testdata/newMonitor.json",
		},
		{
			name: "IMAP port",
			input: Monitor{
				FriendlyName:  "My IMAP port monitor",
				URL:           "example.com",
				Type:          TypePort,
				SubType:       SubTypeIMAP,
				AlertContacts: []string{"3", "5", "7"},
			},
			requestFile:  "testdata/requestNewImapPortMonitor.json",
			responseFile: "testdata/newMonitor.json",
		},
		{
			name: "Custom port",
			input: Monitor{
				FriendlyName:  "My custom port monitor",
				URL:           "example.com",
				Type:          TypePort,
				SubType:       SubTypeCustomPort,
				Port:          8443,
				AlertContacts: []string{"3", "5", "7"},
			},
			requestFile:  "testdata/requestNewCustomPortMonitor.json",
			responseFile: "testdata/newMonitor.json",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := New("dummy")
			// force test coverage of the client's dump functionality
			client.Debug = ioutil.Discard
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
				want, err := ioutil.ReadFile(tc.requestFile)
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
				data, err := os.Open(tc.responseFile)
				if err != nil {
					t.Fatal(err)
				}
				defer data.Close()
				io.Copy(w, data)
			}))
			defer ts.Close()
			client.HTTPClient = ts.Client()
			client.URL = ts.URL
			got, err := client.CreateMonitor(tc.input)
			if err != nil {
				t.Error(err)
			}
			var want int64 = 777810874
			if !cmp.Equal(want, got) {
				t.Error(cmp.Diff(want, got))
			}
		})
	}
}

func TestGetAccountDetails(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getAccountDetails.json")
	defer ts.Close()
	client.HTTPClient = ts.Client()
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

func TestAllAlertContacts(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getAlertContacts.json")
	defer ts.Close()
	client.HTTPClient = ts.Client()
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
	got, err := client.AllAlertContacts()
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
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	want := Monitor{
		ID:           777749809,
		FriendlyName: "Google",
		URL:          "http://www.google.com",
		Type:         TypeHTTP,
		Port:         80,
		Status:       StatusUnknown,
	}
	got, err := client.GetMonitor(want.ID)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMonitorsPages(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyMap := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
			t.Fatal(err)
		}
		var datafile string
		switch bodyMap["offset"] {
		case "0":
			datafile = "testdata/getMonitorsPage1.json"
		case "50":
			datafile = "testdata/getMonitorsPage2.json"
		default:
			t.Fatalf("unexpected offset %s", bodyMap["offset"])
		}
		data, err := os.Open(datafile)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	monitors, err := client.AllMonitors()
	if err != nil {
		t.Error(err)
	}
	if len(monitors) != 100 {
		t.Fatalf("Wanted 100 monitors, but got %d", len(monitors))
	}
	for i, m := range monitors {
		want := fmt.Sprintf("monitor-%d", i+1)
		got := m.FriendlyName
		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestGetMonitors(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/getMonitors.json")
	defer ts.Close()
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	want := []Monitor{
		{
			ID:           777749809,
			FriendlyName: "Google",
			URL:          "http://www.google.com",
			Type:         TypeHTTP,
			Port:         80,
			Status:       StatusUnknown,
		},
		{
			ID:           777712827,
			FriendlyName: "My Web Page",
			URL:          "http://mywebpage.com/",
			Type:         TypeHTTP,
			Status:       StatusUp,
		},
		{
			ID:           777559666,
			FriendlyName: "My FTP Server",
			URL:          "ftp.mywebpage.com",
			Type:         TypePort,
			SubType:      SubTypeFTP,
			Port:         21,
			Status:       StatusUp,
		},
		{
			ID:           781397847,
			FriendlyName: "PortTest",
			URL:          "mywebpage.com",
			Type:         TypePort,
			SubType:      SubTypeCustomPort,
			Port:         8000,
			Status:       StatusUnknown,
		},
	}
	got, err := client.AllMonitors()
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
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	want := []Monitor{
		{
			ID:           777712827,
			FriendlyName: "My Web Page",
			URL:          "http://mywebpage.com/",
			Type:         TypeHTTP,
			Status:       StatusUp,
		},
	}
	got, err := client.SearchMonitors("My Web Page")
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
	client.HTTPClient = ts.Client()
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
	client.HTTPClient = ts.Client()
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

func TestEnsure(t *testing.T) {
	t.Parallel()
	client := New("dummy")
	ts := cannedResponseServer(t, "testdata/ensure.json")
	defer ts.Close()
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	mon := Monitor{
		ID:           777712827,
		FriendlyName: "My Web Page",
		URL:          "http://mywebpage.com/",
		Type:         TypeHTTP,
	}
	// The client will do a SearchMonitors and get a canned response from
	// the test server containing no matches. It will now try to create the
	// monitor, and the test server will just respond with an empty body and
	// OK. The resulting monitor will have an ID of 0.
	want := int64(0)
	got, err := client.EnsureMonitor(mon)
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
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	var want int64 = 777810874
	if err := client.DeleteMonitor(want); err != nil {
		t.Error(err)
	}
}

func TestRenderMonitor(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name     string
		input    Monitor
		wantFile string
	}{
		{
			name: "Simple HTTP",
			input: Monitor{
				ID:            777749809,
				FriendlyName:  "Google",
				URL:           "http://www.google.com",
				Type:          TypeHTTP,
				Port:          0,
				AlertContacts: []string{"3", "5", "7"},
				Status:        StatusUp,
			},
			wantFile: "testdata/monitor_http.txt",
		},
		{
			name: "Keyword exists",
			input: Monitor{
				ID:           777749810,
				FriendlyName: "Google",
				URL:          "http://www.google.com",
				Type:         TypeKeyword,
				KeywordType:  KeywordExists,
				KeywordValue: "bogus",
				Port:         80,
				Status:       StatusMaybeDown,
			},
			wantFile: "testdata/monitor_keyword.txt",
		},
		{
			name: "Keyword not exists",
			input: Monitor{
				ID:           777749811,
				FriendlyName: "Google",
				URL:          "http://www.google.com",
				Type:         TypeKeyword,
				KeywordType:  KeywordNotExists,
				KeywordValue: "bogus",
				Port:         80,
				Status:       StatusUnknown,
			},
			wantFile: "testdata/monitor_keyword_notexists.txt",
		},
		{
			name: "Subtype",
			input: Monitor{
				ID:           777749812,
				FriendlyName: "Google",
				URL:          "http://www.google.com",
				Type:         TypePort,
				SubType:      SubTypeFTP,
				Port:         80,
				Status:       StatusPaused,
			},
			wantFile: "testdata/monitor_subtype.txt",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wantBytes, err := ioutil.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatal(err)
			}
			want := string(wantBytes)
			got := render(monitorTemplate, tc.input)
			if !cmp.Equal(want, got) {
				t.Error(cmp.Diff(want, got))
			}

		})
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
		Type: TypeHTTP,
	}
	want := "HTTP"
	got := m.FriendlyType()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFriendlySubType(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name string
		mon  Monitor
		want string
	}{
		{
			name: "HTTPS",
			mon: Monitor{
				Type:    TypePort,
				SubType: SubTypeHTTPS,
			},
			want: "HTTPS (443)",
		},
		{
			name: "Custom port",
			mon: Monitor{
				Type:    TypePort,
				SubType: SubTypeCustomPort,
				Port:    8080,
			},
			want: "Custom port (8080)",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.mon.FriendlySubType()
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestFriendlyKeywordType(t *testing.T) {
	t.Parallel()
	m := Monitor{
		Type:        TypeKeyword,
		KeywordType: KeywordExists,
	}
	want := "Exists"
	got := m.FriendlyKeywordType()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// cannedResponseServer returns a test TLS server which responds to any request
// with a specified file of canned JSON data.
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
