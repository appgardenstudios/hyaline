package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Validate(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "validate",
		Usage: "Validate configuration",
		Subcommands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Validate a Hyaline configuration file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path of the config file to validate",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path of the output",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ValidateConfig(&action.ValidateConfigArgs{
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
