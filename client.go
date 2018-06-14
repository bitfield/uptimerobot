package uptimerobot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient represents an http.Client, or a mock equivalent.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents an UptimeRobot client.
type Client struct {
	apiKey string
	http   HTTPClient
}

// Error represents an API error.
type Error map[string]interface{}

// Response represents an API response.
type Response struct {
	Stat    string  `json:"stat"`
	Account Account `json:"account"`
	Error   Error   `json:"error"`
}

// Account represents an UptimeRobot account.
type Account struct {
	Email           string `json:"email"`
	MonitorLimit    int    `json:"monitor_limit"`
	MonitorInterval int    `json:"monitor_interval"`
	UpMonitors      int    `json:"up_monitors"`
	DownMonitors    int    `json:"down_monitors"`
	PausedMonitors  int    `json:"paused_monitors"`
}

// New takes an UptimeRobot API key and returns a Client pointer.
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

// GetAccountDetails returns an Account representing the account details.
func (c *Client) GetAccountDetails() (Account, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.uptimerobot.com",
		Path:   "/v2/getAccountDetails",
	}
	form := url.Values{}
	form.Add("api_key", c.apiKey)
	form.Add("format", "json")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	if err != nil {
		return Account{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return Account{}, err
	}
	defer resp.Body.Close()
	r := struct {
		Stat    string  `json:"stat"`
		Account Account `json:"account"`
		Error
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Account{}, err
	}
	if r.Stat != "ok" {
		e, _ := json.MarshalIndent(r.Error, "", " ")
		return Account{}, fmt.Errorf("API error: %s", e)
	}
	return r.Account, nil
}
