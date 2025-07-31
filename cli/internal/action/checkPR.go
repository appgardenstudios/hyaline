package action

import (
	"encoding/json"
	"errors"
	"hyaline/internal/check"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CheckPRArgs struct {
	Config        string
	Documentation string
	PullRequest   string
	Issues        []string
	Output        string
}

func CheckPR(args *CheckPRArgs) error {
	slog.Info("Checking PR",
		"config", args.Config,
		"documentation", args.Documentation,
		"pull-request", args.PullRequest,
		"issues", args.Issues,
		"output", args.Output)

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckPR could not load the config", "error", err)
		return err
	}

	// Ensure check options are set as they are required for this action to run
	if cfg.Check == nil {
		slog.Debug("action.CheckPR did not find check options")
		err = errors.New("the check pr command requires check options be set in the config")
		return err
	}

	// Ensure GitHub token is available
	if cfg.GitHub.Token == "" {
		return errors.New("github token required to retrieve pull-request information")
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckPR could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckPR detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Get Pull Request
	pr, err := github.GetPullRequest(args.PullRequest, cfg.GitHub.Token)
	if err != nil {
		slog.Debug("action.CheckPR could not get pull request", "pull-request", args.PullRequest, "error", err)
		return err
	}
	slog.Info("Retrieved pull-request", "pull-request", args.PullRequest)

	// Get Issue(s)
	issues := []*github.Issue{}
	if len(args.Issues) > 0 {
		for _, issue := range args.Issues {
			issue, err := github.GetIssue(issue, cfg.GitHub.Token)
			if err != nil {
				slog.Debug("action.CheckPR could not get issue", "issue", issue, "error", err)
				return err
			}
			issues = append(issues, issue)
		}
		slog.Info("Retrieved issues", "issues", strings.Join(args.Issues, ", "))
	}

	// Get Documents
	docDB, err := sqlite.InitInput(args.Documentation)
	if err != nil {
		slog.Debug("action.CheckPR could not initialize documentation db", "documentation", args.Documentation, "error", err)
		return err
	}
	documents, err := docs.GetFilteredDocs(&cfg.Check.Documentation, docDB)
	if err != nil {
		slog.Debug("action.CheckPR could not get filtered documents", "error", err)
		return err
	}
	slog.Info("Retrieved filtered documents", "documents", len(documents))

	// Get PR Files
	filteredFiles, changedFiles, err := code.GetFilteredPR(args.PullRequest, cfg.GitHub.Token, &cfg.Check.Code)
	if err != nil {
		slog.Debug("action.CheckPR could not get filtered PR files", "error", err)
		return err
	}
	slog.Info("Retrieved filtered files from PR", "files", len(filteredFiles))

	// Check Diff
	results, err := check.Diff(filteredFiles, documents, pr, issues, cfg.Check, &cfg.LLM)
	if err != nil {
		slog.Debug("action.CheckPR could not check diff", "error", err)
		return err
	}
	slog.Info("Got results", "results", len(results))

	// Format results (reuse existing CheckDiffRecommendation format for compatibility)
	updateSource := cfg.Check.Options.DetectDocumentationUpdates.Source
	recommendations := []CheckDiffRecommendation{}
	for _, result := range results {
		changed := false
		if updateSource == result.Source {
			_, changed = changedFiles[result.Document]
		}
		recommendations = append(recommendations, CheckDiffRecommendation{
			Source:         result.Source,
			Document:       result.Document,
			Section:        result.Section,
			Recommendation: "Consider reviewing and updating this documentation",
			Reasons:        result.Reasons,
			Changed:        changed,
		})
	}
	sort.Sort(CheckDiffRecommendationSort(recommendations))
	output := CheckDiffOutput{
		Recommendations: recommendations,
	}

	// Output the results
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		slog.Debug("action.CheckPR could not marshal json", "error", err)
		return err
	}
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.CheckPR could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.CheckPR could not write output file", "error", err)
		return err
	}
	slog.Info("Output recommendations", "recommendations", len(recommendations), "output", outputAbsPath)

	return nil
}