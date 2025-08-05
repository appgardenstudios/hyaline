package code

import (
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

type FilteredFile struct {
	Filename         string
	OriginalFilename string
	Action           Action
	Contents         []byte
	OriginalContents []byte
	Diff             string
}

type Action string

const (
	ActionInsert Action = "Insert"
	ActionModify Action = "Modify"
	ActionRename Action = "Rename"
	ActionDelete Action = "Delete"
)

func GetFilteredDiff(r *git.Repository, head plumbing.Hash, base plumbing.Hash, cfg *config.CheckCode) (filteredFiles []FilteredFile, changedFiles map[string]struct{}, err error) {
	changedFiles = make(map[string]struct{})

	// Get our diff
	diff, err := repo.GetDiff(r, head, base)
	if err != nil {
		slog.Debug("code.GetFilteredDiff could not get diff", "error", err)
		return
	}

	// Examine each change in the diff
	for _, change := range diff {
		slog.Debug("code.GetFilteredDiff processing diff", "diff", change.String())
		var action merkletrie.Action
		action, err = change.Action()
		if err != nil {
			slog.Debug("code.GetFilteredDiff could not retrieve action for diff", "error", err, "diff", change)
			return
		}
		var from *object.File
		var to *object.File
		from, to, err = change.Files()
		if err != nil {
			slog.Debug("code.GetFilteredDiff could not retrieve files for diff", "error", err, "diff", change)
			return
		}

		// Handle filtering based on action since from/to presence is dependent on the action
		switch action {
		case merkletrie.Insert:
			changedFiles[change.To.Name] = struct{}{}
			if config.PathIsIncluded(change.To.Name, cfg.Include, cfg.Exclude) {
				var toBytes []byte
				toBytes, err = repo.GetBlobBytes(to.Blob)
				if err != nil {
					slog.Debug("code.GetFilteredDiff could not retrieve to blob from insert diff", "error", err)
					return
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					Filename: change.To.Name,
					Action:   ActionInsert,
					Contents: toBytes,
				})
			}
		case merkletrie.Modify:
			changedFiles[change.From.Name] = struct{}{}
			changedFiles[change.To.Name] = struct{}{}
			if config.PathIsIncluded(change.To.Name, cfg.Include, cfg.Exclude) {
				var fromBytes []byte
				fromBytes, err = repo.GetBlobBytes(from.Blob)
				if err != nil {
					slog.Debug("code.GetFilteredDiff could not retrieve from blob from modify diff", "error", err)
					return
				}
				var toBytes []byte
				toBytes, err = repo.GetBlobBytes(to.Blob)
				if err != nil {
					slog.Debug("code.GetFilteredDiff could not retrieve to blob from modify diff", "error", err)
					return
				}
				var action Action
				if change.From.Name != change.To.Name {
					action = ActionRename
				} else {
					action = ActionModify
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					Filename:         change.To.Name,
					OriginalFilename: change.From.Name,
					Action:           action,
					Contents:         toBytes,
					OriginalContents: fromBytes,
				})
			}
		case merkletrie.Delete:
			changedFiles[change.From.Name] = struct{}{}
			if config.PathIsIncluded(change.From.Name, cfg.Include, cfg.Exclude) {
				var fromBytes []byte
				fromBytes, err = repo.GetBlobBytes(from.Blob)
				if err != nil {
					slog.Debug("code.GetFilteredDiff could not retrieve from blob from delete diff", "error", err)
					return
				}
				filteredFiles = append(filteredFiles, FilteredFile{
					OriginalFilename: change.From.Name,
					Action:           ActionDelete,
					OriginalContents: fromBytes,
				})
			}
		}
	}

	return
}
