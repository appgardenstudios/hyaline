package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"hyaline/internal/config"
	"hyaline/internal/github"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
)

type UpdatePRArgs struct {
	Config          string
	PullRequest     string
	Comment         string
	Sha             string
	Recommendations string
	Output          string
}

type UpdatePRComment struct {
	Sha             string                          `json:"sha"`
	Recommendations []UpdatePRCommentRecommendation `json:"recommendations"`
	RawData         string                          `json:"rawData"`
}

type UpdatePRCommentRecommendation struct {
	Checked  bool     `json:"checked"`
	System   string   `json:"system"`
	Source   string   `json:"source"`
	Document string   `json:"document"`
	Section  []string `json:"section"`
	Reasons  []string `json:"reasons"`
}

type UpdatePRCommentRecommendationSort []UpdatePRCommentRecommendation

func (c UpdatePRCommentRecommendationSort) Len() int {
	return len(c)
}
func (c UpdatePRCommentRecommendationSort) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c UpdatePRCommentRecommendationSort) Less(i, j int) bool {
	if c[i].System < c[j].System {
		return true
	}
	if c[i].System > c[j].System {
		return false
	}
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
	return strings.Join(c[i].Section, "#") < strings.Join(c[j].Section, "#")
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
	cfg, err := config.Load(args.Config)
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
	var checkChangeOutput CheckChangeOutput
	err = json.Unmarshal(recsData, &checkChangeOutput)
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
	var comment *UpdatePRComment
	if args.Comment == "" {
		comment, err = updatePRAddComment(args.Sha, checkChangeOutput.Recommendations, args.PullRequest, cfg.GitHub.Token)
	} else {
		comment, err = updatePRUpdateComment(args.Sha, checkChangeOutput.Recommendations, args.PullRequest, args.Comment, cfg.GitHub.Token)
	}
	if err != nil {
		slog.Debug("action.UpdatePR could not update or add comment", "comment", args.Comment, "error", err)
		return err
	}

	// Write comment data to file
	jsonData, err := json.MarshalIndent(comment, "", "  ")
	if err != nil {
		slog.Debug("action.UpdatePR could not marshal json", "error", err)
		return err
	}
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.UpdatePR could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.UpdatePR could not write output file", "error", err)
		return err
	}

	return nil
}

func updatePRUpdateComment(sha string, newRecs []CheckChangeOutputEntry, pr string, comment string, token string) (*UpdatePRComment, error) {
	// Get comment
	existingComment, err := github.GetComment(comment, token)
	if err != nil {
		slog.Debug("action.updatePRFromComment could not get existing comment", "error", err)
		return nil, err
	}

	// Get existing recs from the comment
	existingRecs, err := parsePRComment(existingComment)
	if err != nil {
		slog.Debug("action.updatePRFromComment could not extract data from existing comment", "error", err)
		return nil, err
	}

	// Merge new recommendations with the ones from the comment
	mergedRecs := mergeRecs(newRecs, existingRecs)

	// Format New Raw Data
	rawData, err := formatRawData(&mergedRecs)
	if err != nil {
		slog.Debug("action.updatePRFromComment could not format raw data", "error", err)
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
		slog.Debug("action.updatePRFromComment could not update comment", "pr", pr, "error", err)
		return nil, err
	}

	return &updatedComment, nil
}

// Merge new recs into existing recs and return the resulting list (sorted)
func mergeRecs(newRecs []CheckChangeOutputEntry, existingRecs []UpdatePRCommentRecommendation) (mergedRecs []UpdatePRCommentRecommendation) {
	// Copy over existing recs as is
	mergedRecs = append(mergedRecs, existingRecs...)

	// Add new recs
	for _, newRec := range newRecs {
		// See if this rec already exists
		match := false
		for index, existingRec := range existingRecs {
			if newRecMatchesExisting(&newRec, &existingRec) {
				match = true
				// Do not overwrite existing rec's checked status so we preserve any "manual" checks
				if !existingRec.Checked {
					mergedRecs[index].Checked = newRec.Changed
				}
				// Always merge reasons
				mergedRecs[index].Reasons = mergeReasons(&newRec.Reasons, &existingRec.Reasons)
				break
			}
		}
		// If it does not, add it
		if !match {
			mergedRecs = append(mergedRecs, UpdatePRCommentRecommendation{
				Checked:  newRec.Changed,
				System:   newRec.System,
				Source:   newRec.DocumentationSource,
				Document: newRec.Document,
				Section:  newRec.Section,
				Reasons:  newRec.Reasons,
			})
		}
	}

	// Sort
	sort.Sort(UpdatePRCommentRecommendationSort(mergedRecs))

	return
}

func mergeReasons(newReasons *[]string, existingReasons *[]string) (mergedReasons []string) {
	mergedReasons = append(mergedReasons, *existingReasons...)

	for _, newReason := range *newReasons {
		if !slices.Contains(mergedReasons, newReason) {
			mergedReasons = append(mergedReasons, newReason)
		}
	}

	return
}

func newRecMatchesExisting(newRec *CheckChangeOutputEntry, existingRec *UpdatePRCommentRecommendation) bool {
	if newRec.System != existingRec.System {
		return false
	}
	if newRec.DocumentationSource != existingRec.Source {
		return false
	}
	if newRec.Document != existingRec.Document {
		return false
	}

	return reflect.DeepEqual(newRec.Section, existingRec.Section)
}

func updatePRAddComment(sha string, recommendations []CheckChangeOutputEntry, pr string, token string) (*UpdatePRComment, error) {
	// Format recs
	var recs []UpdatePRCommentRecommendation
	for _, rec := range recommendations {
		recs = append(recs, UpdatePRCommentRecommendation{
			Checked:  rec.Changed,
			System:   rec.System,
			Source:   rec.DocumentationSource,
			Document: rec.Document,
			Section:  rec.Section,
			Reasons:  rec.Reasons,
		})
	}

	// Format raw data
	rawData, err := formatRawData(&recs)
	if err != nil {
		slog.Debug("action.updatePRAddComment could not format raw data", "error", err)
		return nil, err
	}

	// Create comment
	comment := UpdatePRComment{
		Sha:             sha,
		Recommendations: recs,
		RawData:         rawData,
	}

	// Format comment
	formattedComment := formatPRComment(&comment)

	// Add comment to PR
	err = github.AddComment(pr, formattedComment, token)
	if err != nil {
		slog.Debug("action.updatePRAddComment could not add comment", "pr", pr, "error", err)
		return nil, err
	}

	return &comment, nil
}

func formatRawData(recommendations *[]UpdatePRCommentRecommendation) (string, error) {
	data, err := json.Marshal(recommendations)
	if err != nil {
		slog.Debug("action.formatRawData could not marshal json", "error", err)
		return "", err
	}

	return fmt.Sprintf("<![CDATA[ %s ]]>", string(data)), nil
}

const CDATA_START = "<![CDATA[ "
const CDATA_END = " ]]>"
const RECOMMENDATIONS_START = "### Recommendations"

func parsePRComment(comment string) (recs []UpdatePRCommentRecommendation, err error) {
	// Get CData
	dataStart := strings.Index(comment, CDATA_START)
	if dataStart == -1 {
		err = errors.New("could not find start of CDATA block")
		return
	}
	dataEnd := strings.LastIndex(comment, CDATA_END)
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
	recommendationsStart := strings.Index(comment, RECOMMENDATIONS_START)
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

func formatPRComment(comment *UpdatePRComment) string {
	var md strings.Builder

	// Note: The comment MUST start with "# Hyaline" filled with 0-width spaces
	md.WriteString("# H\u200By\u200Ba\u200Bl\u200Bi\u200Bn\u200Be PR Check\n")
	md.WriteString(fmt.Sprintf("**ref**: %s\n", html.EscapeString(comment.Sha)))
	md.WriteString("\n")

	// Note: This starting line always needs to be present because we use it as a sentinel for getting the check marks
	md.WriteString(fmt.Sprintf("%s\n", RECOMMENDATIONS_START))
	if len(comment.Recommendations) > 0 {
		md.WriteString("**Hyaline automatically detects documentation changes** and will check off items as you update them. ")
		md.WriteString("Review the recommendations below and manually update any remaining items as needed.\n")
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
			md.WriteString(fmt.Sprintf("- [%s] **%s**%s in `%s/%s`", checked, html.EscapeString(rec.Document), html.EscapeString(sections), html.EscapeString(rec.System), html.EscapeString(rec.Source)))
			md.WriteString(fmt.Sprintf("<details><summary>Reasons</summary><ul><li>%s</li></ul></details>", reasons))
			md.WriteString("\n")
		}
	} else {
		md.WriteString("Hyaline did not find any documentation related to the contents of this PR. If you are aware of documentation that should have been updated please update it and let your Hyaline administrator know about this message. Thanks!\n")
	}

	// Add raw data
	md.WriteString("\n")
	md.WriteString(fmt.Sprintf("%s\n", comment.RawData))

	return md.String()
}
