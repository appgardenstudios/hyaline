package action

import (
	"hyaline/internal/config"
	"log/slog"
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
	_, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("ExtractChange could not load the config", "error", err)
		return err
	}

	// Create/Scaffold SQLite
	// TODO

	// Get/Insert System
	// TODO

	// Extract/Insert Code
	// TODO

	// Extract/Insert Docs
	// TODO

	slog.Info("Extraction complete")
	return nil
}
