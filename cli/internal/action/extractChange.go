package action

import (
	"database/sql"
	"errors"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ExtractChangeArgs struct {
	Config           string
	System           string
	Base             string
	Head             string
	BaseRef          string
	HeadRef          string
	CodeIDs          []string
	DocumentationIDs []string
	PullRequest      string
	Issues           []string
	Output           string
}

func ExtractChange(args *ExtractChangeArgs) error {
	slog.Info("Extracting changed code and docs")
	slog.Debug("action.ExtractChange Args", slog.Group("args",
		"config", args.Config,
		"system", args.System,
		"base", args.Base,
		"head", args.Head,
		"baseRef", args.BaseRef,
		"headRef", args.HeadRef,
		"codeIDs", args.CodeIDs,
		"documentationIDs", args.DocumentationIDs,
		"pullRequest", args.PullRequest,
		"issues", args.Issues,
		"output", args.Output,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.ExtractChange could not load the config", "error", err)
		return err
	}

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.ExtractChange could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		slog.Debug("action.ExtractChange detected that output db already exists", "absPath", absPath)
		return errors.New("output file already exists")
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.ExtractChange could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer db.Close()
	err = sqlite.CreateSchema(db)
	if err != nil {
		slog.Debug("action.ExtractChange could not create the current schema", "error", err)
		return err
	}

	// Get/Insert System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.ExtractChange could not locate the system", "system", args.System, "error", err)
		return err
	}
	err = sqlite.InsertSystem(sqlite.System{
		ID: system.ID,
	}, db)
	if err != nil {
		slog.Debug("action.ExtractChange could not insert the system", "error", err)
		return err
	}
	slog.Debug("action.ExtractChange system inserted")

	// Determine our set of code/documentation IDs to extract
	// (default to extracting everything if no code AND documentation IDs are passed in)
	codeIDs := append([]string{}, args.CodeIDs...)
	documentationIDs := append([]string{}, args.DocumentationIDs...)
	if len(codeIDs) == 0 && len(documentationIDs) == 0 {
		// Extract all code and documentation ids
		for _, c := range system.CodeSources {
			codeIDs = append(codeIDs, c.ID)
		}
		for _, d := range system.DocumentationSources {
			documentationIDs = append(documentationIDs, d.ID)
		}
	}

	// Extract/Insert Pull Request (if present)
	if args.PullRequest != "" {
		if cfg.GitHub.Token == "" {
			return errors.New("github token required to retrieve pull-request information")
		}
		err = github.InsertPullRequest(args.PullRequest, cfg.GitHub.Token, system.ID, db)
		if err != nil {
			slog.Debug("action.ExtractChange could not insert the system", "error", err)
			return err
		}
		slog.Debug("action.ExtractChange pull request inserted")
	}

	// Extract/Insert Issues
	if len(args.Issues) > 0 {
		if cfg.GitHub.Token == "" {
			return errors.New("github token required to retrieve issue information")
		}
		for _, issue := range args.Issues {
			err = github.InsertIssue(issue, cfg.GitHub.Token, system.ID, db)
			if err != nil {
				slog.Debug("action.ExtractChange could not insert the system", "error", err)
				return err
			}
		}
		slog.Debug("action.ExtractChange issues inserted")
	}

	// Extract/Insert Code
	err = code.ExtractChange(system, args.Head, args.HeadRef, args.Base, args.BaseRef, codeIDs, db)
	if err != nil {
		slog.Debug("action.ExtractChange could not extract code", "error", err)
		return err
	}
	slog.Debug("action.ExtractChange code inserted")

	// Extract/Insert Docs
	err = docs.ExtractChange(system, args.Head, args.HeadRef, args.Base, args.BaseRef, documentationIDs, db)
	if err != nil {
		slog.Debug("action.ExtractChange could not extract docs", "error", err)
		return err
	}
	slog.Debug("action.ExtractChange docs inserted")

	slog.Info("Extraction complete")
	return nil
}
