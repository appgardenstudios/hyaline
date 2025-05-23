package tool

import (
	"database/sql"
	"fmt"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

const listFilesName = "list_files"

type listFilesSchema struct {
}

func ListFiles(systemID string, currentDB *sql.DB) *llm.Tool {
	return &llm.Tool{
		Name:        listFilesName,
		Description: "List available code files",
		Schema:      Reflector.Reflect(&listFilesSchema{}),
		Callback: func(rawInput string) (bool, string, error) {
			slog.Debug("tool.ListFiles - called")

			// Get all files for this system
			files, err := sqlite.GetAllSystemFilesForSystem(systemID, currentDB)
			if err != nil {
				slog.Debug("tool.ListDocuments - could not retrieve files", "systemID", systemID, "error", err)
				return true, "", err
			}

			// Compile our list of files
			var result strings.Builder
			result.WriteString("<files>\n")
			for _, file := range files {
				result.WriteString(fmt.Sprintf("  <file name=\"%s/%s\">\n", file.CodeID, file.ID))
			}
			result.WriteString("</files>\n")

			// Return the result stating we are not done
			return false, result.String(), nil
		},
	}
}
