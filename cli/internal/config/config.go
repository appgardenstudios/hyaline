package config

type Config struct {
	LLM     LLM       `yaml:"llm,omitempty"`
	GitHub  GitHub    `yaml:"github,omitempty"`
	Systems []System  `yaml:"systems,omitempty"`
	Rules   []RuleSet `yaml:"rules,omitempty"`
}

func (c *Config) GetSystem(id string) (system System, found bool) {
	for _, s := range c.Systems {
		if s.ID == id {
			return s, true
		}
	}

	return
}

func (c *Config) GetRuleSet(id string) (ruleSet RuleSet, found bool) {
	for _, r := range c.Rules {
		if r.ID == id {
			return r, true
		}
	}

	return
}

type LLM struct {
	Provider LLMProvider `yaml:"provider,omitempty"`
	Model    string      `yaml:"model,omitempty"`
	Key      string      `yaml:"key,omitempty"`
}

type LLMProvider string

func (p LLMProvider) String() string {
	return string(p)
}

func (p LLMProvider) IsValidLLMProvider() bool {
	switch p {
	case LLMProviderAnthropic, LLMProviderTesting:
		return true
	default:
		return false
	}
}

const (
	LLMProviderAnthropic LLMProvider = "anthropic"
	LLMProviderTesting   LLMProvider = "testing"
)

type GitHub struct {
	Token string `yaml:"token,omitempty"`
}

type System struct {
	ID                   string                `yaml:"id,omitempty"`
	CodeSources          []CodeSource          `yaml:"code,omitempty"`
	DocumentationSources []DocumentationSource `yaml:"documentation,omitempty"`
}

func (s *System) GetDocumentationSource(id string) (doc DocumentationSource, found bool) {
	for _, d := range s.DocumentationSources {
		if d.ID == id {
			return d, true
		}
	}

	return
}

type Extractor struct {
	Type    ExtractorType    `yaml:"type,omitempty"`
	Options ExtractorOptions `yaml:"options,omitempty"`
	Include []string         `yaml:"include,omitempty"`
	Exclude []string         `yaml:"exclude,omitempty"`
}

type ExtractorType string

func (e ExtractorType) String() string {
	return string(e)
}

func (e ExtractorType) IsValidCodeExtractor() bool {
	switch e {
	case ExtractorFs, ExtractorGit:
		return true
	default:
		return false
	}
}

func (e ExtractorType) IsValidDocExtractor() bool {
	switch e {
	case ExtractorFs, ExtractorGit, ExtractorHttp:
		return true
	default:
		return false
	}
}

const (
	ExtractorFs   ExtractorType = "fs"
	ExtractorGit  ExtractorType = "git"
	ExtractorHttp ExtractorType = "http"
)

// TODO there should be a better way rather than crunching everything together
type ExtractorOptions struct {
	Path     string             `yaml:"path,omitempty"`
	Repo     string             `yaml:"repo,omitempty"`
	Branch   string             `yaml:"branch,omitempty"`
	Clone    bool               `yaml:"clone,omitempty"`
	HTTPAuth GitHTTPAuthOptions `yaml:"httpAuth,omitempty"`
	SSHAuth  GitSSHAuthOptions  `yaml:"sshAuth,omitempty"`
	BaseURL  string             `yaml:"baseUrl,omitempty"`
	Start    string             `yaml:"start,omitempty"`
	Headers  map[string]string  `yaml:"headers,omitempty"`
}

type FsOptions struct {
	Path string `yaml:"path,omitempty"`
}

type GitOptions struct {
	Repo     string             `yaml:"repo,omitempty"`
	Branch   string             `yaml:"branch,omitempty"`
	Path     string             `yaml:"path,omitempty"`
	Clone    bool               `yaml:"clone,omitempty"`
	HTTPAuth GitHTTPAuthOptions `yaml:"httpAuth,omitempty"`
	SSHAuth  GitSSHAuthOptions  `yaml:"sshAuth,omitempty"`
}

type HttpOptions struct {
	BaseURL string            `yaml:"baseUrl,omitempty"`
	Start   string            `yaml:"start,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type GitHTTPAuthOptions struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type GitSSHAuthOptions struct {
	User     string `yaml:"user,omitempty"`
	PEM      string `yaml:"pem,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type CodeSource struct {
	ID        string    `yaml:"id,omitempty"`
	Extractor Extractor `yaml:"extractor,omitempty"`
}

type DocumentationSource struct {
	ID        string         `yaml:"id,omitempty"`
	Type      DocType        `yaml:"type,omitempty"`
	HTML      DocHTMLOptions `yaml:"html,omitempty"`
	Extractor Extractor      `yaml:"extractor,omitempty"`
	Rules     []string       `yaml:"rules,omitempty"`
}

type DocType string

func (dt DocType) String() string {
	return string(dt)
}

func (dt DocType) IsValid() bool {
	switch dt {
	case DocTypeMarkdown, DocTypeHTML:
		return true
	default:
		return false
	}
}

const (
	DocTypeMarkdown DocType = "md"
	DocTypeHTML     DocType = "html"
)

type DocHTMLOptions struct {
	Selector string `yaml:"selector,omitempty"`
}

type RuleSet struct {
	ID        string         `yaml:"id,omitempty"`
	Documents []RuleDocument `yaml:"documents,omitempty"`
}

type RuleDocument struct {
	Path     string                `yaml:"path,omitempty"`
	Purpose  string                `yaml:"purpose,omitempty"`
	Required bool                  `yaml:"required,omitempty"`
	Ignore   bool                  `yaml:"ignore,omitempty"`
	UpdateIf UpdateIf              `yaml:"updateIf,omitempty"`
	Sections []RuleDocumentSection `yaml:"sections,omitempty"`
}

type RuleDocumentSection struct {
	ID       string                `yaml:"id,omitempty"`
	Purpose  string                `yaml:"purpose,omitempty"`
	Required bool                  `yaml:"required,omitempty"`
	Ignore   bool                  `yaml:"ignore,omitempty"`
	UpdateIf UpdateIf              `yaml:"updateIf,omitempty"`
	Sections []RuleDocumentSection `yaml:"sections,omitempty"`
}

type UpdateIf struct {
	Touched  []string `yaml:"touched,omitempty"`
	Added    []string `yaml:"added,omitempty"`
	Modified []string `yaml:"modified,omitempty"`
	Deleted  []string `yaml:"deleted,omitempty"`
	Renamed  []string `yaml:"renamed,omitempty"`
}
