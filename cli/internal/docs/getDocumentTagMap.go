package docs

import "hyaline/internal/sqlite"

// GetDocumentTagMap creates a map of document keys to their associated tags
func GetDocumentTagMap(documentTags []sqlite.DOCUMENTTAG) map[string][]FilteredTag {
	tagMap := make(map[string][]FilteredTag)
	for _, tag := range documentTags {
		key := tag.SourceID + "/" + tag.DocumentID
		filteredTag := FilteredTag{
			Key:   tag.TagKey,
			Value: tag.TagValue,
		}
		tagMap[key] = append(tagMap[key], filteredTag)
	}
	return tagMap
}
