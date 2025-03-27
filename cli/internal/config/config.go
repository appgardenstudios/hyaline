package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

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

type Code struct {
	ID        string   `yaml:"id"`
	Extractor string   `yaml:"extractor"`
	Path      string   `yaml:"path"`
	Include   []string `yaml:"include"`
	Exclude   []string `yaml:"exclude"`
}

type Doc struct {
	ID        string   `yaml:"id"`
	Type      string   `yaml:"type"`
	Extractor string   `yaml:"extractor"`
	Path      string   `yaml:"path"`
	Include   []string `yaml:"include"`
	Exclude   []string `yaml:"exclude"`
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
	data = []byte(os.ExpandEnv(string(data)))

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

func validate(cfg *Config) (err error) {
	for _, system := range cfg.Systems {
		// Ensure that system/code combinations are unique
		codeIDs := map[string]struct{}{}
		for _, code := range system.Code {
			if _, ok := codeIDs[code.ID]; ok {
				err = errors.New("duplicate code id detected: " + system.ID + " > " + code.ID)
				slog.Debug("config.Validate found duplicate code id", "system", system.ID, "code", code.ID, "error", err)
				return
			}
			codeIDs[code.ID] = struct{}{}
		}

		// Ensure that system/docs combinations are unique
		docIDs := map[string]struct{}{}
		for _, doc := range system.Docs {
			if _, ok := docIDs[doc.ID]; ok {
				err = errors.New("duplicate docs id detected: " + system.ID + " > " + doc.ID)
				slog.Debug("config.Validate found duplicate docs id", "system", system.ID, "doc", doc.ID, "error", err)
				return
			}
			docIDs[doc.ID] = struct{}{}
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
