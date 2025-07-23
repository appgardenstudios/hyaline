package action

import (
	"database/sql"
	"hyaline/internal/serve/mcp"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ServeMCPArgs struct {
	Documentation string
}

func ServeMCP(args *ServeMCPArgs, version string) error {
	slog.Info("Starting MCP server")
	slog.Debug("action.ServeMCP Args", slog.Group("args",
		"documentation", args.Documentation,
		"version", version,
	))

	// Get absolute path to database
	absPath, err := filepath.Abs(args.Documentation)
	if err != nil {
		slog.Debug("action.ServeMCP could not get an absolute path for documentation", "documentation", args.Documentation, "error", err)
		return err
	}

	// Open SQLite database
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.ServeMCP could not open SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()

	// Create and start MCP server
	server, err := mcp.NewServer(db, version)
	if err != nil {
		slog.Debug("action.ServeMCP could not create MCP server", "error", err)
		return err
	}

	// Start server using stdio transport
	err = server.ServeStdio()
	if err != nil {
		slog.Debug("action.ServeMCP server error", "error", err)
		return err
	}

	slog.Info("MCP server stopped")
	return nil
}