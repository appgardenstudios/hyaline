package action

import (
	"log/slog"

	_ "modernc.org/sqlite"
)

type ExtractDocumentationArgs struct {
	Config string
	Output string
}

func ExtractDocumentation(args *ExtractDocumentationArgs) error {
	slog.Info("Extracting documentation", "config", args.Config, "output", args.Output)

	return nil
}
