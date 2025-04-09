package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/code"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type MergeArgs struct {
	Inputs []string
	Output string
}

func Merge(args *MergeArgs) error {
	slog.Info("Merging data sets")
	slog.Debug("action.Merge Args", slog.Group("args",
		"inputs", args.Inputs,
		"output", args.Output,
	))

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.Merge could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		slog.Debug("action.Merge detected that output db already exists", "absPath", absPath)
		return errors.New("output file already exists")
	}
	outputDB, err := sql.Open("sqlite", absPath)
	if err != nil {
		slog.Debug("action.Merge could not open a new SQLite DB", "dataSourceName", absPath, "error", err)
		return err
	}
	defer outputDB.Close()
	err = sqlite.CreateSchema(outputDB)
	if err != nil {
		slog.Debug("action.Merge could not create the current schema", "error", err)
		return err
	}

	for _, input := range args.Inputs {
		slog.Info("Merging " + input)

		// Open input database
		inputAbsPath, err := filepath.Abs(input)
		if err != nil {
			slog.Debug("action.Merge could not get an absolute path for input", "input", input, "error", err)
			return err
		}
		inputDB, err := sql.Open("sqlite", inputAbsPath)
		if err != nil {
			slog.Debug("action.Merge could not open input SQLite DB", "dataSourceName", inputAbsPath, "error", err)
			return err
		}
		slog.Debug("action.Merge opened input database", "input", input, "path", inputAbsPath)
		defer inputDB.Close()

		// Get systems
		systems, err := sqlite.GetSystems(inputDB)
		if err != nil {
			slog.Debug("action.Merge could not get systems", "input", input, "error", err)
			return err
		}

		// Merge each system
		for _, system := range *systems {
			slog.Info(fmt.Sprintf("Merging system %s from %s", system.ID, input))

			err = code.Merge(system.ID, inputDB, outputDB)
			if err != nil {
				slog.Debug("action.Merge could not merge code", "input", input, "error", err)
				return err
			}

			err = docs.Merge(system.ID, inputDB, outputDB)
			if err != nil {
				slog.Debug("action.Merge could not merge code", "input", input, "error", err)
				return err
			}

			err = github.MergePullRequests(system.ID, inputDB, outputDB)
			if err != nil {
				slog.Debug("action.Merge could not merge pull requests", "input", input, "error", err)
				return err
			}

			err = github.MergeIssues(system.ID, inputDB, outputDB)
			if err != nil {
				slog.Debug("action.Merge could not merge issues", "input", input, "error", err)
				return err
			}
		}
	}

	return nil
}
