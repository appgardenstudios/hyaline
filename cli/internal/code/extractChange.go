package code

import (
	"context"
	"database/sql"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

		// Ensure repo is clean
		// TODO

		// Ensure we are on the head branch
		ref, err := repo.Head()
		if err != nil {
			slog.Debug("code.ExtractChange could not retrieve head", "error", err)
			return err
		}
		slog.Debug("code.ExtractChange repo head", "ref", ref.Name())
		// TODO compare refs

		// Get a list of files change between head and base
		headRef, err := repo.ResolveRevision(plumbing.Revision(head))
		if err != nil {
			slog.Debug("code.ExtractChange could not resolve head", "error", err)
			return err
		}
		headCommit, err := repo.CommitObject(*headRef)
		if err != nil {
			slog.Debug("code.ExtractChange could not get head commit", "error", err)
			return err
		}
		headTree, err := headCommit.Tree()
		if err != nil {
			slog.Debug("code.ExtractChange could not get head tree", "error", err)
			return err
		}
		baseRef, err := repo.ResolveRevision(plumbing.Revision(base))
		if err != nil {
			slog.Debug("code.ExtractChange could not resolve base", "error", err)
			return err
		}
		baseCommit, err := repo.CommitObject(*baseRef)
		if err != nil {
			slog.Debug("code.ExtractChange could not get base commit", "error", err)
			return err
		}
		baseTree, err := baseCommit.Tree()
		if err != nil {
			slog.Debug("code.ExtractChange could not get base tree", "error", err)
			return err
		}
		diffs, err := object.DiffTreeWithOptions(context.Background(), baseTree, headTree, object.DefaultDiffTreeOptions)
		if err != nil {
			slog.Debug("code.ExtractChange could not get diff", "error", err)
			return err
		}

		// TODO
		for _, diff := range diffs {
			fmt.Println(diff.String())
		}

	}

	return
}
