package tools

import (
	"context"
	"fmt"
	"hyaline/internal/mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func ListDocumentsTool() mcp.Tool {
	return mcp.NewTool("list_documents",
		mcp.WithDescription("List all documents at or under the specified URI path, or all documents if no URI provided. Document URIs follow this pattern: `document://system/<system-id>/<documentation-id>/<document-path>`"),
		mcp.WithString("document_uri",
			mcp.Description("The URI path to list documents from (can be partial). Format: document://system/<system-id>/<documentation-id>/<document-path>. If not provided, lists all documents."),
		),
	)
}

func HandleListDocuments(_ context.Context, request mcp.CallToolRequest, documentationData *utils.DocumentationData) (*mcp.CallToolResult, error) {

	// Get the optional document_uri parameter
	documentURIStr := ""
	if uriArg, exists := request.GetArguments()["document_uri"]; exists {
		if uriStr, ok := uriArg.(string); ok {
			documentURIStr = uriStr
		}
	}

	// Parse the URI if provided
	var documentURI *utils.DocumentURI
	if documentURIStr != "" {
		parsedURI, err := utils.NewDocumentURI(documentURIStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid URI format: %s", err.Error())), nil
		}

		documentURI = parsedURI
	}

	// Build the document list in XML format
	results := utils.ProcessDocuments(documentationData, documentURI, false)

	if results.Total == 0 {
		return mcp.NewToolResultText("No documents found matching the specified URI."), nil
	}

	var response strings.Builder
	response.WriteString("The <systems> XML structure contains all available systems, documentation sources, documents, and sections with their corresponding document URIs.\n\n")
	response.WriteString(results.Result.String())

	return mcp.NewToolResultText(response.String()), nil
}
