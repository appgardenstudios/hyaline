package utils

import (
	"context"
	"database/sql"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"sort"
)

// Tags represents a collection of tag keys with their associated values
type Tags map[string][]string

// Add adds a value to the specified tag key
func (t Tags) Add(key, value string) {
	t[key] = append(t[key], value)
}

// Keys returns all tag keys in sorted order for deterministic iteration
func (t Tags) Keys() []string {
	keys := make([]string, 0, len(t))
	for key := range t {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// NewTags creates a new Tags instance
func NewTags() Tags {
	return make(Tags)
}

type Source struct {
	sqlite.SOURCE
	Documents []Document
}

type Document struct {
	sqlite.DOCUMENT
	Tags     Tags
	Sections []Section
}

type Section struct {
	sqlite.SECTION
	Tags Tags
}

// DocumentationData holds all documentation data in memory for fast access
type DocumentationData struct {
	Sources []Source
}

// LoadAllData loads all documentation data from the database into memory
func LoadAllData(db *sql.DB) (*DocumentationData, error) {
	slog.Debug("serve.mcp.data.LoadAllData starting")

	ctx := context.Background()
	queries := sqlite.New(db)

	// Load sources
	sqliteSources, err := queries.GetAllSources(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load sources: %w", err)
	}

	sources := make([]Source, 0, len(sqliteSources))

	// For each source, load documents
	for _, sqliteSource := range sqliteSources {
		source := Source{
			SOURCE: sqliteSource,
		}

		// Load documents for this source
		sqliteDocuments, err := queries.GetDocumentsForSource(ctx, source.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load documents for source %s: %w", source.ID, err)
		}

		documents := make([]Document, 0, len(sqliteDocuments))

		// For each document, load sections and tags
		for _, sqliteDoc := range sqliteDocuments {
			document := Document{
				DOCUMENT: sqliteDoc,
				Tags:     NewTags(),
			}

			// Load document tags
			documentTags, err := queries.GetDocumentTags(ctx, sqlite.GetDocumentTagsParams{
				SourceID:   source.ID,
				DocumentID: document.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to load tags for document %s: %w", document.ID, err)
			}

			// Group tags by key
			for _, tag := range documentTags {
				document.Tags.Add(tag.TagKey, tag.TagValue)
			}

			// Load sections
			sqliteSections, err := queries.GetSectionsForDocument(ctx, sqlite.GetSectionsForDocumentParams{
				SourceID:   source.ID,
				DocumentID: document.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to load sections for document %s: %w", document.ID, err)
			}

			sections := make([]Section, 0, len(sqliteSections))

			// Load section tags
			for _, sqliteSection := range sqliteSections {
				section := Section{
					SECTION: sqliteSection,
					Tags:    NewTags(),
				}

				// Load section tags
				sectionTags, err := queries.GetSectionTags(ctx, sqlite.GetSectionTagsParams{
					SourceID:   source.ID,
					DocumentID: document.ID,
					SectionID:  section.ID,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to load tags for section %s: %w", section.ID, err)
				}

				// Group tags by key
				for _, tag := range sectionTags {
					section.Tags.Add(tag.TagKey, tag.TagValue)
				}

				sections = append(sections, section)
			}

			// Sections are already sorted by PEER_ORDER, ID from the query
			document.Sections = sections
			documents = append(documents, document)
		}

		// Documents are already sorted by ID from the query
		source.Documents = documents
		sources = append(sources, source)
	}

	// Sources are already sorted by ID from the query
	data := &DocumentationData{
		Sources: sources,
	}

	slog.Debug("serve.mcp.data.LoadAllData complete", "sourceCount", len(sources))
	return data, nil
}
