package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestTestCheckDiff(t *testing.T) {
	goldenPath := "./_golden/check-diff.json"
	outputPath := fmt.Sprintf("./_output/check-diff-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "diff",
		"--config", "./_input/check-diff/hyaline.yml",
		"--documentation", "./_input/check-diff/documentation.sqlite",
		"--path", "../../../hyaline-example",
		"--base", "main",
		"--head", "feat-1",
		"--pull-request", "appgardenstudios/hyaline-example/1",
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
