package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Systems []System `yaml:"systems"`
}

type System struct {
	ID   string `yaml:"id"`
	Code []Code `yaml:"code"`
	Docs []Doc  `yaml:"docs"`
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

func Load(path string) (cfg *Config, err error) {
	slog.Debug("Load config starting")
	// Read file from disk
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Debug("Load could not get an absolute path from the provided path", "path", path, "error", err)
		return
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		slog.Debug("Load could not read config file from disk", "error", err)
		return
	}

	// Parse file into struct
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		slog.Debug("Load could not unmarshal yaml config", "error", err)
		return
	}

	// Validate
	slog.Debug("Load config validating")
	codeCombinations := map[string]struct{}{}
	docsCombinations := map[string]struct{}{}
	for _, system := range cfg.Systems {
		// Ensure that system/code combinations are unique
		for _, code := range system.Code {
			id := system.ID + "-" + code.ID
			if _, ok := codeCombinations[id]; ok {
				slog.Debug("Load found duplicate system/code combination", "error", err)
				err = errors.New("duplicate system/code combination detected: " + system.ID + " > " + code.ID)
				return
			}
			codeCombinations[id] = struct{}{}
		}

		// Ensure that system/docs combinations are unique
		for _, doc := range system.Docs {
			id := system.ID + "-" + doc.ID
			if _, ok := docsCombinations[id]; ok {
				slog.Debug("Load found duplicate system/docs combination", "error", err)
				err = errors.New("duplicate system/docs combination detected: " + system.ID + " > " + doc.ID)
				return
			}
			docsCombinations[id] = struct{}{}
		}
	}

	slog.Debug("Load config complete")
	return
}
