package config

import (
	"errors"
	"log/slog"
)

func ValidateLLM(cfg *Config) (err error) {
	if cfg.LLM.Provider != "" && !cfg.LLM.Provider.IsValidLLMProvider() {
		err = errors.New("invalid llm provider detected: " + cfg.LLM.Provider.String())
		slog.Debug("config.Validate found invalid llm provider", "provider", cfg.LLM.Provider.String(), "error", err)
		return
	}

	return
}
