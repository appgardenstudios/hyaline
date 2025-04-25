package action

import (
	"log/slog"

	_ "modernc.org/sqlite"
)

type CheckChangeArgs struct {
	Config  string
	Current string
	Change  string
	System  string
	Output  string
}

func CheckChange(args *CheckChangeArgs) error {
	slog.Info("Checking changed code and docs")
	slog.Debug("action.CheckChange Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"change", args.Change,
		"system", args.System,
		"output", args.Output,
	))

	// TODO

	return nil
}
