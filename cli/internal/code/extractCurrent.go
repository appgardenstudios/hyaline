package code

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mattn/go-zglob"
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
	absPath += string(os.PathSeparator)
	slog.Debug("code.ExtractCurrentFs extracting code from path", "absPath", absPath)

	// Our set of files (as a map so we don't get dupes)
	files := map[string]struct{}{}

	// Loop through our includes and get files
	for _, include := range c.Include {
		slog.Debug("code.ExtractCurrentFs extracting code using include", "include", include, "code", c.ID)

		// Construct our includePattern and get matches
		includePattern := filepath.Join(absPath, include)
		matches, err := zglob.Glob(includePattern)
		if err != nil {
			slog.Debug("code.ExtractCurrentFs could not find code files with glob", "glob", includePattern, "error", err)
			return err
		}

		// Loop through files and add those that aren't in our excludes
		for _, file := range matches {
			// See if we have a match for at least one of our excludes
			match := false
			for _, exclude := range c.Exclude {
				excludePattern := filepath.Join(absPath, exclude)
				match, err = zglob.Match(excludePattern, file)
				if err != nil {
					slog.Debug("code.ExtractCurrentFs could not match exclude", "excludePattern", excludePattern, "file", file, "error", err)
					return err
				}
				if match {
					slog.Debug("code.ExtractCurrentFs file excluded", "file", file, "excludePattern", excludePattern)
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
		contents, err := os.ReadFile(file)
		if err != nil {
			slog.Debug("code.ExtractCurrentFs could not read code file", "error", err, "file", file)
			return err
		}
		relativePath := strings.TrimPrefix(file, absPath)
		err = sqlite.InsertFile(sqlite.File{
			ID:       relativePath,
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
	var r *git.Repository
	if c.GitOptions.Clone {
		// Ensure remote repo is passed in
		if c.GitOptions.Repo == "" {
			err = errors.New("git.repo is required to be set if git.clone is true")
			return
		}
		if c.GitOptions.Path != "" {
			// Clone to disk
			var absPath string
			absPath, err = filepath.Abs(c.GitOptions.Path)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not determine absolute path", "error", err, "path", c.GitOptions.Path)
				return
			}
			r, err = git.PlainClone(absPath, false, &git.CloneOptions{
				URL: c.GitOptions.Repo,
			})
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could clone repo", "error", err, "path", c.GitOptions.Path, "repo", c.GitOptions.Repo)
				return
			}
		} else {
			// Clone into a memory fs
			r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
				URL: c.GitOptions.Repo,
			})
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could clone repo", "error", err, "path", c.GitOptions.Path, "repo", c.GitOptions.Repo)
				return
			}
		}
	} else {
		if c.GitOptions.Path == "" {
			err = errors.New("git.path must be set if git.clone is false")
			return
		} else {
			// Open repo already on disk
			var absPath string
			absPath, err = filepath.Abs(c.GitOptions.Path)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not determine absolute path", "error", err, "path", c.GitOptions.Path)
				return
			}
			r, err = git.PlainOpen(absPath)
			if err != nil {
				slog.Debug("code.ExtractCurrentGit could not open git repo", "error", err, "path", c.GitOptions.Path)
				return
			}
		}
	}

	// Fetch (if needed)
	if c.GitOptions.Fetch {
		remoteName := "origin"
		if c.GitOptions.RemoteName != "" {
			remoteName = c.GitOptions.RemoteName
		}
		r.Fetch(&git.FetchOptions{
			RemoteName: remoteName,
		})
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

	// Validate includes/excludes (move to config validation)
	for _, include := range c.Include {
		if !doublestar.ValidatePattern(include) {
			slog.Debug("code.ExtractCurrentGit could not validate include", "include", include)
			err = errors.New("include pattern invalid: " + include)
			return
		}
	}

	// Get and save code files
	tree.Files().ForEach(func(f *object.File) error {
		for _, include := range c.Include {
			if doublestar.MatchUnvalidated(include, f.Name) {
				for _, exclude := range c.Exclude {
					if doublestar.MatchUnvalidated(exclude, f.Name) {
						continue
					}
				}
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
