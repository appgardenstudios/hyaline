package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Export(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: "Export data from hyaline databases to various formats",
		Subcommands: []*cli.Command{
			{
				Name:  "llms-txt",
				Usage: "Export documentation to llms.txt format for AI assistants",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "current",
						Required: true,
						Usage:    "Path to the current hyaline database",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path for the output llms.txt file",
					},
					&cli.StringFlag{
						Name:     "document-uri",
						Required: false,
						Usage:    "Document URI pattern to filter the export",
					},
					&cli.BoolFlag{
						Name:     "full",
						Required: false,
						Usage:    "Generate llms-full.txt format with complete document content inline",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ExportLlmsTxt(&action.ExportLlmsTxtArgs{
						Current:     cCtx.String("current"),
						Output:      cCtx.String("output"),
						DocumentURI: cCtx.String("document-uri"),
						Full:        cCtx.Bool("full"),
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