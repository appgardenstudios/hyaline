package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestServeMCPReloadDocumentation(t *testing.T) {
	githubToken := os.Getenv("_HYALINE_TEST_GITHUB_TOKEN")

	// Set initial documentation to documentation_1.db
	dispatchWorkflow(t, githubToken, "appgardenstudios/hyaline-example", "set-current-documentation.yml", map[string]interface{}{
		"artifact_source": "documentation_1.db",
	})

	client := setupServeMCPClient(t, "serve mcp --github-repo appgardenstudios/hyaline-example --github-token "+githubToken)
	ctx := context.Background()

	// 1. List documents before reload
	listRequest := mcp.CallToolRequest{}
	listRequest.Params.Name = "list_documents"
	listResponse, err := client.CallTool(ctx, listRequest)
	if err != nil {
		t.Fatalf("expected to call 'list_documents' tool successfully: %v", err)
	}
	listContent, ok := listResponse.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent")
	}

	// Write output for golden file comparison - before reload (documentation_1.db)
	goldenPathBefore := "./_golden/serve-mcp-reload-documentation-list-before.txt"
	outputPathBefore := fmt.Sprintf("./_output/serve-mcp-reload-documentation-list-before-%d.txt", time.Now().UnixMilli())
	err = os.WriteFile(outputPathBefore, []byte(listContent.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPathBefore, outputPathBefore, t)
	}

	compareFiles(goldenPathBefore, outputPathBefore, t)

	// 2. Switch to documentation_2.db and reload
	dispatchWorkflow(t, githubToken, "appgardenstudios/hyaline-example", "set-current-documentation.yml", map[string]interface{}{
		"artifact_source": "documentation_2.db",
	})

	// 3. Reload documentation
	reloadRequest := mcp.CallToolRequest{}
	reloadRequest.Params.Name = "reload_documentation"
	reloadResponse, err := client.CallTool(ctx, reloadRequest)
	if err != nil {
		t.Fatalf("expected to call 'reload_documentation' tool successfully: %v", err)
	}
	reloadContent, ok := reloadResponse.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent")
	}
	if reloadContent.Text != "Documentation reloaded successfully." {
		t.Fatalf("unexpected reload response: %s", reloadContent.Text)
	}

	// 4. List documents after reload
	listResponseAfter, err := client.CallTool(ctx, listRequest)
	if err != nil {
		t.Fatalf("expected to call 'list_documents' tool successfully after reload: %v", err)
	}
	listContentAfter, ok := listResponseAfter.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected content to be of type TextContent after reload")
	}

	// Write output for golden file comparison - after reload (documentation_2.db)
	goldenPathAfter := "./_golden/serve-mcp-reload-documentation-list-after.txt"
	outputPathAfter := fmt.Sprintf("./_output/serve-mcp-reload-documentation-list-after-%d.txt", time.Now().UnixMilli())
	err = os.WriteFile(outputPathAfter, []byte(listContentAfter.Text), 0644)
	if err != nil {
		t.Fatalf("expected to write output file: %v", err)
	}

	if *update {
		updateGolden(goldenPathAfter, outputPathAfter, t)
	}

	compareFiles(goldenPathAfter, outputPathAfter, t)
}
