package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestServeMCPListDocumentsAll(t *testing.T) {
	// Test listing all documents with no URI parameter
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"

	goldenPath := "./_golden/serve-mcp-list-documents-all.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-all-%d.txt", time.Now().UnixMilli())

	callServeMCPServer(t, "./_input/serve-mcp/documentation.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPListDocumentsSingle(t *testing.T) {
	// Test listing a specific document
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://mcp-test/docs/index.html",
	}

	goldenPath := "./_golden/serve-mcp-list-documents-single.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-single-%d.txt", time.Now().UnixMilli())

	callServeMCPServer(t, "./_input/serve-mcp/documentation.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPListDocumentsMultiple(t *testing.T) {
	// Test listing documents that match a specific doc prefix
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://mcp-test/docs",
	}

	goldenPath := "./_golden/serve-mcp-list-documents-multiple.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-multiple-%d.txt", time.Now().UnixMilli())

	callServeMCPServer(t, "./_input/serve-mcp/documentation.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPListDocumentsWithTags(t *testing.T) {
	// Test listing documents with tag filtering: multiple tags AND multiple values
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://mcp-test?category=overview,tutorial&audience=developer",
	}

	goldenPath := "./_golden/serve-mcp-list-documents-tags.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-tags-%d.txt", time.Now().UnixMilli())

	callServeMCPServer(t, "./_input/serve-mcp/documentation.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPListDocumentsNotFound(t *testing.T) {
	// Test listing a non-existent document
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://mcp-test/nonexistent.md",
	}

	goldenPath := "./_golden/serve-mcp-list-documents-notfound.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-notfound-%d.txt", time.Now().UnixMilli())

	callServeMCPServer(t, "./_input/serve-mcp/documentation.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPListDocumentsInvalidURI(t *testing.T) {
	client := setupServeMCPClient(t, ServeMCPClientOptions{DBPath: "./_input/serve-mcp/documentation.sqlite"})
	ctx := context.Background()

	// Test with invalid URI format
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "invalid://bad/uri",
	}

	response, err := client.CallTool(ctx, request)

	// Check for an error response
	if err != nil {
		t.Fatalf("expected call to succeed with error content: %v", err)
	}

	// The response should contain an error message
	if len(response.Content) == 0 {
		t.Fatal("expected response to have content")
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be TextContent")
	}

	// Write output for inspection
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-documents-invalid-%d.txt", time.Now().UnixMilli())
	err = os.WriteFile(outputPath, []byte(textContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	t.Logf("Invalid URI response: %s", textContent.Text)
}
