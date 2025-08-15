package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
