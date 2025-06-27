package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestMCPListDocumentsAll(t *testing.T) {
	// Test listing all documents with no URI parameter
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"

	goldenPath := "./_golden/mcp-list-documents-all.txt"
	outputPath := fmt.Sprintf("./_output/mcp-list-documents-all-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPListDocumentsSystemLevel(t *testing.T) {
	// Test listing documents at system level
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://system/mcp-test",
	}

	goldenPath := "./_golden/mcp-list-documents-system.txt"
	outputPath := fmt.Sprintf("./_output/mcp-list-documents-system-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPListDocumentsDocumentationLevel(t *testing.T) {
	// Test listing documents at documentation level
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "document://system/mcp-test/docs-fs",
	}

	goldenPath := "./_golden/mcp-list-documents-docs.txt"
	outputPath := fmt.Sprintf("./_output/mcp-list-documents-docs-%d.txt", time.Now().UnixMilli())

	callMCPServer(t, "./_input/mcp/current.sqlite", request, outputPath)

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestMCPListDocumentsInvalidURI(t *testing.T) {
	client := setupMCPClient(t, "./_input/mcp/current.sqlite")
	ctx := context.Background()

	// Test list_documents with invalid URI format
	request := mcp.CallToolRequest{}
	request.Params.Name = "list_documents"
	request.Params.Arguments = map[string]any{
		"document_uri": "invalid://bad/uri",
	}

	response, err := client.CallTool(ctx, request)
	if err != nil {
		t.Fatalf("expected to call 'list_documents' tool successfully: %v", err)
	}
	if !response.IsError {
		t.Fatal("expected result to be an error for invalid URI")
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent")
	}

	goldenPath := "./_golden/mcp-list-documents-error.txt"
	outputPath := fmt.Sprintf("./_output/mcp-list-documents-error-%d.txt", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, []byte(textContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
