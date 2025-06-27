package hyaline

import (
	"hyaline/internal/action"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func MCP(logLevel *slog.LevelVar, version string) *cli.Command {
	return &cli.Command{
		Name:  "mcp",
		Usage: "Model Context Protocol server for documentation access",
		Subcommands: []*cli.Command{
			{
				Name:  "stdio",
				Usage: "Start MCP server using standard I/O transport",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "current",
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
					err := action.MCPStdio(&action.MCPStdioArgs{
						Current: cCtx.String("current"),
						Version: version,
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
