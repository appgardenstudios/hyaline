package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	goldenPath := "./_golden/merge.sqlite"
	outputPath := fmt.Sprintf("./_output/merge-%d.db", time.Now().UnixMilli())
	args := []string{
		"merge",
		"--input", "./_input/merge/current.sqlite",
		"--input", "./_input/merge/current-copy.sqlite",
		"--input", "./_input/merge/change.sqlite",
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

	compareDBs(goldenPath, outputPath, t)
}
