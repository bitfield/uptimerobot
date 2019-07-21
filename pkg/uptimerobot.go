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

// MonitorSubTypes maps a numeric monitor subtype to the name of the monitor subtype.
var MonitorSubTypes = map[float64]string{
	1:  "HTTP (80)",
	2:  "HTTPS (443)",
	3:  "FTP (21)",
	4:  "SMTP (25)",
	5:  "POP3 (110)",
	6:  "IMAP (143)",
	99: "Custom Port",
}

// StatusPause is the status value which sets a monitor to paused status when calling EditMonitor.
var StatusPause = "0"

// StatusResume is the status value which sets a monitor to resumed (unpaused) status when calling EditMonitor.
var StatusResume = "1"

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

// Pagination represents the pagination info of an API response.
type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

// Response represents an API response.
type Response struct {
	Stat          string         `json:"stat"`
	Account       Account        `json:"account"`
	Monitors      []Monitor      `json:"monitors"`
	Monitor       Monitor        `json:"monitor"`
	AlertContacts []AlertContact `json:"alert_contacts"`
	Error         Error          `json:"error,omitempty"`
	Pagination    Pagination     `json:"pagination"`
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
	// keyword_type, sub_type, and port are returned as either an integer
	// (if set) or an empty string (if unset), which Go's JSON library won't
	// parse for integer fields: https://github.com/golang/go/issues/22182
	SubType       interface{} `json:"sub_type"`
	KeywordType   interface{} `json:"keyword_type"`
	Port          interface{} `json:"port"`
	KeywordValue  string      `json:"keyword_value"`
	AlertContacts []string    `json:"alert_contacts"`
}

const monitorTemplate = `ID: {{ .ID }}
Name: {{ .FriendlyName }}
URL: {{ .URL }}
Type: {{ .FriendlyType }}
Subtype: {{ .FriendlySubType }}
Keyword type: {{ .FriendlyKeywordType }}
Keyword value: {{ .KeywordValue }}`

// String returns a pretty-printed version of the monitor.
func (m Monitor) String() string {
	return render(monitorTemplate, m)
}

// FriendlyType returns a human-readable name for the monitor type.
func (m Monitor) FriendlyType() string {
	name, ok := MonitorTypes[m.Type]
	if !ok {
		return fmt.Sprintf("%v", m.Type)
	}
	return name
}

// FriendlySubType returns a human-readable name for the monitor subtype,
// including the port number.
func (m Monitor) FriendlySubType() string {
	subType, ok := m.SubType.(float64)
	if !ok {
		return ""
	}
	name, ok := MonitorSubTypes[subType]
	if !ok {
		return fmt.Sprintf("%v", m.SubType)
	}
	if name == "Custom Port" {
		return fmt.Sprintf("%s (%v)", name, m.Port)
	}
	return name
}

// FriendlyKeywordType returns a human-readable name for the monitor keyword type.
func (m Monitor) FriendlyKeywordType() string {
	switch m.KeywordType {
	case 1.0:
		return "exists"
	case 2.0:
		return "not exists"
	default:
		return fmt.Sprintf("%v", m.KeywordType)
	}
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

// GetMonitorByID takes an int64 representing the ID number of an existing monitor,
// and returns the corresponding monitor with the rest of its metadata, or an
// error if the operation failed.
func (c *Client) GetMonitorByID(ID int64) (Monitor, error) {
	r := Response{}
	p := Params{
		"monitors": fmt.Sprintf("%d", ID),
	}
	if err := c.MakeAPICall("getMonitors", &r, p); err != nil {
		return Monitor{}, err
	}
	if len(r.Monitors) == 0 {
		return Monitor{}, fmt.Errorf("monitor %d not found", ID)
	}
	return r.Monitors[0], nil
}

// GetMonitors returns a slice of Monitors representing the existing monitors.
func (c *Client) GetMonitors() (monitors []Monitor, err error) {
	offset := 0
	limit := 50

	for {
		r := Response{}
		params := Params{
			"offset": strconv.Itoa(offset),
			"limit":  strconv.Itoa(limit),
		}
		if err := c.MakeAPICall("getMonitors", &r, params); err != nil {
			break
		}

		monitors = append(monitors, r.Monitors...)

		if r.Error != nil {
			fmt.Println(fmt.Sprintf("%v", r.Error))
			err = fmt.Errorf(fmt.Sprintf("%v", r.Error))
			break
		}
		offset = r.Pagination.Offset + limit
		total := r.Pagination.Total
		condition := offset+limit < total
		if !condition {
			break
		}
	}

	return monitors, err
}

type GetMonitorsResponse struct {
	Monitors chan Monitor
	Error    chan error
}

func (c *Client) GetMonitorsChan() *GetMonitorsResponse {
	out := GetMonitorsResponse{
		make(chan Monitor),
		make(chan error),
	}
	offset := 0
	limit := 50

	go func() {
		defer close(out.Monitors)
		for {
			r := Response{}
			params := Params{
				"offset": strconv.Itoa(offset),
				"limit":  strconv.Itoa(limit),
			}
			if err := c.MakeAPICall("getMonitors", &r, params); err != nil {
				out.Error <- err
				break
			}

			for _, m := range r.Monitors {
				out.Monitors <- m
			}

			if r.Error != nil {
				fmt.Println(fmt.Sprintf("%v", r.Error))
				out.Error <- fmt.Errorf(fmt.Sprintf("%v", r.Error))
				break
			}
			offset = r.Pagination.Offset + limit
			total := r.Pagination.Total
			condition := offset+limit < total
			if !condition {
				break
			}
		}
	}()

	return &out
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

// EnsureMonitor takes a Monitor and creates a new UptimeRobot monitor with the
// specified details, if a monitor for the same URL does not already exist. It
// returns a Monitor with the ID field set to the ID of the newly created
// monitor or the existing monitor if it already existed, or an error if the
// operation failed.
func (c *Client) EnsureMonitor(m Monitor) (Monitor, error) {
	monitors, err := c.GetMonitorsBySearch(m.URL)
	if err != nil {
		return Monitor{}, err
	}
	if len(monitors) == 0 {
		new, err := c.NewMonitor(m)
		if err != nil {
			return Monitor{}, err
		}
		return new, nil
	}
	return monitors[0], nil
}

// PauseMonitor takes a Monitor with the ID field set, and attempts to set the
// monitor status to paused via the API. It returns a Monitor with the ID field
// set to the ID of the monitor, or an error if the operation failed.
func (c *Client) PauseMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	p := Params{
		"id":     strconv.FormatInt(m.ID, 10),
		"status": StatusPause,
	}
	if err := c.MakeAPICall("editMonitor", &r, p); err != nil {
		return Monitor{}, err
	}
	return r.Monitor, nil
}

// StartMonitor takes a Monitor with the ID field set, and attempts to set the
// monitor status to resumed (unpaused) via the API. It returns a Monitor with
// the ID field set to the ID of the monitor, or an error if the operation
// failed.
func (c *Client) StartMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	p := Params{
		"id":     strconv.FormatInt(m.ID, 10),
		"status": StatusResume,
	}
	if err := c.MakeAPICall("editMonitor", &r, p); err != nil {
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
		requestDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return fmt.Errorf("error dumping HTTP request: %v", err)
		}
		fmt.Fprintln(c.Debug, string(requestDump))
		fmt.Fprintln(c.Debug)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	if c.Debug != nil {
		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return fmt.Errorf("error dumping HTTP response: %v", err)
		}
		fmt.Fprintln(c.Debug, string(responseDump))
		fmt.Fprintln(c.Debug)
	}
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
