package checks

import "testing"

func TestContentMatchesRegex(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		pattern     string
		expectedPass bool
		expectError  bool
	}{
		{
			name:        "simple match",
			content:     "This is a README file",
			pattern:     "README",
			expectedPass: true,
			expectError:  false,
		},
		{
			name:        "no match",
			content:     "This is a CHANGELOG file",
			pattern:     "README",
			expectedPass: false,
			expectError:  false,
		},
		{
			name:        "case-insensitive match",
			content:     "This is a readme file",
			pattern:     "(?i)README",
			expectedPass: true,
			expectError:  false,
		},
		{
			name:        "invalid regex",
			content:     "content",
			pattern:     "[invalid",
			expectedPass: false,
			expectError:  false, // We handle this as a failed check, not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass, message := ContentMatchesRegex(tt.content, tt.pattern)
			
			if pass != tt.expectedPass {
				t.Errorf("ContentMatchesRegex() pass = %v, expected %v", pass, tt.expectedPass)
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