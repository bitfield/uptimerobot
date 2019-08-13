package uptimerobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"text/template"
	"time"
)

// Client represents an UptimeRobot client. If the Debug field is set to
// an io.Writer, then the client will dump API requests to it instead of
// calling the real API.
type Client struct {
	apiKey string
	http   *http.Client
	URL    string
	Debug  io.Writer
}

// New takes an UptimeRobot API key and returns a Client.
func New(apiKey string) Client {
	client := Client{
		apiKey: apiKey,
		URL:    "https://api.uptimerobot.com",
		http:   &http.Client{Timeout: 10 * time.Second},
	}
	if os.Getenv("UPTIMEROBOT_DEBUG") != "" {
		client.Debug = os.Stdout
	}
	return client
}

// Error represents an API error.
type Error map[string]interface{}

// Response represents an API response.
type Response struct {
	Stat          string         `json:"stat"`
	Account       Account        `json:"account"`
	Monitors      []Monitor      `json:"monitors"`
	Monitor       Monitor        `json:"monitor"`
	AlertContacts []AlertContact `json:"alert_contacts"`
	Error         Error          `json:"error"`
}

// GetAccountDetails returns an Account representing the account details.
func (c *Client) GetAccountDetails() (Account, error) {
	r := Response{}
	if err := c.MakeAPICall("getAccountDetails", &r, []byte{}); err != nil {
		return Account{}, err
	}
	return r.Account, nil
}

// GetMonitorByID takes an int64 representing the ID number of an existing monitor,
// and returns the corresponding monitor with the rest of its metadata, or an
// error if the operation failed.
func (c *Client) GetMonitorByID(ID int64) (Monitor, error) {
	r := Response{}
	data := []byte(fmt.Sprintf("{\"monitors\": \"%d\"}", ID))
	if err := c.MakeAPICall("getMonitors", &r, data); err != nil {
		return Monitor{}, err
	}
	if len(r.Monitors) == 0 {
		return Monitor{}, fmt.Errorf("monitor %d not found", ID)
	}
	return r.Monitors[0], nil
}

// GetMonitors returns a slice of Monitors representing the existing monitors.
func (c *Client) GetMonitors() (monitors []Monitor, err error) {
	r := Response{}
	if err := c.MakeAPICall("getMonitors", &r, []byte{}); err != nil {
		return monitors, err
	}
	return r.Monitors, nil
}

// GetMonitorsBySearch returns a slice of Monitors whose FriendlyName or URL
// match the search string.
func (c *Client) GetMonitorsBySearch(s string) (monitors []Monitor, err error) {
	r := Response{}
	data := []byte(`{"search": "` + s + `"}`)
	if err := c.MakeAPICall("getMonitors", &r, data); err != nil {
		return monitors, err
	}
	return r.Monitors, nil
}

// GetAlertContacts returns all the AlertContacts associated with the account.
func (c *Client) GetAlertContacts() (contacts []AlertContact, err error) {
	r := Response{}
	if err := c.MakeAPICall("getAlertContacts", &r, []byte{}); err != nil {
		return contacts, err
	}
	return r.AlertContacts, nil
}

// NewMonitor takes a Monitor and creates a new UptimeRobot monitor with the
// specified details. It returns a Monitor with the ID field set to the ID of
// the newly created monitor, or an error if the operation failed.
func (c *Client) NewMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	data, err := json.Marshal(m)
	if err != nil {
		return Monitor{}, err
	}
	if err := c.MakeAPICall("newMonitor", &r, data); err != nil {
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
	data := []byte(fmt.Sprintf("{\"id\": \"%d\",\"status\": %d}", m.ID, StatusPaused))
	if err := c.MakeAPICall("editMonitor", &r, data); err != nil {
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
	data := []byte(fmt.Sprintf("{\"id\": \"%d\",\"status\": %d}", m.ID, StatusResumed))
	if err := c.MakeAPICall("editMonitor", &r, data); err != nil {
		return Monitor{}, err
	}
	return r.Monitor, nil
}

// DeleteMonitor takes a Monitor with the ID field set, and deletes the
// corresponding monitor. It returns a Monitor with the ID field set to the ID
// of the deleted monitor, or an error if the operation failed.
func (c *Client) DeleteMonitor(m Monitor) (Monitor, error) {
	r := Response{}
	data := []byte(fmt.Sprintf("{\"id\": \"%d\"}", m.ID))
	if err := c.MakeAPICall("deleteMonitor", &r, data); err != nil {
		return Monitor{}, err
	}
	return r.Monitor, nil
}

// MakeAPICall calls the UptimeRobot API with the specified verb and data, and
// stores the returned data in the Response struct.
func (c *Client) MakeAPICall(verb string, r *Response, data []byte) error {
	data, err := decorateRequestData(data, c.apiKey)
	if err != nil {
		return err
	}
	requestURL := c.URL + "/v2/" + verb
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("content-type", "application/json")
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
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	resp.Body.Close()
	respString := string(respBytes)
	resp.Body = ioutil.NopCloser(strings.NewReader(respString))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status %d: %q", resp.StatusCode, respString)
	}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding error for %q: %v", respString, err)
	}
	if r.Stat != "ok" {
		e, _ := json.MarshalIndent(r.Error, "", " ")
		return fmt.Errorf("API error: %s", e)
	}
	return nil
}

// decorateRequestData takes JSON data representing an API request, and adds the
// required 'api_key' and 'format' fields to it.
func decorateRequestData(data []byte, apiKey string) ([]byte, error) {
	// Create a temporary map and unmarshal the data into it
	tmp := map[string]interface{}{}
	// Skip unmarshaling empty data
	if len(data) > 0 {
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return []byte{}, fmt.Errorf("unmarshaling request data: %v", err)
		}
	}
	// Add in the necessary request fields
	tmp["api_key"] = apiKey
	tmp["format"] = "json"
	// Marshal it back into string form
	data, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return []byte{}, fmt.Errorf("remarshaling cleaned-up request data: %v", err)
	}
	return data, nil
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
