package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckChangeWithSuggest(t *testing.T) {
	goldenPath := "./_golden/check-change-with-suggest-results.json"
	outputPath := fmt.Sprintf("./_output/check-change-with-suggest-%d.json", time.Now().UnixMilli())
	args := []string{
		"check", "change",
		"--config", "./_input/check-change-with-suggest/config.yml",
		"--current", "./_input/check-change-with-suggest/current.sqlite",
		"--change", "./_input/check-change-with-suggest/change.sqlite",
		"--system", "check-change",
		"--output", outputPath,
		"--suggest",
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
