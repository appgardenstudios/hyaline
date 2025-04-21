package github

import (
	"context"
	"database/sql"
	"hyaline/internal/sqlite"

	"github.com/google/go-github/v71/github"
)

func InsertPullRequest(ref string, token string, systemID string, db *sql.DB) (err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get PR
	client := github.NewClient(nil).WithAuthToken(token)
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, id)
	if err != nil {
		return
	}

	// Insert PR
	err = sqlite.InsertPullRequest(sqlite.PullRequest{
		ID:       ref,
		SystemID: systemID,
		Title:    *pr.Title,
		Body:     *pr.Body,
	}, db)
	if err != nil {
		return
	}

	return
}
