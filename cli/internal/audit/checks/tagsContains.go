package checks

import (
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"regexp"
)

// TagsContains validates that required tag key-value pairs exist
func TagsContains(actualTags []docs.FilteredTag, requiredTags []config.DocumentationFilterTag) (bool, string) {
	for _, required := range requiredTags {
		found := false
		for _, actual := range actualTags {
			// Use regex matching for both key and value
			keyMatches, err := regexp.MatchString(required.Key, actual.Key)
			if err != nil {
				return false, fmt.Sprintf("Invalid regex pattern for tag key: %s", required.Key)
			}
			
			valueMatches, err := regexp.MatchString(required.Value, actual.Value)
			if err != nil {
				return false, fmt.Sprintf("Invalid regex pattern for tag value: %s", required.Value)
			}
			
			if keyMatches && valueMatches {
				found = true
				break
			}
		}
		
		if !found {
			return false, fmt.Sprintf("Required tag with key pattern '%s' and value pattern '%s' not found.", required.Key, required.Value)
		}
	}
	return true, ""
}