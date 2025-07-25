package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

func GetPullRequest(ref string, token string) (body *string, err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get PR
	client := github.NewClient(nil).WithAuthToken(token)
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, int(id))
	if err != nil {
		return
	}

	// Get body
	body = pr.Body

	return
}
