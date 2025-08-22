package config

import (
	"log/slog"
)

const sourceIDRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const metadataTagKeyRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const metadataTagValueRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const auditRuleIDRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`

func Validate(cfg *Config) (err error) {
	// Validate LLM
	err = ValidateLLM(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid llm", "error", err)
		return
	}

	// Verify extract
	err = ValidateExtract(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid extract", "error", err)
		return
	}

	// Verify check
	err = ValidateCheck(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid check", "error", err)
		return
	}

	// Verify audit
	err = ValidateAudit(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid audit", "error", err)
		return
	}

	return
}
