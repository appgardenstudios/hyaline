package llm

import (
	"fmt"
	"hyaline/internal/config"
)

func callLLM(systemPrompt string, userPrompt string, cfg *config.LLM) (result string, err error) {
	switch cfg.Provider {
	case config.LLMProviderAnthropic:
		return callAnthropic(systemPrompt, userPrompt, cfg)
	case config.LLMProviderTesting:
		return "LLM TEST RESPONSE", nil
	default:
		err = fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	return
}
