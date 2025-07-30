package config

import (
	"fmt"
	"regexp"
	"strconv"
)

func validateAudit(cfg *Config) error {
	// If audit was not defined in the config don't check anything.
	// Audit is not always required, so actions requiring the config need to check for nil themselves
	if cfg.Audit == nil {
		return nil
	}

	// Check that we have at least one rule
	if len(cfg.Audit.Rules) == 0 {
		return fmt.Errorf("audit.rules must contain at least one entry, none found")
	}

	// Validate each rule
	for i, rule := range cfg.Audit.Rules {
		if err := validateAuditRule(fmt.Sprintf("audit.rules[%d]", i), &rule, i); err != nil {
			return err
		}
	}

	return nil
}

func validateAuditRule(location string, rule *AuditRule, index int) error {
	// Set default ID if not provided
	if rule.ID == "" {
		rule.ID = strconv.Itoa(index)
	}

	// Check documentation filters
	if len(rule.Documentation) == 0 {
		return fmt.Errorf("%s.documentation must contain at least one entry, none found", location)
	}
	for i, filter := range rule.Documentation {
		if err := validateDocumentationFilter(fmt.Sprintf("%s.documentation[%d]", location, i), &filter); err != nil {
			return err
		}
	}

	// Check ignore filters
	for i, filter := range rule.Ignore {
		if err := validateDocumentationFilter(fmt.Sprintf("%s.ignore[%d]", location, i), &filter); err != nil {
			return err
		}
	}

	// Check that at least one check type is specified
	if err := validateAuditChecks(fmt.Sprintf("%s.checks", location), &rule.Checks); err != nil {
		return err
	}

	return nil
}

func validateAuditChecks(location string, checks *AuditChecks) error {
	hasAtLeastOneCheck := false

	// Content checks
	if checks.Content.Exists {
		hasAtLeastOneCheck = true
	}
	if checks.Content.MinLength > 0 {
		hasAtLeastOneCheck = true
	} else if checks.Content.MinLength < 0 {
		return fmt.Errorf("%s.content.min-length must be non-negative, found: %d", location, checks.Content.MinLength)
	}
	if checks.Content.MatchesRegex != "" {
		hasAtLeastOneCheck = true
		if _, err := regexp.Compile(checks.Content.MatchesRegex); err != nil {
			return fmt.Errorf("%s.content.matches-regex must be a valid regex pattern, found: %s, error: %v", location, checks.Content.MatchesRegex, err)
		}
	}
	if checks.Content.MatchesPrompt != "" {
		hasAtLeastOneCheck = true
	}
	if checks.Content.MatchesPurpose {
		hasAtLeastOneCheck = true
	}

	// Purpose checks
	if checks.Purpose.Exists {
		hasAtLeastOneCheck = true
	}

	// Tags checks
	if len(checks.Tags.Contains) > 0 {
		hasAtLeastOneCheck = true
		for i, tag := range checks.Tags.Contains {
			if _, err := regexp.Compile(tag.Key); err != nil {
				return fmt.Errorf("%s.tags.contains[%d].key must be a valid regex pattern, found: %s, error: %v", location, i, tag.Key, err)
			}
			if _, err := regexp.Compile(tag.Value); err != nil {
				return fmt.Errorf("%s.tags.contains[%d].value must be a valid regex pattern, found: %s, error: %v", location, i, tag.Value, err)
			}
		}
	}

	if !hasAtLeastOneCheck {
		return fmt.Errorf("%s must specify at least one check type", location)
	}

	return nil
}
