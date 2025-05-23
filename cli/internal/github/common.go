package github

import (
	"errors"
	"strconv"
	"strings"
)

// Parse a GitHub PR or Issue reference (OWNER/REPO/ID) and return the various parts
func parseReference(ref string) (owner string, repo string, id int, err error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 3 {
		err = errors.New("reference must contain 3 parts: OWNER/REPO/ID")
		return
	}
	owner = parts[0]
	repo = parts[1]
	id, err = strconv.Atoi(parts[2])

	return
}
