package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func validateCheck(cfg *Config) error {
	// If check was not defined in the config don't check anything.
	// Check is not always required, so actions requiring the config need to check for nil themselves
	if cfg.Check == nil {
		return nil
	}

	// Check code
	if len(cfg.Check.Code.Include) == 0 {
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
	if len(cfg.Check.Documentation.Include) == 0 {
		return errors.New("check.documentation.include must contain at least one entry, none found")
	}
	for i, include := range cfg.Check.Documentation.Include {
		if err := validateCheckDocumentationFilter(fmt.Sprintf("check.documentation.include[%d]", i), include); err != nil {
			return err
		}
	}
	for i, exclude := range cfg.Check.Documentation.Exclude {
		if err := validateCheckDocumentationFilter(fmt.Sprintf("check.documentation.exclude[%d]", i), exclude); err != nil {
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
		if err := validateCheckDocumentationFilter(fmt.Sprintf("%s[%d].documentation", location, i), entry.Documentation); err != nil {
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

func validateCheckDocumentationFilter(location string, filter CheckDocumentationFilter) error {
	// URI trumps source/document/section
	if filter.URI != "" {
		if !strings.HasPrefix(filter.URI, "document://") {
			return fmt.Errorf("%s.uri must start with document://, found: %s", location, filter.URI)
		}
		source, remainder, found := strings.Cut(strings.TrimPrefix(filter.URI, "document://"), "/")
		if !found {
			return fmt.Errorf("%s.uri must contain at least one /, found: %s", location, filter.URI)
		}
		if source == "" || !doublestar.ValidatePattern(source) {
			return fmt.Errorf("%s.uri must contain a valid source pattern, found: %s in %s", location, source, filter.URI)
		}
		document, section, found := strings.Cut(remainder, "#")
		if document == "" || !doublestar.ValidatePattern(document) {
			return fmt.Errorf("%s.uri must contain a valid document pattern, found: %s in %s", location, document, filter.URI)
		}
		if found {
			if section == "" || !doublestar.ValidatePattern(section) {
				return fmt.Errorf("%s.uri must contain a valid section pattern, found: %s in %s", location, section, filter.URI)
			}
		}
	} else {
		if filter.Source == "" || !doublestar.ValidatePattern(filter.Source) {
			return fmt.Errorf("%s.source must be a valid pattern, found: %s", location, filter.Source)
		}
		if filter.Document != "" {
			if !doublestar.ValidatePattern(filter.Document) {
				return fmt.Errorf("%s.document must be a valid pattern, found: %s", location, filter.Document)
			}
		}
		if filter.Section != "" {
			if filter.Document == "" {
				return fmt.Errorf("%s.document must be set if %s.section is set", location, location)
			}
			if !doublestar.ValidatePattern(filter.Section) {
				return fmt.Errorf("%s.section must be a valid pattern, found: %s", location, filter.Section)
			}
		}
	}

	// Check tags
	keyRegex := regexp.MustCompile(metadataTagKeyRegex)
	valueRegex := regexp.MustCompile(metadataTagValueRegex)
	for i, tag := range filter.Tags {
		if !keyRegex.MatchString(tag.Key) {
			return fmt.Errorf("%s.tags[%d].key must match regex /%s/, found: %s", location, i, metadataTagKeyRegex, tag.Key)
		}
		if !valueRegex.MatchString(tag.Value) {
			return fmt.Errorf("%s.tags[%d].value must match regex /%s/, found: %s", location, i, metadataTagValueRegex, tag.Value)
		}
	}

	return nil
}
