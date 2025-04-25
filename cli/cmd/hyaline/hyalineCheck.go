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
				Required: false,
				Usage:    "Path to the change data set",
			},
			&cli.StringFlag{
				Name:     "system",
				Required: true,
				Usage:    "ID of the system to extract",
			},
			&cli.BoolFlag{
				Name:  "recommend",
				Usage: "Include a recommended action when a check does not pass",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// Set log level
			if cCtx.Bool("debug") {
				logLevel.Set(slog.LevelDebug)
			}

			// Execute action
			err := action.Check(&action.CheckArgs{
				Config:    cCtx.String("config"),
				Current:   cCtx.String("current"),
				Change:    cCtx.String("change"),
				System:    cCtx.String("system"),
				Recommend: cCtx.Bool("recommend"),
			})
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			return nil
		},
	}
}
