package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// setupMCPClient creates and initializes an MCP client for testing
func setupMCPClient(t *testing.T, dbPath string) *mcpClient.Client {
	// Get absolute path to database
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		t.Fatalf("expected to get absolute path for database: %v", err)
	}

	// Build the hyaline binary path relative to this test file
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected to get current working directory: %v", err)
	}
	binaryPath := filepath.Join(dir, "../hyaline-e2e")

	// Create MCP client using stdio transport to run hyaline mcp stdio
	args := []string{
		"mcp", "stdio",
		"--current", absDBPath,
	}

	// Set coverage environment
	envVars := []string{"GOCOVERDIR=../.coverdata"}
	
	t.Logf("Starting MCP client with args: %v", args)
	client, err := mcpClient.NewStdioMCPClient(binaryPath, envVars, args...)
	if err != nil {
		t.Fatalf("expected to create MCP client successfully: %v", err)
	}

	t.Cleanup(func() {
		// Ignore broken pipe errors during cleanup as they're expected
		_ = client.Close()
	})

	// Initialize the client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := mcp.InitializeRequest{}
	request.Params.ProtocolVersion = "2024-11-05"
	request.Params.ClientInfo = mcp.Implementation{
		Name:    "hyaline-e2e-test-client",
		Version: "0.0.1",
	}

	result, err := client.Initialize(ctx, request)
	if err != nil {
		t.Fatalf("failed to initialize MCP client: %v", err)
	}
	if result.ServerInfo.Name != "Hyaline Documentation Server" {
		t.Fatalf("unexpected server name, got %s, expected Hyaline Documentation Server", result.ServerInfo.Name)
	}

	return client
}