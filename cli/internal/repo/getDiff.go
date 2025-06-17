package repo

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff(r *git.Repository, head string, base string) (diff object.Changes, err error) {
	// Get a list of files change between head and base
	headRef, err := r.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		// If resolution fails and branch doesn't already have a remote prefix, try with origin/
		if !strings.Contains(head, "/") {
			slog.Debug("repo.GetDiff could not resolve head, trying with origin/", "error", err, "head", head)
			headRef, err = r.ResolveRevision(plumbing.Revision("origin/" + head))
			if err != nil {
				slog.Debug("repo.GetDiff could not resolve head with origin/ prefix", "error", err, "head", "origin/"+head)
				return
			}
		} else {
			slog.Debug("repo.GetDiff could not resolve head", "error", err, "head", head)
			return
		}
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
		// If resolution fails and branch doesn't already have a remote prefix, try with origin/
		if !strings.Contains(base, "/") {
			slog.Debug("repo.GetDiff could not resolve base, trying with origin/", "error", err, "base", base)
			baseRef, err = r.ResolveRevision(plumbing.Revision("origin/" + base))
			if err != nil {
				slog.Debug("repo.GetDiff could not resolve base with origin/ prefix", "error", err, "base", "origin/"+base)
				return
			}
		} else {
			slog.Debug("repo.GetDiff could not resolve base", "error", err, "base", base)
			return
		}
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
	diff, err = object.DiffTreeWithOptions(context.Background(), baseTree, headTree, object.DefaultDiffTreeOptions)
	if err != nil {
		slog.Debug("repo.GetDiff could not get diff", "error", err, "head", head, "base", base)
		return
	}

	return
}
