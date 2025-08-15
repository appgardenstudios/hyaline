package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Extract(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "extract",
		Usage: "Extract code, documentation, and other metadata",
		Subcommands: []*cli.Command{
			{
				Name:  "documentation",
				Usage: "Extract documentation into a current data set",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path of the sqlite database to create",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ExtractDocumentation(&action.ExtractDocumentationArgs{
						Config: cCtx.String("config"),
						Output: cCtx.String("output"),
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
