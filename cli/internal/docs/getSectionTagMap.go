package docs

import "hyaline/internal/sqlite"

// GetSectionTagMap creates a map of section keys to their associated tags
func GetSectionTagMap(sectionTags []sqlite.SECTIONTAG) map[string][]FilteredTag {
	tagMap := make(map[string][]FilteredTag)
	for _, tag := range sectionTags {
		key := tag.SourceID + "/" + tag.DocumentID + "#" + tag.SectionID
		filteredTag := FilteredTag{
			Key:   tag.TagKey,
			Value: tag.TagValue,
		}
		tagMap[key] = append(tagMap[key], filteredTag)
	}
	return tagMap
}
