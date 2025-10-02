package action

import (
	"fmt"
	"hyaline/internal/github"
	"hyaline/internal/io"
	"hyaline/internal/serve/mcp"
	"hyaline/internal/serve/mcp/utils"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type ServeMCPArgs struct {
	Documentation      string
	GitHubRepo         string
	GitHubArtifact     string
	GitHubArtifactPath string
	GitHubToken        string
}

func ServeMCP(args *ServeMCPArgs, version string) error {
	slog.Info("Starting MCP server")
	slog.Debug("action.ServeMCP Args", slog.Group("args",
		"documentation", args.Documentation,
		"githubRepo", args.GitHubRepo,
		"githubArtifact", args.GitHubArtifact,
		"githubArtifactPath", args.GitHubArtifactPath,
		"version", version,
	))

	// Get absolute path to database
	var absPath string
	var err error
	var tempDir string

	if args.GitHubRepo != "" {
		// Create a temporary directory to store the documentation
		tempDir, err = os.MkdirTemp("", "hyaline-docs-*")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(tempDir)

		zipPath, err := github.DownloadLatestArtifact(args.GitHubRepo, args.GitHubArtifact, args.GitHubToken, tempDir)
		if err != nil {
			return fmt.Errorf("failed to download artifact: %w", err)
		}

		// Unzip the artifact
		unzipDir := filepath.Join(tempDir, "unzipped")
		err = io.Unzip(zipPath, unzipDir)
		if err != nil {
			return fmt.Errorf("failed to unzip artifact: %w", err)
		}

		// Join the unzipped directory with the artifact path
		absPath = filepath.Join(unzipDir, args.GitHubArtifactPath)
	} else {
		absPath, err = filepath.Abs(args.Documentation)
		if err != nil {
			slog.Debug("action.ServeMCP could not get an absolute path for documentation", "documentation", args.Documentation, "error", err)
			return err
		}
	}

	db, close, err := sqlite.InitInput(absPath)
	if err != nil {
		slog.Debug("action.ServeMCP could not initialize SQLite DB", "path", absPath, "error", err)
		return err
	}
	defer close()

	// Create and start MCP server
	server, err := mcp.NewServer(db, version, utils.ServerOptions{
		GitHubRepo:         args.GitHubRepo,
		GitHubArtifact:     args.GitHubArtifact,
		GitHubArtifactPath: args.GitHubArtifactPath,
		GitHubToken:        args.GitHubToken,
		DocumentationPath:  args.Documentation,
	})
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
