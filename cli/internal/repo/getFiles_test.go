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

func TestGetFiles(t *testing.T) {
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

	// Test that GetFiles can resolve "main" with fallback to "origin/main"
	fileCount := 0
	err = GetFiles("main", r, func(f *object.File) error {
		fileCount++
		return nil
	})
	if err != nil {
		t.Errorf("GetFiles failed to resolve 'main' with origin/ fallback: %v", err)
	}
	if fileCount != 1 {
		t.Errorf("Expected 1 file, got %d", fileCount)
	}

	// Test that GetFiles works with explicit origin/ prefix
	err = GetFiles("origin/main", r, func(f *object.File) error {
		return nil
	})
	if err != nil {
		t.Errorf("GetFiles failed to resolve 'origin/main': %v", err)
	}
}
