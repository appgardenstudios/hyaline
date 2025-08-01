package github

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v71/github"
)

// GetFileContents retrieves the contents of a file at a specific commit SHA
func GetFileContents(owner, repo, filename, sha, token string) ([]byte, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filename, &github.RepositoryContentGetOptions{
		Ref: sha,
	})
	if err != nil {
		return nil, err
	}

	// Handle directory case (should not happen in PR file list, but be safe)
	if fileContent.GetType() != "file" {
		slog.Debug("github.GetFileContents expected file but got different type", "type", fileContent.GetType(), "filename", filename)
		return nil, nil
	}

	// GetContent() automatically handles base64 decoding
	content, err := fileContent.GetContent()
	if err != nil {
		slog.Debug("github.GetFileContents could not get content", "error", err, "filename", filename)
		return nil, err
	}
	if content == "" {
		// Handle large files or binary files that don't have inline content
		slog.Debug("github.GetFileContents file has no inline content", "filename", filename)
		return nil, nil
	}

	// Convert string to []byte - GetContent() already handles base64 decoding
	return []byte(content), nil
}