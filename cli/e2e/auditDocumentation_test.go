package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestAuditDocumentation(t *testing.T) {
	goldenPath := "./_golden/audit-documentation-results.json"
	outputPath := fmt.Sprintf("./_output/audit-documentation-%d.json", time.Now().UnixMilli())
	args := []string{
		"audit", "documentation",
		"--config", "./_input/audit-documentation/hyaline.yml",
		"--documentation", "./_input/audit-documentation/documentation.sqlite",
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