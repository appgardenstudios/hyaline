package llm

import (
	"context"
	"fmt"
	"hyaline/internal/config"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func callAnthropic(systemPrompt string, userPrompt string, tools []*Tool, cfg *config.LLM) (result string, err error) {
	if cfg.Key == "" {
		slog.Warn("Calling anthropic without a key being set")
	}
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.Key),
	)

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
	}

	// Tools
	toolParams := make([]anthropic.ToolUnionParam, len(tools))
	for i, tool := range tools {
		toolParams[i] = anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: tool.Schema.Properties,
				},
			},
		}
	}
	var toolChoice anthropic.ToolChoiceUnionParam
	if len(tools) > 0 {
		toolChoice = anthropic.ToolChoiceUnionParam{
			OfToolChoiceAny: &anthropic.ToolChoiceAnyParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		}
	}

	// Loop and call llm until we don't have any outstanding tool calls left OR
	// a tool call signals that we are done
	for {
		// Call anthropic with the message(s)
		var message *anthropic.Message
		message, err = client.Messages.New(context.Background(), anthropic.MessageNewParams{
			Model:      cfg.Model,
			MaxTokens:  1024,
			System:     []anthropic.TextBlockParam{{Text: systemPrompt}},
			Messages:   messages,
			Tools:      toolParams,
			ToolChoice: toolChoice,
		})
		if err != nil {
			slog.Debug("llm.callAnthropic errored when sending a new message", "error", err)
			return
		}

		// Add new message(s) to the list
		messages = append(messages, message.ToParam())

		// Initialize our toolResults
		toolResults := []anthropic.ContentBlockParamUnion{}

		// Initialize done sentinel to false
		// This is used to short-circuit the call loop in cases where a tool call
		// signals that we should stop.
		done := false

		// Process block(s) in our response
		for _, block := range message.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.TextBlock:
				result = block.Text
			case anthropic.ToolUseBlock:
				// Get our tool using the block name
				tool := getTool(block.Name, tools)
				if tool == nil {
					err = fmt.Errorf("invalid tool name received: %s", block.Name)
					slog.Debug("llm.callAnthropic received an invalid tool name", "name", block.Name, "error", err)
					return
				}

				// Call the tool
				stop, response, err := tool.Callback(variant.JSON.Input.Raw())

				// Handle if the tool requests that we stop now rather than loop
				if stop {
					done = true
				}

				// Handle if there was an error
				isError := false
				if err != nil {
					isError = true
					slog.Debug("llm.callAnthropic received a tool error", "tool", block.Name, "error", err)
				}

				// Add our result to the tool results
				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, response, isError))
			}
		}

		// If we don't have any tool results OR a tool said stop, we are done
		if len(toolResults) == 0 || done {
			break
		}

		// Add tool result messages to our list of messages
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}

	return
}

func getTool(name string, tools []*Tool) *Tool {
	for _, tool := range tools {
		if tool.Name == name {
			return tool
		}
	}

	return nil
}
