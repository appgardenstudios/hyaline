package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestMCPGetDocumentsAll(t *testing.T) {
	// Test getting all documents with no URI parameter
	request := mcp.CallToolRequest{}
	request.Params.Name = "get_documents"

	goldenPath := "./_golden/mcp-get-documents-all.txt"
	outputPath := fmt.Sprintf("./_output/mcp-get-documents-all-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPGetDocumentsSingle(t *testing.T) {
	// Test getting a specific document
	request := mcp.CallToolRequest{}
	request.Params.Name = "get_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://system/mcp-test/docs-fs/docs/index.html",
	}

	goldenPath := "./_golden/mcp-get-documents-single.txt"
	outputPath := fmt.Sprintf("./_output/mcp-get-documents-single-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPGetDocumentsMultiple(t *testing.T) {
	// Test getting documents that match a specific doc from git sources
	request := mcp.CallToolRequest{}
	request.Params.Name = "get_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://system/mcp-test/docs-fs/docs",
	}

	goldenPath := "./_golden/mcp-get-documents-multiple.txt"
	outputPath := fmt.Sprintf("./_output/mcp-get-documents-multiple-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPGetDocumentsNotFound(t *testing.T) {
	// Test getting a non-existent document
	request := mcp.CallToolRequest{}
	request.Params.Name = "get_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://system/mcp-test/docs-fs/nonexistent.md",
	}

	goldenPath := "./_golden/mcp-get-documents-notfound.txt"
	outputPath := fmt.Sprintf("./_output/mcp-get-documents-notfound-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPGetDocumentsInvalidURI(t *testing.T) {
	client := setupMCPClient(t, "./_input/mcp/current.sqlite")
	ctx := context.Background()

	// Test with invalid URI format
	request := mcp.CallToolRequest{}
	request.Params.Name = "get_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "invalid://bad/uri",
	}

	response, err := client.CallTool(ctx, request)
	if err != nil {
		t.Fatalf("expected to call 'get_documents' tool successfully: %v", err)
	}
	if !response.IsError {
		t.Fatal("expected result to be an error for invalid URI")
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent")
	}

	goldenPath := "./_golden/mcp-get-documents-error.txt"
	outputPath := fmt.Sprintf("./_output/mcp-get-documents-error-%d.txt", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, []byte(textContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
