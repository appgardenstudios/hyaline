package tool

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

const getDocumentName = "get_document"

type getDocumentSchema struct {
	Name string `json:"name" jsonschema:"title=Document name,description=The name of the document,example=app/README.md,example=app/docs/overview.html"`
}

func GetDocument(systemID string, currentDB *sql.DB) *llm.Tool {
	return &llm.Tool{
		Name:        getDocumentName,
		Description: "Get a document by name. If the document does not exist this tool will return (Document Not Found)",
		Schema:      Reflector.Reflect(&getDocumentSchema{}),
		Callback: func(rawInput string) (bool, string, error) {
			slog.Debug("tool.GetDocument - called")

			// Parse input
			var input getDocumentSchema
			err := json.Unmarshal([]byte(rawInput), &input)
			if err != nil {
				slog.Debug("tool.GetDocument - could not parse tool call input, invalid json", "input", rawInput, "error", err)
				return true, "", err
			}

			// Split out documentationID from the document name
			parts := strings.Split(input.Name, "/")
			if len(parts) < 2 {
				slog.Debug("tool.GetDocument - could not parse documentation ID from name", "inputName", input.Name, "error", err)
				return true, "", err
			}
			documentationID := parts[0]
			documentID := strings.Join(parts[1:], "/")

			// Get document from the db
			document, err := sqlite.GetDocument(documentID, documentationID, systemID, currentDB)
			if err != nil {
				slog.Debug("tool.GetDocument - could not retrieve document", "documentID", documentID, "documentationID", documentationID, "systemID", systemID, "error", err)
				return true, "", err
			}

			// Handle document not found
			if document == nil {
				return false, "(Document Not Found)", nil
			}

			// Format result
			var result strings.Builder
			result.WriteString("<document>\n")
			result.WriteString(fmt.Sprintf("  <document_name>%s/%s</document_name>\n", document.DocumentationID, document.ID))
			result.WriteString("  <document_content>\n")
			result.WriteString(document.ExtractedData)
			result.WriteString("\n")
			result.WriteString("  </document_content>\n")
			result.WriteString("</document>\n")

			// Return the result stating we are not done
			return false, result.String(), nil
		},
	}
}
