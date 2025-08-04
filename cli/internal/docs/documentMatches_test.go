package docs

import (
	"hyaline/internal/config"
	"testing"
)

func TestDocumentMatches(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		sourceID   string
		tags       []FilteredTag
		filter     *config.DocumentationFilter
		expected   bool
	}{
		{
			name:       "basic source match",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source: "backend",
			},
			expected: true,
		},
		{
			name:       "source and document match",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
			},
			expected: true,
		},
		{
			name:       "source no match",
			documentID: "README.md",
			sourceID:   "frontend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source: "backend",
			},
			expected: false,
		},
		{
			name:       "document no match",
			documentID: "CHANGELOG.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
			},
			expected: false,
		},
		{
			name:       "section filter excludes document",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
			},
			expected: false,
		},
		{
			name:       "wildcard source match",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source: "*",
			},
			expected: true,
		},
		{
			name:       "wildcard document match",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "*.md",
			},
			expected: true,
		},
		{
			name:       "tag match success",
			documentID: "README.md",
			sourceID:   "backend",
			tags: []FilteredTag{
				{Key: "type", Value: "guide"},
			},
			filter: &config.DocumentationFilter{
				Source: "backend",
				Tags: []config.DocumentationFilterTag{
					{Key: "type", Value: "guide"},
				},
			},
			expected: true,
		},
		{
			name:       "tag match failure",
			documentID: "README.md",
			sourceID:   "backend",
			tags: []FilteredTag{
				{Key: "type", Value: "reference"},
			},
			filter: &config.DocumentationFilter{
				Source: "backend",
				Tags: []config.DocumentationFilterTag{
					{Key: "type", Value: "guide"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DocumentMatches(tt.documentID, tt.sourceID, tt.tags, tt.filter)
			if result != tt.expected {
				t.Errorf("DocumentMatches() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
