package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Audit(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "audit",
		Usage: "Audit documentation for compliance with rules",
		Subcommands: []*cli.Command{
			{
				Name:  "documentation",
				Usage: "Audit documentation against configurable rule checks",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "documentation",
						Required: true,
						Usage:    "Path to the documentation database",
					},
					&cli.StringFlag{
						Name:     "new-flag",
						Required: true,
						Usage:    "The path to the new flag",
					},
					&cli.StringSliceFlag{
						Name:     "source",
						Required: false,
						Usage:    "Only audit specific source ID(s). Can be specified multiple times.",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path to write the audit results JSON file",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.AuditDocumentation(&action.AuditDocumentationArgs{
						Config:        cCtx.String("config"),
						Documentation: cCtx.String("documentation"),
						Sources:       cCtx.StringSlice("source"),
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
