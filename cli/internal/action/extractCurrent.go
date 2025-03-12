package action

import (
	"database/sql"
	"errors"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ExtractCurrentArgs struct {
	Config string
	System string
	Output string
}

func ExtractCurrent(args *ExtractCurrentArgs) error {
	slog.Info("Extracting current code and docs")
	slog.Debug("ExtractCurrent Args", slog.Group("args",
		"config", args.Config,
		"system", args.System,
		"output", args.Output,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("ExtractCurrent could not load the config", "error", err)
		return err
	}

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("ExtractCurrent could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		slog.Debug("ExtractCurrent detected that output db already exists", "absPath", absPath)
		return errors.New("output file already exists")
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("ExtractCurrent could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()
	err = sqlite.CreateCurrentSchema(db)
	if err != nil {
		slog.Debug("ExtractCurrent could not create the current schema", "error", err)
		return err
	}

	slog.Debug("ExtractCurrent starting extraction")

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("ExtractCurrent could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Insert System
	err = sqlite.InsertCurrentSystem(sqlite.CurrentSystem{
		ID: system.ID,
	}, db)
	if err != nil {
		slog.Debug("ExtractCurrent could not insert the system", "error", err)
		return err
	}
	slog.Debug("ExtractCurrent system inserted")

	// Extract/Insert Code
	err = code.ExtractCurrent(system, db)
	if err != nil {
		slog.Debug("ExtractCurrent could not extract code", "error", err)
		return err
	}
	slog.Debug("ExtractCurrent code inserted")

	// Extract/Insert Docs
	err = docs.ExtractCurrent(system, db)
	if err != nil {
		slog.Debug("ExtractCurrent could not extract docs", "error", err)
		return err
	}
	slog.Debug("ExtractCurrent docs inserted")

	slog.Info("Extraction complete")
	return nil
}
