package docs

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Tags represents a collection of tag keys with their associated values
type Tags map[string][]string

// Add adds a value to the specified tag key
func (t Tags) Add(key, value string) {
	t[key] = append(t[key], value)
}

// Keys returns all tag keys in sorted order
func (t Tags) Keys() []string {
	keys := make([]string, 0, len(t))
	for key := range t {
		keys = append(keys, key)
	}
	// Sort for deterministic order
	sort.Strings(keys)
	return keys
}

// NewTags creates a new Tags instance
func NewTags() Tags {
	return make(Tags)
}

// DocumentURI represents a parsed document URI with tag filtering support
type DocumentURI struct {
	SourceID     string
	DocumentPath string
	Tags         Tags // key -> values (comma-separated values become array)
	Section      string
}

// String returns the URI as a string
func (documentURI *DocumentURI) String() string {
	result := "document://"
	if documentURI.SourceID != "" {
		result += documentURI.SourceID
		if documentURI.DocumentPath != "" {
			result += "/" + documentURI.DocumentPath
		}
	}

	// Add query parameters for tags
	if len(documentURI.Tags) > 0 {
		params := url.Values{}
		for key, values := range documentURI.Tags {
			params.Set(key, strings.Join(values, ","))
		}
		result += "?" + params.Encode()
	}

	// Add section fragment
	if documentURI.Section != "" {
		result += "#" + documentURI.Section
	}

	return result
}

// NewDocumentURI parses a document URI string into components
// Format: document://<source-id>/<document-id>[?<key>=<value>][#<section>]
func NewDocumentURI(uriStr string) (*DocumentURI, error) {
	// Expected format: document://<source-id>/<document-id>[?<key>=<value>][#<section>]
	if !strings.HasPrefix(uriStr, "document://") {
		return nil, fmt.Errorf("invalid URI format: must start with 'document://'")
	}

	// Remove the scheme
	path := strings.TrimPrefix(uriStr, "document://")

	// Extract section fragment if present
	section := ""
	if idx := strings.Index(path, "#"); idx != -1 {
		section = path[idx+1:]
		path = path[:idx]
	}

	// Extract query parameters if present
	tags := NewTags()
	if idx := strings.Index(path, "?"); idx != -1 {
		queryStr := path[idx+1:]
		path = path[:idx]

		// Parse query parameters
		params, err := url.ParseQuery(queryStr)
		if err != nil {
			return nil, fmt.Errorf("invalid query parameters: %w", err)
		}

		// Convert to our tag format (comma-separated values to array)
		for key, values := range params {
			if len(values) > 0 && values[0] != "" {
				// Split comma-separated values
				tagValues := strings.Split(values[0], ",")
				for i, v := range tagValues {
					tagValues[i] = strings.TrimSpace(v)
				}
				tags[key] = tagValues
			}
		}
	}

	// Split the path into source and document parts
	parts := strings.SplitN(path, "/", 2)

	result := &DocumentURI{
		Tags:    tags,
		Section: section,
	}

	if len(parts) > 0 && parts[0] != "" {
		result.SourceID = parts[0]
	}
	if len(parts) > 1 && parts[1] != "" {
		result.DocumentPath = parts[1]
	}

	return result, nil
}

// MatchesTags checks if a document/section matches the URI's tag filters
// When multiple values for same tag: match ANY value
// When multiple tags: match at least one value for EACH tag
func (documentURI *DocumentURI) MatchesTags(tags Tags) bool {
	// No tag filters means match all
	if len(documentURI.Tags) == 0 {
		return true
	}

	// Check each required tag
	for requiredKey, requiredValues := range documentURI.Tags {
		values, hasKey := tags[requiredKey]
		if !hasKey || len(values) == 0 {
			return false
		}

		// Check if document/section has at least one of the required values for this tag
		found := false
		for _, requiredValue := range requiredValues {
			for _, value := range values {
				if requiredValue == value {
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			return false
		}
	}

	// All required tags matched
	return true
}