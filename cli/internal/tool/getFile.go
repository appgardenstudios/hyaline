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

const getFileName = "get_file"

type getFileSchema struct {
	Name string `json:"name" jsonschema:"title=File name,description=The name of the file,example=app/server.js,example=app/route/router.js"`
}

func GetFile(systemID string, currentDB *sql.DB) *llm.Tool {
	return &llm.Tool{
		Name:        getFileName,
		Description: "Get a file by name. If the file does not exist this tool will return (File Not Found)",
		Schema:      Reflector.Reflect(&getFileSchema{}),
		Callback: func(rawInput string) (bool, string, error) {
			slog.Debug("tool.GetFile - called")

			// Parse input
			var input getFileSchema
			err := json.Unmarshal([]byte(rawInput), &input)
			if err != nil {
				slog.Debug("tool.GetFile - could not parse tool call input, invalid json", "input", rawInput, "error", err)
				return true, "", err
			}

			// Split out codeID from the file name
			parts := strings.Split(input.Name, "/")
			if len(parts) < 2 {
				slog.Debug("tool.GetFile - could not parse code ID from name", "inputName", input.Name, "error", err)
				return true, "", err
			}
			codeID := parts[0]
			fileID := strings.Join(parts[1:], "/")

			// Get document from the db
			file, err := sqlite.GetFile(fileID, codeID, systemID, currentDB)
			if err != nil {
				slog.Debug("tool.GetFile - could not retrieve document", "fileID", fileID, "codeID", codeID, "systemID", systemID, "error", err)
				return true, "", err
			}

			// Handle document not found
			if file == nil {
				return false, "(File Not Found)", nil
			}

			// Format result
			var result strings.Builder
			result.WriteString("<file>\n")
			result.WriteString(fmt.Sprintf("  <file_name>%s/%s</file_name>\n", file.CodeID, file.ID))
			result.WriteString("  <file_content>\n")
			result.WriteString(file.RawData)
			result.WriteString("\n")
			result.WriteString("  </file_content>\n")
			result.WriteString("</file>\n")

			// Return the result stating we are not done
			return false, result.String(), nil
		},
	}
}
