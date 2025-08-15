package tools

import (
	"context"
	"fmt"
	"hyaline/internal/docs"
	"hyaline/internal/serve/mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetDocumentsTool() mcp.Tool {
	return mcp.NewTool("get_documents",
		mcp.WithDescription("Get the contents of documents matching the specified URI, or all documents if no URI provided. Document URIs follow this pattern: `document://<source-id>/<document-id>[?<key>=<value>][#<section>]` where query parameters filter by tags (multiple values for same key are comma-separated)"),
		mcp.WithString("document_uri",
			mcp.Description("The URI specifying which documents to retrieve (can be partial). Format: document://<source-id>/<document-id>[?<key>=<value>]. Query parameters filter results by tags. If not provided, retrieves all documents."),
		),
	)
}

func HandleGetDocuments(_ context.Context, request mcp.CallToolRequest, documentationData *utils.DocumentationData) (*mcp.CallToolResult, error) {

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
		// No URI provided, match all documents
		documentURI = &docs.DocumentURI{}
	}

	// Process matching documents with content
	results := utils.ProcessDocuments(documentationData, documentURI, true)

	// Check if any documents were found
	resultStr := results.Result.String()
	if !strings.Contains(resultStr, "<document>") {
		return mcp.NewToolResultText("No documents found matching the specified URI."), nil
	}

	var response strings.Builder
	response.WriteString("The <documents> XML structure contains all requested documents and sections. Each <document> has the <document_content> which contains the contents of the document, along with <purpose> and <tags> metadata when available. Sections also include their <purpose> and <tags> when applicable.\n\n")
	response.WriteString(resultStr)

	return mcp.NewToolResultText(response.String()), nil
}
