package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func Serve(logLevel *slog.LevelVar, version string) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run a server for Hyaline",
		Subcommands: []*cli.Command{
			{
				Name:  "mcp",
				Usage: "Start MCP server using standard I/O transport",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "documentation",
						Required: false,
						Usage:    "Local filesystem path to the SQLite database containing documentation. Required when `--github-repo` is not set.",
					},
					&cli.StringFlag{
						Name:     "github-repo",
						Required: false,
						Usage:    "The path of the hyaline-github-app-config repo in GitHub (e.g. `owner/repo`). When set, downloads documentation from the repo's artifacts. Either `--documentation` or `--github-repo` is required.",
					},
					&cli.StringFlag{
						Name:  "github-artifact",
						Value: "_current-documentation",
						Usage: "The name of the documentation artifact in the hyaline-github-app-config repo",
					},
					&cli.StringFlag{
						Name:  "github-artifact-path",
						Value: "documentation.db",
						Usage: "The path to the SQLite database within the GitHub artifact",
					},
					&cli.StringFlag{
						Name: "github-token",
						// Only intended for internal testing purposes
						EnvVars:  []string{"_HYALINE_TEST_GITHUB_TOKEN"},
						Required: false,
						Usage:    "A GitHub Personal Access Token to read action artifacts from the hyaline-github-app-config repo. Required when using `--github-repo`. Consider setting this using an environment variable (e.g. `--github-token $HYALINE_SERVE_MCP_GITHUB_TOKEN`).",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Validate arguments
					if cCtx.String("documentation") == "" && cCtx.String("github-repo") == "" {
						return cli.Exit("One of --documentation or --github-repo must be specified", 1)
					}
					if cCtx.String("documentation") != "" && cCtx.String("github-repo") != "" {
						return cli.Exit("Cannot specify both --documentation and --github-repo", 1)
					}
					if cCtx.String("github-repo") != "" && cCtx.String("github-token") == "" {
						return cli.Exit("--github-token is required when using --github-repo", 1)
					}

					// Execute action
					err := action.ServeMCP(&action.ServeMCPArgs{
						Documentation:      cCtx.String("documentation"),
						GitHubRepo:         cCtx.String("github-repo"),
						GitHubArtifact:     cCtx.String("github-artifact"),
						GitHubArtifactPath: cCtx.String("github-artifact-path"),
						GitHubToken:        cCtx.String("github-token"),
					}, version)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
		},
	}
}
