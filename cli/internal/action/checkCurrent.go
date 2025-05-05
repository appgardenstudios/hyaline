package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type CheckCurrentArgs struct {
	Config  string
	Current string
	System  string
	Output  string
}

func CheckCurrent(args *CheckCurrentArgs) error {
	slog.Info("Checking current docs")
	slog.Debug("action.CheckCurrent Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
		"output", args.Output,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckCurrent could not load the config", "error", err)
		return err
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckCurrent could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckCurrent detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open Current DB
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.CheckCurrent could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.CheckCurrent could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.CheckCurrent opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Get system
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.CheckChange could not locate the system", "system", args.System, "error", err)
		return err
	}

	fmt.Println(system.ID)

	return nil
}
