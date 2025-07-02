package repo

import (
	"context"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff(r *git.Repository, head plumbing.Hash, base plumbing.Hash) (diff object.Changes, err error) {
	slog.Debug("repo.GetDiff getting diff", "head", head.String(), "base", base.String())

	headCommit, err := r.CommitObject(head)
	if err != nil {
		slog.Debug("repo.GetDiff could not get head commit", "error", err, "head", head.String())
		return
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get head tree", "error", err)
		return
	}

	baseCommit, err := r.CommitObject(base)
	if err != nil {
		slog.Debug("repo.GetDiff could not get base commit", "error", err, "base", base.String())
		return
	}

	baseTree, err := baseCommit.Tree()
	if err != nil {
		slog.Debug("repo.GetDiff could not get base tree", "error", err)
		return
	}

	diff, err = object.DiffTreeWithOptions(context.Background(), baseTree, headTree, object.DefaultDiffTreeOptions)
	if err != nil {
		slog.Debug("repo.GetDiff could not get diff", "error", err)
		return
	}

	return
}
