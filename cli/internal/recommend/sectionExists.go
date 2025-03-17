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

type sectionExistsDocument struct {
	Source  string
	Content string
}

func sectionExistsChangeFilesContains(file string, changeFiles *[]sqlite.GetChangeCodeFileRow) bool {
	if changeFiles == nil {
		return false
	}

	for _, changeFile := range *changeFiles {
		if changeFile.ID == file {
			return true
		}
	}

	return false
}

func SectionExists(sectionExists bool, section string, document string, systemID string, current *sql.DB, change *sql.DB, llmOpts config.LLM) (string, error) {
	// Note that this function should really use a string builder
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."

	// Create our existence action for use below
	var existenceAction string
	if !sectionExists {
		existenceAction = "create"
	} else {
		existenceAction = "add content to"
	}

	// Get documents for the prompt
	potentialFiles := []string{"package.json", "Makefile"}
	currentFiles, err := sqlite.GetCurrentCodeFiles(potentialFiles, systemID, current)
	if err != nil {
		slog.Debug("recommend.SectionExists could not get current code files", "error", err, "files", currentFiles)
		return "", err
	}
	slog.Debug("recommend.SectionExists got current files", "error", err, "count", len(*currentFiles))
	var changeFiles *[]sqlite.GetChangeCodeFileRow
	if change != nil {
		changeFiles, err = sqlite.GetChangeCodeFiles(potentialFiles, systemID, change)
		if err != nil {
			slog.Debug("recommend.SectionExists could not get change code files", "error", err, "files", currentFiles)
			return "", err
		}
		slog.Debug("recommend.SectionExists got change files", "error", err, "count", len(*changeFiles))
	}

	// Merge current and change documents
	// Adding files in changeFiles that are not deleted
	files := []sectionExistsDocument{}
	for _, file := range *currentFiles {
		// Add files that are not in change files
		if !sectionExistsChangeFilesContains(file.ID, changeFiles) {
			files = append(files, sectionExistsDocument{
				Source:  file.ID,
				Content: file.RawData,
			})
		}
	}
	if changeFiles != nil {
		for _, file := range *changeFiles {
			// Add files that are not deleted
			if file.Action != "Delete" {
				files = append(files, sectionExistsDocument{
					Source:  file.ID,
					Content: file.RawData,
				})
			}
		}
	}

	// If there are no files, return a generic action
	if len(files) < 1 {
		return fmt.Sprintf(`You should %s the section '%s' in '%s'. The section should describe how to run the project locally, along with any pre-requisites.`, existenceAction, section, document), nil
	}

	// Format our documents for the prompt
	formattedDocuments := []string{}
	for i, doc := range files {
		formattedDocuments = append(formattedDocuments, fmt.Sprintf(`  <document index="%d">
    <source>%s</source>
    <document_content>
%s
    <document_content>
  </document>`, i+1, doc.Source, strings.TrimSpace(doc.Content)))
	}
	documents := fmt.Sprintf(`<documents>
%s
</documents>`, strings.Join(formattedDocuments, "\n"))

	// Generate userPrompt
	userPrompt := fmt.Sprintf(`%s

Given the documents above, generate documentation describing how to run this %s project locally for development. Be clear, accurate, and show console commands where appropriate. Produce the documentation in the %s format.`, documents, "js", "markdown")
	slog.Debug("recommend.SectionExists prompts generated", "systemPrompt", systemPrompt, "userPrompt", userPrompt)

	// Get the suggestion from Anthropic
	suggestion, err := llm.CallAnthropic(systemPrompt, userPrompt, llmOpts.Model, llmOpts.Key)
	if err != nil {
		slog.Debug("recommend.SectionExists could not call anthropic", "error", err, "systemPrompt", systemPrompt, "userPrompt", userPrompt, "model", llmOpts.Model)
		return "", err
	}

	// Format the action
	action := fmt.Sprintf(`You should %s the section '%s' in '%s'. The contents of that section could contain something like:

%s`, existenceAction, section, document, suggestion)

	return action, nil
}
