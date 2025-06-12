package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestUpdatePR(t *testing.T) {
	goldenPath := "./_golden/update-pr.json"
	outputPath := fmt.Sprintf("./_output/update-pr-%d.json", time.Now().UnixMilli())
	args := []string{
		"update", "pr",
		"--config", "./_input/update-pr/config.yml",
		"--pull-request", "appgardenstudios/hyaline-example/1",
		"--sha", "b4c5c73",
		"--recommendations", "./_input/update-pr/recommendations.json",
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
