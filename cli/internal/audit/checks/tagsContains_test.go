package checks

import (
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"testing"
)

func TestTagsContains(t *testing.T) {
	tests := []struct {
		name         string
		actualTags   []docs.FilteredTag
		requiredTags []config.DocumentationFilterTag
		expectedPass bool
	}{
		{
			name: "exact match",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
				{Key: "level", Value: "beginner"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
			},
			expectedPass: true,
		},
		{
			name: "no match",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "reference"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
			},
			expectedPass: false,
		},
		{
			name: "regex pattern match",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "user-guide"},
				{Key: "version", Value: "v1.2.3"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: ".*guide"},
				{Key: "version", Value: "v\\d+\\.\\d+\\.\\d+"},
			},
			expectedPass: true,
		},
		{
			name: "partial regex match fails",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "reference"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: ".*guide"},
			},
			expectedPass: false,
		},
		{
			name: "multiple required tags - all match",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
				{Key: "level", Value: "beginner"},
				{Key: "category", Value: "tutorial"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
				{Key: "level", Value: "beginner"},
			},
			expectedPass: true,
		},
		{
			name: "multiple required tags - one missing",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
				{Key: "level", Value: "beginner"},
			},
			expectedPass: false,
		},
		{
			name:         "no actual tags",
			actualTags:   []docs.FilteredTag{},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
			},
			expectedPass: false,
		},
		{
			name: "empty required tags",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
			},
			requiredTags: []config.DocumentationFilterTag{},
			expectedPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass, message := TagsContains(tt.actualTags, tt.requiredTags)
			
			if pass != tt.expectedPass {
				t.Errorf("TagsContains() pass = %v, expected %v", pass, tt.expectedPass)
			}
			
			if tt.expectedPass && message != "" {
				t.Errorf("Expected empty message for passing check, got: %s", message)
			}
			
			if !tt.expectedPass && message == "" {
				t.Errorf("Expected non-empty message for failing check")
			}
		})
	}
}