package uptimerobot

// TypeHTTP represents an HTTP monitor.
const TypeHTTP = 1

// TypeKeyword represents a keyword monitor.
const TypeKeyword = 2

// TypePing represents a ping monitor.
const TypePing = 3

// TypePort represents a port monitor.
const TypePort = 4

// SubTypeHTTP represents an HTTP monitor subtype.
const SubTypeHTTP = 1

// SubTypeHTTPS represents an HTTPS monitor subtype.
const SubTypeHTTPS = 2

// SubTypeFTP represents an FTP monitor subtype.
const SubTypeFTP = 3

// SubTypeSMTP represents an SMTP monitor subtype.
const SubTypeSMTP = 4

// SubTypePOP3 represents a POP3 monitor subtype.
const SubTypePOP3 = 5

// SubTypeIMAP represents an IMAP monitor subtype.
const SubTypeIMAP = 6

// SubTypeCustomPort represents a custom port monitor subtype.
const SubTypeCustomPort = 99

// KeywordExists represents a keyword check which is critical if the keyword is
// found.
const KeywordExists = 1

// KeywordNotExists represents a keyword check which is critical if the keyword
// is not found.
const KeywordNotExists = 2

// StatusPaused is the status value which sets a monitor to paused status when
// calling EditMonitor.
const StatusPaused = 0

// StatusResumed is the status value which sets a monitor to resumed (unpaused)
// status when calling EditMonitor.
const StatusResumed = 1

// StatusUnknown is the status value indicating that the monitor status is
// currently unknown.
const StatusUnknown = 1

// StatusUp is the status value indicating that the monitor is currently up.
const StatusUp = 2

// StatusMaybeDown is the status value indicating that the monitor may be down,
// but this has not yet been confirmed.
const StatusMaybeDown = 8

// StatusDown is the status value indicating that the monitor is currently down.
const StatusDown = 9

// AlertContactType is a predefined values from uptimerobot. Need for using AlertContact friendly type.
const (
	AlertContactTypePrimaryEmail = 0
	AlertContactTypeEmailToSms   = 1
	AlertContactTypeEmail        = 2
	AlertContactTypeTwitter      = 3
	AlertContactTypeWebHook      = 5
	AlertContactTypePushbullet   = 6
	AlertContactTypeZaiper       = 7
	AlertContactTypeSms          = 8
	AlertContactTypePushover     = 9
	AlertContactTypeHipChat      = 10
	AlertContactTypeSlack        = 11
	AlertContactTypeVoiceCall    = 14
	AlertContactTypeVictorOps    = 15
	AlertContactTypePagerDuty    = 16
	AlertContactTypeTelegram     = 18
	AlertContactTypeTeams        = 20
	AlertContactTypeGoogleChat   = 21
)
