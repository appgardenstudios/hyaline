package code

import (
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"log/slog"
	"path/filepath"

	"github.com/go-git/go-git/v5"
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

func GetFilteredDiff(path string, head string, headRef string, base string, baseRef string, cfg *config.CheckCode) (filteredFiles []FilteredFile, changedFiles map[string]struct{}, err error) {
	changedFiles = make(map[string]struct{})

	// Open repo already on disk
	var absPath string
	absPath, err = filepath.Abs(path)
	if err != nil {
		slog.Debug("code.GetFilteredDiff could not determine absolute path", "error", err, "path", path)
		return
	}
	slog.Info("Opening repo on disk", "absPath", absPath)
	var r *git.Repository
	r, err = git.PlainOpen(absPath)
	if err != nil {
		slog.Debug("code.GetFilteredDiff could not open git repo", "error", err, "path", path)
		return
	}

	// Resolve head and base references
	resolvedHead, err := repo.ResolveRef(r, head, headRef)
	if err != nil {
		slog.Debug("code.GetFilteredDiff could not resolve head reference", "error", err)
		return
	}
	resolvedBase, err := repo.ResolveRef(r, base, baseRef)
	if err != nil {
		slog.Debug("code.GetFilteredDiff could not resolve base reference", "error", err)
		return
	}

	// Get our diff
	diff, err := repo.GetDiff(r, *resolvedHead, *resolvedBase)
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
