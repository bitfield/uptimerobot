package uptimerobot

import "fmt"

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

// FriendlyType returns a human-readable name for the alert contact type
func (a AlertContact) FriendlyType() string {
	switch a.Type {
	case AlertContactTypePrimaryEmail:
		return "PrimaryEmail"
	case AlertContactTypeEmail:
		return "Email"
	case AlertContactTypeSms:
		return "Sms"
	case AlertContactTypeVoiceCall:
		return "Voice"
	case AlertContactTypeWebHook:
		return "Webhook"
	case AlertContactTypeEmailToSms:
		return "EmailToSms"
	case AlertContactTypeTwitter:
		return "Twitter"
	case AlertContactTypeTelegram:
		return "Telegram"
	case AlertContactTypeSlack:
		return "Slack"
	case AlertContactTypeTeams:
		return "Teams"
	case AlertContactTypeGoogleChat:
		return "GoogleChat"
	case AlertContactTypeHipChat:
		return "HipChat"
	case AlertContactTypePagerDuty:
		return "Pagerduty"
	case AlertContactTypePushbullet:
		return "Pushbullet"
	case AlertContactTypePushover:
		return "Pushover"
	case AlertContactTypeVictorOps:
		return "VictorOps"
	case AlertContactTypeZaiper:
		return "Zaiper"
	default:
		return fmt.Sprintf("%d", a.Type)
	}
}
