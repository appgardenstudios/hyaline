package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

// GetPullRequestFiles retrieves the list of files changed in a GitHub Pull Request
func GetPullRequestFiles(ref string, token string) ([]*github.CommitFile, string, string, error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return nil, "", "", err
	}

	// Get GitHub client
	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	// Get PR details to get base and head commit SHAs
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, int(id))
	if err != nil {
		return nil, "", "", err
	}

	baseSHA := pr.Base.SHA
	headSHA := pr.Head.SHA

	// List all files changed in the PR
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, int(id), nil)
	if err != nil {
		return nil, "", "", err
	}

	return files, *baseSHA, *headSHA, nil
}
