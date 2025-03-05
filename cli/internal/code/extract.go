package code

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
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

func ExtractCurrent(system string, cfg *config.Config, db *sql.DB) (err error) {
	// Find our target system (error if not found)
	var targetSystem *config.System
	for _, s := range cfg.Systems {
		if s.ID == system {
			targetSystem = &s
		}
	}
	if targetSystem == nil {
		// TODO better error message here
		return errors.New("system not found")
	}

	// Process each code source
	for _, c := range targetSystem.Code {
		// Insert Code
		codeId := targetSystem.ID + "-" + c.ID
		err = sqlite.InsertCurrentCode(sqlite.CurrentCode{
			ID:       codeId,
			SystemID: targetSystem.ID,
			Path:     c.Path,
		}, db)

		// Get our absolute path
		absPath, err := filepath.Abs(c.Path)
		if err != nil {
			return err
		}
		absPath += string(os.PathSeparator)

		// Make sure we have a valid preset. If not, skip
		preset, ok := presets[c.Preset]
		if !ok {
			// TODO better context in error message
			fmt.Println("Preset not found")
			continue
		}

		// Our set of files (as a map so we don't get dupes)
		files := map[string]struct{}{}

		// Get files from our fully qualified glob path
		glob := filepath.Join(absPath, preset.Glob)
		matches, err := zglob.Glob(glob)
		if err != nil {
			return err
		}
		for _, file := range matches {
			files[file] = struct{}{}
		}

		// Get files from our set of individual preset files
		for _, addtnlFile := range preset.Files {
			file := filepath.Join(absPath, addtnlFile)
			stat, err := os.Stat(file)
			if err == nil && !stat.IsDir() {
				files[file] = struct{}{}
			}
		}

		// Insert files
		for file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(file, absPath)
			err = sqlite.InsertCurrentFile(sqlite.CurrentFile{
				ID:           relativePath,
				CodeID:       codeId,
				SystemID:     targetSystem.ID,
				RelativePath: relativePath,
				RawData:      string(contents),
			}, db)
			if err != nil {
				return err
			}
		}
	}

	return
}
