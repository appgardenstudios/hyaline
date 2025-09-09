package llm

import (
	"context"
	"fmt"
	"hyaline/internal/config"
	"log/slog"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func callOpenAI(systemPrompt string, userPrompt string, tools []*Tool, cfg *config.LLM) (result string, err error) {
	if cfg.Key == "" {
		slog.Warn("Calling OpenAI without a key being set")
	}
	var clientOptions []option.RequestOption
	clientOptions = append(clientOptions, option.WithAPIKey(cfg.Key))

	// Handle endpoint logic
	if cfg.Provider == config.LLMProviderGitHubModels {
		if cfg.Endpoint == "" {
			slog.Debug("Using default GitHub Models endpoint")
			clientOptions = append(clientOptions, option.WithBaseURL("https://models.github.ai/inference"))
		} else {
			slog.Debug("Using custom GitHub Models endpoint", "endpoint", cfg.Endpoint)
			clientOptions = append(clientOptions, option.WithBaseURL(cfg.Endpoint))
		}
	} else if cfg.Endpoint != "" {
		slog.Debug("Using custom OpenAI endpoint", "endpoint", cfg.Endpoint)
		clientOptions = append(clientOptions, option.WithBaseURL(cfg.Endpoint))
	}

	client := openai.NewClient(clientOptions...)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(userPrompt),
	}

	// Tools
	toolParams := make([]openai.ChatCompletionToolUnionParam, len(tools))
	for i, tool := range tools {
		// Convert jsonschema to OpenAI function parameters format
		functionParams := openai.FunctionParameters{
			"type":       "object",
			"properties": tool.Schema.Properties,
		}
		if tool.Schema.Required != nil {
			functionParams["required"] = tool.Schema.Required
		}

		toolParams[i] = openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
					Parameters:  functionParams,
				},
			},
		}
	}

	// Loop and call llm until we don't have any outstanding tool calls left OR
	// a tool call signals that we are done
	for {
		// Call OpenAI with the message(s)
		var chatCompletion *openai.ChatCompletion
		params := openai.ChatCompletionNewParams{
			Model:    cfg.Model,
			Messages: messages,
		}

		// Add system message
		if systemPrompt != "" {
			params.Messages = append([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
			}, params.Messages...)
		}

		// Add tools if present
		if len(tools) > 0 {
			params.Tools = toolParams
			params.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			}
		}

		chatCompletion, err = client.Chat.Completions.New(context.Background(), params)
		if err != nil {
			slog.Error("llm.callOpenAI errored when sending a new message", "error", err)
			return
		}

		// Add new message to the list
		messages = append(messages, chatCompletion.Choices[0].Message.ToParam())

		// Initialize our toolResults
		toolResults := []openai.ChatCompletionMessageParamUnion{}

		// Initialize done sentinel to false
		// This is used to short-circuit the call loop in cases where a tool call
		// signals that we should stop.
		done := false

		// Process choice in our response
		choice := chatCompletion.Choices[0]
		if choice.Message.Content != "" {
			result = choice.Message.Content
		}

		// Handle tool calls
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Type == "function" {
				// Get our tool using the function name
				tool := getTool(toolCall.Function.Name, tools)
				if tool == nil {
					err = fmt.Errorf("invalid tool name received: %s", toolCall.Function.Name)
					slog.Error("llm.callOpenAI received an invalid tool name", "name", toolCall.Function.Name, "error", err)
					return
				}

				// Call the tool
				slog.Debug("llm.callOpenAI invoking tool", "tool", toolCall.Function.Name)
				stop, response, err := tool.Callback(toolCall.Function.Arguments)

				// Handle if the tool requests that we stop now rather than loop
				if stop {
					done = true
				}

				// Handle if there was an error
				if err != nil {
					slog.Error("llm.callOpenAI received a tool error", "tool", toolCall.Function.Name, "error", err)
					response = fmt.Sprintf("Error: %v", err)
				}

				// Add our result to the tool results
				toolResults = append(toolResults, openai.ToolMessage(response, toolCall.ID))
			}
		}

		// If we don't have any tool results OR a tool said stop, we are done
		if len(toolResults) == 0 || done {
			break
		}

		// Add tool result messages to our list of messages
		messages = append(messages, toolResults...)
	}

	return
}
