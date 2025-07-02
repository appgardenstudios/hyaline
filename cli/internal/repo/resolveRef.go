package repo

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// ResolveAlias resolves a reference using alias resolution (branch/tag names)
// This tries to resolve references as branches, remote branches, or tags in that order
func ResolveAlias(r *git.Repository, reference string) (hash *plumbing.Hash, err error) {
	slog.Debug("repo.ResolveAlias attempting to resolve reference", "reference", reference)

	// If it starts with refs/, pass it through directly
	if strings.HasPrefix(reference, "refs/") {
		slog.Debug("repo.ResolveAlias reference starts with refs/, passing through directly", "reference", reference)
		return r.ResolveRevision(plumbing.Revision(reference))
	}

	// Try as local branch
	localBranchRef := "refs/heads/" + reference
	slog.Debug("repo.ResolveAlias trying local branch", "ref", localBranchRef)
	hash, err = r.ResolveRevision(plumbing.Revision(localBranchRef))
	if err == nil {
		slog.Debug("repo.ResolveAlias resolved as local branch", "reference", reference, "hash", hash.String())
		return
	}
	slog.Debug("repo.ResolveAlias could not resolve as local branch", "error", err)

	// Try remote branches - check if reference starts with any remote name
	remotes, err := r.Remotes()
	if err != nil {
		slog.Debug("repo.ResolveAlias could not list remotes", "error", err)
		return
	}

	for _, remote := range remotes {
		remoteName := remote.Config().Name
		// Check if reference starts with this remote name followed by "/"
		if strings.HasPrefix(reference, remoteName+"/") {
			remoteBranchRef := "refs/remotes/" + reference
			slog.Debug("repo.ResolveAlias trying remote branch", "ref", remoteBranchRef, "remote", remoteName)
			hash, err = r.ResolveRevision(plumbing.Revision(remoteBranchRef))
			if err == nil {
				slog.Debug("repo.ResolveAlias resolved as remote branch", "reference", reference, "hash", hash.String())
				return
			}
			break // Found matching remote, no need to check others
		}
	}

	// If there's only one remote, try it
	if len(remotes) == 1 {
		remoteName := remotes[0].Config().Name
		remoteBranchRef := "refs/remotes/" + remoteName + "/" + reference
		slog.Debug("repo.ResolveAlias trying single remote fallback", "ref", remoteBranchRef, "remote", remoteName)
		hash, err = r.ResolveRevision(plumbing.Revision(remoteBranchRef))
		if err == nil {
			slog.Debug("repo.ResolveAlias resolved as single remote fallback", "reference", reference, "hash", hash.String())
			return
		}
		slog.Debug("repo.ResolveAlias could not resolve as single remote fallback", "error", err)
	}

	// Try as tag
	tagRef := "refs/tags/" + reference
	slog.Debug("repo.ResolveAlias trying tag", "ref", tagRef)
	hash, err = r.ResolveRevision(plumbing.Revision(tagRef))
	if err == nil {
		slog.Debug("repo.ResolveAlias resolved as tag", "reference", reference, "hash", hash.String())
		return
	}
	slog.Debug("repo.ResolveAlias could not resolve as tag", "error", err)

	err = fmt.Errorf("could not resolve reference '%s' as branch or tag", reference)
	return
}

// ResolveExplicitRef resolves a reference by passing it directly to go-git's ResolveRevision
// This allows explicit specification of commit hashes or fully qualified references
func ResolveExplicitRef(r *git.Repository, ref string) (hash *plumbing.Hash, err error) {
	slog.Debug("repo.ResolveExplicitRef resolving explicit ref", "ref", ref)
	hash, err = r.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		slog.Debug("repo.ResolveExplicitRef could not resolve ref", "error", err, "ref", ref)
		return
	}
	slog.Debug("repo.ResolveExplicitRef resolved ref", "ref", ref, "hash", hash.String())
	return
}

// ResolveRef resolves a reference using either alias or explicit resolution
// Exactly one of alias or ref must be provided (non-empty)
func ResolveRef(r *git.Repository, alias string, ref string) (hash *plumbing.Hash, err error) {
	slog.Debug("repo.ResolveRef resolving reference", "alias", alias, "ref", ref)

	// Validate that exactly one is provided
	if alias != "" && ref != "" {
		err = fmt.Errorf("both alias (%s) and ref (%s) provided - exactly one must be specified", alias, ref)
		return
	}
	if alias == "" && ref == "" {
		err = fmt.Errorf("neither alias nor ref provided - exactly one must be specified")
		return
	}

	// If ref is provided, use explicit resolution
	if ref != "" {
		return ResolveExplicitRef(r, ref)
	}

	// If alias is provided, use alias resolution
	return ResolveAlias(r, alias)
}
