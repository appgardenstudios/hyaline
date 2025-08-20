package action

import "log/slog"

type ExportDocumentationArgs struct {
	Documentation string
	Format        string
	Includes      []string
	Excludes      []string
	Output        string
}

func ExportDocumentation(args *ExportDocumentationArgs) error {
	slog.Info("Exporting Documentation",
		"documentation", args.Documentation,
		"format", args.Format,
		"includes", args.Includes,
		"excludes", args.Excludes,
		"output", args.Output,
	)

	return nil
}
