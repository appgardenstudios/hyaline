package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExportDocumentationFs(t *testing.T) {
	goldenPath1 := "./_golden/export-documentation-fs-README.md"
	goldenPath2 := "./_golden/export-documentation-fs-SECURITY.md"
	outputPath := fmt.Sprintf("./_output/export-documentation-fs-%d", time.Now().UnixMilli())
	outputPath1 := outputPath + "/README.md"
	outputPath2 := outputPath + "/hyaline/www/content/security.md"
	args := []string{
		"export", "documentation",
		"--documentation", "./_input/export-documentation-fs/documentation.sqlite",
		"--format", "fs",
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
		updateGolden(goldenPath1, outputPath1, t)
		updateGolden(goldenPath2, outputPath2, t)
	}

	compareFiles(goldenPath1, outputPath1, t)
	compareFiles(goldenPath2, outputPath2, t)
}
