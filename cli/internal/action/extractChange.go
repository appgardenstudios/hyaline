package action

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type ExtractChangeArgs struct {
	Config string
	System string
	Base   string
	Head   string
	Output string
}

func ExtractChange(args *ExtractChangeArgs) error {
	slog.Info("Extracting changed code and docs")
	slog.Debug("ExtractChange Args", slog.Group("args",
		"config", args.Config,
		"system", args.System,
		"base", args.Base,
		"head", args.Head,
		"output", args.Output,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("ExtractChange could not load the config", "error", err)
		return err
	}

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("ExtractChange could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		slog.Debug("ExtractChange detected that output db already exists", "absPath", absPath)
		return errors.New("output file already exists")
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("ExtractChange could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()
	err = sqlite.CreateChangeSchema(db)
	if err != nil {
		slog.Debug("ExtractChange could not create the current schema", "error", err)
		return err
	}

	// Get/Insert System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("ExtractChange could not locate the system", "system", args.System, "error", err)
		return err
	}
	err = sqlite.InsertChangeSystem(sqlite.ChangeSystem{
		ID: system.ID,
	}, db)
	if err != nil {
		slog.Debug("ExtractChange could not insert the system", "error", err)
		return err
	}
	slog.Debug("ExtractChange system inserted")

	// Extract/Insert Code
	// TODO

	// Extract/Insert Docs
	// TODO

	slog.Info("Extraction complete")
	return nil
}
