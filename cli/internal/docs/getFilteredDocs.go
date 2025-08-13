package docs

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
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
	documentTagMap := GetDocumentTagMap(documentTags)
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
	sectionTagMap := GetSectionTagMap(sectionTags)

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

func sectionIsIncluded(sectionID string, documentID string, sourceID string, tags []FilteredTag, cfg *config.CheckDocumentation) bool {
	isIncluded := false

	// Section is included if it matches at least one include
	for _, include := range cfg.Include {
		if SectionMatches(sectionID, documentID, sourceID, tags, &include, false) {
			isIncluded = true
			break
		}
	}

	// If section matched an include AND at least one exclude, it is excluded
	if isIncluded {
		for _, exclude := range cfg.Exclude {
			if SectionMatches(sectionID, documentID, sourceID, tags, &exclude, false) {
				return false
			}
		}
	}

	return isIncluded
}

func tagMatches(tags []FilteredTag, filterTags []config.DocumentationFilterTag) bool {
	// If there are no filter tags we always match
	if len(filterTags) == 0 {
		return true
	}

	// Else we need to find a match for all filtered tag pairs
	for _, filterTag := range filterTags {
		found := false
		for _, tag := range tags {
			if filterTag.Key == tag.Key && filterTag.Value == tag.Value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
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
