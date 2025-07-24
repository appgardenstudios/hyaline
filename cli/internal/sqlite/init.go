package sqlite

import (
	"database/sql"
	_ "embed"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func InitOutput(outputPath string) (q *Queries, err error) {
	// Get absolute path
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		slog.Debug("sqlite.InitOutput could not get an absolute path for output", "output", outputPath, "error", err)
		return
	}

	// Ensure output path does not exist as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		err = errors.New("output file already exists")
		slog.Debug("sqlite.InitOutput detected that output db already exists", "absPath", absPath)
		return
	}

	// Open db
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("sqlite.InitOutput could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return
	}

	// Create tables
	_, err = db.Exec(schema)
	if err != nil {
		slog.Debug("sqlite.InitOutput could not create tables", "error", err)
		return
	}

	// Create sqlc queries struct
	q = New(db)

	return
}

func InitInput(inputPath string) (q *Queries, err error) {
	// Get absolute path
	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		slog.Debug("sqlite.InitInput could not get an absolute path for input", "input", inputPath, "error", err)
		return
	}

	// Check if input file exists
	if _, err = os.Stat(absPath); err != nil {
		slog.Debug("sqlite.InitInput input file does not exist", "input", inputPath, "error", err)
		err = errors.New("input file does not exist")
		return
	}

	// Open db
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("sqlite.InitInput could not open input SQLite DB", "dataSourceName", absPath, "error", err)
		return
	}

	// Create sqlc queries struct
	q = New(db)

	return
}
