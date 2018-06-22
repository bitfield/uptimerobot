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

// Params holds optional parameters for API calls
type Params map[string]string

// Response represents an API response.
type Response struct {
	Stat     string    `json:"stat"`
	Account  Account   `json:"account"`
	Monitors []Monitor `json:"monitors"`
	Error    Error     `json:"error"`
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

// Monitor represents an UptimeRobot monitor.
type Monitor struct {
	ID           int64  `json:"id"`
	FriendlyName string `json:"friendly_name"`
	URL          string `json:"url"`
	Type         int    `json:"type"`
	SubType      string `json:"sub_type"`
	// keyword_type is returned as either an integer or an empty string,
	// which Go doesn't allow: https://github.com/golang/go/issues/22182
	KeywordType  interface{} `json:"keyword_type"`
	KeywordValue string      `json:"keyword_value"`
}

func (m Monitor) String() string {
	return fmt.Sprintf("%s: %s", m.FriendlyName, m.URL)
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
	r := Response{}
	if err := c.MakeAPICall("getAccountDetails", &r, Params{}); err != nil {
		return Account{}, err
	}
	return r.Account, nil
}

// GetMonitors returns a slice of Monitors representing the existing monitors.
func (c *Client) GetMonitors() (monitors []Monitor, err error) {
	r := Response{}
	if err := c.MakeAPICall("getMonitors", &r, Params{}); err != nil {
		return monitors, err
	}
	return r.Monitors, nil
}

// GetMonitorsBySearch returns a slice of Monitors whose FriendlyName or URL
// match the search string.
func (c *Client) GetMonitorsBySearch(s string) (monitors []Monitor, err error) {
	r := Response{}
	p := Params{
		"search": s,
	}
	if err := c.MakeAPICall("getMonitors", &r, p); err != nil {
		return monitors, err
	}
	return r.Monitors, nil
}

// MakeAPICall calls the UptimeRobot API with the specified verb and stores the
// returned data in the Response struct.
func (c *Client) MakeAPICall(verb string, r *Response, params Params) error {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.uptimerobot.com",
		Path:   "/v2/" + verb,
	}
	form := url.Values{}
	form.Add("api_key", c.apiKey)
	form.Add("format", "json")
	for k, v := range params {
		form.Add(k, v)
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Stat != "ok" {
		e, _ := json.MarshalIndent(r.Error, "", " ")
		return fmt.Errorf("API error: %s", e)
	}
	return nil
}
