package docs

import (
	"hyaline/internal/config"
	"testing"
)

func TestSectionMatches(t *testing.T) {
	tests := []struct {
		name       string
		sectionID  string
		documentID string
		sourceID   string
		tags       []FilteredTag
		filter     *config.DocumentationFilter
		strict     bool
		expected   bool
	}{
		{
			name:       "strict mode - empty section filter - returns false",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
			},
			strict:   true,
			expected: false,
		},
		{
			name:       "non-strict mode - empty section filter - matches document",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
			},
			strict:   false,
			expected: true,
		},
		{
			name:       "exact section match",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
			},
			strict:   false,
			expected: true,
		},
		{
			name:       "section no match",
			sectionID:  "Usage",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
			},
			strict:   false,
			expected: false,
		},
		{
			name:       "source no match",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "frontend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
			},
			strict:   false,
			expected: false,
		},
		{
			name:       "document no match",
			sectionID:  "Installation",
			documentID: "CHANGELOG.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
			},
			strict:   false,
			expected: false,
		},
		{
			name:       "wildcard section match",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "*ation",
			},
			strict:   false,
			expected: true,
		},
		{
			name:       "tag match with section",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags: []FilteredTag{
				{Key: "difficulty", Value: "easy"},
			},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
				Tags: []config.DocumentationFilterTag{
					{Key: "difficulty", Value: "easy"},
				},
			},
			strict:   false,
			expected: true,
		},
		{
			name:       "tag no match with section",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags: []FilteredTag{
				{Key: "difficulty", Value: "hard"},
			},
			filter: &config.DocumentationFilter{
				Source:   "backend",
				Document: "README.md",
				Section:  "Installation",
				Tags: []config.DocumentationFilterTag{
					{Key: "difficulty", Value: "easy"},
				},
			},
			strict:   false,
			expected: false,
		},
		{
			name:       "document match with no section - non strict",
			sectionID:  "Installation",
			documentID: "README.md",
			sourceID:   "backend",
			tags:       []FilteredTag{},
			filter: &config.DocumentationFilter{
				Source: "backend",
			},
			strict:   false,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SectionMatches(tt.sectionID, tt.documentID, tt.sourceID, tt.tags, tt.filter, tt.strict)
			if result != tt.expected {
				t.Errorf("SectionMatches() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
