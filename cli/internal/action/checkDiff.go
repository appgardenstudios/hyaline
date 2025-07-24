package action

import "log/slog"

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

	return nil
}
