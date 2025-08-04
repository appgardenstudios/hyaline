package code

import (
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/github"
	"log/slog"
)

// GetFilteredPRFiles retrieves filtered files from a GitHub Pull Request
func GetFilteredPRFiles(pullRequest string, token string, cfg *config.CheckCode) (filteredFiles []FilteredFile, changedFiles map[string]struct{}, err error) {
	changedFiles = make(map[string]struct{})

	// Get PR files
	files, err := github.GetPullRequestFiles(pullRequest, token)
	if err != nil {
		slog.Debug("code.GetFilteredPRFiles could not get PR files", "error", err, "pull-request", pullRequest)
		return
	}
	slog.Info("Retrieved PR files", "files", len(files), "pull-request", pullRequest)

	// Examine each change in the PR
	for _, file := range files {
		slog.Debug("code.GetFilteredPRFiles processing file", "filename", *file.Filename, "status", *file.Status)

		filename := *file.Filename
		var originalFilename string
		if file.PreviousFilename != nil {
			originalFilename = *file.PreviousFilename
		}

		// Process each file based on its status
		switch *file.Status {
		case "added":
			changedFiles[filename] = struct{}{}
			if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
				diff := fmt.Sprintf("--- /dev/null\n+++ b/%s\n", filename)
				if file.Patch != nil {
					diff += *file.Patch
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					Filename: filename,
					Action:   ActionInsert,
					Diff:     diff,
				})
			}
		case "modified", "changed":
			changedFiles[filename] = struct{}{}
			if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
				diff := fmt.Sprintf("--- a/%s\n+++ b/%s\n", filename, filename)
				if file.Patch != nil {
					diff += *file.Patch
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					Filename: filename,
					Action:   ActionModify,
					Diff:     diff,
				})
			}
		case "removed":
			changedFiles[filename] = struct{}{}
			if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
				diff := fmt.Sprintf("--- a/%s\n+++ /dev/null\n", filename)
				if file.Patch != nil {
					diff += *file.Patch
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					OriginalFilename: filename,
					Action:           ActionDelete,
					Diff:             diff,
				})
			}
		case "renamed":
			changedFiles[originalFilename] = struct{}{}
			changedFiles[filename] = struct{}{}
			if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
				diff := fmt.Sprintf("--- a/%s\n+++ b/%s\n", originalFilename, filename)
				if file.Patch != nil {
					diff += *file.Patch
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					Filename:         filename,
					OriginalFilename: originalFilename,
					Action:           ActionRename,
					Diff:             diff,
				})
			}
		case "copied":
			// When PreviousFilename is set, treat as renamed
			// When no PreviousFilename is set, treat as added
			if originalFilename != "" {
				// Has previous filename, treat as renamed
				changedFiles[originalFilename] = struct{}{}
				changedFiles[filename] = struct{}{}
				if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
					diff := fmt.Sprintf("--- a/%s\n+++ b/%s\n", originalFilename, filename)
					if file.Patch != nil {
						diff += *file.Patch
					}
					filteredFiles = append(filteredFiles, FilteredFile{
						Filename:         filename,
						OriginalFilename: originalFilename,
						Action:           ActionRename,
						Diff:             diff,
					})
				}
			} else {
				// No previous filename, treat as added
				changedFiles[filename] = struct{}{}
				if config.PathIsIncluded(filename, cfg.Include, cfg.Exclude) {
					diff := fmt.Sprintf("--- /dev/null\n+++ b/%s\n", filename)
					if file.Patch != nil {
						diff += *file.Patch
					}
					filteredFiles = append(filteredFiles, FilteredFile{
						Filename: filename,
						Action:   ActionInsert,
						Diff:     diff,
					})
				}
			}
		default:
			slog.Debug("code.GetFilteredPRFiles unknown file status", "status", *file.Status, "filename", filename)
			continue
		}
	}

	slog.Info("Filtered PR files", "total", len(files), "filtered", len(filteredFiles), "pull-request", pullRequest)

	return
}
