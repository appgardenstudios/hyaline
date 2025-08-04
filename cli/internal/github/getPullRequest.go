package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

type PullRequest struct {
	Title   string
	Body    string
	HeadSHA string
}

func GetPullRequest(ref string, token string) (pr *PullRequest, err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get PR
	client := github.NewClient(nil).WithAuthToken(token)
	rawPr, _, err := client.PullRequests.Get(context.Background(), owner, repo, int(id))
	if err != nil {
		return
	}

	// Get PR details
	pr = &PullRequest{
		Title:   *rawPr.Title,
		Body:    *rawPr.Body,
		HeadSHA: *rawPr.Head.SHA,
	}

	return
}
