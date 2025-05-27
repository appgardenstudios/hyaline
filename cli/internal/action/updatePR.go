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
	var comment *UpdatePRComment
	if args.Comment == "" {
		comment, err = updatePRAddComment(args.Sha, recommendations, args.PullRequest, cfg.GitHub.Token)
	} else {
		err = updatePRFromComment()
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

	// Write the byte slice to the file
	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.UpdatePR could not write output file", "error", err)
		return err
	}

	return nil
}

func updatePRFromComment() error {
	// TODO
	return nil
}

func updatePRAddComment(sha string, recommendations CheckChangeOutput, pr string, token string) (*UpdatePRComment, error) {
	// Format recs
	var recs []UpdatePRCommentRecommendation
	for _, rec := range recommendations.Recommendations {
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

func formatPRComment(comment *UpdatePRComment) string {
	var md strings.Builder

	md.WriteString("# Hyaline PR Check\n")
	md.WriteString(fmt.Sprintf("**ref**: %s\n", html.EscapeString(comment.Sha)))
	md.WriteString("- [ ] Trigger Re-run\n")
	md.WriteString("\n")

	if len(comment.Recommendations) > 0 {
		md.WriteString("### Review and update (if needed) the following document(s) and/or section(s):\n")
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
			md.WriteString(fmt.Sprintf("- [%s] %s%s in %s/%s", checked, html.EscapeString(rec.Document), html.EscapeString(sections), html.EscapeString(rec.System), html.EscapeString(rec.Source)))
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
