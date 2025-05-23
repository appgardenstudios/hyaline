package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckCurrent(t *testing.T) {
	goldenPath := "./_golden/check-current-results.json"
	outputPath := fmt.Sprintf("./_output/check-current-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "current",
		"--config", "./_input/check-current/config.yml",
		"--current", "./_input/check-current/current.sqlite",
		"--system", "check-current",
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
