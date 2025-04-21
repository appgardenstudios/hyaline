package repo

import (
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff(r *git.Repository, head string, base string) (diff object.Changes, err error) {
	// Get a list of files change between head and base
	headRef, err := r.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		slog.Debug("repo.GetDiff could not resolve head", "error", err, "head", head)
		return
	}
	headCommit, err := r.CommitObject(*headRef)
	if err != nil {
		slog.Debug("repo.GetDiff could not get head commit", "error", err, "head", head)
		return
	}
	headTree, err := headCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get head tree", "error", err, "head", head)
		return
	}
	baseRef, err := r.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		slog.Debug("repo.GetDiff could not resolve base", "error", err, "base", base)
		return
	}
	baseCommit, err := r.CommitObject(*baseRef)
	if err != nil {
		slog.Debug("repo.GetDiff could not get base commit", "error", err, "base", base)
		return
	}
	baseTree, err := baseCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get base tree", "error", err, "base", base)
		return
	}
	// Note that we will eventually want to support renames via object.DiffTreeWithOptions
	diff, err = object.DiffTree(baseTree, headTree)
	if err != nil {
		slog.Debug("repo.GetDiff could not get diff", "error", err, "head", head, "base", base)
		return
	}

	return
}
