package tools

import (
	"context"
	"fmt"
	"hyaline/internal/mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetDocumentsTool creates the get_documents tool definition
func GetDocumentsTool() mcp.Tool {
	return mcp.NewTool("get_documents",
		mcp.WithDescription("Get the contents of documents matching the specified URI, or all documents if no URI provided. Document URIs follow this pattern: `document://system/<system-id>/<documentation-id>/<document-path>`"),
		mcp.WithString("document_uri",
			mcp.Description("The URI specifying which documents to retrieve (can be partial). Format: document://system/<system-id>/<documentation-id>/<document-path>. If not provided, retrieves all documents."),
		),
	)
}

// HandleGetDocuments handles the get_documents tool call
func HandleGetDocuments(_ context.Context, request mcp.CallToolRequest, mcpData *utils.MCPData) (*mcp.CallToolResult, error) {

	// Get the optional document_uri parameter
	documentURIStr := ""
	if uriArg, exists := request.GetArguments()["document_uri"]; exists {
		if uriStr, ok := uriArg.(string); ok {
			documentURIStr = uriStr
		}
	}

	// Parse the URI if provided
	var docURI *utils.DocumentURI
	if documentURIStr != "" {
		parsedURI, err := utils.NewDocumentURI(documentURIStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid URI format: %s", err.Error())), nil
		}

		docURI = parsedURI
	}

	// Process matching documents with content
	results := utils.ProcessDocuments(mcpData, docURI, true)

	if results.Total == 0 {
		return mcp.NewToolResultText("No documents found matching the specified URI."), nil
	}

	// Wrap the content with systems XML tag and description
	var result strings.Builder
	result.WriteString("The <systems> XML structure contains all requested systems, documentation sources, and documents. Each <document> has the <document_content> which are the contents of the document.\n\n")
	result.WriteString("<systems>\n")
	result.WriteString(results.Result.String())
	result.WriteString("</systems>")

	return mcp.NewToolResultText(result.String()), nil
}