package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Check(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Check documentation for issues and errors",
		Subcommands: []*cli.Command{
			{
				Name:  "change",
				Usage: "Extract and create a change data set",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "current",
						Required: true,
						Usage:    "Path to the current data set",
					},
					&cli.StringFlag{
						Name:     "change",
						Required: true,
						Usage:    "Path to the change data set",
					},
					&cli.StringFlag{
						Name:     "system",
						Required: true,
						Usage:    "ID of the system to check",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "Path to write the results to",
					},
					&cli.BoolFlag{
						Name:  "suggest",
						Usage: "Call the llm to suggest what edits to make to the documentation for each recommended update",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.CheckChange(&action.CheckChangeArgs{
						Config:  cCtx.String("config"),
						Current: cCtx.String("current"),
						Change:  cCtx.String("change"),
						System:  cCtx.String("system"),
						Output:  cCtx.String("output"),
						Suggest: cCtx.Bool("suggest"),
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
