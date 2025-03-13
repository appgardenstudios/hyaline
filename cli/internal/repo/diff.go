package repo

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff(path string, head string, base string) (diff object.Changes, err error) {
	// Get our absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Debug("repo.GetDiff could not determine absolute code path", "error", err, "path", path)
		return
	}
	absPath += string(os.PathSeparator)
	slog.Debug("repo.GetDiff extracting code from path", "absPath", absPath)

	// Open our git repo
	repo, err := git.PlainOpen(absPath)
	if err != nil {
		slog.Debug("repo.GetDiff could not open git repo", "error", err, "path", path)
		return
	}

	// Get a list of files change between head and base
	headRef, err := repo.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		slog.Debug("repo.GetDiff could not resolve head", "error", err)
		return
	}
	headCommit, err := repo.CommitObject(*headRef)
	if err != nil {
		slog.Debug("repo.GetDiff could not get head commit", "error", err)
		return
	}
	headTree, err := headCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get head tree", "error", err)
		return
	}
	baseRef, err := repo.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		slog.Debug("repo.GetDiff could not resolve base", "error", err)
		return
	}
	baseCommit, err := repo.CommitObject(*baseRef)
	if err != nil {
		slog.Debug("repo.GetDiff could not get base commit", "error", err)
		return
	}
	baseTree, err := baseCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get base tree", "error", err)
		return
	}
	// Note that we will eventually want to support renames via object.DiffTreeWithOptions
	diffs, err := object.DiffTree(baseTree, headTree)
	if err != nil {
		slog.Debug("repo.GetDiff could not get diff", "error", err)
		return
	}

	return diffs, nil
}
