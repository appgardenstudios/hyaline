package action

import (
	"context"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
)

type MergeDocumentationArgs struct {
	Inputs []string
	Output string
}

func MergeDocumentation(args *MergeDocumentationArgs) error {
	slog.Info("Merging documentation", "inputs", args.Inputs, "output", args.Output)

	// Initialize output database
	outputDB, db, err := sqlite.InitOutput(args.Output)
	if err != nil {
		slog.Debug("action.MergeDocumentation could not initialize output", "error", err)
		return err
	}
	defer db.Close()

	slog.Info(fmt.Sprintf("Merging %d data sets", len(args.Inputs)))

	ctx := context.Background()

	// Process each input database
	for i, input := range args.Inputs {
		slog.Info(fmt.Sprintf("Merging %d of %d", i+1, len(args.Inputs)))

		// Initialize input database
		inputDB, err := sqlite.InitInput(input)
		if err != nil {
			slog.Debug("action.MergeDocumentation could not initialize input", "input", input, "error", err)
			return err
		}

		// Get all sources from input database
		sources, err := inputDB.GetAllSources(ctx)
		if err != nil {
			slog.Debug("action.MergeDocumentation could not get sources", "input", input, "error", err)
			return err
		}

		// Process each source
		for _, source := range sources {
			slog.Debug(fmt.Sprintf("Processing source %s from %s", source.ID, input))

			// Delete all existing records for this source in output database
			err = deleteSourceData(ctx, outputDB, source.ID)
			if err != nil {
				slog.Debug("action.MergeDocumentation could not delete existing source data", "sourceID", source.ID, "error", err)
				return err
			}

			// Copy documents and related data
			err = copySourceData(ctx, inputDB, outputDB, source)
			if err != nil {
				slog.Debug("action.MergeDocumentation could not copy source data", "sourceID", source.ID, "error", err)
				return err
			}
		}
	}

	slog.Info(fmt.Sprintf("Merged %d data sets", len(args.Inputs)))
	return nil
}

func deleteSourceData(ctx context.Context, db *sqlite.Queries, sourceID string) error {
	// Delete in reverse order of foreign key dependencies
	if err := db.DeleteSectionTagsForSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete section tags: %w", err)
	}
	if err := db.DeleteSectionsForSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete sections: %w", err)
	}
	if err := db.DeleteDocumentTagsForSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete document tags: %w", err)
	}
	if err := db.DeleteDocumentsForSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	if err := db.DeleteSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}
	return nil
}

func copySourceData(ctx context.Context, inputDB, outputDB *sqlite.Queries, source sqlite.SOURCE) error {
	// Copy source record
	err := outputDB.InsertSource(ctx, sqlite.InsertSourceParams{
		ID:          source.ID,
		Description: source.Description,
		Crawler:     source.Crawler,
		Root:        source.Root,
	})
	if err != nil {
		slog.Debug("action.MergeDocumentation could not insert source", "sourceID", source.ID, "error", err)
		return err
	}

	// Get and copy all documents for the source
	documents, err := inputDB.GetDocumentsForSource(ctx, source.ID)
	if err != nil {
		return fmt.Errorf("could not get documents for source %s: %w", source.ID, err)
	}

	for _, doc := range documents {
		err = outputDB.InsertDocument(ctx, sqlite.InsertDocumentParams{
			ID:            doc.ID,
			SourceID:      doc.SourceID,
			Type:          doc.Type,
			Purpose:       doc.Purpose,
			RawData:       doc.RawData,
			ExtractedData: doc.ExtractedData,
		})
		if err != nil {
			return fmt.Errorf("could not insert document %s: %w", doc.ID, err)
		}
	}

	// Get and copy all document tags for the source
	docTags, err := inputDB.GetAllDocumentTagsForSource(ctx, source.ID)
	if err != nil {
		return fmt.Errorf("could not get document tags for source %s: %w", source.ID, err)
	}

	for _, tag := range docTags {
		err = outputDB.UpsertDocumentTag(ctx, sqlite.UpsertDocumentTagParams{
			SourceID:   source.ID,
			DocumentID: tag.DocumentID,
			TagKey:     tag.TagKey,
			TagValue:   tag.TagValue,
		})
		if err != nil {
			return fmt.Errorf("could not insert document tag for %s: %w", tag.DocumentID, err)
		}
	}

	// Get and copy all sections for the source
	sections, err := inputDB.GetAllSectionsForSource(ctx, source.ID)
	if err != nil {
		return fmt.Errorf("could not get sections for source %s: %w", source.ID, err)
	}

	for _, section := range sections {
		err = outputDB.InsertSection(ctx, sqlite.InsertSectionParams{
			ID:            section.ID,
			DocumentID:    section.DocumentID,
			SourceID:      section.SourceID,
			ParentID:      section.ParentID,
			PeerOrder:     section.PeerOrder,
			Name:          section.Name,
			Purpose:       section.Purpose,
			ExtractedData: section.ExtractedData,
		})
		if err != nil {
			return fmt.Errorf("could not insert section %s: %w", section.ID, err)
		}
	}

	// Get and copy all section tags for the source
	sectionTags, err := inputDB.GetAllSectionTagsForSource(ctx, source.ID)
	if err != nil {
		return fmt.Errorf("could not get section tags for source %s: %w", source.ID, err)
	}

	for _, tag := range sectionTags {
		err = outputDB.UpsertSectionTag(ctx, sqlite.UpsertSectionTagParams{
			SourceID:   source.ID,
			DocumentID: tag.DocumentID,
			SectionID:  tag.SectionID,
			TagKey:     tag.TagKey,
			TagValue:   tag.TagValue,
		})
		if err != nil {
			return fmt.Errorf("could not insert section tag for %s: %w", tag.SectionID, err)
		}
	}

	return nil
}
