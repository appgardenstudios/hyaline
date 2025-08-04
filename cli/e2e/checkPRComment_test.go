package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckPRComment(t *testing.T) {
	goldenPath := "./_golden/check-pr-comment.json"
	outputPath := fmt.Sprintf("./_output/check-pr-comment-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr-comment/hyaline.yml",
		"--documentation", "./_input/check-pr-comment/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/8",
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