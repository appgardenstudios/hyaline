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
	mcpServer *server.MCPServer
	data      *utils.MCPData
}

// NewServer creates and initializes a new MCP server
func NewServer(db *sql.DB) (*Server, error) {
	slog.Debug("mcp.NewServer starting")

	// Load all data into memory
	mcpData, err := utils.LoadAllData(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// Create MCP server instance
	mcpServer := server.NewMCPServer(
		"Hyaline Documentation Server",
		"1.0.0",
		server.WithToolCapabilities(false), // Tools don't change dynamically
	)

	hyalineMCPServer := &Server{
		mcpServer: mcpServer,
		data:      mcpData,
	}

	// Register tools and prompts
	hyalineMCPServer.registerTools()
	hyalineMCPServer.registerPrompts()

	slog.Debug("mcp.NewServer complete")
	return hyalineMCPServer, nil
}

// ServeStdio starts the MCP server using stdio transport
func (hyalineMCPServer *Server) ServeStdio() error {
	return server.ServeStdio(hyalineMCPServer.mcpServer)
}

// registerTools registers all MCP tools
func (hyalineMCPServer *Server) registerTools() {
	// list_documents tool
	hyalineMCPServer.mcpServer.AddTool(tools.ListDocumentsTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListDocuments(ctx, request, hyalineMCPServer.data)
	})

	// get_documents tool
	hyalineMCPServer.mcpServer.AddTool(tools.GetDocumentsTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetDocuments(ctx, request, hyalineMCPServer.data)
	})
}

// registerPrompts registers all MCP prompts
func (hyalineMCPServer *Server) registerPrompts() {
	// answer_question prompt
	hyalineMCPServer.mcpServer.AddPrompt(prompts.AnswerQuestionPrompt(), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return prompts.HandleAnswerQuestion(ctx, request)
	})

}
