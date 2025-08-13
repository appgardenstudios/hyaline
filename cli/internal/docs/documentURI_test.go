package docs

import "testing"

func TestDocumentURI_String_Document(t *testing.T) {
	tests := []struct {
		name     string
		uri      *DocumentURI
		expected string
	}{
		{
			name: "basic document URI",
			uri: &DocumentURI{
				SourceID:     "backend",
				DocumentPath: "README.md",
			},
			expected: "document://backend/README.md",
		},
		{
			name: "document with special characters",
			uri: &DocumentURI{
				SourceID:     "my-service",
				DocumentPath: "docs/api.md",
			},
			expected: "document://my-service/docs/api.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.uri.String()
			if result != tt.expected {
				t.Errorf("DocumentURI.String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDocumentURI_String_Section(t *testing.T) {
	tests := []struct {
		name     string
		uri      *DocumentURI
		expected string
	}{
		{
			name: "basic section URI",
			uri: &DocumentURI{
				SourceID:     "backend",
				DocumentPath: "README.md",
				Section:      "Installation",
			},
			expected: "document://backend/README.md#Installation",
		},
		{
			name: "nested section",
			uri: &DocumentURI{
				SourceID:     "frontend",
				DocumentPath: "docs/guide.md",
				Section:      "Getting Started/Setup",
			},
			expected: "document://frontend/docs/guide.md#Getting Started/Setup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.uri.String()
			if result != tt.expected {
				t.Errorf("DocumentURI.String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewDocumentURI(t *testing.T) {
	tests := []struct {
		name          string
		uriStr        string
		expectedURI   *DocumentURI
		expectedError bool
	}{
		{
			name:   "basic document URI",
			uriStr: "document://backend/README.md",
			expectedURI: &DocumentURI{
				SourceID:     "backend",
				DocumentPath: "README.md",
				Tags:         NewTags(),
			},
			expectedError: false,
		},
		{
			name:   "document URI with section",
			uriStr: "document://backend/README.md#Installation",
			expectedURI: &DocumentURI{
				SourceID:     "backend",
				DocumentPath: "README.md",
				Section:      "Installation",
				Tags:         NewTags(),
			},
			expectedError: false,
		},
		{
			name:          "invalid URI format",
			uriStr:        "http://example.com",
			expectedURI:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewDocumentURI(tt.uriStr)

			if tt.expectedError {
				if err == nil {
					t.Errorf("NewDocumentURI() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("NewDocumentURI() unexpected error: %v", err)
				return
			}

			if result.SourceID != tt.expectedURI.SourceID {
				t.Errorf("SourceID = %v, expected %v", result.SourceID, tt.expectedURI.SourceID)
			}
			if result.DocumentPath != tt.expectedURI.DocumentPath {
				t.Errorf("DocumentPath = %v, expected %v", result.DocumentPath, tt.expectedURI.DocumentPath)
			}
			if result.Section != tt.expectedURI.Section {
				t.Errorf("Section = %v, expected %v", result.Section, tt.expectedURI.Section)
			}
		})
	}
}
