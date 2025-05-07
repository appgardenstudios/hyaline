package check

import (
	"database/sql"
	"hyaline/internal/config"
	"log/slog"
)

func Completeness(systemID string, documentationID string, document string, section []string, purpose string, contents string, cfg *config.LLM, currentDB *sql.DB) (matches bool, reason string, err error) {
	slog.Debug("check.checkLLM checking completeness", "document", document, "section", section)

	return
}
