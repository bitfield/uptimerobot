package uptimerobot

import (
	"encoding/json"
	"net/http"
	"net/url"
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
	q := u.Query()
	q.Set("api_key", c.apiKey)
	q.Set("format", "json")
	q.Set("noJsonCallback", "1")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
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
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Account{}, err
	}
	return r.Account, nil
}
