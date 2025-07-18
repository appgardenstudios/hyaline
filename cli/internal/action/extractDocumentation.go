package action

import (
	"hyaline/internal/config"
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
	_, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.ExtractDocumentation could not load the config", "error", err)
		return err
	}

	return nil
}
