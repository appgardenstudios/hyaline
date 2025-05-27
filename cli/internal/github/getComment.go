package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

func GetComment(ref string, token string) (body string, err error) {

	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get comment
	client := github.NewClient(nil).WithAuthToken(token)
	comment, _, err := client.Issues.GetComment(context.Background(), owner, repo, id)
	if err != nil {
		return
	}

	// Get body
	body = *comment.Body

	return
}
