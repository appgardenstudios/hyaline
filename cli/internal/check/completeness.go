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

const checkIsComplete = "is_complete"

type checkIsCompleteSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this is complete,description=The reason the document and/or section is complete,example=This section is complete because it lists out all major changes to the API endpoints."`
}

const checkIsNotComplete = "is_not_complete"

type checkIsNotCompleteSchema struct {
	Reason string `json:"reason" jsonschema:"title=The reason this is not complete,description=The reason the document and/or section is not complete,example=This section is not complete because it is missing major changes."`
}

func Completeness(systemID string, documentationID string, document string, section []string, purpose string, contents string, cfg *config.LLM, currentDB *sql.DB) (isComplete bool, reason string, err error) {
	slog.Debug("check.Completeness checking completeness", "document", document, "section", section)

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
	prompt.WriteString(fmt.Sprintf("Given the purpose and the %s above, determine if the %s is complete. ", tagName, tagName))
	prompt.WriteString(fmt.Sprintf("If so, call the tool %s to record that the %s is complete. ", checkPurposeMatch, tagName))
	prompt.WriteString(fmt.Sprintf("If not, call the tool %s to record that the %s is not complete. ", checkPurposeMismatch, tagName))
	prompt.WriteString(fmt.Sprintf("You may use other tools as needed to determine if this %s is complete. ", tagName))

	// Tools
	tools := []*llm.Tool{
		{
			Name:        checkIsComplete,
			Description: "Record that the document and/or section is complete",
			Schema:      tool.Reflector.Reflect(&checkIsCompleteSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("check.Completeness - determined document/section is complete")

				// Parse input
				var result checkIsCompleteSchema
				err := json.Unmarshal([]byte(input), &result)
				if err != nil {
					slog.Debug("check.Completeness - could not parse tool call input, invalid json", "tool", checkPurposeMatch, "input", input, "error", err)
					return true, "", err
				}

				// Mark completeness and reason
				isComplete = true
				reason = result.Reason

				// Report that we are done
				return true, "", nil
			},
		},
		{
			Name:        checkIsNotComplete,
			Description: "Record that the document and/or section is NOT complete",
			Schema:      tool.Reflector.Reflect(&checkIsNotCompleteSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("check.Completeness - determined document/section matches the stated purpose")

				// Parse input
				var result checkIsNotCompleteSchema
				err := json.Unmarshal([]byte(input), &result)
				if err != nil {
					slog.Debug("check.Completeness - could not parse tool call input, invalid json", "tool", checkPurposeMismatch, "input", input, "error", err)
					return true, "", err
				}

				// Mark completeness and reason
				isComplete = false
				reason = result.Reason

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
	// fmt.Println(userPrompt)
	slog.Debug("check.Completeness calling the llm")
	// slog.Debug("check.Completeness calling the llm", "systemPrompt", systemPrompt, "userPrompt", userPrompt)
	_, err = llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("check.PuCompletenessrpose encountered an error when calling the llm", "error", err)
		return
	}

	return
}
