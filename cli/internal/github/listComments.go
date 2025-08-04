package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

type Comment struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
}

func ListComments(ref string, token string) (comments []Comment, err error) {
	// Parse reference
	owner, repo, id, err := parseReference(ref)
	if err != nil {
		return
	}

	// List comments with pagination
	client := github.NewClient(nil).WithAuthToken(token)
	
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		commentList, resp, err := client.Issues.ListComments(context.Background(), owner, repo, int(id), opts)
		if err != nil {
			return nil, err
		}

		for _, comment := range commentList {
			if comment.ID != nil && comment.Body != nil {
				comments = append(comments, Comment{
					ID:   *comment.ID,
					Body: *comment.Body,
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return
}