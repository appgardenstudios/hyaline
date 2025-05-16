package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExtractChangeWithIDs(t *testing.T) {
	goldenPath := "./_golden/extract-change-with-ids.sqlite"
	outputPath := fmt.Sprintf("./_output/extract-change-with-ids-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "change",
		"--config", "./_input/extract-change-with-ids/config.yml",
		"--system", "my-app",
		"--base", "main",
		"--head", "origin/feat-1",
		"--code-id", "app-ssh1",
		"--documentation-id", "md-ssh2",
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

	compareDBs(goldenPath, outputPath, t)
}
