package llm

import (
	"fmt"
	"hyaline/internal/config"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	Name        string
	Description string
	Schema      *jsonschema.Schema
	// Take in JSON string, return stop, response, error
	Callback func(string) (bool, string, error)
}

func CallLLM(systemPrompt string, userPrompt string, tools []*Tool, cfg *config.LLM) (result string, err error) {
	switch cfg.Provider {
	case config.LLMProviderAnthropic:
		return callAnthropic(systemPrompt, userPrompt, tools, cfg)
	case config.LLMProviderTesting:
		return "LLM TEST RESPONSE", nil
	default:
		err = fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	return
}
