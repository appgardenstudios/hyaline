package checks

import (
	"fmt"
	"regexp"
)

// ContentMatchesRegex validates content against a regex pattern
func ContentMatchesRegex(content string, pattern string) (bool, string) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Sprintf("Invalid regex pattern: %v", err)
	}

	matches := regex.MatchString(content)
	if matches {
		return true, ""
	} else {
		return false, fmt.Sprintf("Content does not match the regex pattern: %s", pattern)
	}
}