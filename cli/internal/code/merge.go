package code

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func Merge(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all CODE entries from the source we will be copying
	codeToCopy, err := sqlite.GetAllCode(systemID, source)
	if err != nil {
		return err
	}

	// Copy each CODE/FILE(s) from source to dest
	for _, c := range codeToCopy {
		// Delete any existing CODE and FILE entries
		err := sqlite.DeleteCode(c.ID, systemID, dest)
		if err != nil {
			return err
		}
		err = sqlite.DeleteFile(c.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy CODE
		err = sqlite.InsertCode(*c, dest)
		if err != nil {
			return err
		}

		// Copy FILEs
		files, err := sqlite.GetAllFiles(c.ID, systemID, source)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = sqlite.InsertFile(*file, dest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
