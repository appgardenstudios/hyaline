package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func PathIsIncluded(path string, includes []string, excludes []string) bool {
	for _, include := range includes {
		if doublestar.MatchUnvalidated(include, path) {
			for _, exclude := range excludes {
				if doublestar.MatchUnvalidated(exclude, path) {
					return false
				}
			}
			return true
		}
	}

	return false
}

func validateDocumentationFilter(location string, filter *DocumentationFilter) error {
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
