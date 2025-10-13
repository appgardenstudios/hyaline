package resources

import (
	"context"
	"fmt"
	"hyaline/internal/docs"
	"hyaline/internal/serve/mcp/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateDocumentResources creates MCP resources for all documents in the documentation data
func CreateDocumentResources(data *utils.DocumentationData) []mcp.Resource {
	resources := make([]mcp.Resource, 0)

	// Iterate through all sources and documents to create resources
	for _, source := range data.Sources {
		for _, document := range source.Documents {
			uri := &docs.DocumentURI{
				SourceID:     source.ID,
				DocumentPath: document.ID,
			}

			// Create resource
			// Note: MIME type is always text/markdown because extracted data is always markdown
			resource := mcp.NewResource(
				uri.String(),
				document.ID,
				mcp.WithResourceDescription(document.Purpose),
				mcp.WithMIMEType("text/markdown"),
			)

			resources = append(resources, resource)
		}
	}

	return resources
}

// HandleReadResource handles reading a specific resource by URI
func HandleReadResource(_ context.Context, request mcp.ReadResourceRequest, data *utils.DocumentationData) ([]mcp.ResourceContents, error) {
	// Parse the URI
	documentURI, err := docs.NewDocumentURI(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid URI format: %w", err)
	}

	// Find the matching document
	for _, source := range data.Sources {
		if source.ID != documentURI.SourceID {
			continue
		}

		for _, document := range source.Documents {
			if document.ID != documentURI.DocumentPath {
				continue
			}

			// Found matching document
			uri := &docs.DocumentURI{
				SourceID:     source.ID,
				DocumentPath: document.ID,
			}

			return []mcp.ResourceContents{
				&mcp.TextResourceContents{
					URI: uri.String(),
					// MIME type is always text/markdown because extracted data is always markdown
					MIMEType: "text/markdown",
					Text:     document.ExtractedData,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("resource not found: %s", request.Params.URI)
}
