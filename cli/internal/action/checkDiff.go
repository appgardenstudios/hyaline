package action

import (
	"errors"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"log/slog"
)

type CheckDiffArgs struct {
	Config        string
	Documentation string
	Path          string
	Base          string
	BaseRef       string
	Head          string
	HeadRef       string
	PullRequest   string
	Issues        []string
	Output        string
}

func CheckDiff(args *CheckDiffArgs) error {
	slog.Info("Checking diff",
		"config", args.Config,
		"documentation", args.Documentation,
		"path", args.Path,
		"base", args.Base,
		"base-ref", args.BaseRef,
		"head", args.Head,
		"head-ref", args.HeadRef,
		"pull-request", args.PullRequest,
		"issues", args.Issues,
		"output", args.Output)

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckDiff could not load the config", "error", err)
		return err
	}

	// Ensure check options are set as they are required for this action to run
	if cfg.Check == nil {
		slog.Debug("action.CheckDiff did not find check options")
		err = errors.New("the check diff command requires check options be set in the config")
		return err
	}

	// // Get Pull Request
	// var pr *string
	// if args.PullRequest != "" {
	// 	if cfg.GitHub.Token == "" {
	// 		return errors.New("github token required to retrieve pull-request information")
	// 	}
	// 	pr, err = github.GetPullRequest(args.PullRequest, cfg.GitHub.Token)
	// 	if err != nil {
	// 		slog.Debug("action.CheckDiff could not get pull request", "pull-request", args.PullRequest, "error", err)
	// 		return err
	// 	}
	// 	slog.Debug("action.CheckDiff retrieved pull-request", "pull-request", *pr) // TODO clean up
	// }

	// // Get Issue(s)
	// issues := []*string{}
	// if len(args.Issues) > 0 {
	// 	if cfg.GitHub.Token == "" {
	// 		return errors.New("github token required to retrieve issue information")
	// 	}
	// 	for _, issue := range args.Issues {
	// 		body, err := github.GetIssue(issue, cfg.GitHub.Token)
	// 		if err != nil {
	// 			slog.Debug("action.CheckDiff could not get issue", "issue", issue, "error", err)
	// 			return err
	// 		}
	// 		issues = append(issues, body)
	// 	}
	// 	slog.Debug("action.CheckDiff retrieved issues", "issues", issues) // TODO clean up
	// }

	// // Get Documents
	// docDB, err := sqlite.InitInput(args.Documentation)
	// if err != nil {
	// 	slog.Debug("action.CheckDiff could not initialize documentation db", "documentation", args.Documentation, "error", err)
	// 	return err
	// }
	// documents, err := docs.GetFilteredDocs(&cfg.Check.Documentation, docDB)
	// if err != nil {
	// 	slog.Debug("action.CheckDiff could not get filtered documents", "error", err)
	// 	return err
	// }
	// slog.Debug("action.CheckDiff retrieved documents", "documents", documents) // TODO clean up

	// Get Diff
	files, err := code.GetFilteredDiff(args.Path, args.Head, args.HeadRef, args.Base, args.BaseRef, &cfg.Check.Code)
	if err != nil {
		slog.Debug("action.CheckDiff could not get filtered diff", "error", err)
		return err
	}
	slog.Debug("action.CheckDiff retrieved files from diff", "files", files) // TODO clean up

	return nil
}
