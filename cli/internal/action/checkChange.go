package action

import (
	"database/sql"
	"errors"
	"hyaline/internal/check"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type CheckChangeArgs struct {
	Config  string
	Current string
	Change  string
	System  string
	Output  string
}

func CheckChange(args *CheckChangeArgs) error {
	slog.Info("Checking changed code and docs")
	slog.Debug("action.CheckChange Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"change", args.Change,
		"system", args.System,
		"output", args.Output,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckChange could not load the config", "error", err)
		return err
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckChange could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckChange detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open Current DB
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.CheckChange could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.CheckChange could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.CheckChange opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Open Change DB
	var changeDB *sql.DB
	if args.Change != "" {
		changeAbsPath, err := filepath.Abs(args.Change)
		if err != nil {
			slog.Debug("action.CheckChange could not get an absolute path for change", "change", args.Change, "error", err)
			return err
		}
		changeDB, err = sql.Open("sqlite", changeAbsPath)
		if err != nil {
			slog.Debug("action.CheckChange could not open change SQLite DB", "dataSourceName", changeAbsPath, "error", err)
			return err
		}
		slog.Debug("action.CheckChange opened change database", "change", args.Change, "path", changeAbsPath)
		defer changeDB.Close()
	}

	// Get system
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.CheckChange could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Get the full set of ruleDocuments that apply to this system
	rules := map[string]*config.Rule{}
	for _, doc := range system.Docs {
		for _, ruleID := range doc.Rules {
			rules[ruleID] = config.GetRule(cfg.Rules, ruleID)
		}
	}
	ruleDocs := []config.RuleDocument{}
	for _, rule := range rules {
		ruleDocs = append(ruleDocs, rule.Documents...)
	}

	// Get the set of documents/sections that need to be updated for each code change
	// TODO
	for _, c := range system.Code {
		files, err := sqlite.GetAllFiles(c.ID, system.ID, changeDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get files for codeSource", "codeSource", c.ID, "system", args.System, "error", err)
			return err
		}
		for _, file := range files {
			// TODO add documents/sections back to our master list
			check.Change(file, ruleDocs)
		}
	}

	// Merge sets of files into a master list
	// TODO

	// Loop through rules and respect any updateIfs
	// TODO

	// Loop through documents that have been updated and annotate those on the list
	// TODO

	// Output the results
	// TODO

	return nil
}
