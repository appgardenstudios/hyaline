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

	// If the value contains the 2 character sequence "\"+"n", replace it with a newline character.
	if strings.Contains(val, "\\n") {
		val = strings.ReplaceAll(val, "\\n", "\n")
	}

	// If the value contains a newline character, escape the entire string and enclose it in "" so
	// that the yaml parser interprets the \n as newline characters when parsing it.
	if strings.Contains(val, "\n") {
		// Strip out carriage returns (just in case)
		val = strings.ReplaceAll(val, "\r", "")
		// Escape all newlines and double quotes
		val = strings.ReplaceAll(val, "\"", "\\\"")
		val = strings.ReplaceAll(val, "\n", "\\n")
		return fmt.Sprintf("\"%s\"", val)
	}

	return val
}

func validate(cfg *Config) (err error) {
	// Validate LLM
	if cfg.LLM.Provider != "" && !cfg.LLM.Provider.IsValidLLMProvider() {
		err = errors.New("invalid llm provider detected: " + cfg.LLM.Provider.String())
		slog.Debug("config.Validate found invalid llm provider", "provider", cfg.LLM.Provider.String(), "error", err)
		return
	}

	// Validate Systems
	for _, system := range cfg.Systems {

		// Validate code block
		codeIDs := map[string]struct{}{}
		for _, code := range system.CodeSources {
			// Ensure that system/code combinations are unique
			if _, ok := codeIDs[code.ID]; ok {
				err = errors.New("duplicate code id detected: " + system.ID + " > " + code.ID)
				slog.Debug("config.Validate found duplicate code id", "system", system.ID, "code", code.ID, "error", err)
				return
			}
			codeIDs[code.ID] = struct{}{}

			// Ensure extractor is valid
			if !code.Extractor.Type.IsValidCodeExtractor() {
				err = errors.New("invalid code extractor detected: " + system.ID + " > " + code.ID + " > " + code.Extractor.Type.String())
				slog.Debug("config.Validate found invalid code extractor", "extractor", code.Extractor.Type.String(), "system", system.ID, "code", code.ID, "error", err)
				return
			}

			// Ensure include patterns are valid
			for _, include := range code.Extractor.Include {
				if !doublestar.ValidatePattern(include) {
					err = errors.New("invalid code include pattern detected: " + system.ID + " > " + code.ID + " > " + include)
					slog.Debug("config.Validate found invalid include pattern", "include", include, "system", system.ID, "code", code.ID, "error", err)
					return
				}
			}

			// Ensure exclude patterns are valid
			for _, exclude := range code.Extractor.Exclude {
				if !doublestar.ValidatePattern(exclude) {
					err = errors.New("invalid code exclude pattern detected: " + system.ID + " > " + code.ID + " > " + exclude)
					slog.Debug("config.Validate found invalid exclude pattern", "exclude", exclude, "system", system.ID, "code", code.ID, "error", err)
					return
				}
			}
		}

		// Validate docs block
		docIDs := map[string]struct{}{}
		for _, doc := range system.DocumentationSources {
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
			if !doc.Extractor.Type.IsValidDocExtractor() {
				err = errors.New("invalid doc extractor detected: " + system.ID + " > " + doc.ID + " > " + doc.Extractor.Type.String())
				slog.Debug("config.Validate found invalid doc extractor", "extractor", doc.Extractor.Type.String(), "system", system.ID, "doc", doc.ID, "error", err)
				return
			}

			// Ensure include patterns are valid
			for _, include := range doc.Extractor.Include {
				if !doublestar.ValidatePattern(include) {
					err = errors.New("invalid doc include pattern detected: " + system.ID + " > " + doc.ID + " > " + include)
					slog.Debug("config.Validate found invalid doc include", "include", include, "system", system.ID, "doc", doc.ID, "error", err)
					return
				}
			}

			// Ensure exclude patterns are valid
			for _, exclude := range doc.Extractor.Exclude {
				if !doublestar.ValidatePattern(exclude) {
					err = errors.New("invalid doc exclude pattern detected: " + system.ID + " > " + doc.ID + " > " + exclude)
					slog.Debug("config.Validate found invalid doc exclude", "exclude", exclude, "system", system.ID, "doc", doc.ID, "error", err)
					return
				}
			}
		}
	}

	// Validate desiredDocuments
	// TODO
	ruleIDs := map[string]struct{}{}
	for _, rule := range cfg.CommonDocuments {
		if _, ok := ruleIDs[rule.ID]; ok {
			err = errors.New("duplicate rule id detected: " + rule.ID)
			slog.Debug("config.Validate found duplicate rule id", "rule", rule.ID, "error", err)
			return
		}
		ruleIDs[rule.ID] = struct{}{}
	}

	return
}
