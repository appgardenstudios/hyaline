package llm

import (
	"context"
	"hyaline/internal/config"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func callAnthropic(systemPrompt string, userPrompt string, cfg *config.LLM) (result string, err error) {
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.Key),
	)
	var message *anthropic.Message
	message, err = client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: userPrompt},
			}},
		}},
		Model: cfg.Model,
	})

	if err == nil {
		result = message.Content[len(message.Content)-1].Text
	}

	return
}
