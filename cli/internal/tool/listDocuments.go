package tool

import (
	"database/sql"
	"fmt"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

const listDocumentsName = "list_documents"

type listDocumentsSchema struct {
}

func ListDocuments(systemID string, currentDB *sql.DB) *llm.Tool {
	return &llm.Tool{
		Name:        listDocumentsName,
		Description: "List available documentation documents",
		Schema:      Reflector.Reflect(&listDocumentsSchema{}),
		Callback: func(rawInput string) (bool, string, error) {
			slog.Debug("tool.ListDocuments - called")

			// Get all documents from DB
			documents, err := sqlite.GetAllSystemDocumentsForSystem(systemID, currentDB)
			if err != nil {
				slog.Debug("tool.ListDocuments - could not retrieve documents", "systemID", systemID, "error", err)
				return true, "", err
			}

			// Compile our list of documents
			var result strings.Builder
			result.WriteString("<documents>\n")
			for _, document := range documents {
				result.WriteString(fmt.Sprintf("  <document name=\"%s/%s\">\n", document.DocumentationID, document.ID))
			}
			result.WriteString("</documents>\n")

			// Return the result stating we are not done
			return false, result.String(), nil
		},
	}
}
