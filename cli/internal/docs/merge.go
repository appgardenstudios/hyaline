package docs

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func Merge(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all DOCUMENTATION entries from the source we will be copying
	documentationToCopy, err := sqlite.GetAllSystemDocumentation(systemID, source)
	if err != nil {
		return err
	}

	// Copy each CODE/FILE(s) from source to dest
	for _, d := range documentationToCopy {
		// Delete any existing DOCUMENTATION, DOCUMENT, and SECTION entries
		err := sqlite.DeleteSystemDocumentation(d.ID, systemID, dest)
		if err != nil {
			return err
		}
		err = sqlite.DeleteSystemDocument(d.ID, systemID, dest)
		if err != nil {
			return err
		}
		err = sqlite.DeleteSystemSection(d.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy DOCUMENTATION
		err = sqlite.InsertSystemDocumentation(*d, dest)
		if err != nil {
			return err
		}

		// Copy DOCUMENT(s)
		documents, err := sqlite.GetAllSystemDocument(d.ID, systemID, source)
		if err != nil {
			return err
		}
		for _, document := range documents {
			err = sqlite.InsertSystemDocument(*document, dest)
			if err != nil {
				return err
			}
		}

		// Copy SECTION(s)
		sections, err := sqlite.GetAllSystemSection(d.ID, systemID, source)
		if err != nil {
			return err
		}
		for _, section := range sections {
			err = sqlite.InsertSystemSection(*section, dest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
