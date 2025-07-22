package prompts

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCheckCurrentAPIDocs(t *testing.T) {
	scenarioName := "check-current-api-docs"

	expected, err := parseGoldenExpectation(scenarioName, t)
	if err != nil {
		t.Fatalf("Failed to load golden expectation: %v", err)
	}

	runFunc := func(iteration int) (*CheckResult, error) {
		outputPath := fmt.Sprintf("./_output/%s-output-run%d-%s.json", scenarioName, iteration, time.Now().Format("20060102-150405"))
		args := []string{
			"check", "current",
			"--config", "./_input/check-current-api-docs/config.yml",
			"--system", "url-shortener",
			"--current", "./_input/check-current-api-docs/current.sqlite",
			"--output", outputPath,
			"--check-purpose",
			"--check-completeness",
		}

		output, err := runHyalineBinary(args, t)
		if err != nil {
			return nil, fmt.Errorf("hyaline command failed: %v\nOutput: %s", err, string(output))
		}

		outputData, err := os.ReadFile(outputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read output file %s: %v", outputPath, err)
		}
		result, err := parseCheckCurrentResult(outputData, t)
		if err != nil {
			return nil, fmt.Errorf("failed to parse check current result: %v", err)
		}

		return result, nil
	}
	RunBenchmark(3, scenarioName, runFunc, expected, t)
}
