package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

type Issue struct {
	Title string
	Body  string
}

func GetIssue(ref string, token string) (issue *Issue, err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get Issue
	client := github.NewClient(nil).WithAuthToken(token)
	rawIssue, _, err := client.Issues.Get(context.Background(), owner, repo, int(id))
	if err != nil {
		return
	}

	// Get body
	issue = &Issue{
		Title: *rawIssue.Title,
		Body:  *rawIssue.Body,
	}

	return
}
