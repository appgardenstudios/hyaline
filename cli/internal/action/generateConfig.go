package action

import (
	"fmt"
	"hyaline/internal/config"
	"log/slog"
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
	// TODO

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.GenerateConfig could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Make a copy of the rules so we don't overwrite anything
	rules := cfg.Rules

	// Loop through docs in our current system and generate a config for each
	for _, d := range system.Docs {
		// Get a list of Documents from the db for this doc ID
		// TODO

		// Loop through each document to generate rules for it
		// TODO

		// // Get the corresponding rule for this document
		// // TODO

		// // TODO if rule does not exist, create it
		// // TODO check sections against the rule and create sections that don't exist
		// // TODO handle adding document to new rule set or updating existing document in rules
	}

	// Get documentation for system
	// TODO

	fmt.Println("system", system)
	fmt.Println("rules", cfg.Rules)

	return nil
}
