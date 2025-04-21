package config

type Config struct {
	LLM     LLM      `yaml:"llm"`
	GitHub  GitHub   `yaml:"github"`
	Systems []System `yaml:"systems"`
}

type LLM struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Key      string `yaml:"key"`
}

type GitHub struct {
	Token string `yaml:"token"`
}

type System struct {
	ID     string  `yaml:"id"`
	Code   []Code  `yaml:"code"`
	Docs   []Doc   `yaml:"docs"`
	Checks []Check `yaml:"checks"`
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
	Path string `yaml:"path"`
}

type GitOptions struct {
	Repo     string             `yaml:"repo"`
	Branch   string             `yaml:"branch"`
	Path     string             `yaml:"path"`
	Clone    bool               `yaml:"clone"`
	HTTPAuth GitHTTPAuthOptions `yaml:"httpAuth"`
	SSHAuth  GitSSHAuthOptions  `yaml:"sshAuth"`
}

type HttpOptions struct {
	BaseURL string            `yaml:"baseUrl"`
	Start   string            `yaml:"start"`
	Headers map[string]string `yaml:"headers"`
}

type GitHTTPAuthOptions struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type GitSSHAuthOptions struct {
	User     string `yaml:"user"`
	PEM      string `yaml:"pem"`
	Password string `yaml:"password"`
}

type Code struct {
	ID         string     `yaml:"id"`
	Extractor  Extractor  `yaml:"extractor"`
	FsOptions  FsOptions  `yaml:"fs"`
	GitOptions GitOptions `yaml:"git"`
	Include    []string   `yaml:"include"`
	Exclude    []string   `yaml:"exclude"`
}

type Doc struct {
	ID          string         `yaml:"id"`
	Type        DocType        `yaml:"type"`
	HTML        DocHTMLOptions `yaml:"html"`
	Extractor   Extractor      `yaml:"extractor"`
	FsOptions   FsOptions      `yaml:"fs"`
	GitOptions  GitOptions     `yaml:"git"`
	HttpOptions HttpOptions    `yaml:"http"`
	Include     []string       `yaml:"include"`
	Exclude     []string       `yaml:"exclude"`
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
	Selector string `yaml:"selector"`
}

type Check struct {
	ID          string                 `yaml:"id"`
	Description string                 `yaml:"description"`
	Rule        string                 `yaml:"rule"`
	Options     map[string]interface{} `yaml:"options"`
}
