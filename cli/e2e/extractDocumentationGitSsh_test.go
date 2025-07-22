package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExtractDocumentationGitSsh(t *testing.T) {
	goldenPath := "./_golden/extract-documentation-git-ssh.sqlite"
	outputPath := fmt.Sprintf("./_output/extract-documentation-git-ssh-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "documentation",
		"--config", "./_input/extract-documentation-git-ssh/hyaline.yml",
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
