package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExtractCurrentFs(t *testing.T) {
	outputPath := fmt.Sprintf("./_output/extract-current-fs-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "current",
		"--config", "./_input/extract-current-fs/config.yml",
		"--system", "my-app",
		"--output", outputPath,
	}

	stdOutStdErr, err := runBinary(args, t)
	t.Log(string(stdOutStdErr))
	if err != nil {
		t.Fatal(err)
	}

	compareDBs("./_golden/extract-current-fs.sqlite", outputPath, t)
}
