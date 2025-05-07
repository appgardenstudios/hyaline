package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckCurrentWithCompleteness(t *testing.T) {
	goldenPath := "./_golden/check-current-with-completeness-results.json"
	outputPath := fmt.Sprintf("./_output/check-current-with-completeness-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "current",
		"--config", "./_input/check-current-with-completeness/config.yml",
		"--current", "./_input/check-current-with-completeness/current.sqlite",
		"--system", "check-current",
		"--output", outputPath,
		"--check-completeness",
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
