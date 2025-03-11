package recommend

import (
	"database/sql"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

func SectionExists(current *sql.DB) (string, error) {
	// Note that this should really use a string builder
	system := "You are a senior technical writer who writes clear and accurate system documentation."

	// Get documents for the prompt
	potentialFiles := []string{"package.json", "Makefile"}
	files, err := sqlite.GetCodeFile(potentialFiles, current)
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
  </document>`, i, doc.ID, strings.TrimSpace(doc.RawData)))
	}
	documents := fmt.Sprintf(`<documents>
%s
</documents>`, strings.Join(formattedDocuments, "\n"))

	// Generate prompt
	// TODO what to do if there are 0 documents?
	prompt := fmt.Sprintf(`%s

Given the documents above, generate documentation describing how to run this %s project locally for development. Be clear, accurate, and show console commands where appropriate. Produce the documentation in the %s format.`, documents, "js", "markdown")
	slog.Debug("SectionExists prompt generated", "prompt", prompt, "system", system)

	fmt.Println("---")
	fmt.Println(prompt)
	fmt.Println("---")

	return "Action", nil
}
