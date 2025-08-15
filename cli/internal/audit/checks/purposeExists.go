package checks

import "strings"

// PurposeExists validates that a purpose is defined and non-empty
func PurposeExists(purpose string) (bool, string) {
	if strings.TrimSpace(purpose) == "" {
		return false, "Purpose is not defined or is empty."
	}
	return true, ""
}
