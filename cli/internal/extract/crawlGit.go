package extract

import (
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func crawlGit(cfg *config.ExtractCrawler, cb extractorCallback) error {
	// Initialize go-git repo (on disk or in mem)
	var r *git.Repository
	r, err := repo.GetRepo(cfg.Options)
	if err != nil {
		slog.Debug("extract.crawlGit could not get repo", "error", err)
		return err
	}

	// Determine branch to extract from
	branch := "main"
	if cfg.Options.Branch != "" {
		branch = cfg.Options.Branch
	}

	// Resolve branch ref
	ref, err := repo.ResolveAlias(r, branch)
	if err != nil {
		slog.Debug("extract.crawlGit could not resolve branch", "error", err, "branch", branch)
		return err
	}

	// Get files from branch using ref
	err = repo.GetFiles(*ref, r, func(f *object.File) error {
		if config.PathIsIncluded(f.Name, cfg.Include, cfg.Exclude) {
			// Get contents of file and call extractor callback
			var bytes []byte
			bytes, err = repo.GetBlobBytes(f.Blob)
			if err != nil {
				slog.Debug("extract.crawlGit could not get blob bytes", "error", err)
				return err
			}
			return cb(f.Name, bytes)
		}
		return nil
	})
	if err != nil {
		slog.Debug("extract.crawlGit could not get files from branch", "error", err, "branch", branch)
		return err
	}

	return nil
}
