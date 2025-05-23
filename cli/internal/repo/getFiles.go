package repo

import (
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetFiles(branch string, r *git.Repository, cb func(*object.File) error) (err error) {
	ref, err := r.ResolveRevision(plumbing.Revision(branch))
	if err != nil {
		slog.Debug("repo.GetFiles could not resolve head", "error", err)
		return
	}
	commit, err := r.CommitObject(*ref)
	if err != nil {
		slog.Debug("repo.GetFiles could not get head commit", "error", err)
		return
	}
	tree, err := commit.Tree()
	if err != nil {
		slog.Debug("repo.GetFiles could not get head tree", "error", err)
		return
	}

	// Call cb for each file
	err = tree.Files().ForEach(cb)

	return
}
