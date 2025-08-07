package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestCheckPRNewEmptyComment(t *testing.T) {
	goldenPath := "./_golden/check-pr-new-empty-comment.json"
	outputPath := fmt.Sprintf("./_output/check-pr-new-empty-comment-%d.json", time.Now().UnixMilli())
	outputCurrentPath := fmt.Sprintf("./_output/check-pr-new-empty-comment-current-%d.json", time.Now().UnixMilli())
	outputPreviousPath := fmt.Sprintf("./_output/check-pr-new-empty-comment-previous-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr-new-empty-comment/hyaline.yml",
		"--documentation", "./_input/check-pr-new-empty-comment/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/9",
		"--issue", "appgardenstudios/hyaline-example/2",
		"--issue", "appgardenstudios/hyaline-example/3",
		"--output", outputPath,
		"--output-current", outputCurrentPath,
		"--output-previous", outputPreviousPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the comment we just added
	// gh -R appgardenstudios/hyaline-example pr comment 9 --delete-last --yes
	cmd := exec.Command("gh", "-R", "appgardenstudios/hyaline-example", "pr", "comment", "9", "--delete-last", "--yes")
	cmd.Env = os.Environ()
	stdOutStdErr, err = cmd.CombinedOutput()
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
		updateGolden("./_golden/check-pr-new-empty-comment-current.json", outputCurrentPath, t)
		updateGolden("./_golden/check-pr-new-empty-comment-previous.json", outputPreviousPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
	compareFiles("./_golden/check-pr-new-empty-comment-current.json", outputCurrentPath, t)
	compareFiles("./_golden/check-pr-new-empty-comment-previous.json", outputPreviousPath, t)
}
