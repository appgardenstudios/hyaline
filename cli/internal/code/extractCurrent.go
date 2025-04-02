package code

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

func ExtractCurrent(system *config.System, db *sql.DB) (err error) {
	// Process each code source
	for _, c := range system.Code {
		slog.Debug("code.ExtractCurrent extracting code", "system", system.ID, "code", c.ID)
		// Insert Code
		err = sqlite.InsertCode(sqlite.Code{
			ID:       c.ID,
			SystemID: system.ID,
			Path:     c.FsOptions.Path,
		}, db)
		if err != nil {
			slog.Debug("code.ExtractCurrent could not insert code", "error", err, "code", c.ID)
			return err
		}

		// Extract based on the extractor
		switch c.Extractor {
		case config.ExtractorFs:
			err = ExtractCurrentFs(system.ID, &c, db)
			if err != nil {
				slog.Debug("code.ExtractCurrent could not extract code using fs extractor", "error", err, "code", c.ID)
				return
			}
		case config.ExtractorGit:
			err = ExtractCurrentGit(system.ID, &c, db)
			if err != nil {
				slog.Debug("code.ExtractCurrent could not extract code using git extractor", "error", err, "code", c.ID)
				return
			}
		default:
			slog.Debug("code.ExtractCurrent unknown extractor", "extractor", c.Extractor.String(), "code", c.ID)
			return errors.New("Unknown Extractor '" + c.Extractor.String() + "' for code " + c.ID)
		}
	}

	return
}

func ExtractCurrentFs(systemID string, c *config.Code, db *sql.DB) (err error) {
	// Get our absolute path
	absPath, err := filepath.Abs(c.FsOptions.Path)
	if err != nil {
		slog.Debug("code.ExtractCurrentFs could not determine absolute code path", "error", err, "path", c.FsOptions.Path)
		return err
	}
	slog.Debug("code.ExtractCurrentFs extracting code from path", "absPath", absPath)

	// Get our root FS
	// We use a root FS so symlinks and relative paths don't escape our path
	// https://pkg.go.dev/os@go1.24.1#Root
	root, err := os.OpenRoot(absPath)
	if err != nil {
		slog.Debug("code.ExtractCurrentFs could not open fs root", "error", err, "path", c.FsOptions.Path)
		return err
	}
	fsRoot := root.FS()

	// Our set of files (as a map so we don't get dupes)
	files := map[string]struct{}{}

	// Loop through our includes and get files
	for _, include := range c.Include {
		slog.Debug("code.ExtractCurrentFs extracting code using include", "include", include, "code", c.ID)

		// Get matched files
		matches, err := doublestar.Glob(fsRoot, include)
		if err != nil {
			slog.Debug("code.ExtractCurrentFs could not find code files with include", "include", include, "error", err)
			return err
		}

		// Loop through files and add those that aren't in our excludes
		for _, file := range matches {
			// See if we have a match for at least one of our excludes
			match := false
			for _, exclude := range c.Exclude {
				match = doublestar.MatchUnvalidated(exclude, file)
				if match {
					slog.Debug("code.ExtractCurrentFs file excluded", "file", file, "exclude", exclude)
					break
				}
			}
			if !match {
				files[file] = struct{}{}
			}
		}
	}

	// Insert files
	for file := range files {
		contents, err := fs.ReadFile(fsRoot, file)
		if err != nil {
			slog.Debug("code.ExtractCurrentFs could not read code file", "error", err, "file", file)
			return err
		}

		err = sqlite.InsertFile(sqlite.File{
			ID:       file,
			CodeID:   c.ID,
			SystemID: systemID,
			RawData:  string(contents),
		}, db)
		if err != nil {
			slog.Debug("code.ExtractCurrentFs could not insert file", "error", err)
			return err
		}
	}

	return
}

func ExtractCurrentGit(systemID string, c *config.Code, db *sql.DB) (err error) {
	// Initialize go-git (on disk or in mem)
	r, err := getRepo(c.GitOptions)
	if err != nil {
		slog.Debug("code.ExtractCurrentGit could not get repo", "error", err)
		return
	}

	// Extract files from branch
	branch := "main"
	if c.GitOptions.Branch != "" {
		branch = c.GitOptions.Branch
	}
	ref, err := r.ResolveRevision(plumbing.Revision(branch))
	if err != nil {
		slog.Debug("code.ExtractCurrentGit could not resolve head", "error", err)
		return
	}
	commit, err := r.CommitObject(*ref)
	if err != nil {
		slog.Debug("code.ExtractCurrentGit could not get head commit", "error", err)
		return
	}
	tree, err := commit.Tree()
	if err != nil {
		slog.Debug("code.ExtractCurrentGit could not get head tree", "error", err)
		return
	}

	// Get and save code files
	tree.Files().ForEach(func(f *object.File) error {
		for _, include := range c.Include {
			if doublestar.MatchUnvalidated(include, f.Name) {
				// If excluded, skip this file and continue
				excluded := false
				for _, exclude := range c.Exclude {
					if doublestar.MatchUnvalidated(exclude, f.Name) {
						excluded = true
						break
					}
				}
				if excluded {
					continue
				}

				// Get contents of file and insert into db
				var bytes []byte
				bytes, err = repo.GetBlobBytes(f.Blob)
				if err != nil {
					slog.Debug("code.ExtractCurrentGit could not get blob bytes", "error", err)
					return err
				}
				err = sqlite.InsertFile(sqlite.File{
					ID:       f.Name,
					CodeID:   c.ID,
					SystemID: systemID,
					RawData:  string(bytes),
				}, db)
				if err != nil {
					slog.Debug("code.ExtractCurrentFs could not insert file", "error", err)
					return err
				}
			}
		}
		return nil
	})

	return
}

func getRepo(options config.GitOptions) (r *git.Repository, err error) {
	if options.Clone {
		// Ensure remote repo is passed in
		if options.Repo == "" {
			err = errors.New("git.repo is required to be set if git.clone is true")
			return
		}
		// Create cloneOptions
		cloneOptions := &git.CloneOptions{
			URL: options.Repo,
		}

		// Add http auth if password is set
		if options.HTTPAuth.Password != "" {
			username := "git"
			if options.HTTPAuth.Username != "" {
				username = options.HTTPAuth.Username
			}
			cloneOptions.Auth = &http.BasicAuth{
				Username: username,
				Password: options.HTTPAuth.Password,
			}
		}

		// Add ssh auth if PEM is set
		if options.SSHAuth.PEM != "" {
			user := "git"
			if options.SSHAuth.User != "" {
				user = options.SSHAuth.User
			}
			var keys *ssh.PublicKeys
			keys, err = ssh.NewPublicKeys(user, []byte(options.SSHAuth.PEM), options.SSHAuth.Password)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not parse git ssh PEM", "error", err)
				return
			}
			cloneOptions.Auth = keys
		}

		if options.Path != "" {
			// Clone to disk
			var absPath string
			absPath, err = filepath.Abs(options.Path)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not determine absolute path", "error", err, "path", options.Path)
				return
			}
			r, err = git.PlainClone(absPath, false, cloneOptions)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could clone repo", "error", err, "path", options.Path, "repo", options.Repo)
				return
			}
		} else {
			// Clone into a memory fs
			r, err = git.Clone(memory.NewStorage(), nil, cloneOptions)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could clone repo", "error", err, "path", options.Path, "repo", options.Repo)
				return
			}
		}
	} else {
		if options.Path == "" {
			err = errors.New("git.path must be set if git.clone is false")
			return
		} else {
			// Open repo already on disk
			var absPath string
			absPath, err = filepath.Abs(options.Path)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not determine absolute path", "error", err, "path", options.Path)
				return
			}
			r, err = git.PlainOpen(absPath)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not open git repo", "error", err, "path", options.Path)
				return
			}
		}
	}

	return
}
