package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	goldenPath := "./_golden/validate-config-output.json"
	outputPath := fmt.Sprintf("./_output/validate-config-%d.json", time.Now().UnixMilli())
	args := []string{
		"validate", "config",
		"--config", "./_input/validate-config/hyaline.yml",
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
