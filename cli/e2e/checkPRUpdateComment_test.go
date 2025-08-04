package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckPRUpdateComment(t *testing.T) {
	goldenPath := "./_golden/check-pr-update-comment.json"
	outputPath := fmt.Sprintf("./_output/check-pr-update-comment-%d.json", time.Now().UnixMilli())

	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr-update-comment/hyaline.yml",
		"--documentation", "./_input/check-pr-update-comment/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/8",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log("Check PR output:", string(stdOutStdErr))
	if err != nil {
		t.Fatal("Check PR failed:", err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}
