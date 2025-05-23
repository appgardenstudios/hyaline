package code

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func Merge(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all CODE entries from the source we will be copying
	codeToCopy, err := sqlite.GetAllSystemCode(systemID, source)
	if err != nil {
		return err
	}

	// Copy each CODE/FILE(s) from source to dest
	for _, c := range codeToCopy {
		// Delete any existing CODE and FILE entries
		err := sqlite.DeleteSystemCode(c.ID, systemID, dest)
		if err != nil {
			return err
		}
		err = sqlite.DeleteSystemFile(c.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy CODE
		err = sqlite.InsertSystemCode(*c, dest)
		if err != nil {
			return err
		}

		// Copy FILEs
		files, err := sqlite.GetAllSystemFiles(c.ID, systemID, source)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = sqlite.InsertSystemFile(*file, dest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
