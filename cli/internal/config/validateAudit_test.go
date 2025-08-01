package config

import (
	"strings"
	"testing"
)

// buildTestConfig generates a Config struct for testing based on AuditRule parameters
func buildTestConfig(rules ...AuditRule) *Config {
	return &Config{
		Audit: &Audit{
			Rules: rules,
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
			name:        "empty rules",
			cfg:         buildTestConfig(),
			expectError: true,
			errorMsg:    "audit.rules must contain at least one entry",
		},
		{
			name: "valid basic rule",
			cfg: buildTestConfig(AuditRule{
				ID:          "test-rule",
				Description: "Test rule",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: false,
		},
		{
			name: "rule with auto-generated ID",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: false,
		},
		{
			name: "rule without documentation filters",
			cfg: buildTestConfig(AuditRule{
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "documentation must contain at least one entry",
		},
		{
			name: "rule without checks",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
			}),
			expectError: true,
			errorMsg:    "must specify at least one check type",
		},
		{
			name: "valid URI filter",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{URI: "document://backend/README.md#Installation"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						MinLength: 100,
					},
				},
			}),
			expectError: false,
		},
		{
			name: "invalid URI without document://",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{URI: "backend/README.md"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "uri must start with document://",
		},
		{
			name: "invalid regex pattern",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						MatchesRegex: "[invalid",
					},
				},
			}),
			expectError: true,
			errorMsg:    "must be a valid regex pattern",
		},
		{
			name: "negative min-length",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						MinLength: -1,
					},
				},
			}),
			expectError: true,
			errorMsg:    "min-length must be non-negative",
		},
		{
			name: "valid tags check",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Tags: AuditTagsChecks{
						Contains: []DocumentationFilterTag{
							{Key: "type", Value: "guide"},
						},
					},
				},
			}),
			expectError: false,
		},
		{
			name: "valid regex tag patterns",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Tags: AuditTagsChecks{
						Contains: []DocumentationFilterTag{
							{Key: "type.*", Value: "(guide|tutorial)"},
						},
					},
				},
			}),
			expectError: false,
		},
		{
			name: "invalid regex tag key",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Tags: AuditTagsChecks{
						Contains: []DocumentationFilterTag{
							{Key: "[invalid", Value: "guide"},
						},
					},
				},
			}),
			expectError: true,
			errorMsg:    "key must be a valid regex pattern",
		},
		{
			name: "invalid regex tag value",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Tags: AuditTagsChecks{
						Contains: []DocumentationFilterTag{
							{Key: "type", Value: "(guide|incomplete"},
						},
					},
				},
			}),
			expectError: true,
			errorMsg:    "value must be a valid regex pattern",
		},
		{
			name: "section without document",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend", Section: "Installation"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "document must be set if",
		},
		{
			name: "all check types",
			cfg: buildTestConfig(AuditRule{
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists:         true,
						MinLength:      100,
						MatchesRegex:   ".*README.*",
						MatchesPrompt:  "This is a test prompt",
						MatchesPurpose: true,
					},
					Purpose: AuditPurposeChecks{
						Exists: true,
					},
					Tags: AuditTagsChecks{
						Contains: []DocumentationFilterTag{
							{Key: "type", Value: "guide"},
						},
					},
				},
			}),
			expectError: false,
		},
		{
			name: "valid rule ID",
			cfg: buildTestConfig(AuditRule{
				ID: "valid-rule_123",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: false,
		},
		{
			name: "invalid rule ID with spaces",
			cfg: buildTestConfig(AuditRule{
				ID: "rule with spaces",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "must match pattern",
		},
		{
			name: "invalid rule ID starting with hyphen",
			cfg: buildTestConfig(AuditRule{
				ID: "-invalid",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "must match pattern",
		},
		{
			name: "invalid rule ID with special characters",
			cfg: buildTestConfig(AuditRule{
				ID: "rule@special",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "must match pattern",
		},
		{
			name: "invalid rule ID too long",
			cfg: buildTestConfig(AuditRule{
				ID: "this-is-a-very-long-rule-id-that-exceeds-the-maximum-length-of-64-characters",
				Documentation: []DocumentationFilter{
					{Source: "backend"},
				},
				Checks: AuditChecks{
					Content: AuditContentChecks{
						Exists: true,
					},
				},
			}),
			expectError: true,
			errorMsg:    "must match pattern",
		},
		{
			name: "duplicate rule IDs",
			cfg: buildTestConfig(
				AuditRule{
					ID: "duplicate-id",
					Documentation: []DocumentationFilter{
						{Source: "backend"},
					},
					Checks: AuditChecks{
						Content: AuditContentChecks{
							Exists: true,
						},
					},
				},
				AuditRule{
					ID: "duplicate-id",
					Documentation: []DocumentationFilter{
						{Source: "frontend"},
					},
					Checks: AuditChecks{
						Content: AuditContentChecks{
							Exists: true,
						},
					},
				},
			),
			expectError: true,
			errorMsg:    "is not unique",
		},
		{
			name: "empty rule IDs are allowed to coexist",
			cfg: buildTestConfig(
				AuditRule{
					ID: "",
					Documentation: []DocumentationFilter{
						{Source: "backend"},
					},
					Checks: AuditChecks{
						Content: AuditContentChecks{
							Exists: true,
						},
					},
				},
				AuditRule{
					ID: "",
					Documentation: []DocumentationFilter{
						{Source: "frontend"},
					},
					Checks: AuditChecks{
						Content: AuditContentChecks{
							Exists: true,
						},
					},
				},
			),
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
