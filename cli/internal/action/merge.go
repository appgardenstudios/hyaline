package action

import (
	"database/sql"
	"errors"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type MergeArgs struct {
	Inputs []string
	Output string
}

func Merge(args *MergeArgs) error {
	slog.Info("Merging data sets")
	slog.Debug("action.Merge Args", slog.Group("args",
		"inputs", args.Inputs,
		"output", args.Output,
	))

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.Merge could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		slog.Debug("action.Merge detected that output db already exists", "absPath", absPath)
		return errors.New("output file already exists")
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.Merge could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()
	err = sqlite.CreateSchema(db)
	if err != nil {
		slog.Debug("action.Merge could not create the current schema", "error", err)
		return err
	}

	// Merge Systems

	return nil
}
