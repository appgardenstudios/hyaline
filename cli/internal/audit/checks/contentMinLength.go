package checks

import "fmt"

// ContentMinLength validates that content meets minimum length requirement
func ContentMinLength(content string, minLength int) (bool, string) {
	contentLength := len(content)
	pass := contentLength >= minLength

	if !pass {
		return false, fmt.Sprintf("Content length is %d, minimum required is %d.", contentLength, minLength)
	}

	return true, ""
}
