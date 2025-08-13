package checks

import (
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"regexp"
)

// TagsContains validates that required tag key-value pairs exist
func TagsContains(actualTags []docs.FilteredTag, requiredTags []config.DocumentationFilterTag) (bool, string, error) {
	for _, required := range requiredTags {
		found := false
		for _, actual := range actualTags {
			// Use regex matching for both key and value
			keyMatches, err := regexp.MatchString(required.Key, actual.Key)
			if err != nil {
				return false, "", fmt.Errorf("invalid regex pattern for tag key '%s': %w", required.Key, err)
			}

			valueMatches, err := regexp.MatchString(required.Value, actual.Value)
			if err != nil {
				return false, "", fmt.Errorf("invalid regex pattern for tag value '%s': %w", required.Value, err)
			}

			if keyMatches && valueMatches {
				found = true
				break
			}
		}

		if !found {
			return false, fmt.Sprintf("Required tag with key pattern '%s' and value pattern '%s' not found.", required.Key, required.Value), nil
		}
	}
	return true, "", nil
}
