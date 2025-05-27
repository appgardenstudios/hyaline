package action

import (
	"encoding/json"
	"errors"
	"hyaline/internal/config"
	"log/slog"
	"os"
	"path/filepath"
)

type UpdatePRArgs struct {
	Config          string
	PullRequest     string
	Comment         string
	Sha             string
	Recommendations string
	Output          string
}

type comment struct {
	Sha             string
	Recommendations []string // TODO
	RawData         string   // TODO
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

	// Load Config
	_, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.UpdatePR could not load the config", "error", err)
		return err
	}

	// Load recommendations
	recsAbsPath, err := filepath.Abs(args.Recommendations)
	if err != nil {
		slog.Debug("action.UpdatePR could not get an absolute path for recommendations", "recommendations", args.Recommendations, "error", err)
		return err
	}
	recsData, err := os.ReadFile(recsAbsPath)
	if err != nil {
		slog.Debug("action.UpdatePR could not read recommendations", "recommendations", args.Recommendations, "error", err)
		return err
	}
	var recommendations CheckChangeOutput
	err = json.Unmarshal(recsData, &recommendations)
	if err != nil {
		slog.Debug("action.UpdatePR could not load recommendations", "recommendations", args.Recommendations, "error", err)
		return err
	}

	// Ensure output location does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.UpdatePR could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.UpdatePR detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Handle updating or adding the comment
	if args.Comment == "" {
		err = updatePRFromComment()
	} else {
		err = updatePRAddComment(recommendations, args.PullRequest)
	}
	if err != nil {
		slog.Debug("action.UpdatePR could not update or add comment", "comment", args.Comment, "error", err)
		return err
	}

	// Output recommendations and metadata TODO update this comment with the name of the struct we are outputting
	// TODO

	return nil
}

func updatePRFromComment() error {
	return nil
}

func updatePRAddComment(recommendations CheckChangeOutput, pr string) error {
	// Create comment
	// TODO

	// Format comment
	// TODO

	// Add comment
	// TODO

	return nil
}
