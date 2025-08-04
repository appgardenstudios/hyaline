package audit

import (
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/sqlite"
	"testing"
)

func TestMatchesAnyFilter(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		sourceID   string
		sectionID  string
		tags       []docs.FilteredTag
		filters    []config.DocumentationFilter
		expected   bool
	}{
		{
			name:       "document matches source filter",
			documentID: "README.md",
			sourceID:   "backend",
			sectionID:  "",
			tags:       []docs.FilteredTag{},
			filters: []config.DocumentationFilter{
				{
					Source: "backend",
				},
			},
			expected: true,
		},
		{
			name:       "document doesn't match source filter",
			documentID: "README.md",
			sourceID:   "frontend",
			sectionID:  "",
			tags:       []docs.FilteredTag{},
			filters: []config.DocumentationFilter{
				{
					Source: "backend",
				},
			},
			expected: false,
		},
		{
			name:       "section matches specific filter",
			documentID: "README.md",
			sourceID:   "backend",
			sectionID:  "Installation",
			tags:       []docs.FilteredTag{},
			filters: []config.DocumentationFilter{
				{
					Source:   "backend",
					Document: "README.md",
					Section:  "Installation",
				},
			},
			expected: true,
		},
		{
			name:       "section with strict matching",
			documentID: "README.md",
			sourceID:   "backend",
			sectionID:  "Installation",
			tags:       []docs.FilteredTag{},
			filters: []config.DocumentationFilter{
				{
					Source:   "backend",
					Document: "README.md",
				},
			},
			expected: false, // Should not match because we use strict=true for sections
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesAnyFilter(tt.documentID, tt.sourceID, tt.sectionID, tt.tags, tt.filters)
			if result != tt.expected {
				t.Errorf("matchesAnyFilter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFilterBySource(t *testing.T) {
	documents := []sqlite.DOCUMENT{
		{ID: "doc1", SourceID: "backend"},
		{ID: "doc2", SourceID: "frontend"},
		{ID: "doc3", SourceID: "backend"},
	}

	sources := []string{"backend"}
	filtered := filterBySource(documents, sources, getDocumentSourceID)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(filtered))
	}

	for _, doc := range filtered {
		if doc.SourceID != "backend" {
			t.Errorf("Expected all documents to have sourceID 'backend', got %s", doc.SourceID)
		}
	}
}
