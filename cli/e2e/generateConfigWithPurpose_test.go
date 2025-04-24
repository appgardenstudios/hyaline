package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateConfigWithPurpose(t *testing.T) {
	goldenPath := "./_golden/generate-config-with-purpose.yml"
	outputPath := fmt.Sprintf("./_output/generate-config-with-purpose-%d.yml", time.Now().UnixMilli())
	args := []string{
		"generate", "config",
		"--config", "./_input/generate-config-with-purpose/config.yml",
		"--current", "./_input/generate-config-with-purpose/current.sqlite",
		"--system", "generate-config",
		"--output", outputPath,
		"--include-purpose",
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
