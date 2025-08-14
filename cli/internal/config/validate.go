package config

import (
	"errors"
	"log/slog"
)

const sourceIDRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const metadataTagKeyRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const metadataTagValueRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const auditRuleIDRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`

func validate(cfg *Config) (err error) {
	// Validate LLM
	if cfg.LLM.Provider != "" && !cfg.LLM.Provider.IsValidLLMProvider() {
		err = errors.New("invalid llm provider detected: " + cfg.LLM.Provider.String())
		slog.Debug("config.Validate found invalid llm provider", "provider", cfg.LLM.Provider.String(), "error", err)
		return
	}

	// Verify extract
	err = validateExtract(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid extract", "error", err)
		return
	}

	// Verify check
	err = validateCheck(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid check", "error", err)
		return
	}

	// Verify audit
	err = validateAudit(cfg)
	if err != nil {
		slog.Debug("config.Validate found invalid audit", "error", err)
		return
	}

	return
}
