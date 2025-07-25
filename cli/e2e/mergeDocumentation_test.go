package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestMergeDocumentation(t *testing.T) {
	// Run the merge command
	goldenPath := "./_golden/merge-documentation.sqlite"
	outputPath := fmt.Sprintf("./_output/merge-documentation-%d.db", time.Now().UnixMilli())

	mergeArgs := []string{
		"merge", "documentation",
		"--input", "./_input/merge-documentation/input-1.sqlite",
		"--input", "./_input/merge-documentation/input-2.sqlite",
		"--input", "./_input/merge-documentation/input-3.sqlite",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(mergeArgs, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		updateGolden(goldenPath, outputPath, t)
	}

	compareDBs(goldenPath, outputPath, t)
}
