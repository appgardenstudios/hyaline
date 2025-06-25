package utils

import (
	"database/sql"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
)

// DocumentationData holds all documentation data in memory for fast access
type DocumentationData struct {
	Systems       map[string]*sqlite.System                                // systemID -> System
	Documentation map[string]map[string]*sqlite.SystemDocumentation        // systemID -> docID -> Documentation
	Documents     map[string]map[string]map[string]*sqlite.SystemDocument  // systemID -> docID -> documentID -> Document
	Sections      map[string]map[string]map[string][]*sqlite.SystemSection // systemID -> docID -> documentID -> []*Section
}

// LoadAllData loads all documentation data from the database into memory
func LoadAllData(db *sql.DB) (*DocumentationData, error) {
	slog.Debug("data.LoadAllData starting")

	data := &DocumentationData{
		Systems:       make(map[string]*sqlite.System),
		Documentation: make(map[string]map[string]*sqlite.SystemDocumentation),
		Documents:     make(map[string]map[string]map[string]*sqlite.SystemDocument),
		Sections:      make(map[string]map[string]map[string][]*sqlite.SystemSection),
	}

	// Load systems
	systems, err := sqlite.GetAllSystem(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load systems: %w", err)
	}

	for _, sys := range systems {
		data.Systems[sys.ID] = sys
		data.Documentation[sys.ID] = make(map[string]*sqlite.SystemDocumentation)
		data.Documents[sys.ID] = make(map[string]map[string]*sqlite.SystemDocument)
		data.Sections[sys.ID] = make(map[string]map[string][]*sqlite.SystemSection)

		// Load documentation for this system
		docs, err := sqlite.GetAllSystemDocumentation(sys.ID, db)
		if err != nil {
			return nil, fmt.Errorf("failed to load documentation for system %s: %w", sys.ID, err)
		}

		for _, doc := range docs {
			data.Documentation[sys.ID][doc.ID] = doc
			data.Documents[sys.ID][doc.ID] = make(map[string]*sqlite.SystemDocument)
			data.Sections[sys.ID][doc.ID] = make(map[string][]*sqlite.SystemSection)

			// Load documents for this documentation
			documents, err := sqlite.GetAllSystemDocument(doc.ID, sys.ID, db)
			if err != nil {
				return nil, fmt.Errorf("failed to load documents for documentation %s: %w", doc.ID, err)
			}

			for _, document := range documents {
				data.Documents[sys.ID][doc.ID][document.ID] = document

				// Load sections for this document
				sections, err := sqlite.GetAllSystemSectionsForDocument(document.ID, doc.ID, sys.ID, db)
				if err != nil {
					return nil, fmt.Errorf("failed to load sections for document %s: %w", document.ID, err)
				}

				data.Sections[sys.ID][doc.ID][document.ID] = sections
			}
		}
	}

	slog.Debug("data.LoadAllData complete")
	return data, nil
}
