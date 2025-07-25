package docs

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"

	"github.com/bmatcuk/doublestar/v4"
)

type FilteredDoc struct {
	Document *sqlite.DOCUMENT // Document can be nil if a section matches but a document does not
	Sections []*sqlite.SECTION
}

type FilteredTag struct {
	Key   string
	Value string
}

func GetFilteredDocs(cfg *config.CheckDocumentation, db *sqlite.Queries) (docs []*FilteredDoc, err error) {
	docMap := make(map[string]*FilteredDoc)

	// Get and filter all docs from database
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
	for _, document := range documents {
		tags := documentTagMap[document.SourceID+"/"+document.ID]
		if documentIsIncluded(document.ID, document.SourceID, tags, cfg) {
			docMap[document.SourceID+"/"+document.ID] = &FilteredDoc{
				Document: &document,
			}
		}
	}

	// Get and filter all sections from the database
	sections, err := db.GetAllSections(context.Background())
	if err != nil {
		slog.Debug("docs.GetFilteredDocs could not get all sections", "error", err)
		return
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
	for _, section := range sections {
		tags := documentTagMap[section.SourceID+"/"+section.DocumentID+"#"+section.ID]
		if sectionIsIncluded(section.ID, section.DocumentID, section.SourceID, tags, cfg) {
			filteredDoc, ok := docMap[section.SourceID+"/"+section.DocumentID]
			if ok {
				filteredDoc.Sections = append(filteredDoc.Sections, &section)
			} else {
				docMap[section.SourceID+"/"+section.DocumentID] = &FilteredDoc{
					Sections: []*sqlite.SECTION{&section},
				}
			}
		}
	}

	// Loop through map and populate docs array
	for _, doc := range docMap {
		docs = append(docs, doc)
	}

	return
}

func documentIsIncluded(documentID string, sourceID string, tags []FilteredTag, cfg *config.CheckDocumentation) bool {
	isIncluded := false

	// Document is included if it matches at least one include
	for _, include := range cfg.Include {
		if documentMatches(documentID, sourceID, tags, &include) {
			isIncluded = true
			break
		}
	}

	// If document matched an include AND at least one exclude, it is excluded
	if isIncluded {
		for _, exclude := range cfg.Exclude {
			if documentMatches(documentID, sourceID, tags, &exclude) {
				return false
			}
		}
	}

	return isIncluded
}

func documentMatches(documentID string, sourceID string, tags []FilteredTag, filter *config.CheckDocumentationFilter) bool {
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
		if sectionMatches(sectionID, documentID, sourceID, tags, &include) {
			isIncluded = true
			break
		}
	}

	// If section matched an include AND at least one exclude, it is excluded
	if isIncluded {
		for _, exclude := range cfg.Exclude {
			if sectionMatches(sectionID, documentID, sourceID, tags, &exclude) {
				return false
			}
		}
	}

	return isIncluded
}

func sectionMatches(sectionID string, documentID string, sourceID string, tags []FilteredTag, filter *config.CheckDocumentationFilter) bool {
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
