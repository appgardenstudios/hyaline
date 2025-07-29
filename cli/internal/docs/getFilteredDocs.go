package docs

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"

	"github.com/bmatcuk/doublestar/v4"
)

type FilteredDoc struct {
	Document *sqlite.DOCUMENT
	Tags     []FilteredTag
	Sections []FilteredSection
}

type FilteredSection struct {
	Section  *sqlite.SECTION
	Tags     []FilteredTag
	Sections []FilteredSection
}

type FilteredTag struct {
	Key   string
	Value string
}

func GetFilteredDocs(cfg *config.CheckDocumentation, db *sqlite.Queries) (docs []*FilteredDoc, err error) {
	// Get/format data from sqlite
	documents, err := db.GetAllDocuments(context.Background())
	if err != nil {
		slog.Debug("docs.GetFilteredDocs could not get all documents", "error", err)
		return
	}
	documentTags, err := db.GetAllDocumentTags(context.Background())
	if err != nil {
		slog.Debug("docs.GetFilteredDocs could not get all document tags", "error", err)
		return
	}
	documentTagMap := make(map[string][]FilteredTag)
	for _, tag := range documentTags {
		key := tag.SourceID + "/" + tag.DocumentID
		filteredTag := FilteredTag{
			Key:   tag.TagKey,
			Value: tag.TagValue,
		}
		entry, ok := documentTagMap[key]
		if ok {
			entry = append(entry, filteredTag)
			documentTagMap[key] = entry
		} else {
			documentTagMap[key] = []FilteredTag{filteredTag}
		}
	}
	sections, err := db.GetAllSections(context.Background())
	if err != nil {
		slog.Debug("docs.GetFilteredDocs could not get all sections", "error", err)
		return
	}
	documentSectionMap := make(map[string][]*sqlite.SECTION)
	// Note: mapped sections MUST be in peer order for subsequent logic to work
	for _, section := range sections {
		key := section.SourceID + "/" + section.DocumentID
		entry, ok := documentSectionMap[key]
		if ok {
			entry = append(entry, &section)
			documentSectionMap[key] = entry
		} else {
			documentSectionMap[key] = []*sqlite.SECTION{&section}
		}
	}
	sectionTags, err := db.GetAllSectionTags(context.Background())
	if err != nil {
		slog.Debug("docs.GetFilteredDocs could not get all section tags", "error", err)
		return
	}
	sectionTagMap := make(map[string][]FilteredTag)
	for _, tag := range sectionTags {
		key := tag.SourceID + "/" + tag.DocumentID + "#" + tag.SectionID
		filteredTag := FilteredTag{
			Key:   tag.TagKey,
			Value: tag.TagValue,
		}
		entry, ok := sectionTagMap[key]
		if ok {
			entry = append(entry, filteredTag)
			sectionTagMap[key] = entry
		} else {
			sectionTagMap[key] = []FilteredTag{filteredTag}
		}
	}

	// Get and filter all docs from database
	for _, document := range documents {
		documentKey := document.SourceID + "/" + document.ID
		if documentIsIncluded(document.ID, document.SourceID, documentTagMap[documentKey], cfg) {
			// Add document and ALL sections if document included
			docs = append(docs, &FilteredDoc{
				Document: &document,
				Tags:     documentTagMap[documentKey],
				Sections: getDocumentSections(documentSectionMap[documentKey], ""),
			})
		} else {
			// Else get filtered sections and add document if there is at least 1 section included
			// Note: filteredSections will contain only sections that match along with their parents up to the root
			filteredSections := getFilteredDocumentSections(documentSectionMap[documentKey], "", sectionTagMap, cfg)
			if len(filteredSections) > 0 {
				docs = append(docs, &FilteredDoc{
					Document: &document,
					Tags:     documentTagMap[documentKey],
					Sections: filteredSections,
				})
			}
		}
	}

	return
}

func documentIsIncluded(documentID string, sourceID string, tags []FilteredTag, cfg *config.CheckDocumentation) bool {
	isIncluded := false

	// Document is included if it matches at least one include
	for _, include := range cfg.Include {
		if DocumentMatches(documentID, sourceID, tags, &include) {
			isIncluded = true
			break
		}
	}

	// If document matched an include AND at least one exclude, it is excluded
	if isIncluded {
		for _, exclude := range cfg.Exclude {
			if DocumentMatches(documentID, sourceID, tags, &exclude) {
				return false
			}
		}
	}

	return isIncluded
}

func DocumentMatches(documentID string, sourceID string, tags []FilteredTag, filter *config.CheckDocumentationFilter) bool {
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

func sectionIsIncluded(sectionID string, documentID string, sourceID string, tags []FilteredTag, cfg *config.CheckDocumentation) bool {
	isIncluded := false

	// Section is included if it matches at least one include
	for _, include := range cfg.Include {
		if SectionMatches(sectionID, documentID, sourceID, tags, &include) {
			isIncluded = true
			break
		}
	}

	// If section matched an include AND at least one exclude, it is excluded
	if isIncluded {
		for _, exclude := range cfg.Exclude {
			if SectionMatches(sectionID, documentID, sourceID, tags, &exclude) {
				return false
			}
		}
	}

	return isIncluded
}

func SectionMatches(sectionID string, documentID string, sourceID string, tags []FilteredTag, filter *config.CheckDocumentationFilter) bool {
	source, document, section := filter.GetParts()

	// If section is blank just match on document.
	// If document is a match include this section.
	if section == "" {
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

func tagMatches(tags []FilteredTag, filterTags []config.CheckDocumentationFilterTag) bool {
	// If there are no filter tags we always match regardless
	if len(filterTags) == 0 {
		return true
	}

	// Else we need to match at least one tag pair
	for _, filterTag := range filterTags {
		for _, tag := range tags {
			if filterTag.Key == tag.Key && filterTag.Value == tag.Value {
				return true
			}
		}
	}

	return false
}

// Note: This requires documentSections to be in peer order
func getDocumentSections(documentSections []*sqlite.SECTION, parent string) (sections []FilteredSection) {
	for _, section := range documentSections {
		if section.ParentID == parent {
			sections = append(sections, FilteredSection{
				Section:  section,
				Sections: getDocumentSections(documentSections, section.ID),
			})
		}
	}
	return
}

// Note: This requires documentSections to be in peer order
func getFilteredDocumentSections(documentSections []*sqlite.SECTION, parent string, sectionTagMap map[string][]FilteredTag, cfg *config.CheckDocumentation) (sections []FilteredSection) {
	for _, section := range documentSections {
		if section.ParentID == parent {
			childSections := getFilteredDocumentSections(documentSections, section.ID, sectionTagMap, cfg)
			tags := sectionTagMap[section.SourceID+"/"+section.DocumentID+"#"+section.ID]
			// If a child matches always include this section
			if len(childSections) > 0 {
				sections = append(sections, FilteredSection{
					Section:  section,
					Tags:     tags,
					Sections: childSections,
				})
			} else {
				// Else include this section if it should be included
				if sectionIsIncluded(section.ID, section.DocumentID, section.SourceID, tags, cfg) {
					sections = append(sections, FilteredSection{
						Section: section,
						Tags:    tags,
					})
				}
			}
		}
	}
	return
}
