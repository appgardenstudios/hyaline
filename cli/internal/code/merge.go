package code

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func Merge(systemID string, from *sql.DB, to *sql.DB) error {
	codeToCopy, err := sqlite.GetAllCode(systemID, from)
	if err != nil {
		return err
	}

	for _, code := range codeToCopy {
		// Delete any existing CODE and FILE entries
		err := sqlite.DeleteCodeAndFiles(code.ID, systemID, to)
		if err != nil {
			return err
		}

		// Copy CODE
		code, err := sqlite.GetCode(code.ID, systemID, from)
		if err != nil {
			return err
		}
		if code != nil {
			err = sqlite.InsertCode(*code, to)
			if err != nil {
				return err
			}
		}

		// Copy FILEs
		files, err := sqlite.GetAllFiles(code.ID, systemID, from)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = sqlite.InsertFile(*file, to)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
