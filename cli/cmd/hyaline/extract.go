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
				Name:  "change",
				Usage: "Extract and create a change data set",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "system",
						Required: true,
						Usage:    "ID of the system to extract",
					},
					&cli.StringFlag{
						Name:     "base",
						Required: true,
						Usage:    "Base branch (where changes will be applied)",
					},
					&cli.StringFlag{
						Name:     "head",
						Required: true,
						Usage:    "Head branch (which changes will be applied)",
					},
					&cli.StringSliceFlag{
						Name:     "code-id",
						Required: false,
						Usage:    "IDs of the code source(s) that will be extracted",
					},
					&cli.StringSliceFlag{
						Name:     "documentation-id",
						Required: false,
						Usage:    "IDs of the documentation source(s) that will be extracted",
					},
					&cli.StringFlag{
						Name:     "pull-request",
						Required: false,
						Usage:    "GitHub Pull Request to include in the change (OWNER/REPO/PR_NUMBER)",
					},
					&cli.StringSliceFlag{
						Name:     "issue",
						Required: false,
						Usage:    "GitHub Issue to include in the change (OWNER/REPO/PR_NUMBER). Accepts multiple issues by setting multiple times.",
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
					err := action.ExtractChange(&action.ExtractChangeArgs{
						Config:           cCtx.String("config"),
						System:           cCtx.String("system"),
						Base:             cCtx.String("base"),
						Head:             cCtx.String("head"),
						CodeIDs:          cCtx.StringSlice("code-id"),
						DocumentationIDs: cCtx.StringSlice("documentation-id"),
						PullRequest:      cCtx.String("pull-request"),
						Issues:           cCtx.StringSlice("issue"),
						Output:           cCtx.String("output"),
					})
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
			{
				Name:  "current",
				Usage: "Extract and create a current data set",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
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
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ExtractCurrent(&action.ExtractCurrentArgs{
						Config: cCtx.String("config"),
						System: cCtx.String("system"),
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
