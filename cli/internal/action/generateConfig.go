package action

import (
	"database/sql"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type GenerateConfigArgs struct {
	Config         string
	Current        string
	System         string
	Output         string
	IncludePurpose bool
}

func GenerateConfig(args *GenerateConfigArgs) error {
	slog.Info("Generating Config")
	slog.Debug("action.GenerateConfig Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
		"output", args.Output,
		"include-purpose", args.IncludePurpose,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.GenerateConfig could not load the config", "error", err)
		return err
	}

	// Open current db
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.Check could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.Check could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.Check opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.GenerateConfig could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Make a copy of the rules so we don't overwrite anything
	// TODO rules := cfg.Rules

	// Loop through docs in our current system and generate a config for each
	for _, d := range system.Docs {
		// Get a list of Documents from the db for this doc ID
		documents, err := sqlite.GetAllDocument(d.ID, system.ID, currentDB)
		if err != nil {
			slog.Debug("action.GenerateConfig could not get documents from current db", "document", d.ID, "system", system.ID, "error", err)
			return err
		}

		// Loop through each document to generate rules for it
		for _, doc := range documents {
			// Get the corresponding rules for this document
			// TODO

			// TODO if rule does not exist, create it
			// TODO check sections against the rule and create sections that don't exist
			// TODO handle adding document to new rule set or updating existing document in rules

		}

	}

	fmt.Println("system", system)
	fmt.Println("rules", cfg.Rules)

	return nil
}
