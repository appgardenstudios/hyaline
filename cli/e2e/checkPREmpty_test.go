package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestCheckPREmpty(t *testing.T) {
	goldenPath := "./_golden/check-pr-empty.json"
	outputPath := fmt.Sprintf("./_output/check-pr-empty-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr-empty/hyaline.yml",
		"--documentation", "./_input/check-pr-empty/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/1",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the comment we just added
	// gh -R appgardenstudios/hyaline-example pr comment 1 --delete-last --yes
	cmd := exec.Command("gh", "-R", "appgardenstudios/hyaline-example", "pr", "comment", "1", "--delete-last", "--yes")
	cmd.Env = os.Environ()
	stdOutStdErr, err = cmd.CombinedOutput()
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareFiles(goldenPath, outputPath, t)
}