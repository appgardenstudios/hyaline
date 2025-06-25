package utils

import (
	"fmt"
	"strings"
)

// DocumentURI represents a parsed document URI
type DocumentURI struct {
	SystemID        string
	DocumentationID string
	DocumentPath    string
}

// String returns the URI as a string
func (d *DocumentURI) String() string {
	result := "document://system"
	if d.SystemID != "" {
		result += "/" + d.SystemID
		if d.DocumentationID != "" {
			result += "/" + d.DocumentationID
			if d.DocumentPath != "" {
				result += "/" + d.DocumentPath
			}
		}
	}
	return result
}

// NewDocumentURI parses a document URI string into components
func NewDocumentURI(uriStr string) (*DocumentURI, error) {
	// Expected format: document://system/<system-id>/<documentation-id>/<document-path>
	if !strings.HasPrefix(uriStr, "document://system") {
		return nil, fmt.Errorf("invalid URI format: must start with 'document://system'")
	}

	// Remove the scheme
	path := strings.TrimPrefix(uriStr, "document://system")
	if path == "" {
		// URI is just "document://system"
		return &DocumentURI{}, nil
	}

	// Remove any fragment (section ID) - we no longer support section IDs in URIs
	if idx := strings.Index(path, "#"); idx != -1 {
		path = path[:idx]
	}

	// Split the path into components
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		return &DocumentURI{}, nil
	}

	result := &DocumentURI{}

	if len(parts) >= 1 && parts[0] != "" {
		result.SystemID = parts[0]
	}
	if len(parts) >= 2 && parts[1] != "" {
		result.DocumentationID = parts[1]
	}
	if len(parts) >= 3 {
		result.DocumentPath = strings.Join(parts[2:], "/")
	}

	return result, nil
}
