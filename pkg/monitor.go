package uptimerobot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Monitor represents an Uptime Robot monitor.
type Monitor struct {
	ID            int64    `json:"id,omitempty"`
	FriendlyName  string   `json:"friendly_name"`
	URL           string   `json:"url"`
	Type          int      `json:"type"`
	SubType       int      `json:"sub_type,omitempty"`
	KeywordType   int      `json:"keyword_type,omitempty"`
	Port          int      `json:"port"`
	KeywordValue  string   `json:"keyword_value,omitempty"`
	AlertContacts []string `json:"alert_contacts,omitempty"`
	Status        Status   `json:"status"`
}

const monitorTemplate = `ID: {{ .ID }}
Name: {{ .FriendlyName }}
URL: {{ .URL -}}
{{ if .Port }}{{ printf "\nPort: %d" .Port }}{{ end -}}
{{ if .Type }}{{ printf "\nType: %s" .FriendlyType }}{{ end -}}
Status: {{ .Status }}
{{ if .SubType }}{{ printf "\nSubtype: %s" .FriendlySubType }}{{ end -}}
{{ if .KeywordType }}{{ printf "\nKeywordType: %s" .FriendlyKeywordType }}{{ end -}}
{{ if .KeywordValue }}{{ printf "\nKeyword: %s" .KeywordValue }}{{ end }}`

// String returns a pretty-printed version of the monitor.
func (m Monitor) String() string {
	return render(monitorTemplate, m)
}

// FriendlyType returns a human-readable name for the monitor type.
func (m Monitor) FriendlyType() string {
	switch m.Type {
	case TypeHTTP:
		return "HTTP"
	case TypeKeyword:
		return "Keyword"
	case TypePing:
		return "Ping"
	case TypePort:
		return "Port"
	default:
		return fmt.Sprintf("%v", m.Type)
	}
}

// FriendlySubType returns a human-readable name for the monitor subtype,
// including the port number.
func (m Monitor) FriendlySubType() string {
	switch m.SubType {
	case SubTypeHTTP:
		return "HTTP (80)"
	case SubTypeHTTPS:
		return "HTTPS (443)"
	case SubTypeFTP:
		return "FTP (21)"
	case SubTypeSMTP:
		return "SMTP (25)"
	case SubTypePOP3:
		return "POP3 (110)"
	case SubTypeIMAP:
		return "IMAP (143)"
	case SubTypeCustomPort:
		return fmt.Sprintf("Custom port (%d)", m.Port)
	default:
		return fmt.Sprintf("%v", m.SubType)
	}
}

// FriendlyKeywordType returns a human-readable name for the monitor keyword type.
func (m Monitor) FriendlyKeywordType() string {
	switch m.KeywordType {
	case KeywordExists:
		return "Exists"
	case KeywordNotExists:
		return "NotExists"
	default:
		return fmt.Sprintf("%v", m.KeywordType)
	}
}

// MarshalJSON converts a Monitor struct into its string JSON representation,
// handling the special encoding of the alert_contacts field.
func (m Monitor) MarshalJSON() ([]byte, error) {
	// Use a temporary type definition to avoid infinite recursion when
	// marshaling
	type MonitorAlias Monitor
	ma := MonitorAlias(m)
	data, err := json.Marshal(ma)
	if err != nil {
		return []byte{}, err
	}
	// Create a temporary map and unmarshal the data into it
	tmp := map[string]interface{}{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return []byte{}, err
	}
	contacts := make([]string, len(m.AlertContacts))
	for i, c := range m.AlertContacts {
		contacts[i] = c + "_0_0"
	}
	tmp["alert_contacts"] = strings.Join(contacts, "-")
	// Marshal the cleaned-up data back to JSON again
	data, err = json.Marshal(tmp)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// UnmarshalJSON converts a JSON monitor representation to a Monitor struct,
// handling the Uptime Robot API's invalid encoding of integer zeros as empty
// strings.
func (m *Monitor) UnmarshalJSON(data []byte) error {
	// We need a custom unmarshaler because keyword_type, sub_type, and port
	// are returned as either a quoted integer (if set) or an empty string
	// (if unset), which Go's JSON library won't parse for integer fields:
	// https://github.com/golang/go/issues/22182
	//
	// Create a temporary map and unmarshal the data into it
	raw := map[string]interface{}{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	// Check and clean up any problematic fields
	fields := []string{
		"sub_type",
		"keyword_type",
		"port",
	}
	for _, f := range fields {
		// If the field is empty string, that means zero.
		if raw[f] == "" {
			raw[f] = 0
		}
		// Otherwise, try to convert it to int.
		if s, ok := raw[f].(string); ok {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			raw[f] = v
		}
	}
	// Marshal the cleaned-up data back to JSON
	data, err = json.Marshal(raw)
	if err != nil {
		return err
	}
	// Use a temporary type definition to avoid infinite recursion when unmarshaling
	type MonitorAlias Monitor
	var ma MonitorAlias
	if err := json.Unmarshal(data, &ma); err != nil {
		return err
	}
	// Finally, convert the temporary type back to a Monitor
	*m = Monitor(ma)
	return nil
}
