package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckPRCommentEmpty(t *testing.T) {
	goldenPath := "./_golden/check-pr-comment-empty.json"
	outputPath := fmt.Sprintf("./_output/check-pr-comment-empty-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr-comment-empty/hyaline.yml",
		"--documentation", "./_input/check-pr-comment-empty/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/1",
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