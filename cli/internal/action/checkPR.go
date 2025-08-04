package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"hyaline/internal/check"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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

	// Get recommendations
	recommendations, err := getRecommendations(filteredFiles, documents, pr, issues, changedFiles, cfg.Check, &cfg.LLM)
	if err != nil {
		slog.Debug("action.CheckPR could not get recommendations", "error", err)
		return err
	}

	// Update PR comment using SHA from PR head
	err = updatePRComment(pr, pr.HeadSHA, recommendations, args.PullRequest, cfg.GitHub.Token)
	if err != nil {
		slog.Debug("action.CheckPR could not update PR comment", "error", err)
		return err
	}
	slog.Info("Updated PR comment", "sha", pr.HeadSHA)

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

func getRecommendations(filteredFiles []code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue, changedFiles map[string]struct{}, checkConfig *config.Check, llmConfig *config.LLM) ([]CheckDiffRecommendation, error) {
	// Check Diff
	results, err := check.Diff(filteredFiles, documents, pr, issues, checkConfig, llmConfig)
	if err != nil {
		slog.Debug("getRecommendations could not check diff", "error", err)
		return nil, err
	}
	slog.Info("Got results", "results", len(results))

	// Format results (reuse existing CheckDiffRecommendation format for compatibility)
	updateSource := checkConfig.Options.DetectDocumentationUpdates.Source
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

	return recommendations, nil
}

func updatePRComment(pr *github.PullRequest, sha string, recommendations []CheckDiffRecommendation, prRef string, token string) error {
	// Find existing Hyaline comment
	commentID, err := findHyalineComment(prRef, token)
	if err != nil {
		slog.Debug("updatePRComment could not find Hyaline comment", "error", err)
		return err
	}

	// Convert CheckDiffRecommendations to format expected by comment functions
	updateRecs := []UpdatePRCommentRecommendation{}
	for _, rec := range recommendations {
		updateRecs = append(updateRecs, UpdatePRCommentRecommendation{
			Checked:  rec.Changed,
			System:   "checkPR", // Use a consistent system identifier
			Source:   rec.Source,
			Document: rec.Document,
			Section:  rec.Section,
			Reasons:  rec.Reasons,
		})
	}

	var comment *UpdatePRComment
	if commentID == 0 {
		// No existing comment, create a new one
		comment, err = checkPRAddComment(sha, updateRecs, prRef, token)
	} else {
		// Update existing comment
		commentRef := fmt.Sprintf("%s/%d", prRef, commentID)
		comment, err = checkPRUpdateComment(sha, updateRecs, prRef, commentRef, token)
	}
	
	if err != nil {
		slog.Debug("updatePRComment could not update or add comment", "commentID", commentID, "error", err)
		return err
	}

	slog.Info("Updated PR comment", "commentID", commentID, "recommendations", len(comment.Recommendations))
	return nil
}

func checkPRUpdateComment(sha string, newRecs []UpdatePRCommentRecommendation, pr string, comment string, token string) (*UpdatePRComment, error) {
	// Get comment
	existingComment, err := github.GetComment(comment, token)
	if err != nil {
		slog.Debug("checkPRUpdateComment could not get existing comment", "error", err)
		return nil, err
	}

	// Get existing recs from the comment
	existingRecs, err := parsePRComment(existingComment)
	if err != nil {
		slog.Debug("checkPRUpdateComment could not extract data from existing comment", "error", err)
		return nil, err
	}

	// Merge new recommendations with the ones from the comment
	mergedRecs := mergeRecsForCheckPR(newRecs, existingRecs)

	// Format New Raw Data
	rawData, err := formatRawData(&mergedRecs)
	if err != nil {
		slog.Debug("checkPRUpdateComment could not format raw data", "error", err)
		return nil, err
	}

	// Create comment
	updatedComment := UpdatePRComment{
		Sha:             sha,
		Recommendations: mergedRecs,
		RawData:         rawData,
	}

	// Format comment
	formattedComment := formatPRComment(&updatedComment)

	// Update comment
	err = github.UpdateComment(comment, formattedComment, token)
	if err != nil {
		slog.Debug("checkPRUpdateComment could not update comment", "pr", pr, "error", err)
		return nil, err
	}

	return &updatedComment, nil
}

func checkPRAddComment(sha string, recommendations []UpdatePRCommentRecommendation, pr string, token string) (*UpdatePRComment, error) {
	// Format raw data
	rawData, err := formatRawData(&recommendations)
	if err != nil {
		slog.Debug("checkPRAddComment could not format raw data", "error", err)
		return nil, err
	}

	// Create comment
	comment := UpdatePRComment{
		Sha:             sha,
		Recommendations: recommendations,
		RawData:         rawData,
	}

	// Format comment
	formattedComment := formatPRComment(&comment)

	// Add comment to PR
	err = github.AddComment(pr, formattedComment, token)
	if err != nil {
		slog.Debug("checkPRAddComment could not add comment", "pr", pr, "error", err)
		return nil, err
	}

	return &comment, nil
}

func mergeRecsForCheckPR(newRecs []UpdatePRCommentRecommendation, existingRecs []UpdatePRCommentRecommendation) (mergedRecs []UpdatePRCommentRecommendation) {
	// Initialize to empty slice to avoid returning nil
	mergedRecs = []UpdatePRCommentRecommendation{}
	// Copy over existing recs as is
	mergedRecs = append(mergedRecs, existingRecs...)

	// Add new recs
	for _, newRec := range newRecs {
		// See if this rec already exists
		match := false
		for index, existingRec := range existingRecs {
			if newRecMatchesExistingForCheckPR(&newRec, &existingRec) {
				match = true
				// Do not overwrite existing rec's checked status so we preserve any "manual" checks
				if !existingRec.Checked {
					mergedRecs[index].Checked = newRec.Checked
				}
				// Always merge reasons
				mergedRecs[index].Reasons = mergeReasonsForCheckPR(&newRec.Reasons, &existingRec.Reasons)
				break
			}
		}
		// If it does not, add it
		if !match {
			mergedRecs = append(mergedRecs, newRec)
		}
	}

	// Sort
	sort.Sort(UpdatePRCommentRecommendationSort(mergedRecs))

	return
}

func mergeReasonsForCheckPR(newReasons *[]string, existingReasons *[]string) (mergedReasons []string) {
	mergedReasons = append(mergedReasons, *existingReasons...)

	for _, newReason := range *newReasons {
		if !slices.Contains(mergedReasons, newReason) {
			mergedReasons = append(mergedReasons, newReason)
		}
	}

	return
}

func newRecMatchesExistingForCheckPR(newRec *UpdatePRCommentRecommendation, existingRec *UpdatePRCommentRecommendation) bool {
	if newRec.System != existingRec.System {
		return false
	}
	if newRec.Source != existingRec.Source {
		return false
	}
	if newRec.Document != existingRec.Document {
		return false
	}

	return reflect.DeepEqual(newRec.Section, existingRec.Section)
}

// findHyalineComment finds an existing Hyaline comment in a PR
// Returns the comment ID if found, or 0 if not found
func findHyalineComment(ref string, token string) (commentID int64, err error) {
	// Get all comments for the PR
	comments, err := github.ListComments(ref, token)
	if err != nil {
		return 0, err
	}

	// Search for a comment that starts with the Hyaline header (with zero-width spaces)
	hyalineHeader := "# H\u200By\u200Ba\u200Bl\u200Bi\u200Bn\u200Be"
	
	for _, comment := range comments {
		if strings.HasPrefix(comment.Body, hyalineHeader) {
			return comment.ID, nil
		}
	}

	// No Hyaline comment found
	return 0, nil
}