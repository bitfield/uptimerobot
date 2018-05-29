package uptimerobot

import (
	"net/http"
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

func TestGetAccountDetails(t *testing.T) {
	u := New("dummy")
	u.http = &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// do whatever you want
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, nil
		},
	}

	u.GetAccountDetails()
}
