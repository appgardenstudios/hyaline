package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckChange(t *testing.T) {
	goldenPath := "./_golden/check-change-results.json"
	outputPath := fmt.Sprintf("./_output/check-change-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "change",
		"--config", "./_input/check-change/config.yml",
		"--current", "./_input/check-change/current.sqlite",
		"--change", "./_input/check-change/change.sqlite",
		"--system", "check-change",
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
