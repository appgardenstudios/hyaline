package rule

type Result struct {
	System      string      `json:"system"`
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Rule        string      `json:"rule"`
	Options     interface{} `json:"options"`
	Pass        bool        `json:"pass"`
	Severity    string      `json:"severity,omitempty"`
	Message     string      `json:"message,omitempty"`
}
