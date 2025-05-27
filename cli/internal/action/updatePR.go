package action

import "log/slog"

type UpdatePRArgs struct {
	Config          string
	PullRequest     string
	Comment         string
	Sha             string
	Recommendations string
	Output          string
}

func UpdatePR(args *UpdatePRArgs) error {
	slog.Info("Update PR")
	slog.Debug("action.UpdatePR Args", slog.Group("args",
		"config", args.Config,
		"pullRequest", args.PullRequest,
		"comment", args.Comment,
		"sha", args.Sha,
		"recommendations", args.Recommendations,
		"output", args.Output,
	))

	return nil
}
