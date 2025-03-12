package code

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func ExtractChange(system *config.System, head string, base string, db *sql.DB) (err error) {
	// Process each code source
	for _, c := range system.Code {
		slog.Debug("code.ExtractChange extracting code", "system", system, "code", c.ID)

		// Insert Code
		codeId := system.ID + "-" + c.ID
		err = sqlite.InsertChangeCode(sqlite.ChangeCode{
			ID:       codeId,
			SystemID: system.ID,
			Path:     c.Path,
		}, db)

		// Get our absolute path
		absPath, err := filepath.Abs(c.Path)
		if err != nil {
			slog.Debug("code.ExtractChange could not determine absolute code path", "error", err, "path", c.Path)
			return err
		}
		absPath += string(os.PathSeparator)
		slog.Debug("code.ExtractChange extracting code from path", "absPath", absPath)

		// Make sure we have a valid preset. If not, skip
		preset, ok := presets[c.Preset]
		if !ok {
			slog.Info("Code Preset Not Found. Skipping...", "system", system, "code", c.ID, "preset", c.Preset)
			continue
		}
		slog.Debug("code.ExtractChange extracting code using preset", "presetID", c.Preset, "preset", preset)

		// Open our git repo
		repo, err := git.PlainOpen(absPath)
		if err != nil {
			slog.Debug("code.ExtractChange could not open git repo", "error", err, "path", c.Path)
			return err
		}

		// Ensure we are on the head branch
		ref, err := repo.Head()
		if err != nil {
			slog.Debug("code.ExtractChange could not retrieve head", "error", err)
			return err
		}
		slog.Debug("code.ExtractChange repo head", "ref", ref.Name())
		// TODO compare refs

		// Get a list of files change between head and base
		// TODO get the set of

	}

	return
}
