package code

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

		// Get our absolute path
		absPath, err := filepath.Abs(c.FsOptions.Path)
		if err != nil {
			slog.Debug("code.ExtractCurrent could not determine absolute code path", "error", err, "path", c.FsOptions.Path)
			return err
		}
		absPath += string(os.PathSeparator)
		slog.Debug("code.ExtractCurrent extracting code from path", "absPath", absPath)

		// Our set of files (as a map so we don't get dupes)
		files := map[string]struct{}{}

		// Loop through our includes and get files
		for _, include := range c.Include {
			slog.Debug("code.ExtractCurrent extracting code using include", "include", include, "code", c.ID)

			// Construct our includePattern and get matches
			includePattern := filepath.Join(absPath, include)
			matches, err := zglob.Glob(includePattern)
			if err != nil {
				slog.Debug("code.ExtractCurrent could not find code files with glob", "glob", includePattern, "error", err)
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
						slog.Debug("code.ExtractCurrent could not match exclude", "excludePattern", excludePattern, "file", file, "error", err)
						return err
					}
					if match {
						slog.Debug("code.ExtractCurrent file excluded", "file", file, "excludePattern", excludePattern)
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
				slog.Debug("code.ExtractCurrent could not read code file", "error", err, "file", file)
				return err
			}
			relativePath := strings.TrimPrefix(file, absPath)
			err = sqlite.InsertFile(sqlite.File{
				ID:       relativePath,
				CodeID:   c.ID,
				SystemID: system.ID,
				RawData:  string(contents),
			}, db)
			if err != nil {
				slog.Debug("code.ExtractCurrent could not insert file", "error", err)
				return err
			}
		}
	}

	return
}
