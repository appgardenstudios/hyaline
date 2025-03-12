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
	ID        string `yaml:"id"`
	Extractor string `yaml:"extractor"`
	Path      string `yaml:"path"`
	Preset    string `yaml:"preset"`
}

type Doc struct {
	ID        string `yaml:"id"`
	Type      string `yaml:"type"`
	Extractor string `yaml:"extractor"`
	Path      string `yaml:"path"`
	Glob      string `yaml:"glob"`
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
	data, err := os.ReadFile(absPath)
	if err != nil {
		slog.Debug("config.Load could not read config file from disk", "error", err)
		return
	}

	// Replace any env references ($KEY or ${KEY} with the contents of KEY from env)
	data = []byte(os.ExpandEnv(string(data)))

	// Parse file into struct
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		slog.Debug("config.Load could not unmarshal yaml config", "error", err)
		return
	}

	// Validate
	slog.Debug("config.Load config validating")
	codeCombinations := map[string]struct{}{}
	docsCombinations := map[string]struct{}{}
	for _, system := range cfg.Systems {
		// Ensure that system/code combinations are unique
		for _, code := range system.Code {
			id := system.ID + "-" + code.ID
			if _, ok := codeCombinations[id]; ok {
				slog.Debug("config.Load found duplicate system/code combination", "error", err)
				err = errors.New("duplicate system/code combination detected: " + system.ID + " > " + code.ID)
				return
			}
			codeCombinations[id] = struct{}{}
		}

		// Ensure that system/docs combinations are unique
		for _, doc := range system.Docs {
			id := system.ID + "-" + doc.ID
			if _, ok := docsCombinations[id]; ok {
				slog.Debug("config.Load found duplicate system/docs combination", "error", err)
				err = errors.New("duplicate system/docs combination detected: " + system.ID + " > " + doc.ID)
				return
			}
			docsCombinations[id] = struct{}{}
		}
	}

	slog.Debug("config.Load config complete")
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
