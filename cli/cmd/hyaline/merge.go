package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Merge(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merge data sets",
		Subcommands: []*cli.Command{
			{
				Name:  "documentation",
				Usage: "Merge 2 or more documentation data sets into a single output database",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "input",
						Required: true,
						Usage:    "Path of the sqlite databases to merge. At least 2 inputs are required",
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

					// Validate at least 2 inputs
					inputs := cCtx.StringSlice("input")
					if len(inputs) < 2 {
						return cli.Exit("At least 2 input databases are required", 1)
					}

					// Execute action
					err := action.MergeDocumentation(&action.MergeDocumentationArgs{
						Inputs: inputs,
						Output: cCtx.String("output"),
					})
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
		},
		// For backwards compatibility, support the old merge command without subcommand
		// Don't mark flags as required here to avoid the validation issue
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "input",
				Usage: "Path of the sqlite database to merge",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Path of the sqlite database to create",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// If input and output are provided, run the legacy merge
			if cCtx.IsSet("input") && cCtx.IsSet("output") {
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
			}

			// Note: Technically, this isn't a subcommand, but this is the best way to get
			// a help message printed
			cli.ShowSubcommandHelp(cCtx)
			return cli.Exit("\nError: --input and --output are required", 1)
		},
	}

}
