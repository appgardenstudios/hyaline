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

func TestGetDiff(t *testing.T) {
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

	// Create a reference for origin/main
	ref := plumbing.NewHashReference("refs/remotes/origin/main", hash)
	err = r.Storer.SetReference(ref)
	if err != nil {
		t.Fatal(err)
	}

	// Test that GetDiff can resolve head "main" with fallback to "origin/main"
	_, err = GetDiff(r, "main", "origin/main")
	if err != nil {
		t.Errorf("GetDiff failed to resolve head 'main' with origin/ fallback: %v", err)
	}

	// Test that GetDiff can resolve base "main" with fallback to "origin/main"
	_, err = GetDiff(r, "origin/main", "main")
	if err != nil {
		t.Errorf("GetDiff failed to resolve base 'main' with origin/ fallback: %v", err)
	}
}
