package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExportDocumentationJson(t *testing.T) {
	goldenPath := "./_golden/export-documentation-export.json"
	outputPath := fmt.Sprintf("./_output/export-documentation-json-%d.json", time.Now().UnixMilli())
	args := []string{
		"export", "documentation",
		"--documentation", "./_input/export-documentation-json/documentation.sqlite",
		"--format", "json",
		"--exclude", "document://*/cli/*",
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
