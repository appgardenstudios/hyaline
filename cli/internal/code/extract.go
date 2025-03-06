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

type Preset struct {
	Glob  string
	Files []string
}

var presets = map[string]Preset{
	"js": {
		Glob:  "./**/*.js",
		Files: []string{"./package.json", "./Makefile"},
	},
}

func ExtractCurrent(system *config.System, db *sql.DB) (err error) {
	// Process each code source
	for _, c := range system.Code {
		slog.Debug("ExtractCurrent extracting code", "system", system, "code", c.ID)
		// Insert Code
		codeId := system.ID + "-" + c.ID
		err = sqlite.InsertCurrentCode(sqlite.CurrentCode{
			ID:       codeId,
			SystemID: system.ID,
			Path:     c.Path,
		}, db)

		// Get our absolute path
		absPath, err := filepath.Abs(c.Path)
		if err != nil {
			slog.Debug("ExtractCurrent could not determine absolute code path", "error", err, "path", c.Path)
			return err
		}
		absPath += string(os.PathSeparator)
		slog.Debug("ExtractCurrent extracting code from path", "absPath", absPath)

		// Make sure we have a valid preset. If not, skip
		preset, ok := presets[c.Preset]
		if !ok {
			slog.Info("Code Preset Not Found. Skipping...", "system", system, "code", c.ID, "preset", c.Preset)
			continue
		}
		slog.Debug("ExtractCurrent extracting code using preset", "presetID", c.Preset, "preset", preset)

		// Our set of files (as a map so we don't get dupes)
		files := map[string]struct{}{}

		// Get files from our fully qualified glob path
		glob := filepath.Join(absPath, preset.Glob)
		matches, err := zglob.Glob(glob)
		if err != nil {
			slog.Debug("ExtractCurrent could not find code files with glob", "error", err)
			return err
		}
		for _, file := range matches {
			files[file] = struct{}{}
		}
		slog.Debug("ExtractCurrent found the following code file matches using glob", "glob", glob, "matches", matches)

		// Get files from our set of individual preset files
		for _, addtnlFile := range preset.Files {
			file := filepath.Join(absPath, addtnlFile)
			stat, err := os.Stat(file)
			if err == nil && !stat.IsDir() {
				files[file] = struct{}{}
			}
		}
		slog.Debug("ExtractCurrent will insert the following code files (glob plus additional)", "glob", glob, "additionalFiles", preset.Files, "files", files)

		// Insert files
		for file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				slog.Debug("ExtractCurrent could not read code file", "error", err, "file", file)
				return err
			}
			relativePath := strings.TrimPrefix(file, absPath)
			err = sqlite.InsertCurrentFile(sqlite.CurrentFile{
				ID:           relativePath,
				CodeID:       codeId,
				SystemID:     system.ID,
				RelativePath: relativePath,
				RawData:      string(contents),
			}, db)
			if err != nil {
				slog.Debug("ExtractCurrent could not insert file", "error", err)
				return err
			}
		}
	}

	return
}
