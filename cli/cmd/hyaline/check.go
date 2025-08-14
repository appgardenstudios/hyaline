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
				Name:  "diff",
				Usage: "Check a diff for issues",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "documentation",
						Required: true,
						Usage:    "Path to the current documentation data set",
					},
					&cli.StringFlag{
						Name:     "path",
						Required: false,
						Usage:    "Path to the root of the repository to check",
					},
					&cli.StringFlag{
						Name:     "base",
						Required: false,
						Usage:    "Base branch (where changes will be applied). Either --base or --base-ref must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "base-ref",
						Required: false,
						Usage:    "Base reference (explicit commit hash or fully qualified reference). Either --base-ref or --base must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "head",
						Required: false,
						Usage:    "Head branch (which changes will be applied). Either --head or --head-ref must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "head-ref",
						Required: false,
						Usage:    "Head reference (explicit commit hash or fully qualified reference). Either --head-ref or --head must be provided, but not both.",
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
						Usage:    "Path to write the results to",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Helper function to show help and exit with error
					showHelpAndExit := func(message string) error {
						cli.ShowSubcommandHelp(cCtx)
						return cli.Exit("\nError: "+message, 1)
					}

					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Validate mutual exclusivity and required arguments
					base := cCtx.String("base")
					baseRef := cCtx.String("base-ref")
					head := cCtx.String("head")
					headRef := cCtx.String("head-ref")

					// Validate base arguments
					if base != "" && baseRef != "" {
						return showHelpAndExit("--base and --base-ref are mutually exclusive")
					}
					if base == "" && baseRef == "" {
						return showHelpAndExit("either --base or --base-ref is required")
					}

					// Validate head arguments
					if head != "" && headRef != "" {
						return showHelpAndExit("--head and --head-ref are mutually exclusive")
					}
					if head == "" && headRef == "" {
						return showHelpAndExit("either --head or --head-ref is required")
					}

					// Execute action
					err := action.CheckDiff(&action.CheckDiffArgs{
						Config:        cCtx.String("config"),
						Documentation: cCtx.String("documentation"),
						Path:          cCtx.String("path"),
						Base:          base,
						BaseRef:       baseRef,
						Head:          head,
						HeadRef:       headRef,
						PullRequest:   cCtx.String("pull-request"),
						Issues:        cCtx.StringSlice("issue"),
						Output:        cCtx.String("output"),
					})
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			}, {
				Name:  "pr",
				Usage: "Check a pull request for issues and add recommendations as a comment on the PR",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "documentation",
						Required: true,
						Usage:    "Path to the current documentation data set",
					},
					&cli.StringFlag{
						Name:     "pull-request",
						Required: true,
						Usage:    "GitHub Pull Request to check (OWNER/REPO/PR_NUMBER)",
					},
					&cli.StringSliceFlag{
						Name:     "issue",
						Required: false,
						Usage:    "GitHub Issue to include in the change (OWNER/REPO/ISSUE_NUMBER). Accepts multiple issues by setting multiple times.",
					},
					&cli.StringFlag{
						Name:     "output",
						Required: false,
						Usage:    "Path to write the combined (current and previous merged together) recommendations to (optional)",
					},
					&cli.StringFlag{
						Name:     "output-current",
						Required: false,
						Usage:    "Path to write the current recommendations to (optional)",
					},
					&cli.StringFlag{
						Name:     "output-previous",
						Required: false,
						Usage:    "Path to write the previous recommendations to (optional)",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.CheckPR(&action.CheckPRArgs{
						Config:         cCtx.String("config"),
						Documentation:  cCtx.String("documentation"),
						PullRequest:    cCtx.String("pull-request"),
						Issues:         cCtx.StringSlice("issue"),
						Output:         cCtx.String("output"),
						OutputCurrent:  cCtx.String("output-current"),
						OutputPrevious: cCtx.String("output-previous"),
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
