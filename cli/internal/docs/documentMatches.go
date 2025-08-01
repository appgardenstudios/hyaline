package docs

import (
	"hyaline/internal/config"

	"github.com/bmatcuk/doublestar/v4"
)

func DocumentMatches(documentID string, sourceID string, tags []FilteredTag, filter *config.DocumentationFilter) bool {
	source, document, section := filter.GetParts()

	// Document does not match filters that specify sections
	if section != "" {
		return false
	}
	matches := false

	// Document matches if we match source OR source and document (if document is set)
	if doublestar.MatchUnvalidated(source, sourceID) {
		if document != "" {
			if doublestar.MatchUnvalidated(document, documentID) {
				matches = true
			}
		} else {
			matches = true
		}
	}

	// If document matches and there are tags it must match at least one tag
	if matches {
		matches = tagMatches(tags, filter.Tags)
	}

	return matches
}