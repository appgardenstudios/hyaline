package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckPRUpdatePR(t *testing.T) {
	checkOutputPath := fmt.Sprintf("./_output/check-pr-update-pr-check-%d.json", time.Now().UnixMilli())
	updateOutputPath := fmt.Sprintf("./_output/check-pr-update-pr-update-%d.json", time.Now().UnixMilli())
	
	// First run check pr to generate recommendations
	checkArgs := []string{
		"check", "pr",
		"--config", "./_input/check-pr/hyaline.yml",
		"--documentation", "./_input/check-pr/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/8",
		"--output", checkOutputPath,
	}

	stdOutStdErr, err := runBinary(checkArgs, t)
	t.Log("Check PR output:", string(stdOutStdErr))
	if err != nil {
		t.Fatal("Check PR failed:", err)
	}

	// Then run update pr using the generated recommendations
	updateArgs := []string{
		"update", "pr",
		"--config", "./_input/update-pr/config.yml",
		"--pull-request", "appgardenstudios/hyaline-example/8",
		"--sha", "4eefcddf5a56daff4734039dbb89129877a20dd5",
		"--recommendations", checkOutputPath,
		"--output", updateOutputPath,
	}

	stdOutStdErr, err = runBinary(updateArgs, t)
	t.Log("Update PR output:", string(stdOutStdErr))
	if err != nil {
		t.Fatal("Update PR failed:", err)
	}

	// Note: This is a basic workflow test
	// The full spec calls for more comprehensive workflow tests but this demonstrates the integration
}