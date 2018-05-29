package uptimerobot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	called bool
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.called = true
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

func TestGetAccountDetails(t *testing.T) {
	c := New("dummy")
	want := "/getAccountDetails"
	mockClient := MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {

			u := req.URL
			if u.Path != want {
				t.Errorf("GetAccountDetails called %q, want %q", u.Path, want)
			}
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			}, nil
		},
	}
	c.http = &mockClient
	c.GetAccountDetails()
	if !mockClient.called {
		t.Error("GetAccountDetails didn't make HTTP request")
	}
}
