package uptimerobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// MonitorTypes maps an integer monitor type to the name of the monitor type.
var MonitorTypes = map[int]string{
	1: "HTTP",
	2: "keyword",
	3: "ping",
	4: "port",
}

// HTTPClient represents an http.Client, or a mock equivalent.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents an UptimeRobot client. If the Debug field is set to
// an io.Writer, then the client will dump API requests to it instead of
// calling the real API.
type Client struct {
	apiKey string
	http   HTTPClient
	Debug  io.Writer
}

// Error represents an API error.
type Error map[string]interface{}

// Params holds optional parameters for API calls.
type Params map[string]string

// Response represents an API response.
type Response struct {
	Stat          string         `json:"stat"`
	Account       Account        `json:"account"`
	Monitors      []Monitor      `json:"monitors"`
	Monitor       Monitor        `json:"monitor"`
	AlertContacts []AlertContact `json:"alert_contacts"`
	Error         Error          `json:"error"`
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

const accountTemplate = `Email: {{ .Email }}
Monitor limit: {{ .MonitorLimit }}
Monitor interval: {{ .MonitorInterval }}
Up monitors: {{ .UpMonitors }}
Down monitors: {{ .DownMonitors }}
Paused monitors: {{ .PausedMonitors }}`

// String returns a pretty-printed version of the account details.
func (a Account) String() string {
	return render(accountTemplate, a)
}

// AlertContact represents an alert contact.
type AlertContact struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendly_name"`
	Type         int    `json:"type"`
	Status       int    `json:"status"`
	Value        string `json:"value"`
}

const alertContactTemplate = `ID: {{ .ID }}
Name: {{ .FriendlyName }}
Type: {{ .Type }}
Status: {{ .Status }}
Value: {{ .Value }}`

// String returns a pretty-printed version of the alert contact.
func (a AlertContact) String() string {
	return render(alertContactTemplate, a)
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
	KeywordType   interface{} `json:"keyword_type"`
	KeywordValue  string      `json:"keyword_value"`
	AlertContacts []string    `json:"alert_contacts"`
}

const monitorTemplate = `ID: {{ .ID }}
Name: {{ .FriendlyName }}
URL: {{ .URL }}
Type: {{ .FriendlyType }}
Subtype: {{ .SubType }}
Keyword type: {{ .KeywordType }}
Keyword value: {{ .KeywordValue }}`

// String returns a pretty-printed version of the monitor.
func (m Monitor) String() string {
	return render(monitorTemplate, m)
}

// FriendlyType returns a human-readable name for the monitor type.
func (m Monitor) FriendlyType() string {
	name, ok := MonitorTypes[m.Type]
	if !ok {
		log.Fatalf("Unknown monitor type %d", m.Type)
	}
	return name
}

// New takes an UptimeRobot API key and returns a Client.
func New(apiKey string) Client {
	return Client{
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

// GetAlertContacts returns all the AlertContacts associated with the account.
func (c *Client) GetAlertContacts() (contacts []AlertContact, err error) {
	r := Response{}
	if err := c.MakeAPICall("getAlertContacts", &r, Params{}); err != nil {
		return contacts, err
	}
	return r.AlertContacts, nil
}

// NewMonitor takes a Monitor and creates a new UptimeRobot monitor with the
// specified details. It returns a Monitor with the ID field set to the ID of
// the newly created monitor, or an error if the operation failed.
func (c *Client) NewMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	p := Params{
		"friendly_name":  m.FriendlyName,
		"url":            m.URL,
		"type":           strconv.Itoa(m.Type),
		"alert_contacts": buildAlertContactList(m.AlertContacts),
	}
	if err := c.MakeAPICall("newMonitor", &r, p); err != nil {
		return Monitor{}, err
	}
	return r.Monitor, nil
}

// DeleteMonitor takes a Monitor with the ID field set, and deletes the
// corresponding monitor. It returns a Monitor with the ID field set to the ID
// of the deleted monitor, or an error if the operation failed.
func (c *Client) DeleteMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	p := Params{
		"id": strconv.FormatInt(m.ID, 10),
	}
	if err := c.MakeAPICall("deleteMonitor", &r, p); err != nil {
		return Monitor{}, err
	}
	return r.Monitor, nil
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
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if c.Debug != nil {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return fmt.Errorf("error dumping HTTP request: %v", err)
		}
		fmt.Fprintln(c.Debug, string(dump))
		return nil
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	if r.Stat != "ok" {
		e, _ := json.MarshalIndent(r.Error, "", " ")
		return fmt.Errorf("API error: %s", e)
	}
	return nil
}

// render takes a template and a data value, and returns the string result of
// executing the template in the context of the value.
func render(templateName string, value interface{}) string {
	var output bytes.Buffer
	tmpl, err := template.New("").Parse(templateName)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(&output, value)
	if err != nil {
		log.Fatal(err)
	}
	return output.String()
}

// MonitorType returns the monitor type number associated with the given type name.
func MonitorType(t string) int {
	for number, name := range MonitorTypes {
		if name == t {
			return number
		}
	}
	log.Fatalf("unknown monitor type %q", t)
	return 0
}

// buildAlertContactList constructs a string in the right format to pass to the
// 'new monitor' API to set alert contacts on a monitor.
func buildAlertContactList(contactIDs []string) string {
	contacts := make([]string, len(contactIDs))
	for i, c := range contactIDs {
		contacts[i] = c + "_0_0"
	}
	return strings.Join(contacts, "-")
}
