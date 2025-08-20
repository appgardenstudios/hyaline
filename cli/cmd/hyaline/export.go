package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Export(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: "Export documentation",
		Subcommands: []*cli.Command{
			{
				Name:  "documentation",
				Usage: "Export documentation into a variety of formats",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "documentation",
						Required: true,
						Usage:    "Path to the current documentation data set",
					},
					&cli.StringFlag{
						Name:     "format",
						Required: true,
						Usage:    "Format to use when exporting (one of fs, llmsfulltxt, json, sqlite)",
					},
					&cli.StringSliceFlag{
						Name:     "include",
						Required: false,
						Usage:    "Document URIs to include. Accepts multiple includes by setting multiple times.",
					},
					&cli.StringSliceFlag{
						Name:     "exclude",
						Required: false,
						Usage:    "Document URIs to include. Accepts multiple includes by setting multiple times.",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path to use when exporting the documentation",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ExportDocumentation(&action.ExportDocumentationArgs{
						Documentation: cCtx.String("documentation"),
						Format:        cCtx.String("format"),
						Includes:      cCtx.StringSlice("include"),
						Excludes:      cCtx.StringSlice("exclude"),
						Output:        cCtx.String("output"),
					})
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
		},
	}
}
