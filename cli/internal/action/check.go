package action

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hyaline/internal/check"
	"hyaline/internal/config"
	"hyaline/internal/rule"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type CheckArgs struct {
	Config    string
	Current   string
	Change    string
	System    string
	Recommend bool
}

func Check(args *CheckArgs) error {
	slog.Info("Checking current docs")
	slog.Debug("action.Check Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"change", args.Change,
		"system", args.System,
		"recommend", args.Recommend,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.Check could not load the config", "error", err)
		return err
	}

	// Open current data set database
	absPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.Check could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.Check could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer currentDB.Close()

	// Open change data set database (if passed in)
	// var changeDB *sql.DB
	if args.Change != "" {
		// TODO load change db
	}

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.Check could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Run checks and count failures
	results := []*rule.Result{}
	failed := 0
	for _, c := range system.Checks {
		slog.Info("Running check " + c.ID)
		result, err := check.Run(c, system.ID, currentDB, args.Recommend, cfg.LLM)
		if err != nil {
			slog.Debug("action.Check could not run", "check", c.ID, "error", err)
			return err
		}
		if !result.Pass {
			failed++
		}
		results = append(results, result)
	}

	// Print out checks
	data := struct {
		Results []*rule.Result `json:"results"`
	}{results}
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Debug("action.Check could not marshal results", "error", err)
		return err
	}
	fmt.Println(string(output))

	// If >0 failed, return an error so the program error code != 0
	if failed > 0 {
		return fmt.Errorf("%d checks failed", failed)
	}

	return nil
}
