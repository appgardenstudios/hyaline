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
						Required: true,
						Usage:    "Path to the SQLite database containing documentation",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.ServeMCP(&action.ServeMCPArgs{
						Documentation: cCtx.String("documentation"),
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
