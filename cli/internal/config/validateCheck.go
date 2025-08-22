package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

func validateCheck(cfg *Config) error {
	// If check was not defined in the config don't check anything.
	// Check is not always required, so actions requiring the config need to check for nil themselves
	if cfg.Check == nil {
		return nil
	}

	// Check code
	if !cfg.Check.Disabled && len(cfg.Check.Code.Include) == 0 {
		return errors.New("check.code.include must contain at least one entry, none found")
	}
	for i, include := range cfg.Check.Code.Include {
		if include == "" || !doublestar.ValidatePattern(include) {
			return fmt.Errorf("check.code.include[%d] must be a valid pattern, found: %s", i, include)
		}
	}
	for i, exclude := range cfg.Check.Code.Exclude {
		if exclude == "" || !doublestar.ValidatePattern(exclude) {
			return fmt.Errorf("check.code.exclude[%d] must be a valid pattern, found: %s", i, exclude)
		}
	}

	// Check documentation
	if !cfg.Check.Disabled && len(cfg.Check.Documentation.Include) == 0 {
		return errors.New("check.documentation.include must contain at least one entry, none found")
	}
	for i, include := range cfg.Check.Documentation.Include {
		if err := validateDocumentationFilter(fmt.Sprintf("check.documentation.include[%d]", i), &include); err != nil {
			return err
		}
	}
	for i, exclude := range cfg.Check.Documentation.Exclude {
		if err := validateDocumentationFilter(fmt.Sprintf("check.documentation.exclude[%d]", i), &exclude); err != nil {
			return err
		}
	}

	// Check options
	if cfg.Check.Options.DetectDocumentationUpdates.Source != "" {
		source := cfg.Check.Options.DetectDocumentationUpdates.Source
		if !regexp.MustCompile(sourceIDRegex).MatchString(source) {
			return fmt.Errorf("extract.options.detectDocumentationUpdates.source must match regex /%s/, found: %s", sourceIDRegex, source)
		}
	}
	if err := validateCheckUpdateIf("check.options.updateIf.touched", cfg.Check.Options.UpdateIf.Touched); err != nil {
		return err
	}
	if err := validateCheckUpdateIf("check.options.updateIf.added", cfg.Check.Options.UpdateIf.Added); err != nil {
		return err
	}
	if err := validateCheckUpdateIf("check.options.updateIf.modified", cfg.Check.Options.UpdateIf.Modified); err != nil {
		return err
	}
	if err := validateCheckUpdateIf("check.options.updateIf.deleted", cfg.Check.Options.UpdateIf.Deleted); err != nil {
		return err
	}
	if err := validateCheckUpdateIf("check.options.updateIf.renamed", cfg.Check.Options.UpdateIf.Renamed); err != nil {
		return err
	}

	return nil
}

func validateCheckUpdateIf(location string, entries []CheckOptionsUpdateIfEntry) error {
	for i, entry := range entries {
		if err := validateCheckCodeFilter(fmt.Sprintf("%s[%d].code", location, i), entry.Code); err != nil {
			return err
		}
		if err := validateDocumentationFilter(fmt.Sprintf("%s[%d].documentation", location, i), &entry.Documentation); err != nil {
			return err
		}
	}

	return nil
}

func validateCheckCodeFilter(location string, filter CheckCodeFilter) error {
	if filter.Path == "" || !doublestar.ValidatePattern(filter.Path) {
		return fmt.Errorf("%s.path must be a valid pattern, found: %s", location, filter.Path)
	}

	return nil
}
