package check

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hyaline/internal/config"
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

func Purpose(systemID string, documentationID string, document string, section []string, purpose string, contents string, cfg *config.LLM, currentDB *sql.DB) (matches bool, reason string, err error) {
	slog.Debug("check.checkLLM checking purpose", "document", document, "section", section)

	// Generate the system and user prompt
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	// Document/Section Info/Contents

	// Note that the document name has the documentationID prepended to
	// make it unique across documentation sources for the system.
	// The common tools in the `tools` package expect this syntax as well.
	documentName := fmt.Sprintf("%s/%s", documentationID, document)
	tagName := "document"
	name := documentName
	if len(section) > 0 {
		tagName = "section"
		name = section[len(section)-1]
		prompt.WriteString(fmt.Sprintf("This section belongs to the document %s.\n", documentName))
	}
	prompt.WriteString(fmt.Sprintf("The contents of the %s are give in <%s>.\n\n", tagName, tagName))
	prompt.WriteString(fmt.Sprintf("<%s>\n", tagName))
	prompt.WriteString(fmt.Sprintf("  <%s_name>%s</%s_name>\n", tagName, name, tagName))
	prompt.WriteString(fmt.Sprintf("  <%s_content>\n", tagName))
	prompt.WriteString(contents)
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
	prompt.WriteString(fmt.Sprintf("You may use other tools as needed to determine if this %s matches the stated purpose above. ", tagName))

	// Tools
	tools := []*llm.Tool{
		{
			Name:        checkPurposeMatch,
			Description: "Record that the document and/or section matches the stated purpose",
			Schema:      tool.Reflector.Reflect(&checkPurposeMatchSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("check.Purpose - determined document/section matches the stated purpose")

				// Parse input
				var match checkPurposeMatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("check.Purpose - could not parse tool call input, invalid json", "tool", checkPurposeMatch, "input", input, "error", err)
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
				slog.Debug("check.Purpose - determined document/section matches the stated purpose")

				// Parse input
				var match checkPurposeMismatchSchema
				err := json.Unmarshal([]byte(input), &match)
				if err != nil {
					slog.Debug("check.Purpose - could not parse tool call input, invalid json", "tool", checkPurposeMismatch, "input", input, "error", err)
					return true, "", err
				}

				// Mark match and reason
				matches = false
				reason = match.Reason

				// Report that we are done
				return true, "", nil
			},
		},
		tool.ListFiles(systemID, currentDB),
		tool.GetFile(systemID, currentDB),
		tool.ListDocuments(systemID, currentDB),
		tool.GetDocument(systemID, currentDB),
	}

	// Call LLM
	userPrompt := prompt.String()
	fmt.Println(userPrompt)
	slog.Debug("check.Purpose calling the llm")
	// slog.Debug("check.Purpose calling the llm", "systemPrompt", systemPrompt, "userPrompt", userPrompt)
	_, err = llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("check.Purpose encountered an error when calling the llm", "error", err)
		return
	}

	return
}
