package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Update(logLevel *slog.LevelVar) *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update various items such as pull requests",
		Subcommands: []*cli.Command{
			{
				Name:  "pr",
				Usage: "Update a GitHub PR",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "pull-request",
						Required: true,
						Usage:    "GitHub Pull Request to use (OWNER/REPO/PR_NUMBER)",
					},
					&cli.StringFlag{
						Name:     "comment",
						Required: false,
						Usage:    "GitHub Pull Request Comment to update (OWNER/REPO/COMMENT_NUMBER)",
					},
					&cli.StringFlag{
						Name:     "sha",
						Required: false,
						Usage:    "SHA to add to the comment",
					},
					&cli.StringFlag{
						Name:     "recommendations",
						Required: true,
						Usage:    "Path to the recommendations to use (output of check change)",
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
					err := action.UpdatePR(&action.UpdatePRArgs{
						Config:          cCtx.String("config"),
						PullRequest:     cCtx.String("pull-request"),
						Comment:         cCtx.String("comment"),
						Sha:             cCtx.String("sha"),
						Recommendations: cCtx.String("recommendations"),
						Output:          cCtx.String("output"),
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
