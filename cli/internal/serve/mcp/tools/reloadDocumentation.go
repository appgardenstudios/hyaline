package tools

import (
	"context"
	"fmt"
	"hyaline/internal/github"
	"hyaline/internal/io"
	"hyaline/internal/serve/mcp/utils"
	"hyaline/internal/sqlite"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
)

func ReloadDocumentationTool() mcp.Tool {
	return mcp.NewTool("reload_documentation",
		mcp.WithDescription("Reload the documentation dataset."),
	)
}

func HandleReloadDocumentation(_ context.Context, request mcp.CallToolRequest, opts utils.ServerOptions) (*mcp.CallToolResult, *utils.DocumentationData, error) {
	var absPath string

	// If GitHub repository is configured, download from GitHub
	if opts.GitHubRepo != "" {
		// Check if GitHub token is configured
		if opts.GitHubToken == "" {
			return mcp.NewToolResultError("GitHub token is not configured."), nil, nil
		}

		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "hyaline-docs-reload-*")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create temp dir: %s", err.Error())), nil, nil
		}
		defer os.RemoveAll(tempDir)

		// Download latest artifact
		zipPath, err := github.DownloadLatestArtifact(opts.GitHubRepo, opts.GitHubArtifact, opts.GitHubToken, tempDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to download artifact: %s", err.Error())), nil, nil
		}

		// Unzip the artifact
		unzipDir := filepath.Join(tempDir, "unzipped")
		err = io.Unzip(zipPath, unzipDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to unzip artifact: %s", err.Error())), nil, nil
		}

		// Join the unzipped directory with the GitHub artifact path
		absPath = filepath.Join(unzipDir, opts.GitHubArtifactPath)
	} else {
		// Use the documentation path from the filesystem
		absPath = opts.DocumentationPath
	}

	// Initialize database
	db, close, err := sqlite.InitInput(absPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to initialize database: %s", err.Error())), nil, nil
	}
	defer close()

	// Load documentation data
	documentationData, err := utils.LoadAllData(db)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load data: %s", err.Error())), nil, nil
	}

	return mcp.NewToolResultText("Documentation reloaded successfully."), documentationData, nil
}
