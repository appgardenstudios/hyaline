package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateConfig(t *testing.T) {
	goldenPath := "./_golden/generate-config.yml"
	outputPath := fmt.Sprintf("./_output/generate-config-%d.yml", time.Now().UnixMilli())
	args := []string{
		"generate", "config",
		"--config", "./_input/generate-config/config.yml",
		"--current", "./_input/generate-config/current.sqlite",
		"--system", "generate-config",
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
