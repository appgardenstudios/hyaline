package checks

import (
	"encoding/json"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/llm"
	"hyaline/internal/tool"
	"log/slog"
	"strings"
)

const checkPromptMatch = "prompt_match"

type checkPromptMatchSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this criteria is met,description=The reason the document and/or section meets the given criteria,example=This section meets the given criteria because it contains detailed installation instructions."`
}

const checkPromptMismatch = "prompt_mismatch"

type checkPromptMismatchSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this criteria is not met,description=The reason the document and/or section does not meet the given criteria,example=This section does not meet the given criteria because it lacks specific installation steps."`
}

// ContentMatchesPrompt uses LLM to validate content against a custom prompt
func ContentMatchesPrompt(sourceID, documentID string, sectionID string, prompt string, content string, cfg *config.LLM) (bool, string, error) {
	slog.Debug("audit.checks.ContentMatchesPrompt starting", "document", documentID, "section", sectionID)

	var matches bool
	var reason string

	// Generate the system and user prompt
	systemPrompt := "You are a senior technical writer who writes clear and accurate documentation."
	var userPromptBuilder strings.Builder

	tagName := "document"
	if sectionID != "" {
		tagName = "section"
	}

	documentURI := &docs.DocumentURI{
		SourceID:     sourceID,
		DocumentPath: documentID,
		Section:      sectionID,
	}

	userPromptBuilder.WriteString(fmt.Sprintf("The contents of the %s are given in <%s>.\n\n", tagName, tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("<%s>\n", tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("  <%s_uri>%s</%s_uri>\n", tagName, documentURI.String(), tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("  <%s_content>\n", tagName))
	userPromptBuilder.WriteString(content)
	userPromptBuilder.WriteString("\n")
	userPromptBuilder.WriteString(fmt.Sprintf("  </%s_content>\n", tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("</%s>\n", tagName))
	userPromptBuilder.WriteString("\n\n")

	// Criteria
	userPromptBuilder.WriteString(fmt.Sprintf("The criteria for this %s is: %s\n\n", tagName, prompt))
	// Instructions
	userPromptBuilder.WriteString(fmt.Sprintf("Given the criteria and the %s above, determine if the %s meets the given criteria. ", tagName, tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("If so, call the tool %s to record that the %s meets the given criteria. ", checkPromptMatch, tagName))
	userPromptBuilder.WriteString(fmt.Sprintf("If not, call the tool %s to record that the %s does not meet the given criteria. ", checkPromptMismatch, tagName))

	// Tools for LLM response
	tools := []*llm.Tool{
		{
			Name:        checkPromptMatch,
			Description: "Record that the content matches the given criteria",
			Schema:      tool.Reflector.Reflect(&checkPromptMatchSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("audit.checks.ContentMatchesPrompt - determined content matches the given criteria")

				// Parse input
				var match checkPromptMatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("audit.checks.ContentMatchesPrompt - could not parse tool call input, invalid json", "tool", checkPromptMatch, "input", input, "error", err)
					return true, "", err
				}

				// Mark match and reason
				matches = true
				reason = match.Reason

				// Report that we are done
				return true, "", nil
			},
		},
		{
			Name:        checkPromptMismatch,
			Description: "Record that the content does not match the given criteria",
			Schema:      tool.Reflector.Reflect(&checkPromptMismatchSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("audit.checks.ContentMatchesPrompt - determined content does not match the given criteria")

				// Parse input
				var match checkPromptMismatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("audit.checks.ContentMatchesPrompt - could not parse tool call input, invalid json", "tool", checkPromptMismatch, "input", input, "error", err)
					return true, "", err
				}

				// Mark match and reason
				matches = false
				reason = match.Reason

				// Report that we are done
				return true, "", nil
			},
		},
	}

	// Call LLM
	userPrompt := userPromptBuilder.String()
	slog.Debug("audit.checks.ContentMatchesPrompt calling the llm")
	_, err := llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("audit.checks.ContentMatchesPrompt encountered an error when calling the llm", "error", err)
		return false, "", err
	}

	return matches, reason, nil
}
