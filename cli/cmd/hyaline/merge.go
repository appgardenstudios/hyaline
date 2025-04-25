package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Merge(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merge 2 or more data sets into a single output database",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "input",
				Required: true,
				Usage:    "Path of the sqlite database to merge",
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
			err := action.Merge(&action.MergeArgs{
				Inputs: cCtx.StringSlice("input"),
				Output: cCtx.String("output"),
			})
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			return nil
		},
	}

}
