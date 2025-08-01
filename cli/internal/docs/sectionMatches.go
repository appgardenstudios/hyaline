package docs

import (
	"hyaline/internal/config"

	"github.com/bmatcuk/doublestar/v4"
)

func SectionMatches(sectionID string, documentID string, sourceID string, tags []FilteredTag, filter *config.DocumentationFilter, strict bool) bool {
	source, document, section := filter.GetParts()

	// If section is blank just match on document.
	// If document is a match include this section.
	if section == "" {
		// If strict matching is enabled and no section is set, the section is not a match.
		if strict && section == "" {
			return false
		}
		// Document matches if we match source OR source and document (if document is set)
		if doublestar.MatchUnvalidated(source, sourceID) {
			if document != "" {
				if doublestar.MatchUnvalidated(document, documentID) {
					return tagMatches(tags, filter.Tags)
				}
			} else {
				return tagMatches(tags, filter.Tags)
			}
		}
		return false
	}

	// Else see if all 3 match
	if doublestar.MatchUnvalidated(source, sourceID) &&
		doublestar.MatchUnvalidated(document, documentID) &&
		doublestar.MatchUnvalidated(section, sectionID) {
		// If everything matches we also have to match tags
		return tagMatches(tags, filter.Tags)
	}
	return false
}
