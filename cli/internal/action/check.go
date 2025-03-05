package action

import (
	"database/sql"
	"hyaline/internal/config"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type CheckArgs struct {
	Config  string
	Current string
	System  string
}

func Check(args *CheckArgs) error {
	slog.Info("Checking current docs")
	slog.Debug("Check Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("Check could not load the config", "error", err)
		return err
	}

	// Open current data set database
	absPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("Check could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("Check could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("Check could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Run checks
	for _, c := range system.Checks {
		slog.Info("Running check " + c.ID)
		// TODO run the check
	}

	return nil
}
