package config

import (
	"strings"
	"testing"
)

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
			name: "valid basic rule",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							ID:          "test-rule",
							Description: "Test rule",
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									Exists: true,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "rule with auto-generated ID",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									Exists: true,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "rule without documentation filters",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									Exists: true,
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "documentation must contain at least one entry",
		},
		{
			name: "rule without checks",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "must specify at least one check type",
		},
		{
			name: "valid URI filter",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									URI: "document://backend/README.md#Installation",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									MinLength: 100,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid URI without document://",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									URI: "backend/README.md",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									Exists: true,
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "uri must start with document://",
		},
		{
			name: "invalid regex pattern",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									MatchesRegex: "[invalid",
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "must be a valid regex pattern",
		},
		{
			name: "negative min-length",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									MinLength: -1,
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "min-length must be non-negative",
		},
		{
			name: "valid tags check",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Tags: AuditTagsChecks{
									Contains: []DocumentationFilterTag{
										{
											Key:   "type",
											Value: "guide",
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid regex tag patterns",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Tags: AuditTagsChecks{
									Contains: []DocumentationFilterTag{
										{
											Key:   "type.*",
											Value: "(guide|tutorial)",
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid regex tag key",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Tags: AuditTagsChecks{
									Contains: []DocumentationFilterTag{
										{
											Key:   "[invalid",
											Value: "guide",
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "key must be a valid regex pattern",
		},
		{
			name: "invalid regex tag value",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
							},
							Checks: AuditChecks{
								Tags: AuditTagsChecks{
									Contains: []DocumentationFilterTag{
										{
											Key:   "type",
											Value: "(guide|incomplete",
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "value must be a valid regex pattern",
		},
		{
			name: "section without document",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source:  "backend",
									Section: "Installation",
								},
							},
							Checks: AuditChecks{
								Content: AuditContentChecks{
									Exists: true,
								},
							},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "document must be set if",
		},
		{
			name: "all check types",
			cfg: &Config{
				Audit: &Audit{
					Rules: []AuditRule{
						{
							Documentation: []DocumentationFilter{
								{
									Source: "backend",
								},
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
										{
											Key:   "type",
											Value: "guide",
										},
									},
								},
							},
						},
					},
				},
			},
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
