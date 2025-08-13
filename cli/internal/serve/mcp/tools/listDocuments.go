package tools

import (
	"context"
	"fmt"
	"hyaline/internal/docs"
	"hyaline/internal/serve/mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func ListDocumentsTool() mcp.Tool {
	return mcp.NewTool("list_documents",
		mcp.WithDescription("List all documents at or under the specified URI path, or all documents if no URI provided. Document URIs follow this pattern: `document://<source-id>/<document-id>[?<key>=<value>][#<section>]` where query parameters filter by tags (multiple values for same key are comma-separated)"),
		mcp.WithString("document_uri",
			mcp.Description("The URI path to list documents from (can be partial). Format: document://<source-id>/<document-id>[?<key>=<value>]. Query parameters filter results by tags. If not provided, lists all documents."),
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
	var documentURI *docs.DocumentURI
	if documentURIStr != "" {
		parsedURI, err := docs.NewDocumentURI(documentURIStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid URI format: %s", err.Error())), nil
		}
		documentURI = parsedURI
	} else {
		// No URI provided, list all documents
		documentURI = &docs.DocumentURI{}
	}

	// Build the document list in XML format
	results := utils.ProcessDocuments(documentationData, documentURI, false)

	// Check if any documents were found
	resultStr := results.Result.String()
	if !strings.Contains(resultStr, "<document>") {
		return mcp.NewToolResultText("No documents found matching the specified URI."), nil
	}

	var response strings.Builder
	response.WriteString("The <documents> XML structure contains all available documents and sections with their corresponding document URIs, along with <purpose> and <tags> metadata when available.\n\n")
	response.WriteString(resultStr)

	return mcp.NewToolResultText(response.String()), nil
}
