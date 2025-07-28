package extract

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

func extractMd(id string, sourceID string, rawData []byte, db *sqlite.Queries) error {
	// Clean up raw data
	extractedData := strings.TrimSpace(string(rawData))
	extractedData = strings.ReplaceAll(extractedData, "\r", "")

	// Insert document
	err := db.InsertDocument(context.Background(), sqlite.InsertDocumentParams{
		ID:            id,
		SourceID:      sourceID,
		Type:          config.DocTypeMarkdown.String(),
		Purpose:       "",
		RawData:       string(rawData),
		ExtractedData: extractedData,
	})
	if err != nil {
		slog.Debug("extract.extractMd could not insert document", "error", err)
		return err
	}

	// Extract/insert sections
	err = extractSections(id, sourceID, extractedData, db)
	if err != nil {
		slog.Debug("extract.extractMd could not extract sections", "error", err)
		return err
	}

	return nil
}
