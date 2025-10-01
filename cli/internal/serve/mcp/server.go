package mcp

import (
	"context"
	"fmt"
	"hyaline/internal/serve/mcp/prompts"
	"hyaline/internal/serve/mcp/tools"
	"hyaline/internal/serve/mcp/utils"
	"hyaline/internal/sqlite"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP server with in-memory data
type Server struct {
	mcpServer          *server.MCPServer
	documentationData  *utils.DocumentationData
	githubRepo         string
	githubArtifact     string
	githubToken        string
	githubArtifactPath string
	filesystemDocPath  string
}

// NewServer creates and initializes a new MCP server
func NewServer(db *sqlite.Queries, version string, githubRepo string, githubArtifact string, githubArtifactPath string, githubToken string, filesystemDocPath string) (*Server, error) {
	slog.Debug("serve.mcp.NewServer starting")

	// Load all data into memory
	documentationData, err := utils.LoadAllData(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// Create MCP server instance
	mcpServer := server.NewMCPServer(
		"Hyaline Documentation Server",
		version,
		server.WithToolCapabilities(false), // Tools don't change dynamically
	)

	hyalineMCPServer := &Server{
		mcpServer:          mcpServer,
		documentationData:  documentationData,
		githubRepo:         githubRepo,
		githubArtifact:     githubArtifact,
		githubArtifactPath: githubArtifactPath,
		githubToken:        githubToken,
		filesystemDocPath:  filesystemDocPath,
	}

	// Register tools and prompts
	hyalineMCPServer.registerTools()
	hyalineMCPServer.registerPrompts()

	slog.Debug("serve.mcp.NewServer complete")
	return hyalineMCPServer, nil
}

func (hyalineMCPServer *Server) ServeStdio() error {
	return server.ServeStdio(hyalineMCPServer.mcpServer)
}

func (hyalineMCPServer *Server) registerTools() {
	hyalineMCPServer.mcpServer.AddTool(tools.ListDocumentsTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListDocuments(ctx, request, hyalineMCPServer.documentationData)
	})

	hyalineMCPServer.mcpServer.AddTool(tools.GetDocumentsTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetDocuments(ctx, request, hyalineMCPServer.documentationData)
	})

	hyalineMCPServer.mcpServer.AddTool(tools.ReloadDocumentationTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Use the appropriate documentation path based on mode
		documentationPath := hyalineMCPServer.filesystemDocPath
		if hyalineMCPServer.githubRepo != "" {
			documentationPath = hyalineMCPServer.githubArtifactPath
		}

		result, newDocumentationData, err := tools.HandleReloadDocumentation(ctx, request, hyalineMCPServer.githubRepo, hyalineMCPServer.githubArtifact, hyalineMCPServer.githubToken, documentationPath)
		if err != nil {
			return result, err
		}

		// Update the documentation data if reload was successful
		if newDocumentationData != nil {
			hyalineMCPServer.documentationData = newDocumentationData
		}

		return result, nil
	})
}

func (hyalineMCPServer *Server) registerPrompts() {
	hyalineMCPServer.mcpServer.AddPrompt(prompts.AnswerQuestionPrompt(), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return prompts.HandleAnswerQuestion(ctx, request)
	})
}
