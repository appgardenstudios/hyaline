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
		expectError  bool
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
			expectError:  false,
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
			expectError:  false,
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
			expectError:  false,
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
			expectError:  false,
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
			expectError:  false,
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
			expectError:  false,
		},
		{
			name:       "no actual tags",
			actualTags: []docs.FilteredTag{},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "guide"},
			},
			expectedPass: false,
			expectError:  false,
		},
		{
			name: "empty required tags",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
			},
			requiredTags: []config.DocumentationFilterTag{},
			expectedPass: true,
			expectError:  false,
		},
		{
			name: "invalid regex in key pattern",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "[invalid", Value: "guide"},
			},
			expectedPass: false,
			expectError:  true,
		},
		{
			name: "invalid regex in value pattern",
			actualTags: []docs.FilteredTag{
				{Key: "type", Value: "guide"},
			},
			requiredTags: []config.DocumentationFilterTag{
				{Key: "type", Value: "[invalid"},
			},
			expectedPass: false,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass, message, err := TagsContains(tt.actualTags, tt.requiredTags)

			if tt.expectError {
				if err == nil {
					t.Errorf("TagsContains() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("TagsContains() unexpected error: %v", err)
			}

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
