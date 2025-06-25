package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdatePRCommentEmpty(t *testing.T) {
	goldenPath := "./_golden/update-pr-comment-empty.json"
	outputPath := fmt.Sprintf("./_output/update-pr-comment-empty-%d.json", time.Now().UnixMilli())
	args := []string{
		"update", "pr",
		"--config", "./_input/update-pr-comment-empty/config.yml",
		"--pull-request", "appgardenstudios/hyaline-example/1",
		"--comment", "appgardenstudios/hyaline-example/2981785605",
		"--sha", "b4c5c73",
		"--recommendations", "./_input/update-pr-comment-empty/recommendations.json",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
