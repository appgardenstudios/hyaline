package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestCheckPRUpdatePR(t *testing.T) {
	outputPath := fmt.Sprintf("./_output/check-pr-update-pr-%d.json", time.Now().UnixMilli())
	
	// Run check pr to both generate recommendations and update PR comment
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr/hyaline.yml",
		"--documentation", "./_input/check-pr/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/8",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log("Check PR output:", string(stdOutStdErr))
	if err != nil {
		t.Fatal("Check PR failed:", err)
	}

	// Clean up the comment we just added
	// gh -R appgardenstudios/hyaline-example pr comment 8 --delete-last --yes
	cmd := exec.Command("gh", "-R", "appgardenstudios/hyaline-example", "pr", "comment", "8", "--delete-last", "--yes")
	cmd.Env = os.Environ()
	stdOutStdErr, err = cmd.CombinedOutput()
	t.Log("Cleanup output:", string(stdOutStdErr))
	if err != nil {
		t.Fatal("Cleanup failed:", err)
	}

	// Note: This demonstrates that check pr is now self-contained and doesn't need update pr
}