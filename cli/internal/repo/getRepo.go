package repo

import (
	"errors"
	"hyaline/internal/config"
	"log/slog"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

func GetRepo(options config.ExtractorOptions) (r *git.Repository, err error) {
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

		// Add http auth if set
		if options.Auth.Type == config.ExtractorAuthHTTP {
			username := "git"
			if options.Auth.Options.Username != "" {
				username = options.Auth.Options.Username
			}
			cloneOptions.Auth = &http.BasicAuth{
				Username: username,
				Password: options.Auth.Options.Password,
			}
		}

		// Add ssh auth if set
		if options.Auth.Type == config.ExtractorAuthSSH {
			user := "git"
			if options.Auth.Options.User != "" {
				user = options.Auth.Options.User
			}
			var keys *ssh.PublicKeys
			keys, err = ssh.NewPublicKeys(user, []byte(options.Auth.Options.PEM), options.Auth.Options.Password)
			if err != nil {
				slog.Debug("repo.GetRepo could not parse git ssh PEM", "error", err)
				return
			}
			cloneOptions.Auth = keys
		}

		if options.Path != "" {
			// Clone to disk
			var absPath string
			absPath, err = filepath.Abs(options.Path)
			if err != nil {
				slog.Debug("repo.GetRepo could not determine absolute path", "error", err, "path", options.Path)
				return
			}
			slog.Info("Cloning to disk", "absPath", absPath)
			r, err = git.PlainClone(absPath, false, cloneOptions)
			if err != nil {
				slog.Debug("repo.GetRepo could not clone repo to disk", "error", err, "path", options.Path, "repo", options.Repo)
				return
			}
		} else {
			// Clone into a memory fs
			slog.Info("Cloning to memory fs")
			r, err = git.Clone(memory.NewStorage(), nil, cloneOptions)
			if err != nil {
				slog.Debug("repo.GetRepo could not clone repo to memory", "error", err, "path", options.Path, "repo", options.Repo)
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
				slog.Debug("repo.GetRepo could not determine absolute path", "error", err, "path", options.Path)
				return
			}
			slog.Info("Opening repo on disk", "absPath", absPath)
			r, err = git.PlainOpen(absPath)
			if err != nil {
				slog.Debug("repo.GetRepo could not open git repo", "error", err, "path", options.Path)
				return
			}
		}
	}

	return
}
