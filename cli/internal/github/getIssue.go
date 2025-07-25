package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

func GetIssue(ref string, token string) (body *string, err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get Issue
	client := github.NewClient(nil).WithAuthToken(token)
	issue, _, err := client.Issues.Get(context.Background(), owner, repo, int(id))
	if err != nil {
		return
	}

	// Get body
	body = issue.Body

	return
}
