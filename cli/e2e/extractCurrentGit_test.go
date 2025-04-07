package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExtractCurrentGit(t *testing.T) {
	outputPath := fmt.Sprintf("./_output/extract-current-git-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "current",
		"--config", "./_input/extract-current-git/config.yml",
		"--system", "my-app",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	compareDBs("./_golden/extract-current-git.sqlite", outputPath, t)
}
