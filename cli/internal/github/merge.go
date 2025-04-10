package github

import (
	"database/sql"
	"hyaline/internal/sqlite"
)

func MergePullRequests(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all PULL_REQUEST entries from the source we will be copying
	pullRequestsToCopy, err := sqlite.GetAllPullRequest(systemID, source)
	if err != nil {
		return err
	}

	// Copy each PULL_REQUEST from source to dest
	for _, p := range pullRequestsToCopy {
		// Delete the existing PULL_REQUEST (if any)
		err := sqlite.DeletePullRequest(p.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy PULL_REQUEST
		err = sqlite.InsertPullRequest(*p, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func MergeIssues(systemID string, source *sql.DB, dest *sql.DB) error {
	// Get all ISSUE entries from the source we will be copying
	issuesToCopy, err := sqlite.GetAllIssue(systemID, source)
	if err != nil {
		return err
	}

	// Copy each ISSUE from source to dest
	for _, i := range issuesToCopy {
		// Delete the existing ISSUE (if any)
		err := sqlite.DeleteIssue(i.ID, systemID, dest)
		if err != nil {
			return err
		}

		// Copy ISSUE
		err = sqlite.InsertIssue(*i, dest)
		if err != nil {
			return err
		}
	}

	return nil
}
