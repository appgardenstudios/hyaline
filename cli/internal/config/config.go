package config

type Config struct {
	LLM     LLM      `yaml:"llm,omitempty"`
	GitHub  GitHub   `yaml:"github,omitempty"`
	Systems []System `yaml:"systems,omitempty"`
	Rules   []Rule   `yaml:"rules,omitempty"`
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
	ID     string  `yaml:"id,omitempty"`
	Code   []Code  `yaml:"code,omitempty"`
	Docs   []Doc   `yaml:"docs,omitempty"`
	Checks []Check `yaml:"checks,omitempty"`
}

type Extractor string

func (e Extractor) String() string {
	return string(e)
}

func (e Extractor) IsValidCodeExtractor() bool {
	switch e {
	case ExtractorFs, ExtractorGit:
		return true
	default:
		return false
	}
}

func (e Extractor) IsValidDocExtractor() bool {
	switch e {
	case ExtractorFs, ExtractorGit, ExtractorHttp:
		return true
	default:
		return false
	}
}

const (
	ExtractorFs   Extractor = "fs"
	ExtractorGit  Extractor = "git"
	ExtractorHttp Extractor = "http"
)

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

type Code struct {
	ID         string     `yaml:"id,omitempty"`
	Extractor  Extractor  `yaml:"extractor,omitempty"`
	FsOptions  FsOptions  `yaml:"fs,omitempty"`
	GitOptions GitOptions `yaml:"git,omitempty"`
	Include    []string   `yaml:"include,omitempty"`
	Exclude    []string   `yaml:"exclude,omitempty"`
}

type Doc struct {
	ID          string         `yaml:"id,omitempty"`
	Type        DocType        `yaml:"type,omitempty"`
	HTML        DocHTMLOptions `yaml:"html,omitempty"`
	Extractor   Extractor      `yaml:"extractor,omitempty"`
	FsOptions   FsOptions      `yaml:"fs,omitempty"`
	GitOptions  GitOptions     `yaml:"git,omitempty"`
	HttpOptions HttpOptions    `yaml:"http,omitempty"`
	Include     []string       `yaml:"include,omitempty"`
	Exclude     []string       `yaml:"exclude,omitempty"`
	Rules       []string       `yaml:"rules,omitempty"`
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

type Check struct {
	ID          string                 `yaml:"id,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Rule        string                 `yaml:"rule,omitempty"`
	Options     map[string]interface{} `yaml:"options,omitempty"`
}

type Rule struct {
	ID        string         `yaml:"id,omitempty"`
	Documents []RuleDocument `yaml:"documents,omitempty"`
}

type RuleDocument struct {
	Path     string                `yaml:"path,omitempty"`
	Purpose  string                `yaml:"purpose,omitempty"`
	Required bool                  `yaml:"required,omitempty"`
	Ignore   bool                  `yaml:"ignore,omitempty"`
	Sections []RuleDocumentSection `yaml:"sections,omitempty"`
}

type RuleDocumentSection struct {
	ID       string                `yaml:"id,omitempty"`
	Purpose  string                `yaml:"purpose,omitempty"`
	Required bool                  `yaml:"required,omitempty"`
	Ignore   bool                  `yaml:"ignore,omitempty"`
	Sections []RuleDocumentSection `yaml:"sections,omitempty"`
}
