package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type AnthropicMessagesRequest struct {
	Model       string                            `json:"model"`
	MaxTokens   int                               `json:"max_tokens"`
	System      string                            `json:"system"`
	Messages    []AnthropicMessagesRequestMessage `json:"messages"`
	Temperature float64                           `json:"temperature,omitempty"`
}

type AnthropicMessagesRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicMessagesResponse struct {
	ID           string                             `json:"id"`
	Type         string                             `json:"type"`
	Role         string                             `json:"role"`
	Model        string                             `json:"model"`
	Content      []AnthropicMessagesResponseContent `json:"content"`
	StopReason   string                             `json:"stop_reason"`
	StopSequence string                             `json:"stop_sequence"`
	Usage        AnthropicMessagesResponseUsage     `json:"usage"`
}

type AnthropicMessagesResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicMessagesResponseUsage struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	OutputTokens             int `json:"output_tokens"`
}

func CallAnthropic(systemPrompt string, userPrompt string, model string, key string) (action string, err error) {
	// Create our client
	// Note that we will want to create this at a higher scope and re-use it
	client := &http.Client{Timeout: 60 * time.Second}

	// Marshal our request body
	reqData := AnthropicMessagesRequest{
		Model:     model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages: []AnthropicMessagesRequestMessage{
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		slog.Debug("llm.CallAnthropic could not marshal json", "error", err)
		return
	}

	// Create our request
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Debug("llm.CallAnthropic could not create request", "error", err)
		return
	}

	// Add headers
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("llm.CallAnthropic could not make request", "error", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("llm.CallAnthropic could not read response body", "error", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		slog.Debug("llm.CallAnthropic returned with a non-200 status code", "statusCode", resp.StatusCode, "body", string(resBodyBytes))
		return
	}

	// Parse response body
	var resBody AnthropicMessagesResponse
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		slog.Debug("llm.CallAnthropic could not unmarshal response body", "error", err, "body", string(resBodyBytes))
		return
	}

	// Get response text
	// Note that we should probably test for type == text here eventually
	for _, content := range resBody.Content {
		action = content.Text
	}

	return
}
