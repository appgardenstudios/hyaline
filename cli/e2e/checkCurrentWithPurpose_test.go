package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckCurrentWithPurpose(t *testing.T) {
	goldenPath := "./_golden/check-current-with-purpose-results.json"
	outputPath := fmt.Sprintf("./_output/check-current-with-purpose-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "current",
		"--config", "./_input/check-current-with-purpose/config.yml",
		"--current", "./_input/check-current-with-purpose/current.sqlite",
		"--system", "check-current",
		"--output", outputPath,
		"--check-purpose",
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
