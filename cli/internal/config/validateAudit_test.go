package config

import (
	"strings"
	"testing"
)

// TestCase represents the input parameters for generating a test case
type TestCase struct {
	ID               string
	Description      string
	Source           string
	URI              string
	Section          string
	ContentExists    bool
	ContentMinLength int
	ContentRegex     string
	ContentPrompt    string
	ContentPurpose   bool
	PurposeExists    bool
	TagKey           string
	TagValue         string
}

// buildTestConfig generates a Config struct for testing based on simple parameters
func buildTestConfig(tc TestCase) *Config {
	rule := AuditRule{
		ID:          tc.ID,
		Description: tc.Description,
	}

	// Build documentation filters based on what's provided
	if tc.URI != "" {
		rule.Documentation = []DocumentationFilter{{URI: tc.URI}}
	} else if tc.Section != "" {
		rule.Documentation = []DocumentationFilter{{Source: tc.Source, Section: tc.Section}}
	} else if tc.Source != "" {
		rule.Documentation = []DocumentationFilter{{Source: tc.Source}}
	}

	// Build checks based on what's provided
	if tc.ContentExists {
		rule.Checks.Content.Exists = true
	}
	if tc.ContentMinLength != 0 {
		rule.Checks.Content.MinLength = tc.ContentMinLength
	}
	if tc.ContentRegex != "" {
		rule.Checks.Content.MatchesRegex = tc.ContentRegex
	}
	if tc.ContentPrompt != "" {
		rule.Checks.Content.MatchesPrompt = tc.ContentPrompt
	}
	if tc.ContentPurpose {
		rule.Checks.Content.MatchesPurpose = true
	}
	if tc.PurposeExists {
		rule.Checks.Purpose.Exists = true
	}
	if tc.TagKey != "" {
		rule.Checks.Tags.Contains = []DocumentationFilterTag{{Key: tc.TagKey, Value: tc.TagValue}}
	}

	return &Config{
		Audit: &Audit{
			Rules: []AuditRule{rule},
		},
	}
}

func TestValidateAudit(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil audit config",
			cfg:         &Config{},
			expectError: false,
		},
		{
			name: "empty rules",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{},
				},
			},
			expectError: true,
			errorMsg:    "audit.rules must contain at least one entry",
		},
		{
			name:        "valid basic rule",
			cfg:         buildTestConfig(TestCase{ID: "test-rule", Description: "Test rule", Source: "backend", ContentExists: true}),
			expectError: false,
		},
		{
			name:        "rule with auto-generated ID",
			cfg:         buildTestConfig(TestCase{Source: "backend", ContentExists: true}),
			expectError: false,
		},
		{
			name:        "rule without documentation filters",
			cfg:         buildTestConfig(TestCase{ContentExists: true}),
			expectError: true,
			errorMsg:    "documentation must contain at least one entry",
		},
		{
			name:        "rule without checks",
			cfg:         buildTestConfig(TestCase{Source: "backend"}),
			expectError: true,
			errorMsg:    "must specify at least one check type",
		},
		{
			name:        "valid URI filter",
			cfg:         buildTestConfig(TestCase{URI: "document://backend/README.md#Installation", ContentMinLength: 100}),
			expectError: false,
		},
		{
			name:        "invalid URI without document://",
			cfg:         buildTestConfig(TestCase{URI: "backend/README.md", ContentExists: true}),
			expectError: true,
			errorMsg:    "uri must start with document://",
		},
		{
			name:        "invalid regex pattern",
			cfg:         buildTestConfig(TestCase{Source: "backend", ContentRegex: "[invalid"}),
			expectError: true,
			errorMsg:    "must be a valid regex pattern",
		},
		{
			name:        "negative min-length",
			cfg:         buildTestConfig(TestCase{Source: "backend", ContentMinLength: -1}),
			expectError: true,
			errorMsg:    "min-length must be non-negative",
		},
		{
			name:        "valid tags check",
			cfg:         buildTestConfig(TestCase{Source: "backend", TagKey: "type", TagValue: "guide"}),
			expectError: false,
		},
		{
			name:        "valid regex tag patterns",
			cfg:         buildTestConfig(TestCase{Source: "backend", TagKey: "type.*", TagValue: "(guide|tutorial)"}),
			expectError: false,
		},
		{
			name:        "invalid regex tag key",
			cfg:         buildTestConfig(TestCase{Source: "backend", TagKey: "[invalid", TagValue: "guide"}),
			expectError: true,
			errorMsg:    "key must be a valid regex pattern",
		},
		{
			name:        "invalid regex tag value",
			cfg:         buildTestConfig(TestCase{Source: "backend", TagKey: "type", TagValue: "(guide|incomplete"}),
			expectError: true,
			errorMsg:    "value must be a valid regex pattern",
		},
		{
			name:        "section without document",
			cfg:         buildTestConfig(TestCase{Source: "backend", Section: "Installation", ContentExists: true}),
			expectError: true,
			errorMsg:    "document must be set if",
		},
		{
			name: "all check types",
			cfg: buildTestConfig(TestCase{
				Source:           "backend",
				ContentExists:    true,
				ContentMinLength: 100,
				ContentRegex:     ".*README.*",
				ContentPrompt:    "This is a test prompt",
				ContentPurpose:   true,
				PurposeExists:    true,
				TagKey:           "type",
				TagValue:         "guide",
			}),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudit(tt.cfg)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectError && err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', but got: %v", tt.errorMsg, err)
				}
			}
		})
	}
}
