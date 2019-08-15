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

