package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Generate(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate various items such as configuration",
		Subcommands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Generate configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "current",
						Required: true,
						Usage:    "Path to the current sqlite database",
					},
					&cli.StringFlag{
						Name:     "system",
						Required: true,
						Usage:    "ID of the system to extract",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path of the sqlite database to create",
					},
					&cli.BoolFlag{
						Name:  "include-purpose",
						Usage: "If set, will call an LLM to generate the document/section purpose",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.GenerateConfig(&action.GenerateConfigArgs{
						Config:         cCtx.String("config"),
						Current:        cCtx.String("current"),
						System:         cCtx.String("system"),
						Output:         cCtx.String("output"),
						IncludePurpose: cCtx.Bool("include-purpose"),
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
