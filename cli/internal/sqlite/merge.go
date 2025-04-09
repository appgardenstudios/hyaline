package sqlite

import (
	"database/sql"
	"log/slog"
)

func DeleteCodeAndFiles(codeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  CODE
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		slog.Debug("sqlite.DeleteCodeAndFiles delete code statement failure")
		return err
	}

	_, err = stmt.Exec(codeID, systemID)
	if err != nil {
		slog.Debug("sqlite.DeleteCodeAndFiles delete code exec failure")
		return err
	}

	stmt, err = db.Prepare(`
DELETE FROM
  FILE
WHERE
  CODE_ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		slog.Debug("sqlite.DeleteCodeAndFiles delete file statement failure")
		return err
	}

	_, err = stmt.Exec(codeID, systemID)
	if err != nil {
		slog.Debug("sqlite.DeleteCodeAndFiles delete file exec failure")
		return err
	}

	return nil
}
