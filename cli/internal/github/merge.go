package github

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func MergeChanges(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all CHANGE entries from the source we will be copying
	changeToCopy, err := sqlite.GetAllSystemChange(systemID, source)
	if err != nil {
		return err
	}

	// Copy each CHANGE from source to dest
	for _, p := range changeToCopy {
		// Delete the existing CHANGE (if any)
		err := sqlite.DeleteSystemChange(p.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy CHANGE
		err = sqlite.InsertSystemChange(*p, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func MergeTasks(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all TASK entries from the source we will be copying
	tasksToCopy, err := sqlite.GetAllSystemTask(systemID, source)
	if err != nil {
		return err
	}

	// Copy each TASK from source to dest
	for _, i := range tasksToCopy {
		// Delete the existing TASK (if any)
		err := sqlite.DeleteSystemTask(i.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy TASK
		err = sqlite.InsertSystemTask(*i, dest)
		if err != nil {
			return err
		}
	}

	return nil
}
