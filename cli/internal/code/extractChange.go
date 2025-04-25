package code

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

func ExtractChange(system *config.System, head string, base string, db *sql.DB) (err error) {
	// Process each code source
	for _, c := range system.Code {
		// Only extract changed code for git sources
		if c.Extractor != config.ExtractorGit {
			slog.Debug("code.ExtractChange skipping non-git code source", "system", system.ID, "code", c.ID)
			continue
		}
		slog.Debug("code.ExtractChange extracting code", "system", system.ID, "code", c.ID, "head", head, "base", base)

		// Get document path
		path := c.GitOptions.Path
		if path == "" {
			path = c.GitOptions.Repo
		}

		// Insert Code
		err = sqlite.InsertCode(sqlite.Code{
			ID:       c.ID,
			SystemID: system.ID,
			Path:     path,
		}, db)
		if err != nil {
			slog.Debug("code.ExtractChange could not insert code", "error", err, "code", c.ID)
			return err
		}

		// Initialize go-git repo (on disk or in mem)
		var r *git.Repository
		r, err = repo.GetRepo(c.GitOptions)
		if err != nil {
			slog.Debug("code.ExtractChange could not get repo", "error", err)
			return
		}

		// Get our diff
		diff, err := repo.GetDiff(r, head, base)
		if err != nil {
			slog.Debug("code.ExtractChange could not get diff", "error", err)
			return err
		}

		// Load any files in the diff that match our preset
		for _, change := range diff {
			slog.Debug("code.ExtractChange processing diff", "diff", change.String())
			action, err := change.Action()
			if err != nil {
				slog.Debug("code.ExtractChange could not retrieve action for diff", "error", err, "diff", change)
				return err
			}
			_, to, err := change.Files()
			if err != nil {
				slog.Debug("code.ExtractChange could not retrieve files for diff", "error", err, "diff", change)
				return err
			}
			switch action {
			case merkletrie.Insert:
				fallthrough
			case merkletrie.Modify:
				if config.PathIsIncluded(change.To.Name, c.Include, c.Exclude) {
					slog.Debug("code.ExtractChange inserting file", "file", change.To.Name, "action", action)
					bytes, err := repo.GetBlobBytes(to.Blob)
					if err != nil {
						slog.Debug("code.ExtractChange could not retrieve blob from diff", "error", err)
						return err
					}
					err = sqlite.InsertFile(sqlite.File{
						ID:         change.To.Name,
						CodeID:     c.ID,
						SystemID:   system.ID,
						Action:     sqlite.MapAction(action, change.From.Name, change.To.Name),
						OriginalID: change.From.Name,
						RawData:    string(bytes),
					}, db)
					if err != nil {
						slog.Debug("code.ExtractChange could not insert file", "error", err)
						return err
					}
				}
			case merkletrie.Delete:
				if config.PathIsIncluded(change.To.Name, c.Include, c.Exclude) {
					slog.Debug("code.ExtractChange inserting file", "file", change.From.Name, "action", action)
					err = sqlite.InsertFile(sqlite.File{
						ID:         change.From.Name,
						CodeID:     c.ID,
						SystemID:   system.ID,
						Action:     sqlite.ActionDelete,
						OriginalID: "",
						RawData:    "",
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
