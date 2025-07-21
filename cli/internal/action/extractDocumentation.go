package action

import (
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"

	_ "modernc.org/sqlite"
)

type ExtractDocumentationArgs struct {
	Config string
	Output string
}

func ExtractDocumentation(args *ExtractDocumentationArgs) error {
	slog.Info("Extracting documentation", "config", args.Config, "output", args.Output)

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.ExtractDocumentation could not load the config", "error", err)
		return err
	}

	// Ensure extract options are set as they are required for this action to run
	if cfg.Extract == nil {
		slog.Debug("action.ExtractDocumentation did not find extract options")
		err = errors.New("the extract documentation command requires extract options be set in the config")
		return err
	}

	// Initialize our output database
	_, err = sqlite.InitOutput(args.Output)
	if err != nil {
		slog.Debug("action.ExtractDocumentation could not initialize output", "error", err)
		return err
	}

	return nil
}
