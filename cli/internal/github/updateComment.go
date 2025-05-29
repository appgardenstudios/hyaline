package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

func UpdateComment(ref string, body string, token string) (err error) {

	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Update Comment
	client := github.NewClient(nil).WithAuthToken(token)
	_, _, err = client.Issues.EditComment(context.Background(), owner, repo, id, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return
	}

	return
}
