package extract

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"

	"github.com/bmatcuk/doublestar/v4"
)

func addMetadata(sourceID string, cfg []config.ExtractMetadata, db *sqlite.Queries) error {
	ctx := context.Background()

	// Get the list of all inserted documents and sections
	documents, err := db.GetDocumentIDsForSource(ctx, sourceID)
	if err != nil {
		slog.Debug("extract.addMetadata could not retrieve inserted documents", "error", err)
		return err
	}
	sections, err := db.GetSectionIDsForSource(ctx, sourceID)
	if err != nil {
		slog.Debug("extract.addMetadata could not retrieve inserted sections", "error", err)
		return err
	}

	// Loop through metadata and insert/update matching documents/sections
	for i, metadata := range cfg {
		// Only add metadata to sections if section is set
		if metadata.Section == "" {
			for _, document := range documents {
				if doublestar.MatchUnvalidated(metadata.Document, document) {
					// Add purpose to document (if set)
					if metadata.Purpose != "" {
						err = db.UpdateDocumentPurpose(ctx, sqlite.UpdateDocumentPurposeParams{
							Purpose:  metadata.Purpose,
							ID:       document,
							SourceID: sourceID,
						})
						if err != nil {
							slog.Debug("extract.addMetadata could not update document purpose", "metadata", i, "error", err)
							return err
						}
					}

					// Add tags to document
					for j, tag := range metadata.Tags {
						err = db.UpsertDocumentTag(ctx, sqlite.UpsertDocumentTagParams{
							SourceID:   sourceID,
							DocumentID: document,
							TagKey:     tag.Key,
							TagValue:   tag.Value,
						})
						if err != nil {
							slog.Debug("extract.addMetadata could not insert document tag", "metadata", i, "tag", j, "error", err)
							return err
						}
					}
				}
			}
		} else {
			for _, section := range sections {
				if doublestar.MatchUnvalidated(metadata.Document, section.DocumentID) && doublestar.MatchUnvalidated(metadata.Section, section.ID) {
					// Add purpose to document (if set)
					if metadata.Purpose != "" {
						err = db.UpdateSectionPurpose(ctx, sqlite.UpdateSectionPurposeParams{
							Purpose:    metadata.Purpose,
							ID:         section.ID,
							DocumentID: section.DocumentID,
							SourceID:   sourceID,
						})
						if err != nil {
							slog.Debug("extract.addMetadata could not update document purpose", "metadata", i, "error", err)
							return err
						}
					}

					// Add tags to document
					for j, tag := range metadata.Tags {
						err = db.UpsertSectionTag(ctx, sqlite.UpsertSectionTagParams{
							SourceID:   sourceID,
							DocumentID: section.DocumentID,
							SectionID:  section.ID,
							TagKey:     tag.Key,
							TagValue:   tag.Value,
						})
						if err != nil {
							slog.Debug("extract.addMetadata could not insert section tag", "metadata", i, "tag", j, "error", err)
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
