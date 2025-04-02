package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExtractCurrent(t *testing.T) {
	outputPath := fmt.Sprintf("./_output/extract-current-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "current",
		"--config", "./_input/extract-current/config.yml",
		"--system", "my-app",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	compareDBs("./_golden/extract-current.sqlite", outputPath, t)
}
