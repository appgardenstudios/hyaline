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
						Required: false,
						Usage:    "Base branch (where changes will be applied). Either --base or --base-ref must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "head",
						Required: false,
						Usage:    "Head branch (which changes will be applied). Either --head or --head-ref must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "base-ref",
						Required: false,
						Usage:    "Base reference (explicit commit hash or fully qualified reference). Either --base-ref or --base must be provided, but not both.",
					},
					&cli.StringFlag{
						Name:     "head-ref",
						Required: false,
						Usage:    "Head reference (explicit commit hash or fully qualified reference). Either --head-ref or --head must be provided, but not both.",
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
					err := action.ExtractChange(&action.ExtractChangeArgs{
						Config:           cCtx.String("config"),
						System:           cCtx.String("system"),
						Base:             base,
						Head:             head,
						BaseRef:          baseRef,
						HeadRef:          headRef,
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
			{
				Name:  "documentation",
				Usage: "Extract documentation into a current data set",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
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
					err := action.ExtractDocumentation(&action.ExtractDocumentationArgs{
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
