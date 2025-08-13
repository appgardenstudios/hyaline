package repo

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func setupTestRepo(t *testing.T) (*git.Repository, plumbing.Hash) {
	// Create an in-memory repository with worktree
	storage := memory.NewStorage()
	fs := memfs.New()
	r, err := git.Init(storage, fs)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test commit
	wt, err := r.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	// Add a remote
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://example.com/repo.git"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a test file
	file, err := fs.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write([]byte("test content"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Add the file
	_, err = wt.Add("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Create initial commit
	hash, err := wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create main branch reference (git.Init creates master by default)
	mainRef := plumbing.NewHashReference("refs/heads/main", hash)
	err = r.Storer.SetReference(mainRef)
	if err != nil {
		t.Fatal(err)
	}

	return r, hash
}

func TestResolveAlias(t *testing.T) {
	r, hash := setupTestRepo(t)

	// Test resolving already qualified reference
	resolved, err := ResolveAlias(r, "refs/heads/main")
	if err != nil {
		t.Errorf("ResolveAlias failed to resolve refs/heads/main: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveAlias resolved refs/heads/main to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Create a remote branch reference
	remoteRef := plumbing.NewHashReference("refs/remotes/origin/main", hash)
	err = r.Storer.SetReference(remoteRef)
	if err != nil {
		t.Fatal(err)
	}

	// Test resolving remote branch when only one remote exists
	resolved, err = ResolveAlias(r, "main")
	if err != nil {
		t.Errorf("ResolveAlias failed to resolve 'main' with remote fallback: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveAlias resolved 'main' to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Create a tag reference
	tagRef := plumbing.NewHashReference("refs/tags/v1.0.0", hash)
	err = r.Storer.SetReference(tagRef)
	if err != nil {
		t.Fatal(err)
	}

	// Test resolving tag
	resolved, err = ResolveAlias(r, "v1.0.0")
	if err != nil {
		t.Errorf("ResolveAlias failed to resolve tag 'v1.0.0': %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveAlias resolved 'v1.0.0' to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test resolving non-existent reference
	_, err = ResolveAlias(r, "nonexistent")
	if err == nil {
		t.Error("ResolveAlias should have failed for non-existent reference")
	}
}

func TestResolveAliasMultipleRemotes(t *testing.T) {
	r, hash := setupTestRepo(t)

	// Add a second remote
	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{"https://example.com/upstream.git"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create remote branch references for both remotes
	originRef := plumbing.NewHashReference("refs/remotes/origin/features/new-feature", hash)
	err = r.Storer.SetReference(originRef)
	if err != nil {
		t.Fatal(err)
	}

	upstreamRef := plumbing.NewHashReference("refs/remotes/upstream/features/new-feature", hash)
	err = r.Storer.SetReference(upstreamRef)
	if err != nil {
		t.Fatal(err)
	}

	// Test that we can resolve remote branches with explicit remote prefix
	resolved, err := ResolveAlias(r, "origin/features/new-feature")
	if err != nil {
		t.Errorf("ResolveAlias failed to resolve 'origin/features/new-feature': %v", err)
	}
	if resolved == nil {
		t.Error("ResolveAlias returned nil hash for 'origin/features/new-feature'")
		return
	}
	if *resolved != hash {
		t.Errorf("ResolveAlias resolved 'origin/features/new-feature' to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test that we can resolve the other remote too
	resolved, err = ResolveAlias(r, "upstream/features/new-feature")
	if err != nil {
		t.Errorf("ResolveAlias failed to resolve 'upstream/features/new-feature': %v", err)
	}
	if resolved == nil {
		t.Error("ResolveAlias returned nil hash for 'upstream/features/new-feature'")
		return
	}
	if *resolved != hash {
		t.Errorf("ResolveAlias resolved 'upstream/features/new-feature' to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test that plain branch name without remote prefix still fails when multiple remotes exist
	_, err = ResolveAlias(r, "features/new-feature")
	if err == nil {
		t.Error("ResolveAlias should have failed when multiple remotes exist and no local branch found")
	}
}

func TestResolveExplicitRef(t *testing.T) {
	r, hash := setupTestRepo(t)

	// Test resolving explicit commit hash
	resolved, err := ResolveExplicitRef(r, hash.String())
	if err != nil {
		t.Errorf("ResolveExplicitRef failed to resolve commit hash: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveExplicitRef resolved to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test resolving explicit reference
	resolved, err = ResolveExplicitRef(r, "refs/heads/main")
	if err != nil {
		t.Errorf("ResolveExplicitRef failed to resolve refs/heads/main: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveExplicitRef resolved to wrong hash. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test resolving non-existent reference
	_, err = ResolveExplicitRef(r, "nonexistent")
	if err == nil {
		t.Error("ResolveExplicitRef should have failed for non-existent reference")
	}
}

func TestResolveRef(t *testing.T) {
	r, hash := setupTestRepo(t)

	// Create a remote branch reference
	remoteRef := plumbing.NewHashReference("refs/remotes/origin/main", hash)
	err := r.Storer.SetReference(remoteRef)
	if err != nil {
		t.Fatal(err)
	}

	// Test alias resolution
	resolved, err := ResolveRef(r, "main", "")
	if err != nil {
		t.Errorf("ResolveRef failed with alias resolution: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveRef returned wrong hash for alias. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test explicit resolution
	resolved, err = ResolveRef(r, "", hash.String())
	if err != nil {
		t.Errorf("ResolveRef failed with explicit resolution: %v", err)
	}
	if *resolved != hash {
		t.Errorf("ResolveRef returned wrong hash for explicit ref. Expected %s, got %s", hash.String(), resolved.String())
	}

	// Test validation - both provided
	_, err = ResolveRef(r, "main", hash.String())
	if err == nil {
		t.Error("ResolveRef should have failed when both alias and ref are provided")
	}

	// Test validation - neither provided
	_, err = ResolveRef(r, "", "")
	if err == nil {
		t.Error("ResolveRef should have failed when neither alias nor ref are provided")
	}
}
