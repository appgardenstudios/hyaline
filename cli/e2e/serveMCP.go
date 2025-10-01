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

type ServeMCPClientOptions struct {
	DBPath         string
	GitHubRepo     string
	GitHubArtifact string
	GitHubToken    string
}

// setupServeMCPClient creates and initializes an MCP client for testing the new serve mcp command
func setupServeMCPClient(t *testing.T, opts ServeMCPClientOptions) *mcpClient.Client {
	// Build the hyaline binary path relative to this test file
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected to get current working directory: %v", err)
	}
	binaryPath := filepath.Join(dir, "../hyaline-e2e")

	// Build args based on options
	args := []string{"serve", "mcp"}

	if opts.GitHubRepo != "" {
		args = append(args, "--github-repo", opts.GitHubRepo)
		if opts.GitHubArtifact != "" {
			args = append(args, "--github-artifact", opts.GitHubArtifact)
		}
		if opts.GitHubToken != "" {
			args = append(args, "--github-token", opts.GitHubToken)
		}
	} else if opts.DBPath != "" {
		// Get absolute path to database
		absDBPath, err := filepath.Abs(opts.DBPath)
		if err != nil {
			t.Fatalf("expected to get absolute path for database: %v", err)
		}
		args = append(args, "--documentation", absDBPath)
	} else {
		t.Fatal("either DBPath or GitHubRepo must be specified")
	}

	t.Logf("Starting MCP client")
	client, err := mcpClient.NewStdioMCPClient(binaryPath, nil, args...)
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

// callServeMCPServer creates an MCP client, performs a request, and writes output to the provided path
func callServeMCPServer(t *testing.T, dbPath string, request mcp.CallToolRequest, outputPath string) {
	client := setupServeMCPClient(t, ServeMCPClientOptions{DBPath: dbPath})
	ctx := context.Background()

	response, err := client.CallTool(ctx, request)
	if err != nil {
		t.Fatalf("expected to call '%s' tool successfully: %v", request.Params.Name, err)
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent")
	}

	t.Logf("Response content: %s", textContent.Text)

	err = os.WriteFile(outputPath, []byte(textContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}
}
