package github

import (
	"context"

	"github.com/google/go-github/v74/github"
)

func AddComment(ref string, body string, token string) (err error) {

	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Add Comment
	client := github.NewClient(nil).WithAuthToken(token)
	_, _, err = client.Issues.CreateComment(context.Background(), owner, repo, int(id), &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return
	}

	return
}
