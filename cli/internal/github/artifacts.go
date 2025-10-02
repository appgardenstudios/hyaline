package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v74/github"
)

func DownloadLatestArtifact(githubRepo string, artifactName string, githubToken string, destDir string) (string, error) {
	slog.Debug("github.DownloadLatestArtifact starting",
		slog.Group("args",
			"githubRepo", githubRepo,
			"artifactName", artifactName,
			"destDir", destDir,
		),
	)

	parts := strings.Split(githubRepo, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid github repo format: %s", githubRepo)
	}
	owner, repo := parts[0], parts[1]

	client := github.NewClient(nil).WithAuthToken(githubToken)
	artifacts, _, err := client.Actions.ListArtifacts(context.Background(), owner, repo, &github.ListArtifactsOptions{
		Name: &artifactName,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to list artifacts: %w", err)
	}

	if len(artifacts.Artifacts) == 0 {
		return "", fmt.Errorf("artifact not found: %s", artifactName)
	}

	targetArtifact := artifacts.Artifacts[0]

	slog.Debug("github.DownloadLatestArtifact found artifact", "artifactID", targetArtifact.GetID())

	downloadURL, _, err := client.Actions.DownloadArtifact(context.Background(), owner, repo, targetArtifact.GetID(), 3)
	if err != nil {
		return "", fmt.Errorf("failed to get artifact download URL: %w", err)
	}

	slog.Debug("github.DownloadLatestArtifact downloading artifact", "url", downloadURL.String())

	req, err := http.NewRequest("GET", downloadURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download artifact: %w", err)
	}
	defer resp.Body.Close()

	zipPath := filepath.Join(destDir, "artifact.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save zip file: %w", err)
	}

	slog.Debug("github.DownloadLatestArtifact downloaded artifact", "path", zipPath)

	return zipPath, nil
}
