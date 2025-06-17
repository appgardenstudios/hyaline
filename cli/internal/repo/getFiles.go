package repo

import (
	"log/slog"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetFiles(branch string, r *git.Repository, cb func(*object.File) error) (err error) {
	ref, err := r.ResolveRevision(plumbing.Revision(branch))
	if err != nil {
		// If resolution fails and branch doesn't already have a remote prefix, try with origin/
		if !strings.Contains(branch, "/") {
			slog.Debug("repo.GetFiles could not resolve branch, trying with origin/", "error", err, "branch", branch)
			ref, err = r.ResolveRevision(plumbing.Revision("origin/" + branch))
			if err != nil {
				slog.Debug("repo.GetFiles could not resolve branch with origin/ prefix", "error", err, "branch", "origin/"+branch)
				return
			}
		} else {
			slog.Debug("repo.GetFiles could not resolve branch", "error", err, "branch", branch)
			return
		}
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
