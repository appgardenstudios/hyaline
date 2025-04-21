package github

import (
	"context"
	"database/sql"
	"hyaline/internal/sqlite"

	"github.com/google/go-github/v71/github"
)

func InsertIssue(ref string, token string, systemID string, db *sql.DB) (err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// Get Issue
	client := github.NewClient(nil).WithAuthToken(token)
	issue, _, err := client.Issues.Get(context.Background(), owner, repo, id)
	if err != nil {
		return
	}

	// Insert PR
	err = sqlite.InsertIssue(sqlite.Issue{
		ID:       ref,
		SystemID: systemID,
		Title:    *issue.Title,
		Body:     *issue.Body,
	}, db)
	if err != nil {
		return
	}

	return
}
