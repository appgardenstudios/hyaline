package action

import (
	"errors"
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

	return nil
}
