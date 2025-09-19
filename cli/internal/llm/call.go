package llm

import (
	"errors"
	"fmt"
	"hyaline/internal/config"
	"log/slog"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	Name        string
	Description string
	Schema      *jsonschema.Schema
	// Take in JSON input string, return stop, response, error
	Callback func(string) (bool, string, error)
}

type CallLLMHandler func(systemPrompt string, userPrompt string, tools []*Tool, cfg *config.LLM) (string, error)

func CallLLM(systemPrompt string, userPrompt string, tools []*Tool, cfg *config.LLM) (result string, err error) {
	if cfg == nil || cfg.Provider == "" {
		slog.Error("llm configuration must be present to call an llm")
		err = errors.New("llm configuration missing")
		return
	}
	slog.Debug("Calling LLM", "provider", cfg.Provider, "model", cfg.Model)
	switch cfg.Provider {
	case config.LLMProviderAnthropic:
		return callAnthropic(systemPrompt, userPrompt, tools, cfg)
	case config.LLMProviderOpenAI, config.LLMProviderGitHubModels:
		return callOpenAI(systemPrompt, userPrompt, tools, cfg)
	case config.LLMProviderTesting:
		return "LLM TEST RESPONSE", nil
	default:
		err = fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	return
}
