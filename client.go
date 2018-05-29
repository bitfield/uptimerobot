package uptimerobot

import "net/http"

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
		http:   &http.Client{},
	}
}

// GetAccountDetails returns an Account representing the account details.
func (c *Client) GetAccountDetails() (Account, error) {
	return Account{}, nil
}
