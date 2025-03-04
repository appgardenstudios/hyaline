package config

import (
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
	// Read file from disk
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return
	}

	// Parse file into struct
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return
	}

	// Validate config
	// TODO Ensure that system/code/documentation combinations are unique (see below)

	// Return
	return
}
