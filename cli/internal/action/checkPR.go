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
	"hyaline/internal/io"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"strings"
)

const CHECK_PR_HYALINE_HEADER = "# H\u200By\u200Ba\u200Bl\u200Bi\u200Bn\u200Be"
const CHECK_PR_CDATA_START = "<![CDATA[ "
const CHECK_PR_CDATA_END = " ]]>"
const CHECK_PR_RECOMMENDATIONS_START = "### Recommendations"

type CheckPRArgs struct {
	Config         string
	Documentation  string
	PullRequest    string
	Issues         []string
	Output         string
	OutputCurrent  string
	OutputPrevious string
}

func CheckPR(args *CheckPRArgs) error {
	slog.Info("Checking PR",
		"config", args.Config,
		"documentation", args.Documentation,
		"pull-request", args.PullRequest,
		"issues", args.Issues,
		"output", args.Output,
		"output-current", args.OutputCurrent,
		"output-previous", args.OutputPrevious,
	)

	// Load Config
	cfg, err := config.Load(args.Config, true)
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

	// If check is disabled, skip
	if cfg.Check.Disabled {
		slog.Info("Check disabled. Skipping...")
		return nil
	}

	// Ensure GitHub token is available
	if cfg.GitHub.Token == "" {
		return errors.New("github token required to retrieve pull-request information")
	}

	// Initialize output files if provided
	var outputFile, outputCurrentFile, outputPreviousFile *os.File

	if args.Output != "" {
		outputFile, err = io.InitOutput(args.Output)
		if err != nil {
			slog.Debug("action.CheckPR could not initialize output file", "error", err)
			return err
		}
		defer outputFile.Close()
	}

	if args.OutputCurrent != "" {
		outputCurrentFile, err = io.InitOutput(args.OutputCurrent)
		if err != nil {
			slog.Debug("action.CheckPR could not initialize output-current file", "error", err)
			return err
		}
		defer outputCurrentFile.Close()
	}

	if args.OutputPrevious != "" {
		outputPreviousFile, err = io.InitOutput(args.OutputPrevious)
		if err != nil {
			slog.Debug("action.CheckPR could not initialize output-previous file", "error", err)
			return err
		}
		defer outputPreviousFile.Close()
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
	recommendations, fileCheckContextHashes, err := getRecommendations(filteredFiles, documents, pr, issues, changedFiles, cfg.Check, &cfg.LLM)
	if err != nil {
		slog.Debug("action.CheckPR could not get recommendations", "error", err)
		return err
	}
	slog.Info("Retrieved recommendations", "recommendations", len(recommendations))

	// Get existing Hyaline comment
	existingComment, err := findHyalineComment(args.PullRequest, cfg.GitHub.Token)
	if err != nil {
		slog.Debug("action.CheckPR could not search for existing Hyaline comment", "error", err)
		return err
	}

	// Get previous recommendations from existing comment if available
	var previousOutput *CheckOutput
	if existingComment != nil {
		slog.Info("Retrieving previous recommendations from comment", "commentID", existingComment.ID)
		previousOutput, err = parseCheckPRComment(existingComment.Body)
		if err != nil {
			slog.Debug("action.CheckPR could not parse existing comment", "error", err)
			return err
		}
	}

	previousRecommendations := []CheckRecommendation{}
	if previousOutput != nil {
		slog.Info("Found previous recommendations", "recommendations", len(previousOutput.Recommendations))
		previousRecommendations = previousOutput.Recommendations
	}

	// Merge current and previous recommendations
	mergedRecommendations := mergeCheckRecommendations(recommendations, previousRecommendations, fileCheckContextHashes)

	mergedOutput := CheckOutput{
		Recommendations: mergedRecommendations,
		Head:            pr.Head,
		Base:            pr.Base,
	}

	slog.Info("Merged recommendations", "mergedRecommendations", len(mergedRecommendations))

	err = upsertPRComment(args.PullRequest, existingComment, mergedOutput, cfg.GitHub.Token)
	if err != nil {
		slog.Debug("action.CheckPR could not upsert PR comment", "error", err)
		return err
	}

	// Write merged recommendations to output file if provided
	if outputFile != nil {
		output := CheckOutput{
			Recommendations: mergedRecommendations,
			Head:            pr.Head,
			Base:            pr.Base,
		}

		err = io.WriteJSON(outputFile, output)
		if err != nil {
			slog.Debug("action.CheckPR could not write merged recommendations to output file", "error", err)
			return err
		}
		slog.Info("Output merged recommendations", "recommendations", len(mergedRecommendations), "output", args.Output)
	}

	// Write current recommendations to output-current file if provided
	if outputCurrentFile != nil {
		output := CheckOutput{
			Recommendations: recommendations,
			Head:            pr.Head,
			Base:            pr.Base,
		}

		err = io.WriteJSON(outputCurrentFile, output)
		if err != nil {
			slog.Debug("action.CheckPR could not write current recommendations to output file", "error", err)
			return err
		}
		slog.Info("Output current recommendations", "recommendations", len(recommendations), "output", args.OutputCurrent)
	}

	// Write previous recommendations to output-previous file if provided
	if outputPreviousFile != nil {
		if previousOutput == nil {
			// Write empty JSON object when no previous output exists
			_, err = outputPreviousFile.Write([]byte("{}"))
			if err != nil {
				slog.Debug("action.CheckPR could not write empty previous recommendations to output file", "error", err)
				return err
			}
		} else {
			err = io.WriteJSON(outputPreviousFile, previousOutput)
			if err != nil {
				slog.Debug("action.CheckPR could not write previous recommendations to output file", "error", err)
				return err
			}
		}
		slog.Info("Output previous recommendations", "recommendations", len(previousRecommendations), "output", args.OutputPrevious)
	}

	return nil
}

func upsertPRComment(pr string, existingComment *github.Comment, output CheckOutput, token string) error {
	formattedComment := formatCheckPRComment(&output)

	if existingComment != nil {
		slog.Info("Updating existing PR comment", "commentID", existingComment.ID)
		prParts := strings.Split(pr, "/")
		commentRef := fmt.Sprintf("%s/%s/%d", prParts[0], prParts[1], existingComment.ID)

		err := github.UpdateComment(commentRef, formattedComment, token)
		if err != nil {
			slog.Debug("updatePRComment could not update comment", "commentRef", commentRef, "error", err)
			return err
		}
	} else {
		slog.Info("Adding new PR comment")

		// Add comment to PR
		err := github.AddComment(pr, formattedComment, token)
		if err != nil {
			slog.Debug("addPRComment could not add comment", "pr", pr, "error", err)
			return err
		}
	}

	return nil
}

func mergeCheckRecommendations(newRecs []CheckRecommendation, existingRecs []CheckRecommendation, fileCheckContextHashes check.FileCheckContextHashes) (mergedRecs []CheckRecommendation) {
	// Initialize to empty slice to avoid returning nil
	mergedRecs = []CheckRecommendation{}

	// 1. Copy over all new recommendations as-is
	mergedRecs = append(mergedRecs, newRecs...)

	// 2. Merge existing recommendations
	for _, existingRec := range existingRecs {
		// See if this rec already exists
		match := false
		for index, mergedRec := range mergedRecs {
			// If a matching recommendation is found, merge it with the existing one
			if recommendationsMatch(&mergedRec, &existingRec) {
				match = true
				// Merge recommendations
				updatedMergedRec := mergedRec
				updatedMergedRec.Checked = existingRec.Checked
				updatedMergedRec.Reasons = mergeCheckReasons(&mergedRec.Reasons, &existingRec.Reasons, fileCheckContextHashes)
				mergedRecs[index] = updatedMergedRec
				break
			}
		}

		// If no match was found, append it
		if !match {
			mergedRec := existingRec
			// Call mergeCheckReasons with an empty newReasons so that outdated old reasons still get marked outdated
			mergedRec.Reasons = mergeCheckReasons(&[]check.Reason{}, &existingRec.Reasons, fileCheckContextHashes)
			mergedRecs = append(mergedRecs, mergedRec)
		}
	}

	// 3. Mark recommendations as outdated if all their reasons are outdated
	for i := range mergedRecs {
		allReasonsOutdated := true
		for _, reason := range mergedRecs[i].Reasons {
			if !reason.Outdated {
				allReasonsOutdated = false
				break
			}
		}
		mergedRecs[i].Outdated = allReasonsOutdated
	}

	// 4. Sort recommendations
	sortCheckRecommendations(mergedRecs)

	return
}

func mergeCheckReasons(newReasons *[]check.Reason, existingReasons *[]check.Reason, fileCheckContextHashes check.FileCheckContextHashes) (mergedReasons []check.Reason) {
	// 1. Copy over all new reasons as-is
	mergedReasons = append(mergedReasons, *newReasons...)

	// 2. Merge existing reasons
	for _, existingReason := range *existingReasons {
		// See if the reason already exists
		match := false
		for _, mergedReason := range mergedReasons {
			if reasonsMatch(&mergedReason, &existingReason) {
				match = true
				break
			}
		}

		// If no match was found, merge the existing reason
		if !match {
			updatedExistingReason := existingReason
			// Check if the file context hash has changed. If so, or if it doesn't exist, mark as outdated
			if fileContextHash, exists := fileCheckContextHashes[existingReason.Check.File]; exists {
				if fileContextHash[existingReason.Check.Type] != existingReason.Check.ContextHash {
					updatedExistingReason.Outdated = true
				}
			} else {
				updatedExistingReason.Outdated = true
			}
			// Add the existing reason to the merged reasons
			mergedReasons = append(mergedReasons, updatedExistingReason)
		}
	}

	sortCheckReasons(mergedReasons)

	return
}

func recommendationsMatch(recA *CheckRecommendation, recB *CheckRecommendation) bool {
	if recA.Source != recB.Source {
		return false
	}
	if recA.Document != recB.Document {
		return false
	}

	return strings.Join(recA.Section, "/") == strings.Join(recB.Section, "/")
}

func reasonsMatch(reasonA *check.Reason, reasonB *check.Reason) bool {
	if reasonA.Check.File != reasonB.Check.File {
		return false
	}

	return reasonA.Check.Type == reasonB.Check.Type
}

func findHyalineComment(ref string, token string) (*github.Comment, error) {
	// Get all comments for the PR
	comments, err := github.ListComments(ref, token)
	if err != nil {
		return nil, err
	}

	// Search for a comment that starts with the Hyaline header (with zero-width spaces)
	for _, comment := range comments {
		if strings.HasPrefix(comment.Body, CHECK_PR_HYALINE_HEADER) {
			return &comment, nil
		}
	}

	// No Hyaline comment found
	return nil, nil
}

func formatCheckPRRawData(output *CheckOutput) (string, error) {
	data, err := json.Marshal(output)
	if err != nil {
		slog.Debug("formatCheckPRRawData could not marshal json", "error", err)
		return "", err
	}

	return fmt.Sprintf("<![CDATA[ %s ]]>", string(data)), nil
}

func parseCheckPRComment(comment string) (output *CheckOutput, err error) {
	// Get CData
	dataStart := strings.Index(comment, CHECK_PR_CDATA_START)
	if dataStart == -1 {
		err = errors.New("could not find start of CDATA block")
		return
	}
	dataEnd := strings.LastIndex(comment, CHECK_PR_CDATA_END)
	if dataEnd == -1 {
		err = errors.New("could not find end of CDATA block")
		return
	}
	data := comment[dataStart+10 : dataEnd]

	// Parse as CheckOutput
	output = &CheckOutput{}
	err = json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}

	// Get any manual checks and update the corresponding rec
	// Note that this relies on the CData order being the same as the recs list order
	recommendationsStart := strings.Index(comment, CHECK_PR_RECOMMENDATIONS_START)
	if recommendationsStart == -1 {
		err = errors.New("could not find start of recommendations")
		return
	}
	checks := comment[recommendationsStart:dataStart]
	lines := strings.Split(checks, "\n")
	currentRec := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "- [") && currentRec < len(output.Recommendations) {
			// Always pull checked from the list and NOT CData
			if strings.HasPrefix(line, "- [ ]") {
				output.Recommendations[currentRec].Checked = false
			}
			if strings.HasPrefix(line, "- [x]") {
				output.Recommendations[currentRec].Checked = true
			}
			currentRec++
		}
	}

	return
}

func formatCheckPRComment(output *CheckOutput) string {
	var md strings.Builder

	// Note: The comment MUST start with the Hyaline header
	md.WriteString(CHECK_PR_HYALINE_HEADER + " PR Check\n")
	md.WriteString(fmt.Sprintf("**ref**: %s\n", html.EscapeString(output.Head)))
	md.WriteString("\n")

	// Split recommendations into valid and outdated
	var validRecs []CheckRecommendation
	var outdatedRecs []CheckRecommendation

	for _, rec := range output.Recommendations {
		if rec.Outdated {
			outdatedRecs = append(outdatedRecs, rec)
		} else {
			validRecs = append(validRecs, rec)
		}
	}

	// Note: This starting line always needs to be present because we use it as a sentinel for getting the check marks
	md.WriteString(fmt.Sprintf("%s\n", CHECK_PR_RECOMMENDATIONS_START))

	// Render valid recommendations
	if len(validRecs) > 0 {
		md.WriteString("Review the following recommendations and update the corresponding documentation as needed:\n")
		for _, rec := range validRecs {
			checked := " "
			if rec.Checked {
				checked = "x"
			}
			sections := formatSections(rec.Section)
			reasons := formatReasons(rec.Reasons)
			md.WriteString(fmt.Sprintf("- [%s] **%s**%s in `%s`", checked, html.EscapeString(rec.Document), html.EscapeString(sections), html.EscapeString(rec.Source)))
			md.WriteString(fmt.Sprintf("<details><summary>Reasons</summary><ul><li>%s</li></ul></details>", reasons))
			md.WriteString("\n")
		}
		md.WriteString("\nNote: Hyaline will automatically detect documentation updated in this PR and mark corresponding recommendations as reviewed.\n")
	} else if len(outdatedRecs) == 0 {
		md.WriteString("Hyaline did not find any documentation related to the contents of this PR. If you are aware of documentation that should have been updated please update it and let your Hyaline administrator know about this message. Thanks!\n")
	}

	// Render outdated recommendations section if any exist
	if len(outdatedRecs) > 0 {
		md.WriteString("\n<details><summary>Changes have caused the following recommendations to be outdated:</summary>\n\n")
		for _, rec := range outdatedRecs {
			sections := formatSections(rec.Section)
			reasons := formatReasons(rec.Reasons)
			md.WriteString(fmt.Sprintf("- **%s**%s in `%s`", html.EscapeString(rec.Document), html.EscapeString(sections), html.EscapeString(rec.Source)))
			md.WriteString(fmt.Sprintf("<details><summary>Reasons</summary><ul><li>%s</li></ul></details>", reasons))
			md.WriteString("\n")
		}
		md.WriteString("</details>\n")
	}

	// Add raw data
	rawData, _ := formatCheckPRRawData(output)
	md.WriteString("\n")
	md.WriteString(fmt.Sprintf("%s\n", rawData))

	return md.String()
}

func formatSections(sections []string) string {
	if len(sections) > 0 {
		return fmt.Sprintf(" > %s", strings.Join(sections, " > "))
	}
	return ""
}

func formatReasons(reasons []check.Reason) string {
	var cleanReasons []string
	for _, reason := range reasons {
		reasonText := html.EscapeString(reason.Reason)
		if reason.Outdated {
			reasonText = fmt.Sprintf("~~%s~~ (Outdated)", reasonText)
		}
		cleanReasons = append(cleanReasons, reasonText)
	}
	return strings.Join(cleanReasons, "</li><li>")
}
