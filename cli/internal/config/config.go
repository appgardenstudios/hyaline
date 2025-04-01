package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LLM     LLM      `yaml:"llm"`
	Systems []System `yaml:"systems"`
}

type LLM struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Key      string `yaml:"key"`
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

func (e Extractor) IsValid() bool {
	switch e {
	case ExtractorFs, ExtractorGit:
		return true
	default:
		return false
	}
}

const (
	ExtractorFs  Extractor = "fs"
	ExtractorGit Extractor = "git"
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
	ID         string         `yaml:"id"`
	Type       DocType        `yaml:"type"`
	HTML       DocHTMLOptions `yaml:"html"`
	Extractor  Extractor      `yaml:"extractor"`
	FsOptions  FsOptions      `yaml:"fs"`
	GitOptions GitOptions     `yaml:"git"`
	Include    []string       `yaml:"include"`
	Exclude    []string       `yaml:"exclude"`
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

func Load(path string) (cfg *Config, err error) {
	slog.Debug("config.Load config starting")
	// Read file from disk
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Debug("config.Load could not get an absolute path from the provided path", "path", path, "error", err)
		return
	}
	slog.Debug("config.Load resolved absolute path for config", "path", path, "absPath", absPath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		slog.Debug("config.Load could not read config file from disk", "error", err)
		return
	}

	// Replace any env references ($KEY or ${KEY} with the contents of KEY from env)
	data = []byte(os.Expand(string(data), getEscapedEnv))

	// Parse file into the struct
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		slog.Debug("config.Load could not unmarshal yaml config", "error", err)
		return
	}

	// Validate
	err = validate(cfg)
	if err != nil {
		slog.Debug("config.Load found an invalid config", "error", err)
		return
	}

	slog.Debug("config.Load config complete")
	return
}

// Handle cases where an env var contains newlines by escaping them and
// wrapping the value in double quotes so that \n will be expanded back out
// in the final string value (ex. PEM files). This is done so our env var
// substitution does not mess up the yaml config file.
func getEscapedEnv(key string) string {
	val := os.Getenv(key)
	if strings.Contains(val, "\n") {
		// Strip out carriage returns
		val = strings.ReplaceAll(val, "\r", "")
		// Escape all newlines and double quotes
		val = strings.ReplaceAll(val, "\"", "\\\"")
		val = strings.ReplaceAll(val, "\n", "\\n")
		return fmt.Sprintf("\"%s\"", val)
	}

	return val
}

func validate(cfg *Config) (err error) {
	// Validate Systems
	for _, system := range cfg.Systems {

		// Validate code block
		codeIDs := map[string]struct{}{}
		for _, code := range system.Code {
			// Ensure that system/code combinations are unique
			if _, ok := codeIDs[code.ID]; ok {
				err = errors.New("duplicate code id detected: " + system.ID + " > " + code.ID)
				slog.Debug("config.Validate found duplicate code id", "system", system.ID, "code", code.ID, "error", err)
				return
			}
			codeIDs[code.ID] = struct{}{}

			// Ensure extractor is valid
			// TODO

			// Ensure include patterns are valid
			for _, include := range code.Include {
				if !doublestar.ValidatePattern(include) {
					err = errors.New("invalid code include pattern detected: " + system.ID + " > " + code.ID + ">" + include)
					slog.Debug("config.Validate ", "include", include, "system", system.ID, "code", code.ID, "error", err)
					return
				}
			}

			// Ensure exclude patterns are valid
			for _, exclude := range code.Exclude {
				if !doublestar.ValidatePattern(exclude) {
					err = errors.New("invalid code exclude pattern detected: " + system.ID + " > " + code.ID + ">" + exclude)
					slog.Debug("config.Validate ", "exclude", exclude, "system", system.ID, "code", code.ID, "error", err)
					return
				}
			}
		}

		// Validate docs block
		docIDs := map[string]struct{}{}
		for _, doc := range system.Docs {
			// Ensure that system/docs combinations are unique
			if _, ok := docIDs[doc.ID]; ok {
				err = errors.New("duplicate docs id detected: " + system.ID + " > " + doc.ID)
				slog.Debug("config.Validate found duplicate docs id", "system", system.ID, "doc", doc.ID, "error", err)
				return
			}
			docIDs[doc.ID] = struct{}{}

			// Ensure that doc type is valid
			if !doc.Type.IsValid() {
				err = errors.New("invalid doc type '" + doc.Type.String() + "' detected: " + system.ID + " > " + doc.ID)
				slog.Debug("config.Validate found invalid doc type", "system", system.ID, "doc", doc.ID, "type", doc.Type.String(), "error", err)
				return
			}

			// Ensure extractor is valid
			// TODO

			// Ensure include patterns are valid
			for _, include := range doc.Include {
				if !doublestar.ValidatePattern(include) {
					err = errors.New("invalid doc include pattern detected: " + system.ID + " > " + doc.ID + ">" + include)
					slog.Debug("config.Validate ", "include", include, "system", system.ID, "doc", doc.ID, "error", err)
					return
				}
			}

			// Ensure exclude patterns are valid
			for _, exclude := range doc.Exclude {
				if !doublestar.ValidatePattern(exclude) {
					err = errors.New("invalid doc exclude pattern detected: " + system.ID + " > " + doc.ID + ">" + exclude)
					slog.Debug("config.Validate ", "exclude", exclude, "system", system.ID, "doc", doc.ID, "error", err)
					return
				}
			}
		}
	}
	return
}

func GetSystem(system string, cfg *Config) (targetSystem *System, err error) {
	for _, s := range cfg.Systems {
		if s.ID == system {
			targetSystem = &s
		}
	}
	if targetSystem == nil {
		slog.Debug("config.GetSystem target system not found", "system", system)
		err = errors.New("system not found: " + system)
	}

	return
}
