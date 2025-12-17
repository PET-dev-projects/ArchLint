package types

// Severity represents finding severity.
type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
)

// Finding describes a single rule violation or informational message.
type Finding struct {
	RuleID   string         `json:"ruleId"`
	Severity Severity       `json:"severity"`
	Message  string         `json:"message"`
	Path     string         `json:"path"`
	Meta     map[string]any `json:"meta,omitempty"`
}
