package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestServeMCPListResources(t *testing.T) {
	client := setupServeMCPClient(t, "serve mcp --documentation ./_input/serve-mcp/documentation.sqlite")
	ctx := context.Background()

	// List all resources
	request := mcp.ListResourcesRequest{}
	response, err := client.ListResources(ctx, request)
	if err != nil {
		t.Fatalf("expected to list resources successfully: %v", err)
	}

	// Convert response to JSON for comparison
	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("expected to marshal response: %v", err)
	}

	goldenPath := "./_golden/serve-mcp-list-resources-all.json"
	outputPath := fmt.Sprintf("./_output/serve-mcp-list-resources-all-%d.json", time.Now().UnixMilli())

	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
