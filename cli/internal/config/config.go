package config

import (
	"fmt"
	"strings"
)

type Config struct {
	LLM     LLM      `yaml:"llm,omitempty"`
	GitHub  GitHub   `yaml:"github,omitempty"`
	Extract *Extract `yaml:"extract,omitempty"`
	Check   *Check   `yaml:"check,omitempty"`
	Audit   *Audit   `yaml:"audit,omitempty"`
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

type Extractor struct {
	Type    CrawlerType    `yaml:"type,omitempty"`
	Options CrawlerOptions `yaml:"options,omitempty"`
	Include []string       `yaml:"include,omitempty"`
	Exclude []string       `yaml:"exclude,omitempty"`
}

type CrawlerType string

func (e CrawlerType) String() string {
	return string(e)
}

func (e CrawlerType) IsValidCodeExtractor() bool {
	switch e {
	case ExtractorTypeFs, ExtractorTypeGit:
		return true
	default:
		return false
	}
}

func (e CrawlerType) IsValidDocExtractor() bool {
	switch e {
	case ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp:
		return true
	default:
		return false
	}
}

func (e CrawlerType) IsValid() bool {
	switch e {
	case ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp:
		return true
	default:
		return false
	}
}

func (e CrawlerType) PossibleValues() string {
	return fmt.Sprintf("%s, %s, %s", ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp)
}

const (
	ExtractorTypeFs   CrawlerType = "fs"
	ExtractorTypeGit  CrawlerType = "git"
	ExtractorTypeHttp CrawlerType = "http"
)

// Note: there should be a better way rather than crunching everything together
type CrawlerOptions struct {
	Path    string            `yaml:"path,omitempty"`
	Repo    string            `yaml:"repo,omitempty"`
	Branch  string            `yaml:"branch,omitempty"`
	Clone   bool              `yaml:"clone,omitempty"`
	Auth    ExtractorAuth     `yaml:"auth,omitempty"`
	BaseURL string            `yaml:"baseUrl,omitempty"`
	Start   string            `yaml:"start,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type ExtractorAuthType string

func (e ExtractorAuthType) String() string {
	return string(e)
}

const (
	ExtractorAuthHTTP ExtractorAuthType = "http"
	ExtractorAuthSSH  ExtractorAuthType = "ssh"
)

type ExtractorAuth struct {
	Type    ExtractorAuthType    `yaml:"type,omitempty"`
	Options ExtractorAuthOptions `yaml:"options,omitempty"`
}

type ExtractorAuthOptions struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	User     string `yaml:"user,omitempty"`
	PEM      string `yaml:"pem,omitempty"`
}

type CodeSource struct {
	ID        string    `yaml:"id,omitempty"`
	Extractor Extractor `yaml:"extractor,omitempty"`
}

type DocumentationSource struct {
	ID               string           `yaml:"id,omitempty"`
	Type             ExtractorType    `yaml:"type,omitempty"`
	Options          ExtractorOptions `yaml:"options,omitempty"`
	Extractor        Extractor        `yaml:"extractor,omitempty"`
	IncludeDocuments []string         `yaml:"includeDocuments,omitempty"`
	Documents        []Document       `yaml:"documents,omitempty"`
}

type ExtractorType string

func (dt ExtractorType) String() string {
	return string(dt)
}

func (dt ExtractorType) IsValid() bool {
	switch dt {
	case DocTypeMarkdown, DocTypeHTML:
		return true
	default:
		return false
	}
}

func (dt ExtractorType) PossibleValues() string {
	return fmt.Sprintf("%s, %s", DocTypeMarkdown, DocTypeHTML)
}

const (
	DocTypeMarkdown ExtractorType = "md"
	DocTypeHTML     ExtractorType = "html"
)

type ExtractorOptions struct {
	Selector                 string `yaml:"selector,omitempty"`
	DisablePurposeExtraction bool   `yaml:"disablePurposeExtraction,omitempty"`
	PurposeKey               string `yaml:"purposeKey,omitempty"`
}

type DocumentSet struct {
	ID        string     `yaml:"id,omitempty"`
	Documents []Document `yaml:"documents,omitempty"`
}

type Document struct {
	Name     string            `yaml:"name,omitempty"`
	Purpose  string            `yaml:"purpose,omitempty"`
	Required bool              `yaml:"required,omitempty"`
	Ignore   bool              `yaml:"ignore,omitempty"`
	UpdateIf UpdateIf          `yaml:"updateIf,omitempty"`
	Sections []DocumentSection `yaml:"sections,omitempty"`
}

type DocumentSection struct {
	Name     string            `yaml:"name,omitempty"`
	Purpose  string            `yaml:"purpose,omitempty"`
	Required bool              `yaml:"required,omitempty"`
	Ignore   bool              `yaml:"ignore,omitempty"`
	UpdateIf UpdateIf          `yaml:"updateIf,omitempty"`
	Sections []DocumentSection `yaml:"sections,omitempty"`
}

type UpdateIf struct {
	Touched  []UpdateIfEntry `yaml:"touched,omitempty"`
	Added    []UpdateIfEntry `yaml:"added,omitempty"`
	Modified []UpdateIfEntry `yaml:"modified,omitempty"`
	Deleted  []UpdateIfEntry `yaml:"deleted,omitempty"`
	Renamed  []UpdateIfEntry `yaml:"renamed,omitempty"`
}

type UpdateIfEntry struct {
	CodeSource string `yaml:"codeID,omitempty"`
	Glob       string `yaml:"glob,omitempty"`
}

type Extract struct {
	Disabled   bool               `yaml:"disabled,omitempty"`
	Source     ExtractSource      `yaml:"source,omitempty"`
	Crawler    ExtractCrawler     `yaml:"crawler,omitempty"`
	Extractors []ExtractExtractor `yaml:"extractors,omitempty"`
	Metadata   []ExtractMetadata  `yaml:"metadata,omitempty"`
}

type ExtractSource struct {
	ID          string `yaml:"id,omitempty"`
	Description string `yaml:"description,omitempty"`
	Root        string `yaml:"root,omitempty"`
}

type ExtractCrawler struct {
	Type    CrawlerType    `yaml:"type,omitempty"`
	Options CrawlerOptions `yaml:"options,omitempty"`
	Include []string       `yaml:"include,omitempty"`
	Exclude []string       `yaml:"exclude,omitempty"`
}

type ExtractExtractor struct {
	Type    ExtractorType    `yaml:"type,omitempty"`
	Options ExtractorOptions `yaml:"options,omitempty"`
	Include []string         `yaml:"include,omitempty"`
	Exclude []string         `yaml:"exclude,omitempty"`
}

type ExtractMetadata struct {
	Document string               `yaml:"document,omitempty"`
	Section  string               `yaml:"section,omitempty"`
	Purpose  string               `yaml:"purpose,omitempty"`
	Tags     []ExtractMetadataTag `yaml:"tags,omitempty"`
}

type ExtractMetadataTag struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Check struct {
	Disabled      bool               `yaml:"disabled,omitempty"`
	Code          CheckCode          `yaml:"code,omitempty"`
	Documentation CheckDocumentation `yaml:"documentation,omitempty"`
	Options       CheckOptions       `yaml:"options,omitempty"`
}

type CheckCode struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type CheckDocumentation struct {
	Include []DocumentationFilter `yaml:"include,omitempty"`
	Exclude []DocumentationFilter `yaml:"exclude,omitempty"`
}

type CheckOptions struct {
	DetectDocumentationUpdates CheckOptionsDetectDocumentationUpdates `yaml:"detectDocumentationUpdates,omitempty"`
	UpdateIf                   CheckOptionsUpdateIf                   `yaml:"updateIf,omitempty"`
}

type CheckOptionsDetectDocumentationUpdates struct {
	Source string `yaml:"source,omitempty"`
}

type CheckOptionsUpdateIf struct {
	Touched  []CheckOptionsUpdateIfEntry `yaml:"touched,omitempty"`
	Added    []CheckOptionsUpdateIfEntry `yaml:"added,omitempty"`
	Modified []CheckOptionsUpdateIfEntry `yaml:"modified,omitempty"`
	Deleted  []CheckOptionsUpdateIfEntry `yaml:"deleted,omitempty"`
	Renamed  []CheckOptionsUpdateIfEntry `yaml:"renamed,omitempty"`
}

type CheckOptionsUpdateIfEntry struct {
	Code          CheckCodeFilter     `yaml:"code,omitempty"`
	Documentation DocumentationFilter `yaml:"documentation,omitempty"`
}

type CheckCodeFilter struct {
	Path string `yaml:"path,omitempty"`
}

type DocumentationFilter struct {
	Source   string                   `yaml:"source,omitempty"`
	Document string                   `yaml:"document,omitempty"`
	Section  string                   `yaml:"section,omitempty"`
	URI      string                   `yaml:"uri,omitempty"`
	Tags     []DocumentationFilterTag `yaml:"tags,omitempty"`
}

func (filter *DocumentationFilter) GetParts() (source string, document string, section string) {
	if filter.URI != "" {
		var remainder string
		source, remainder, _ = strings.Cut(strings.TrimPrefix(filter.URI, "document://"), "/")
		document, section, _ = strings.Cut(remainder, "#")
	} else {
		source = filter.Source
		document = filter.Document
		section = filter.Section
	}

	return
}

type DocumentationFilterTag struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Audit struct {
	Disabled bool        `yaml:"disabled,omitempty"`
	Rules    []AuditRule `yaml:"rules,omitempty"`
}

type AuditRule struct {
	ID            string                `yaml:"id,omitempty"`
	Description   string                `yaml:"description,omitempty"`
	Documentation []DocumentationFilter `yaml:"documentation,omitempty"`
	Ignore        []DocumentationFilter `yaml:"ignore,omitempty"`
	Checks        AuditChecks           `yaml:"checks,omitempty"`
}

type AuditChecks struct {
	Content AuditContentChecks `yaml:"content,omitempty"`
	Purpose AuditPurposeChecks `yaml:"purpose,omitempty"`
	Tags    AuditTagsChecks    `yaml:"tags,omitempty"`
}

type AuditContentChecks struct {
	Exists         bool   `yaml:"exists,omitempty"`
	MinLength      int    `yaml:"min-length,omitempty"`
	MatchesRegex   string `yaml:"matches-regex,omitempty"`
	MatchesPrompt  string `yaml:"matches-prompt,omitempty"`
	MatchesPurpose bool   `yaml:"matches-purpose,omitempty"`
}

type AuditPurposeChecks struct {
	Exists bool `yaml:"exists,omitempty"`
}

type AuditTagsChecks struct {
	Contains []DocumentationFilterTag `yaml:"contains,omitempty"`
}
