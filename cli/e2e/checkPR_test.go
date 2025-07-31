package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckPR(t *testing.T) {
	goldenPath := "./_golden/check-pr.json"
	outputPath := fmt.Sprintf("./_output/check-pr-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "pr",
		"--config", "./_input/check-pr/hyaline.yml",
		"--documentation", "./_input/check-pr/documentation.sqlite",
		"--pull-request", "appgardenstudios/hyaline-example/8",
		"--issue", "appgardenstudios/hyaline-example/2",
		"--issue", "appgardenstudios/hyaline-example/3",
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