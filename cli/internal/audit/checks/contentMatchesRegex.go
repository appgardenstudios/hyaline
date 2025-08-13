package checks

import (
	"fmt"
	"regexp"
)

// ContentMatchesRegex validates content against a regex pattern
func ContentMatchesRegex(content string, pattern string) (bool, string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, "", err
	}

	matches := regex.MatchString(content)
	if matches {
		return true, "", nil
	} else {
		return false, fmt.Sprintf("Content does not match the regex pattern: %s", pattern), nil
	}
}
