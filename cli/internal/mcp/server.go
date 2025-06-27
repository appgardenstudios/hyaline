package mcp

import (
	"context"
	"database/sql"
	"fmt"
	"hyaline/internal/mcp/prompts"
	"hyaline/internal/mcp/tools"
	"hyaline/internal/mcp/utils"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP server with in-memory data
type Server struct {
	mcpServer         *server.MCPServer
	documentationData *utils.DocumentationData
}

// NewServer creates and initializes a new MCP server
func NewServer(db *sql.DB, version string) (*Server, error) {
	slog.Debug("mcp.NewServer starting")

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
		mcpServer:         mcpServer,
		documentationData: documentationData,
	}

	// Register tools and prompts
	hyalineMCPServer.registerTools()
	hyalineMCPServer.registerPrompts()

	slog.Debug("mcp.NewServer complete")
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
}

func (hyalineMCPServer *Server) registerPrompts() {
	hyalineMCPServer.mcpServer.AddPrompt(prompts.AnswerQuestionPrompt(), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return prompts.HandleAnswerQuestion(ctx, request)
	})
}
