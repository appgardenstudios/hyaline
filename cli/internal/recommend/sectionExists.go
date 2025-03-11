package recommend

import (
	"database/sql"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

func SectionExists(sectionExists bool, section string, document string, systemID string, current *sql.DB, llmOpts config.LLM) (string, error) {
	// Note that this should really use a string builder
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."

	// Get documents for the prompt
	potentialFiles := []string{"package.json", "Makefile"}
	files, err := sqlite.GetCodeFile(potentialFiles, systemID, current)
	if err != nil {
		slog.Debug("SectionExists could not get code files", "error", err, "files", files)
		return "", err
	}
	formattedDocuments := []string{}
	for i, doc := range *files {
		formattedDocuments = append(formattedDocuments, fmt.Sprintf(`  <document index="%d">
    <source>%s</source>
    <document_content>
%s
    <document_content>
  </document>`, i+1, doc.ID, strings.TrimSpace(doc.RawData)))
	}
	documents := fmt.Sprintf(`<documents>
%s
</documents>`, strings.Join(formattedDocuments, "\n"))

	// Generate userPrompt
	// TODO what to do if there are 0 documents?
	userPrompt := fmt.Sprintf(`%s

Given the documents above, generate documentation describing how to run this %s project locally for development. Be clear, accurate, and show console commands where appropriate. Produce the documentation in the %s format.`, documents, "js", "markdown")
	slog.Debug("SectionExists prompts generated", "systemPrompt", systemPrompt, "userPrompt", userPrompt)

	suggestion, err := llm.CallAnthropic(systemPrompt, userPrompt, llmOpts.Model, llmOpts.Key)
	if err != nil {
		slog.Debug("SectionExists could not call anthropic", "error", err, "systemPrompt", systemPrompt, "userPrompt", userPrompt, "model", llmOpts.Model)
		return "", err
	}

	// Format the action
	var existenceAction string
	if sectionExists {
		existenceAction = "create"
	} else {
		existenceAction = "add content to"
	}
	action := fmt.Sprintf(`You should %s the section '%s' in '%s'. The contents of that section could contain something like:

%s`, existenceAction, section, document, suggestion)

	return action, nil
}
