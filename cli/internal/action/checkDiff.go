package action

import (
	"encoding/json"
	"errors"
	"hyaline/internal/check"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/llm"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
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

type CheckOutput struct {
	Recommendations []CheckRecommendation `json:"recommendations"`
	Head            string                `json:"head"`
	Base            string                `json:"base"`
}

type CheckRecommendation struct {
	Source         string         `json:"documentationSource"`
	Document       string         `json:"document"`
	Section        []string       `json:"section,omitempty"`
	Recommendation string         `json:"recommendation"`
	Reasons        []check.Reason `json:"reasons"`
	Changed        bool           `json:"changed"`
	Checked        bool           `json:"checked"`
	Outdated       bool           `json:"outdated"`
}

func sortCheckRecommendations(recommendations []CheckRecommendation) {
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Outdated != recommendations[j].Outdated {
			return !recommendations[i].Outdated
		}
		if recommendations[i].Source < recommendations[j].Source {
			return true
		}
		if recommendations[i].Source > recommendations[j].Source {
			return false
		}
		if recommendations[i].Document < recommendations[j].Document {
			return true
		}
		if recommendations[i].Document > recommendations[j].Document {
			return false
		}
		return strings.Join(recommendations[i].Section, "/") < strings.Join(recommendations[j].Section, "/")
	})
}

func sortCheckReasons(reasons []check.Reason) {
	sort.Slice(reasons, func(i, j int) bool {
		if reasons[i].Outdated != reasons[j].Outdated {
			return !reasons[i].Outdated
		}
		if reasons[i].Check.File < reasons[j].Check.File {
			return true
		}
		if reasons[i].Check.File > reasons[j].Check.File {
			return false
		}
		return reasons[i].Check.Type < reasons[j].Check.Type
	})
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
	cfg, err := config.Load(args.Config, true)
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

	// If check is disabled, skip
	if cfg.Check.Disabled {
		slog.Info("Check disabled. Skipping...")
		return nil
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckDiff could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckDiff detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Get Pull Request
	var pr *github.PullRequest
	if args.PullRequest != "" {
		if cfg.GitHub.Token == "" {
			return errors.New("github token required to retrieve pull-request information")
		}
		pr, err = github.GetPullRequest(args.PullRequest, cfg.GitHub.Token)
		if err != nil {
			slog.Debug("action.CheckDiff could not get pull request", "pull-request", args.PullRequest, "error", err)
			return err
		}
		slog.Info("Retrieved pull-request", "pull-request", args.PullRequest)
	}

	// Get Issue(s)
	issues := []*github.Issue{}
	if len(args.Issues) > 0 {
		if cfg.GitHub.Token == "" {
			return errors.New("github token required to retrieve issue information")
		}
		for _, issue := range args.Issues {
			issue, err := github.GetIssue(issue, cfg.GitHub.Token)
			if err != nil {
				slog.Debug("action.CheckDiff could not get issue", "issue", issue, "error", err)
				return err
			}
			issues = append(issues, issue)
		}
		slog.Info("Retrieved issues", "issues", strings.Join(args.Issues, ", "))
	}

	// Get Documents
	docDB, err := sqlite.InitInput(args.Documentation)
	if err != nil {
		slog.Debug("action.CheckDiff could not initialize documentation db", "documentation", args.Documentation, "error", err)
		return err
	}
	documents, err := docs.GetFilteredDocs(&cfg.Check.Documentation, docDB)
	if err != nil {
		slog.Debug("action.CheckDiff could not get filtered documents", "error", err)
		return err
	}
	slog.Info("Retrieved filtered documents", "documents", len(documents))

	// Open repo and resolve head and base references
	var absPath string
	absPath, err = filepath.Abs(args.Path)
	if err != nil {
		slog.Debug("action.CheckDiff could not determine absolute path", "error", err, "path", args.Path)
		return err
	}
	slog.Info("Opening repo on disk", "absPath", absPath)
	var r *git.Repository
	r, err = git.PlainOpen(absPath)
	if err != nil {
		slog.Debug("action.CheckDiff could not open git repo", "error", err, "path", args.Path)
		return err
	}

	// Resolve head and base references
	resolvedHead, err := repo.ResolveRef(r, args.Head, args.HeadRef)
	if err != nil {
		slog.Debug("action.CheckDiff could not resolve head reference", "error", err)
		return err
	}
	resolvedBase, err := repo.ResolveRef(r, args.Base, args.BaseRef)
	if err != nil {
		slog.Debug("action.CheckDiff could not resolve base reference", "error", err)
		return err
	}

	// Get Diff
	filteredFiles, changedFiles, err := code.GetFilteredDiff(r, *resolvedHead, *resolvedBase, &cfg.Check.Code)
	if err != nil {
		slog.Debug("action.CheckDiff could not get filtered diff", "error", err)
		return err
	}
	slog.Info("Retrieved filtered files from diff", "files", len(filteredFiles))

	// Get recommendations
	recommendations, _, err := getRecommendations(filteredFiles, documents, pr, issues, changedFiles, cfg.Check, &cfg.LLM)
	if err != nil {
		slog.Debug("action.CheckDiff could not get recommendations", "error", err)
		return err
	}

	output := CheckOutput{
		Recommendations: recommendations,
		Head:            (*resolvedHead).String(),
		Base:            (*resolvedBase).String(),
	}

	// Output the results
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		slog.Debug("action.CheckDiff could not marshal json", "error", err)
		return err
	}
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.CheckDiff could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.CheckDiff could not write output file", "error", err)
		return err
	}
	slog.Info("Output recommendations", "recommendations", len(recommendations), "output", outputAbsPath)

	return nil
}

func getRecommendations(filteredFiles []code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue, changedFiles map[string]struct{}, checkConfig *config.Check, llmConfig *config.LLM) ([]CheckRecommendation, check.FileCheckContextHashes, error) {
	// Check Diff
	results, fileCheckContextHashes, err := check.Diff(filteredFiles, documents, pr, issues, checkConfig, llmConfig, llm.CallLLM)
	if err != nil {
		slog.Debug("getRecommendations could not check diff", "error", err)
		return nil, nil, err
	}
	slog.Info("Got results", "results", len(results))

	// Format results
	updateSource := checkConfig.Options.DetectDocumentationUpdates.Source
	recommendations := []CheckRecommendation{}
	for _, result := range results {
		changed := false
		if updateSource == result.Source {
			_, changed = changedFiles[result.Document]
		}
		recommendations = append(recommendations, CheckRecommendation{
			Source:         result.Source,
			Document:       result.Document,
			Section:        result.Section,
			Recommendation: "Consider reviewing and updating this documentation",
			Reasons:        result.Reasons,
			Changed:        changed,
			Checked:        changed,
		})
	}

	sortCheckRecommendations(recommendations)

	return recommendations, fileCheckContextHashes, nil
}
