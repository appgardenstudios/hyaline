package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
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

type CheckPRComment struct {
	Sha             string                         `json:"sha"`
	Recommendations []CheckPRCommentRecommendation `json:"recommendations"`
	RawData         string                         `json:"rawData"`
}

type CheckPRCommentRecommendation struct {
	Checked  bool     `json:"checked"`
	Source   string   `json:"source"`
	Document string   `json:"document"`
	Section  []string `json:"section"`
	Reasons  []string `json:"reasons"`
}

type CheckPRCommentRecommendationSort []CheckPRCommentRecommendation

func (c CheckPRCommentRecommendationSort) Len() int {
	return len(c)
}
func (c CheckPRCommentRecommendationSort) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CheckPRCommentRecommendationSort) Less(i, j int) bool {
	if c[i].Source < c[j].Source {
		return true
	}
	if c[i].Source > c[j].Source {
		return false
	}
	if c[i].Document < c[j].Document {
		return true
	}
	if c[i].Document > c[j].Document {
		return false
	}
	return strings.Join(c[i].Section, "/") < strings.Join(c[j].Section, "/")
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
	filteredFiles, changedFiles, err := code.GetFilteredPRFiles(args.PullRequest, cfg.GitHub.Token, &cfg.Check.Code)
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

	slog.Info("Upserting PR comment", "sha", pr.HeadSHA, "recommendations", len(recommendations))
	comment, err := upsertPRComment(pr, pr.HeadSHA, recommendations, args.PullRequest, cfg.GitHub.Token)
	if err != nil {
		slog.Debug("action.CheckPR could not upsert PR comment", "error", err)
		return err
	}

	// Output the comment to a file
	jsonData, err := json.MarshalIndent(comment, "", "  ")
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
	slog.Info("Output recommendations", "recommendations", len(comment.Recommendations), "output", outputAbsPath)

	return nil
}

func getRecommendations(filteredFiles []code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue, changedFiles map[string]struct{}, checkConfig *config.Check, llmConfig *config.LLM) ([]CheckPRCommentRecommendation, error) {
	// Check Diff
	results, err := check.Diff(filteredFiles, documents, pr, issues, checkConfig, llmConfig)
	if err != nil {
		slog.Debug("getRecommendations could not check diff", "error", err)
		return nil, err
	}
	slog.Info("Got results", "results", len(results))

	// Format results
	updateSource := checkConfig.Options.DetectDocumentationUpdates.Source
	recommendations := []CheckPRCommentRecommendation{}
	for _, result := range results {
		changed := false
		if updateSource == result.Source {
			_, changed = changedFiles[result.Document]
		}
		recommendations = append(recommendations, CheckPRCommentRecommendation{
			Checked:  changed,
			Source:   result.Source,
			Document: result.Document,
			Section:  result.Section,
			Reasons:  result.Reasons,
		})
	}

	return recommendations, nil
}

func upsertPRComment(pr *github.PullRequest, sha string, recommendations []CheckPRCommentRecommendation, prRef string, token string) (*CheckPRComment, error) {
	// Find existing Hyaline comment
	existingComment, err := findHyalineComment(prRef, token)
	if err != nil {
		slog.Debug("upsertPRComment could not search for Hyaline comment", "error", err)
		return nil, err
	}

	var comment *CheckPRComment
	if existingComment == nil {
		slog.Info("Adding new PR comment")
		comment, err = addPRComment(sha, recommendations, prRef, token)

		if err != nil {
			slog.Debug("upsertPRComment could not add a comment", "error", err)
			return nil, err
		}
	} else {
		slog.Info("Updating existing comment", "commentID", existingComment.ID)
		comment, err = updatePRComment(sha, recommendations, prRef, existingComment, token)

		if err != nil {
			slog.Debug("upsertPRComment could not update a comment", "commentID", existingComment.ID, "error", err)
			return nil, err
		}
	}

	return comment, nil
}

func updatePRComment(sha string, newRecs []CheckPRCommentRecommendation, pr string, existingComment *github.Comment, token string) (*CheckPRComment, error) {
	// Get existing recs from the comment
	existingRecs, err := parseCheckPRComment(existingComment.Body)
	if err != nil {
		slog.Debug("updatePRComment could not extract data from existing comment", "error", err)
		return nil, err
	}

	// Merge new recommendations with the ones from the comment
	mergedRecs := mergeCheckPRRecs(newRecs, existingRecs)

	// Format New Raw Data
	rawData, err := formatCheckPRRawData(&mergedRecs)
	if err != nil {
		slog.Debug("updatePRComment could not format raw data", "error", err)
		return nil, err
	}

	// Create comment
	updatedComment := CheckPRComment{
		Sha:             sha,
		Recommendations: mergedRecs,
		RawData:         rawData,
	}

	// Format comment
	formattedComment := formatCheckPRComment(&updatedComment)

	prParts := strings.Split(pr, "/")
	commentRef := fmt.Sprintf("%s/%s/%d", prParts[0], prParts[1], existingComment.ID)
	err = github.UpdateComment(commentRef, formattedComment, token)
	if err != nil {
		slog.Debug("updatePRComment could not update comment", "pr", pr, "error", err)
		return nil, err
	}

	return &updatedComment, nil
}

func addPRComment(sha string, recommendations []CheckPRCommentRecommendation, pr string, token string) (*CheckPRComment, error) {
	// Format raw data
	rawData, err := formatCheckPRRawData(&recommendations)
	if err != nil {
		slog.Debug("addPRComment could not format raw data", "error", err)
		return nil, err
	}

	// Create comment
	comment := CheckPRComment{
		Sha:             sha,
		Recommendations: recommendations,
		RawData:         rawData,
	}

	// Format comment
	formattedComment := formatCheckPRComment(&comment)

	// Add comment to PR
	err = github.AddComment(pr, formattedComment, token)
	if err != nil {
		slog.Debug("addPRComment could not add comment", "pr", pr, "error", err)
		return nil, err
	}

	return &comment, nil
}

func mergeCheckPRRecs(newRecs []CheckPRCommentRecommendation, existingRecs []CheckPRCommentRecommendation) (mergedRecs []CheckPRCommentRecommendation) {
	// Initialize to empty slice to avoid returning nil
	mergedRecs = []CheckPRCommentRecommendation{}
	// Copy over existing recs as is
	mergedRecs = append(mergedRecs, existingRecs...)

	// Add new recs
	for _, newRec := range newRecs {
		// See if this rec already exists
		match := false
		for index, existingRec := range existingRecs {
			if newCheckPRRecMatchesExisting(&newRec, &existingRec) {
				match = true
				// Do not overwrite existing rec's checked status so we preserve any "manual" checks
				if !existingRec.Checked {
					mergedRecs[index].Checked = newRec.Checked
				}
				// Always merge reasons
				mergedRecs[index].Reasons = mergeCheckPRReasons(&newRec.Reasons, &existingRec.Reasons)
				break
			}
		}
		// If it does not, add it
		if !match {
			mergedRecs = append(mergedRecs, newRec)
		}
	}

	// Sort
	sort.Sort(CheckPRCommentRecommendationSort(mergedRecs))

	return
}

func mergeCheckPRReasons(newReasons *[]string, existingReasons *[]string) (mergedReasons []string) {
	mergedReasons = append(mergedReasons, *existingReasons...)

	for _, newReason := range *newReasons {
		if !slices.Contains(mergedReasons, newReason) {
			mergedReasons = append(mergedReasons, newReason)
		}
	}

	return
}

func newCheckPRRecMatchesExisting(newRec *CheckPRCommentRecommendation, existingRec *CheckPRCommentRecommendation) bool {
	if newRec.Source != existingRec.Source {
		return false
	}
	if newRec.Document != existingRec.Document {
		return false
	}

	return reflect.DeepEqual(newRec.Section, existingRec.Section)
}

func findHyalineComment(ref string, token string) (*github.Comment, error) {
	// Get all comments for the PR
	comments, err := github.ListComments(ref, token)
	if err != nil {
		return nil, err
	}

	// Search for a comment that starts with the Hyaline header (with zero-width spaces)
	for _, comment := range comments {
		if strings.HasPrefix(comment.Body, CHECKPR_HYALINE_HEADER) {
			return &comment, nil
		}
	}

	// No Hyaline comment found
	return nil, nil
}

func formatCheckPRRawData(recommendations *[]CheckPRCommentRecommendation) (string, error) {
	data, err := json.Marshal(recommendations)
	if err != nil {
		slog.Debug("formatCheckPRRawData could not marshal json", "error", err)
		return "", err
	}

	return fmt.Sprintf("<![CDATA[ %s ]]>", string(data)), nil
}

const CHECKPR_HYALINE_HEADER = "# H\u200By\u200Ba\u200Bl\u200Bi\u200Bn\u200Be"
const CHECKPR_CDATA_START = "<![CDATA[ "
const CHECKPR_CDATA_END = " ]]>"
const CHECKPR_RECOMMENDATIONS_START = "### Recommendations"

func parseCheckPRComment(comment string) (recs []CheckPRCommentRecommendation, err error) {
	// Get CData
	dataStart := strings.Index(comment, CHECKPR_CDATA_START)
	if dataStart == -1 {
		err = errors.New("could not find start of CDATA block")
		return
	}
	dataEnd := strings.LastIndex(comment, CHECKPR_CDATA_END)
	if dataEnd == -1 {
		err = errors.New("could not find end of CDATA block")
		return
	}
	data := comment[dataStart+10 : dataEnd]
	err = json.Unmarshal([]byte(data), &recs)
	if err != nil {
		return
	}

	// Get any manual checks and update the corresponding rec
	// Note that this relies on the CData order being the same as the recs list order
	recommendationsStart := strings.Index(comment, CHECKPR_RECOMMENDATIONS_START)
	if recommendationsStart == -1 {
		err = errors.New("could not find start of recommendations")
		return
	}
	checks := comment[recommendationsStart:dataStart]
	lines := strings.Split(checks, "\n")
	currentRec := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "- [") && currentRec < len(recs) {
			// Always pull checked from the list and NOT CData
			if strings.HasPrefix(line, "- [ ]") {
				recs[currentRec].Checked = false
			}
			if strings.HasPrefix(line, "- [x]") {
				recs[currentRec].Checked = true
			}
			currentRec++
		}
	}

	return
}

func formatCheckPRComment(comment *CheckPRComment) string {
	var md strings.Builder

	// Note: The comment MUST start with the Hyaline header
	md.WriteString(CHECKPR_HYALINE_HEADER + " PR Check\n")
	md.WriteString(fmt.Sprintf("**ref**: %s\n", html.EscapeString(comment.Sha)))
	md.WriteString("\n")

	// Note: This starting line always needs to be present because we use it as a sentinel for getting the check marks
	md.WriteString(fmt.Sprintf("%s\n", CHECKPR_RECOMMENDATIONS_START))
	if len(comment.Recommendations) > 0 {
		md.WriteString("Review the following recommendations and update the corresponding documentation as needed:\n")
		for _, rec := range comment.Recommendations {
			checked := " "
			if rec.Checked {
				checked = "x"
			}
			sections := ""
			if len(rec.Section) > 0 {
				sections = fmt.Sprintf(" > %s", strings.Join(rec.Section, " > "))
			}
			var cleanReasons []string
			for _, reason := range rec.Reasons {
				cleanReasons = append(cleanReasons, html.EscapeString(reason))
			}
			reasons := strings.Join(cleanReasons, "</li><li>")
			md.WriteString(fmt.Sprintf("- [%s] **%s**%s in `%s`", checked, html.EscapeString(rec.Document), html.EscapeString(sections), html.EscapeString(rec.Source)))
			md.WriteString(fmt.Sprintf("<details><summary>Reasons</summary><ul><li>%s</li></ul></details>", reasons))
			md.WriteString("\n")
		}
		md.WriteString("\nNote: Hyaline will automatically detect documentation updated in this PR and mark corresponding recommendations as reviewed.\n")
	} else {
		md.WriteString("Hyaline did not find any documentation related to the contents of this PR. If you are aware of documentation that should have been updated please update it and let your Hyaline administrator know about this message. Thanks!\n")
	}

	// Add raw data
	md.WriteString("\n")
	md.WriteString(fmt.Sprintf("%s\n", comment.RawData))

	return md.String()
}
