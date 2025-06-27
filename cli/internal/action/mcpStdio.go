package action

import (
	"database/sql"
	"hyaline/internal/mcp"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type MCPStdioArgs struct {
	Current string
	Version string
}

func MCPStdio(args *MCPStdioArgs) error {
	slog.Info("Starting MCP server")
	slog.Debug("action.MCPStdio Args", slog.Group("args",
		"current", args.Current,
		"version", args.Version,
	))

	// Get absolute path to database
	absPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.MCPStdio could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}

	// Open SQLite database
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.MCPStdio could not open SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()

	// Create and start MCP server
	server, err := mcp.NewServer(db, args.Version)
	if err != nil {
		slog.Debug("action.MCPStdio could not create MCP server", "error", err)
		return err
	}

	// Start server using stdio transport
	err = server.ServeStdio()
	if err != nil {
		slog.Debug("action.MCPStdio server error", "error", err)
		return err
	}

	slog.Info("MCP server stopped")
	return nil
}
