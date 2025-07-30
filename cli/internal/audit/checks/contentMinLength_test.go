package checks

import "testing"

func TestContentMinLength(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		minLength    int
		expectedPass bool
	}{
		{
			name:         "content meets minimum length",
			content:      "This is a long enough piece of content for the test",
			minLength:    20,
			expectedPass: true,
		},
		{
			name:         "content exactly meets minimum length",
			content:      "exactly twenty chars",
			minLength:    20,
			expectedPass: true,
		},
		{
			name:         "content shorter than minimum length",
			content:      "short",
			minLength:    20,
			expectedPass: false,
		},
		{
			name:         "empty content",
			content:      "",
			minLength:    1,
			expectedPass: false,
		},
		{
			name:         "zero minimum length",
			content:      "any content",
			minLength:    0,
			expectedPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass, message := ContentMinLength(tt.content, tt.minLength)
			
			if pass != tt.expectedPass {
				t.Errorf("ContentMinLength() pass = %v, expected %v", pass, tt.expectedPass)
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