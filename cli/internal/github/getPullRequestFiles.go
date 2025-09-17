package github

import (
	"context"

	"github.com/google/go-github/v74/github"
)

// GetPullRequestFiles retrieves the list of files changed in a GitHub Pull Request
func GetPullRequestFiles(ref string, token string) ([]*github.CommitFile, error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return nil, err
	}

	// Get GitHub client
	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	// List all files changed in the PR with pagination
	var allFiles []*github.CommitFile
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, int(id), opts)
		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, files...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allFiles, nil
}
