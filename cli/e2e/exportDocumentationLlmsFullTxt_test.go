package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExportDocumentationLlmsFullTxt(t *testing.T) {
	goldenPath := "./_golden/export-documentation-llmsfulltxt.txt"
	outputPath := fmt.Sprintf("./_output/export-documentation-llmsfulltxt-%d.txt", time.Now().UnixMilli())
	args := []string{
		"export", "documentation",
		"--documentation", "./_input/export-documentation-llmsfulltxt/documentation.sqlite",
		"--format", "llmsfulltxt",
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
