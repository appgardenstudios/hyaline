package prompts

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCheckChangeAddFeature(t *testing.T) {
	scenarioName := "check-change-add-feature"

	// Load golden expectations
	expected, err := parseGoldenExpectation(scenarioName, t)
	if err != nil {
		t.Fatalf("Failed to load golden expectation: %v", err)
	}

	// Define the run function for this specific scenario
	runFunc := func(iteration int) (*CheckResult, error) {
		outputPath := fmt.Sprintf("./_output/%s-output-run%d-%s.json", scenarioName, iteration, time.Now().Format("20060102-150405"))
		args := []string{
			"check", "change",
			"--config", "./_input/check-change-add-feature/config.yml",
			"--system", "url-shortener",
			"--current", "./_input/check-change-add-feature/current.sqlite",
			"--change", "./_input/check-change-add-feature/change.sqlite",
			"--output", outputPath,
		}

		output, err := runHyalineBinary(args, t)
		if err != nil {
			return nil, fmt.Errorf("hyaline command failed: %v\nOutput: %s", err, string(output))
		}

		outputData, err := os.ReadFile(outputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read output file %s: %v", outputPath, err)
		}

		result, err := parseCheckChangeResult(outputData, t)
		if err != nil {
			return nil, fmt.Errorf("failed to parse check change result: %v", err)
		}

		return result, nil
	}

	RunBenchmark(3, scenarioName, runFunc, expected, t)
}
