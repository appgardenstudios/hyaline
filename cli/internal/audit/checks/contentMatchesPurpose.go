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

const checkPurposeMatch = "purpose_match"

type checkPurposeMatchSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this purpose is met,description=The reason the document and/or section meets the stated purpose,example=This section meets the stated purpose of listing out all major changes to the API endpoints."`
}

const checkPurposeMismatch = "purpose_mismatch"

type checkPurposeMismatchSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this purpose is not met,description=The reason the document and/or section does not meet the stated purpose,example=This section does not meet the stated purpose of listing out all major changes to the API endpoints because it is missing major changes."`
}

// ContentMatchesPurpose validates content matches its stated purpose using LLM
func ContentMatchesPurpose(sourceID, documentID string, sectionID string, purpose string, content string, cfg *config.LLM) (bool, string, error) {
	slog.Debug("audit.checks.ContentMatchesPurpose starting", "document", documentID, "section", sectionID)

	var matches bool
	var reason string

	// Generate the system and user prompt
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	tagName := "document"
	if sectionID != "" {
		tagName = "section"
	}

	documentURI := &docs.DocumentURI{
		SourceID:     sourceID,
		DocumentPath: documentID,
		Section:      sectionID,
	}

	prompt.WriteString(fmt.Sprintf("The contents of the %s are given in <%s>.\n\n", tagName, tagName))
	prompt.WriteString(fmt.Sprintf("<%s>\n", tagName))
	prompt.WriteString(fmt.Sprintf("  <%s_uri>%s</%s_uri>\n", tagName, documentURI.String(), tagName))
	prompt.WriteString(fmt.Sprintf("  <%s_content>\n", tagName))
	prompt.WriteString(content)
	prompt.WriteString("\n")
	prompt.WriteString(fmt.Sprintf("  </%s_content>\n", tagName))
	prompt.WriteString(fmt.Sprintf("</%s>\n", tagName))
	prompt.WriteString("\n\n")

	// Purpose
	prompt.WriteString(fmt.Sprintf("The purpose of this %s should be \"%s\".\n\n", tagName, purpose))
	// Instructions
	prompt.WriteString(fmt.Sprintf("Given the purpose and the %s above, determine if the %s matches the stated purpose. ", tagName, tagName))
	prompt.WriteString(fmt.Sprintf("If so, call the tool %s to record that the %s matches the stated purpose. ", checkPurposeMatch, tagName))
	prompt.WriteString(fmt.Sprintf("If not, call the tool %s to record that the %s does not match the stated purpose. ", checkPurposeMismatch, tagName))

	// Tools for LLM response
	tools := []*llm.Tool{
		{
			Name:        checkPurposeMatch,
			Description: "Record that the document and/or section matches the stated purpose",
			Schema:      tool.Reflector.Reflect(&checkPurposeMatchSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("audit.checks.ContentMatchesPurpose - determined document/section matches the stated purpose")

				// Parse input
				var match checkPurposeMatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("audit.checks.ContentMatchesPurpose - could not parse tool call input, invalid json", "tool", checkPurposeMatch, "input", input, "error", err)
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
			Name:        checkPurposeMismatch,
			Description: "Record that the document and/or section does not match the stated purpose",
			Schema:      tool.Reflector.Reflect(&checkPurposeMismatchSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("audit.checks.ContentMatchesPurpose - determined document/section does not match the stated purpose")

				// Parse input
				var match checkPurposeMismatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("audit.checks.ContentMatchesPurpose - could not parse tool call input, invalid json", "tool", checkPurposeMismatch, "input", input, "error", err)
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
	userPrompt := prompt.String()
	slog.Debug("audit.checks.ContentMatchesPurpose calling the llm")
	_, err := llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("audit.checks.ContentMatchesPurpose encountered an error when calling the llm", "error", err)
		return false, "", err
	}

	return matches, reason, nil
}
