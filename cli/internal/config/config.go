package config

import (
	"fmt"
	"strings"
)

type Config struct {
	LLM             LLM           `yaml:"llm,omitempty"`
	GitHub          GitHub        `yaml:"github,omitempty"`
	Systems         []System      `yaml:"systems,omitempty"`
	CommonDocuments []DocumentSet `yaml:"commonDocuments,omitempty"`
	Extract         *Extract      `yaml:"extract,omitempty"`
	Check           *Check        `yaml:"check,omitempty"`
}

func (c *Config) GetSystem(id string) (system System, found bool) {
	for _, s := range c.Systems {
		if s.ID == id {
			return s, true
		}
	}

	return
}

func (c *Config) GetCommonDocumentSet(id string) (documentSet DocumentSet, found bool) {
	for _, s := range c.CommonDocuments {
		if s.ID == id {
			return s, true
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
	case ExtractorTypeFs, ExtractorTypeGit:
		return true
	default:
		return false
	}
}

func (e ExtractorType) IsValidDocExtractor() bool {
	switch e {
	case ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp:
		return true
	default:
		return false
	}
}

func (e ExtractorType) IsValid() bool {
	switch e {
	case ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp:
		return true
	default:
		return false
	}
}

func (e ExtractorType) PossibleValues() string {
	return fmt.Sprintf("%s, %s, %s", ExtractorTypeFs, ExtractorTypeGit, ExtractorTypeHttp)
}

const (
	ExtractorTypeFs   ExtractorType = "fs"
	ExtractorTypeGit  ExtractorType = "git"
	ExtractorTypeHttp ExtractorType = "http"
)

// Note: there should be a better way rather than crunching everything together
type ExtractorOptions struct {
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
	ID               string                     `yaml:"id,omitempty"`
	Type             DocType                    `yaml:"type,omitempty"`
	Options          DocumentationSourceOptions `yaml:"options,omitempty"`
	Extractor        Extractor                  `yaml:"extractor,omitempty"`
	IncludeDocuments []string                   `yaml:"includeDocuments,omitempty"`
	Documents        []Document                 `yaml:"documents,omitempty"`
}

func (d *DocumentationSource) GetDocuments(c *Config) (documents []Document) {
	// create a map of added documents so we can see what we have already added
	documentMap := map[string]struct{}{}

	// Add all documents from our documentation source first
	for _, document := range d.Documents {
		_, found := documentMap[document.Name]
		if !found {
			documentMap[document.Name] = struct{}{}
			documents = append(documents, document)
		}

	}

	// Add documents from our common documents so that documents in later sets take priority of those
	// in earlier sets as defined by the order of the commonDocument IDs.
	for i := len(d.IncludeDocuments) - 1; i >= 0; i-- {
		documentSetID := d.IncludeDocuments[i]
		docSet, docSetFound := c.GetCommonDocumentSet(documentSetID)
		if !docSetFound {
			continue
		}

		for _, document := range docSet.Documents {
			_, found := documentMap[document.Name]
			if !found {
				documentMap[document.Name] = struct{}{}
				documents = append(documents, document)
			}
		}
	}

	return
}

func (d *DocumentationSource) GetDocument(c *Config, path string) (document Document, found bool) {
	for _, doc := range d.GetDocuments(c) {
		if doc.Name == path {
			return doc, true
		}
	}

	return
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

func (dt DocType) PossibleValues() string {
	return fmt.Sprintf("%s, %s", DocTypeMarkdown, DocTypeHTML)
}

const (
	DocTypeMarkdown DocType = "md"
	DocTypeHTML     DocType = "html"
)

type DocumentationSourceOptions struct {
	Selector string `yaml:"selector,omitempty"`
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
	Type    ExtractorType    `yaml:"type,omitempty"`    // TODO rename ExtractorType to CrawlerType during cleanup in #206
	Options ExtractorOptions `yaml:"options,omitempty"` // TODO rename ExtractorOptions to CrawlerOptions during cleanup in #206
	Include []string         `yaml:"include,omitempty"`
	Exclude []string         `yaml:"exclude,omitempty"`
}

type ExtractExtractor struct {
	Type    DocType                    `yaml:"type,omitempty"`    // TODO rename DocType to ExtractorType during cleanup in #206
	Options DocumentationSourceOptions `yaml:"options,omitempty"` // TODO rename DocumentationSourceOptions to ExtractorOptions during cleanup in #206
	Include []string                   `yaml:"include,omitempty"`
	Exclude []string                   `yaml:"exclude,omitempty"`
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
	Code          CheckCode          `yaml:"code,omitempty"`
	Documentation CheckDocumentation `yaml:"documentation,omitempty"`
	Options       CheckOptions       `yaml:"options,omitempty"`
}

type CheckCode struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type CheckDocumentation struct {
	Include []CheckDocumentationFilter `yaml:"include,omitempty"`
	Exclude []CheckDocumentationFilter `yaml:"exclude,omitempty"`
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
	Code          CheckCodeFilter          `yaml:"code,omitempty"`
	Documentation CheckDocumentationFilter `yaml:"documentation,omitempty"`
}

type CheckCodeFilter struct {
	Path string `yaml:"path,omitempty"`
}

type CheckDocumentationFilter struct {
	Source   string                        `yaml:"source,omitempty"`
	Document string                        `yaml:"document,omitempty"`
	Section  string                        `yaml:"section,omitempty"`
	URI      string                        `yaml:"uri,omitempty"`
	Tags     []CheckDocumentationFilterTag `yaml:"tags,omitempty"`
}

func (filter *CheckDocumentationFilter) GetParts() (source string, document string, section string) {
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

type CheckDocumentationFilterTag struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}
