package code

import (
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/github"
	"log/slog"
	"strings"
)

// GetFilteredPR retrieves filtered files from a GitHub Pull Request
func GetFilteredPR(pullRequest string, token string, cfg *config.CheckCode) (filteredFiles []FilteredFile, changedFiles map[string]struct{}, err error) {
	changedFiles = make(map[string]struct{})

	// Parse PR reference to get owner/repo for API calls
	parts := strings.Split(pullRequest, "/")
	if len(parts) != 3 {
		slog.Debug("code.GetFilteredPR invalid pull request reference", "pull-request", pullRequest)
		return nil, nil, errors.New("pull request reference must be in format OWNER/REPO/NUMBER")
	}
	owner := parts[0]
	repo := parts[1]

	// Get PR files and commit SHAs
	files, baseSHA, headSHA, err := github.GetPullRequestFiles(pullRequest, token)
	if err != nil {
		slog.Debug("code.GetFilteredPR could not get PR files", "error", err, "pull-request", pullRequest)
		return
	}
	slog.Info("Retrieved PR files", "files", len(files), "pull-request", pullRequest)

	// Examine each change in the PR
	for _, file := range files {
		slog.Debug("code.GetFilteredPR processing file", "filename", *file.Filename, "status", *file.Status)
		
		filename := *file.Filename
		var originalFilename string
		if file.PreviousFilename != nil {
			originalFilename = *file.PreviousFilename
		}

		// Determine action and track changed files
		var action Action
		switch *file.Status {
		case "added":
			action = ActionInsert
			changedFiles[filename] = struct{}{}
		case "modified":
			if originalFilename != "" && originalFilename != filename {
				action = ActionRename
				changedFiles[originalFilename] = struct{}{}
				changedFiles[filename] = struct{}{}
			} else {
				action = ActionModify
				changedFiles[filename] = struct{}{}
			}
		case "removed":
			action = ActionDelete
			if originalFilename != "" {
				changedFiles[originalFilename] = struct{}{}
			} else {
				changedFiles[filename] = struct{}{}
			}
		case "renamed":
			action = ActionRename
			changedFiles[originalFilename] = struct{}{}
			changedFiles[filename] = struct{}{}
		default:
			slog.Debug("code.GetFilteredPR unknown file status", "status", *file.Status, "filename", filename)
			continue
		}

		// Handle filtering based on action since from/to presence is dependent on the action
		var pathToCheck string
		switch action {
		case ActionInsert:
			pathToCheck = filename
		case ActionModify:
			pathToCheck = filename
		case ActionRename:
			pathToCheck = filename
		case ActionDelete:
			pathToCheck = originalFilename
		}

		if config.PathIsIncluded(pathToCheck, cfg.Include, cfg.Exclude) {
			filteredFile := FilteredFile{
				Filename:         filename,
				OriginalFilename: originalFilename,
				Action:           action,
			}

			// Get file contents for non-deleted files
			if action != ActionDelete {
				contents, err := github.GetFileContents(owner, repo, filename, headSHA, token)
				if err != nil {
					slog.Debug("code.GetFilteredPR could not get head file contents", "error", err, "filename", filename, "sha", headSHA)
					return nil, nil, err
				}
				filteredFile.Contents = contents
			}

			// Get original file contents for modified/renamed/deleted files
			if action == ActionModify || action == ActionRename || action == ActionDelete {
				sourceFilename := filename
				if originalFilename != "" {
					sourceFilename = originalFilename
				}
				originalContents, err := github.GetFileContents(owner, repo, sourceFilename, baseSHA, token)
				if err != nil {
					slog.Debug("code.GetFilteredPR could not get base file contents", "error", err, "filename", sourceFilename, "sha", baseSHA)
					return nil, nil, err
				}
				filteredFile.OriginalContents = originalContents
			}

			filteredFiles = append(filteredFiles, filteredFile)
		}
	}

	slog.Info("Filtered PR files", "total", len(files), "filtered", len(filteredFiles), "pull-request", pullRequest)

	return
}