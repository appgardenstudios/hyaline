package code

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/mattn/go-zglob"
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
		glob, err := zglob.New(preset.Glob)
		if err != nil {
			slog.Debug("code.ExtractChange could not instantiate preset glob", "error", err, "glob", preset.Glob)
			return err
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
		// Note that we will eventually want to support renames via object.DiffTreeWithOptions
		diffs, err := object.DiffTree(baseTree, headTree)
		if err != nil {
			slog.Debug("code.ExtractChange could not get diff", "error", err)
			return err
		}

		// Load any files in the diff that match our preset
		for _, diff := range diffs {
			slog.Debug("code.ExtractChange processing diff", "diff", diff.String())
			action, err := diff.Action()
			if err != nil {
				slog.Debug("code.ExtractChange could not retrieve action for diff", "error", err, "diff", diff)
				return err
			}
			from, to, err := diff.Files()
			if err != nil {
				slog.Debug("code.ExtractChange could not retrieve files for diff", "error", err, "diff", diff)
				return err
			}
			switch action {
			case merkletrie.Insert:
				fallthrough
			case merkletrie.Modify:
				// TODO make the files match less fragile
				if glob.Match(to.Name) || slices.Contains(preset.Files, "./"+to.Name) {
					slog.Debug("code.ExtractChange inserting file", "file", to.Name, "action", action)
					bytes, err := getBlobBytes(to.Blob)
					if err != nil {
						slog.Debug("code.ExtractChange could not retrieve blob from diff", "error", err)
						return err
					}
					err = sqlite.InsertChangeFile(sqlite.ChangeFile{
						ID:           to.Name,
						CodeID:       codeId,
						SystemID:     system.ID,
						RelativePath: to.Name,
						Action:       action.String(),
						RawData:      string(bytes),
					}, db)
					if err != nil {
						slog.Debug("code.ExtractChange could not insert file", "error", err)
						return err
					}
				}
			case merkletrie.Delete:
				// TODO make the files match less fragile
				if glob.Match(from.Name) || slices.Contains(preset.Files, "./"+from.Name) {
					slog.Debug("code.ExtractChange inserting file", "file", from.Name, "action", action)
					err = sqlite.InsertChangeFile(sqlite.ChangeFile{
						ID:           from.Name,
						CodeID:       codeId,
						SystemID:     system.ID,
						RelativePath: from.Name,
						Action:       action.String(),
						RawData:      "",
					}, db)
					if err != nil {
						slog.Debug("code.ExtractChange could not insert file", "error", err)
						return err
					}
				}
			}
		}
	}

	return
}

func getBlobBytes(blob object.Blob) (bytes []byte, err error) {
	r, err := blob.Reader()
	if err != nil {
		return
	}
	bytes, err = io.ReadAll(r)
	return
}
