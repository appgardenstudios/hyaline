package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestServeMCPReadResourceValid(t *testing.T) {
	client := setupServeMCPClient(t, "serve mcp --documentation ./_input/serve-mcp/documentation.sqlite")
	ctx := context.Background()

	// Read a specific resource
	request := mcp.ReadResourceRequest{}
	request.Params.URI = "document://mcp-test/docs/index.html"

	response, err := client.ReadResource(ctx, request)
	if err != nil {
		t.Fatalf("expected to read resource successfully: %v", err)
	}

	// Verify we got content back
	if len(response.Contents) == 0 {
		t.Fatal("expected response to have content")
	}

	// Cast to TextResourceContents to access Text field
	textContent, ok := mcp.AsTextResourceContents(response.Contents[0])
	if !ok {
		t.Fatal("expected content to be TextResourceContents")
	}

	// Write the content to output file
	goldenPath := "./_golden/serve-mcp-read-resource-valid.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-read-resource-valid-%d.txt", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, []byte(textContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPReadResourceNotFound(t *testing.T) {
	client := setupServeMCPClient(t, "serve mcp --documentation ./_input/serve-mcp/documentation.sqlite")
	ctx := context.Background()

	// Try to read a non-existent resource
	request := mcp.ReadResourceRequest{}
	request.Params.URI = "document://mcp-test/nonexistent.md"

	_, err := client.ReadResource(ctx, request)

	// We expect an error
	if err == nil {
		t.Fatal("expected error when reading non-existent resource")
	}

	goldenPath := "./_golden/serve-mcp-read-resource-notfound.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-read-resource-notfound-%d.txt", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, []byte(err.Error()), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}

func TestServeMCPReadResourceInvalidURI(t *testing.T) {
	client := setupServeMCPClient(t, "serve mcp --documentation ./_input/serve-mcp/documentation.sqlite")
	ctx := context.Background()

	// Try to read with invalid URI format
	request := mcp.ReadResourceRequest{}
	request.Params.URI = "invalid://bad/uri"

	_, err := client.ReadResource(ctx, request)

	// We expect an error
	if err == nil {
		t.Fatal("expected error when reading with invalid URI")
	}

	goldenPath := "./_golden/serve-mcp-read-resource-invalid.txt"
	outputPath := fmt.Sprintf("./_output/serve-mcp-read-resource-invalid-%d.txt", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, []byte(err.Error()), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
