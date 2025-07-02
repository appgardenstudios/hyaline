package repo

import (
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetFiles(ref plumbing.Hash, r *git.Repository, cb func(*object.File) error) (err error) {
	slog.Debug("repo.GetFiles getting files from ref", "ref", ref.String())

	commit, err := r.CommitObject(ref)
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
