package action

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type MergeDocumentationArgs struct {
	Inputs []string
	Output string
}

func MergeDocumentation(args *MergeDocumentationArgs) error {
	slog.Info(fmt.Sprintf("Merging documentation in this order: %v with output: %s", args.Inputs, args.Output))

	// Initialize output database  
	outputDB, err := sqlite.InitOutput(args.Output)
	if err != nil {
		slog.Debug("action.MergeDocumentation could not initialize output", "error", err)
		return err
	}

	slog.Info(fmt.Sprintf("Merging %d data sets", len(args.Inputs)))

	ctx := context.Background()

	// Process each input database
	for i, input := range args.Inputs {
		slog.Info(fmt.Sprintf("Merging %d of %d", i+1, len(args.Inputs)))

		// Open input database
		inputAbsPath, err := filepath.Abs(input)
		if err != nil {
			slog.Debug("action.MergeDocumentation could not get an absolute path for input", "input", input, "error", err)
			return err
		}

		// Check if input file exists
		if _, err := os.Stat(inputAbsPath); err != nil {
			slog.Debug("action.MergeDocumentation input file does not exist", "input", input, "error", err)
			return fmt.Errorf("input file does not exist: %s", input)
		}

		inputDB, err := sql.Open("sqlite", inputAbsPath)
		if err != nil {
			slog.Debug("action.MergeDocumentation could not open input SQLite DB", "dataSourceName", inputAbsPath, "error", err)
			return err
		}
		defer inputDB.Close()

		inputQueries := sqlite.New(inputDB)

		// Get all sources from input database
		sources, err := inputQueries.GetAllSources(ctx)
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

			// Copy source record
			err = outputDB.InsertSource(ctx, sqlite.InsertSourceParams{
				ID:          source.ID,
				Description: source.Description,
				Crawler:     source.Crawler,
				Root:        source.Root,
			})
			if err != nil {
				slog.Debug("action.MergeDocumentation could not insert source", "sourceID", source.ID, "error", err)
				return err
			}

			// Copy documents
			documents, err := inputQueries.GetDocumentsForSource(ctx, source.ID)
			if err != nil {
				slog.Debug("action.MergeDocumentation could not get documents", "sourceID", source.ID, "error", err)
				return err
			}

			for _, doc := range documents {
				// Insert document
				err = outputDB.InsertDocument(ctx, sqlite.InsertDocumentParams{
					ID:            doc.ID,
					SourceID:      doc.SourceID,
					Type:          doc.Type,
					Purpose:       doc.Purpose,
					RawData:       doc.RawData,
					ExtractedData: doc.ExtractedData,
				})
				if err != nil {
					slog.Debug("action.MergeDocumentation could not insert document", "documentID", doc.ID, "error", err)
					return err
				}

				// Copy document tags
				docTags, err := inputQueries.GetDocumentTags(ctx, sqlite.GetDocumentTagsParams{
					SourceID:   source.ID,
					DocumentID: doc.ID,
				})
				if err != nil {
					slog.Debug("action.MergeDocumentation could not get document tags", "documentID", doc.ID, "error", err)
					return err
				}

				for _, tag := range docTags {
					err = outputDB.UpsertDocumentTag(ctx, sqlite.UpsertDocumentTagParams{
						SourceID:   source.ID,
						DocumentID: doc.ID,
						TagKey:     tag.TagKey,
						TagValue:   tag.TagValue,
					})
					if err != nil {
						slog.Debug("action.MergeDocumentation could not insert document tag", "documentID", doc.ID, "error", err)
						return err
					}
				}

				// Copy sections
				sections, err := inputQueries.GetSectionsForDocument(ctx, sqlite.GetSectionsForDocumentParams{
					SourceID:   source.ID,
					DocumentID: doc.ID,
				})
				if err != nil {
					slog.Debug("action.MergeDocumentation could not get sections", "documentID", doc.ID, "error", err)
					return err
				}

				for _, section := range sections {
					// Insert section
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
						slog.Debug("action.MergeDocumentation could not insert section", "sectionID", section.ID, "error", err)
						return err
					}

					// Copy section tags
					sectionTags, err := inputQueries.GetSectionTags(ctx, sqlite.GetSectionTagsParams{
						SourceID:   source.ID,
						DocumentID: doc.ID,
						SectionID:  section.ID,
					})
					if err != nil {
						slog.Debug("action.MergeDocumentation could not get section tags", "sectionID", section.ID, "error", err)
						return err
					}

					for _, tag := range sectionTags {
						err = outputDB.UpsertSectionTag(ctx, sqlite.UpsertSectionTagParams{
							SourceID:   source.ID,
							DocumentID: doc.ID,
							SectionID:  section.ID,
							TagKey:     tag.TagKey,
							TagValue:   tag.TagValue,
						})
						if err != nil {
							slog.Debug("action.MergeDocumentation could not insert section tag", "sectionID", section.ID, "error", err)
							return err
						}
					}
				}
			}
		}
	}

	slog.Info(fmt.Sprintf("Merged %d data sets", len(args.Inputs)))
	return nil
}

func deleteSourceData(ctx context.Context, db *sqlite.Queries, sourceID string) error {
	// Delete in reverse order of foreign key dependencies
	if err := db.DeleteSectionTagsForSource(ctx, sourceID); err != nil {
		return errors.New("failed to delete section tags: " + err.Error())
	}
	if err := db.DeleteSectionsForSource(ctx, sourceID); err != nil {
		return errors.New("failed to delete sections: " + err.Error())
	}
	if err := db.DeleteDocumentTagsForSource(ctx, sourceID); err != nil {
		return errors.New("failed to delete document tags: " + err.Error())
	}
	if err := db.DeleteDocumentsForSource(ctx, sourceID); err != nil {
		return errors.New("failed to delete documents: " + err.Error())
	}
	if err := db.DeleteSource(ctx, sourceID); err != nil {
		return errors.New("failed to delete source: " + err.Error())
	}
	return nil
}