package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExportDocumentationSqlite(t *testing.T) {
	goldenPath := "./_golden/export-documentation-export.sqlite"
	outputPath := fmt.Sprintf("./_output/export-documentation-sqlite-%d.db", time.Now().UnixMilli())
	args := []string{
		"export", "documentation",
		"--documentation", "./_input/export-documentation-sqlite/documentation.sqlite",
		"--format", "sqlite",
		"--include", "document://*/**/*?type=security",
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
