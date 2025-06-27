package utils

import (
	"database/sql"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"sort"
)

// System holds a system and its associated documentation
type System struct {
	System        *sqlite.System
	Documentation []Documentation
}

// Documentation holds documentation and its associated documents
type Documentation struct {
	Documentation *sqlite.SystemDocumentation
	Documents     []Document
}

// Document holds a document and its associated sections
type Document struct {
	Document *sqlite.SystemDocument
	Sections []*sqlite.SystemSection
}

// DocumentationData holds all documentation data in memory for fast access
type DocumentationData struct {
	Systems []System
}

// LoadAllData loads all documentation data from the database into memory
func LoadAllData(db *sql.DB) (*DocumentationData, error) {
	slog.Debug("data.LoadAllData starting")

	// Load systems
	sqliteSystems, err := sqlite.GetAllSystem(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load systems: %w", err)
	}

	var systems []System
	for _, sys := range sqliteSystems {
		system := System{
			System: sys,
		}

		// Load documentation for this system
		docs, err := sqlite.GetAllSystemDocumentation(sys.ID, db)
		if err != nil {
			return nil, fmt.Errorf("failed to load documentation for system %s: %w", sys.ID, err)
		}

		for _, doc := range docs {
			documentation := Documentation{
				Documentation: doc,
			}

			// Load documents for this documentation
			documents, err := sqlite.GetAllSystemDocument(doc.ID, sys.ID, db)
			if err != nil {
				return nil, fmt.Errorf("failed to load documents for documentation %s: %w", doc.ID, err)
			}

			for _, document := range documents {
				// Load sections for this document
				sections, err := sqlite.GetAllSystemSectionsForDocument(document.ID, doc.ID, sys.ID, db)
				if err != nil {
					return nil, fmt.Errorf("failed to load sections for document %s: %w", document.ID, err)
				}

				doc := Document{
					Document: document,
					Sections: sections,
				}
				documentation.Documents = append(documentation.Documents, doc)
			}

			// Sort documents alphabetically by ID
			sort.Slice(documentation.Documents, func(i, j int) bool {
				return documentation.Documents[i].Document.ID < documentation.Documents[j].Document.ID
			})

			system.Documentation = append(system.Documentation, documentation)
		}

		// Sort documentation alphabetically by ID
		sort.Slice(system.Documentation, func(i, j int) bool {
			return system.Documentation[i].Documentation.ID < system.Documentation[j].Documentation.ID
		})

		systems = append(systems, system)
	}

	// Sort systems alphabetically by ID
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].System.ID < systems[j].System.ID
	})

	data := &DocumentationData{
		Systems: systems,
	}

	slog.Debug("data.LoadAllData complete")
	return data, nil
}
