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

	// Load Config
	// TODO

	// Open Current DB
	// TODO

	// Open Change DB
	// TODO

	// Get system
	// TODO

	// Ensure output file does not exist
	// TODO

	// Get the set of files that need to be updated for each code change
	// TODO

	// Merge sets of files into a master list
	// TODO

	// Loop through documents that have been updated and annotate those on the list
	// TODO

	// Output the results
	// TODO

	return nil
}
