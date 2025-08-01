package checks

import "testing"

func TestPurposeExists(t *testing.T) {
	tests := []struct {
		name        string
		purpose     string
		expectedPass bool
	}{
		{
			name:        "purpose exists",
			purpose:     "This document explains how to install the software",
			expectedPass: true,
		},
		{
			name:        "empty purpose",
			purpose:     "",
			expectedPass: false,
		},
		{
			name:        "whitespace only purpose",
			purpose:     "   \n\t  ",
			expectedPass: false,
		},
		{
			name:        "purpose with content",
			purpose:     "A valid purpose",
			expectedPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass, message := PurposeExists(tt.purpose)
			
			if pass != tt.expectedPass {
				t.Errorf("PurposeExists() pass = %v, expected %v", pass, tt.expectedPass)
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