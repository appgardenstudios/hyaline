package prompts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"hyaline/internal/action"
)

// GoldenExpectation defines what we expect from a benchmark
type GoldenExpectation struct {
	ExpectedRecommendations []ExpectedRecommendation `json:"expectedRecommendations"`
	Description             string                   `json:"description"`
}

// ExpectedRecommendation defines an expected recommendation
type ExpectedRecommendation struct {
	Document  string   `json:"document"`
	Section   []string `json:"section,omitempty"`   // If omitted, matches document-level recommendations
	CheckType string   `json:"checkType,omitempty"` // For check current tests: "MATCHES_PURPOSE" or "COMPLETE"
}

// CheckResult is a unified result structure for both CheckChangeOutput and CheckCurrentOutput
type CheckResult struct {
	Recommendations []Recommendation `json:"recommendations"`
}

// Recommendation represents a single recommendation from either check command
type Recommendation struct {
	Document    string   `json:"document"`
	Section     []string `json:"section"`
	CheckType   string   `json:"checkType,omitempty"` // For check current: "MATCHES_PURPOSE" or "COMPLETE"
	Description string   `json:"description"`         // The actual recommendation/message
}

// RunResult contains the results of comparing actual vs expected
type RunResult struct {
	Scenario         string           `json:"scenario"`
	Summary          string           `json:"summary"`
	BenchmarkResults BenchmarkResults `json:"benchmarkResults"`
	RawResult        *CheckResult     `json:"rawResult"`
	Timestamp        string           `json:"timestamp"`
}

// MultiRunResult contains the results of multiple benchmark runs
type MultiRunResult struct {
	Scenario                string                  `json:"scenario"`
	Summary                 string                  `json:"summary"`
	RunResults              []RunResult             `json:"runResults"`
	AverageBenchmarkResults AverageBenchmarkResults `json:"averageBenchmarkResults"`
	Timestamp               string                  `json:"timestamp"`
}

// AverageBenchmarkResults contains averaged benchmark evaluation results
type AverageBenchmarkResults struct {
	AverageScore         float64 `json:"averageScore"`
	ExpectedMatchesAvg   float64 `json:"expectedMatchesAvg"`
	MissingMatchesAvg    float64 `json:"missingMatchesAvg"`
	UnexpectedMatchesAvg float64 `json:"unexpectedMatchesAvg"`
}

// BenchmarkResults contains exact matching results
type BenchmarkResults struct {
	ExpectedMatches   []string `json:"expectedMatches"`
	MissingMatches    []string `json:"missingMatches"`
	UnexpectedMatches []string `json:"unexpectedMatches"`
	Score             float64  `json:"score"`
}

func runHyalineBinary(args []string, t *testing.T) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current dir: %v", err)
	}

	binaryPath := filepath.Join(dir, "../../hyaline")
	// Run from CLI directory where .env is located
	workingDir := dir
	t.Log("binaryPath", binaryPath)
	t.Log("workingDir", workingDir)

	// Run hyaline binary directly, trusting environment variables are already loaded
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = workingDir

	// Pass through environment
	cmd.Env = os.Environ()

	return cmd.CombinedOutput()
}

// parseGoldenExpectation loads the expected results for a scenario
func parseGoldenExpectation(scenarioName string, t *testing.T) (*GoldenExpectation, error) {
	goldenPath := filepath.Join("_golden", scenarioName+".json")
	data, err := os.ReadFile(goldenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read golden file %s: %v", goldenPath, err)
	}

	var expectation GoldenExpectation
	if err := json.Unmarshal(data, &expectation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal golden file %s: %v", goldenPath, err)
	}

	return &expectation, nil
}

// parseCheckChangeResult parses the JSON output from hyaline check change commands
func parseCheckChangeResult(output []byte, t *testing.T) (*CheckResult, error) {
	var result action.CheckChangeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse check change result: %v", err)
	}

	var entries []Recommendation
	for _, rec := range result.Recommendations {
		entries = append(entries, Recommendation{
			Document:    rec.Document,
			Section:     rec.Section,
			Description: rec.Recommendation,
		})
	}
	return &CheckResult{Recommendations: entries}, nil
}

// parseCheckCurrentResult parses check current JSON output
// Only specific ERROR results are converted to recommendations:
// - Check type must be "MATCHES_PURPOSE" or "COMPLETE" (LLM-based checks)
// - Result must be "ERROR" (actual problems, not PASS/WARN/SKIPPED)
// - Message must NOT be about document/section existence (focus on content quality, not structural issues)
func parseCheckCurrentResult(data []byte, t *testing.T) (*CheckResult, error) {
	var result action.CheckCurrentOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse check current result: %v", err)
	}

	var entries []Recommendation

	for _, entry := range result.Results {
		if entry.Result == "ERROR" &&
			(entry.Check == "MATCHES_PURPOSE" || entry.Check == "COMPLETE") &&
			entry.Message != "This document does not exist" &&
			entry.Message != "This section does not exist" {
			entries = append(entries, Recommendation{
				Document:    entry.Document,
				Section:     entry.Section,
				CheckType:   entry.Check,
				Description: entry.Message,
			})
		}
	}

	return &CheckResult{Recommendations: entries}, nil
}

func evaluateBenchmark(actual *CheckResult, expected *GoldenExpectation) BenchmarkResults {
	var results BenchmarkResults

	// Create maps for efficient lookup
	actualRecs := make(map[string]bool)
	for _, rec := range actual.Recommendations {
		// Join section array into a single string for comparison
		sectionStr := strings.Join(rec.Section, "/")
		checkType := rec.CheckType
		if checkType == "" {
			key := fmt.Sprintf("%s:%s", rec.Document, sectionStr)
			actualRecs[key] = true
		} else {
			key := fmt.Sprintf("%s:%s:%s", rec.Document, sectionStr, checkType)
			actualRecs[key] = true
		}
	}

	expectedRecs := make(map[string]bool)
	for _, exp := range expected.ExpectedRecommendations {
		checkType := exp.CheckType
		var sectionStr string
		if len(exp.Section) == 0 {
			// No section specified - this matches the document directly
			sectionStr = ""
		} else {
			sectionStr = strings.Join(exp.Section, "/")
		}
		if checkType == "" {
			key := fmt.Sprintf("%s:%s", exp.Document, sectionStr)
			expectedRecs[key] = true
		} else {
			key := fmt.Sprintf("%s:%s:%s", exp.Document, sectionStr, checkType)
			expectedRecs[key] = true
		}
	}

	// Check for expected matches
	for expKey := range expectedRecs {
		if actualRecs[expKey] {
			results.ExpectedMatches = append(results.ExpectedMatches, expKey)
		} else {
			results.MissingMatches = append(results.MissingMatches, expKey)
		}
	}

	// Check for unexpected matches
	for actualKey := range actualRecs {
		if !expectedRecs[actualKey] {
			results.UnexpectedMatches = append(results.UnexpectedMatches, actualKey)
		}
	}

	results.Score = calculateScore(expectedRecs, results.MissingMatches, results.UnexpectedMatches)

	return results
}

const unexpectedPenalty = 0.25

// calculateScore computes the evaluation score based on matches and misses
// Score = (total_expected - missing - unexpectedPenalty*unexpected) / total_expected_required
func calculateScore(expectedRecs map[string]bool, missingMatches, unexpectedMatches []string) float64 {
	totalExpectedRequired := len(expectedRecs)
	missingCount := len(missingMatches)
	unexpectedCount := len(unexpectedMatches)

	if totalExpectedRequired > 0 {
		scoreNumerator := float64(totalExpectedRequired) - float64(missingCount) - unexpectedPenalty*float64(unexpectedCount)
		score := scoreNumerator / float64(totalExpectedRequired)
		return score
	}

	// If no expected recommendations, score based on unexpected penalty
	// Perfect score (1.0) when no unexpected, but penalize for unexpected recommendations
	if unexpectedCount == 0 {
		return 1.0 // Perfect score when no required matches and no unexpected ones
	}

	return -unexpectedPenalty * float64(unexpectedCount)
}

// addJSONSection adds a collapsible JSON section to the markdown buffer
func addJSONSection(buf *bytes.Buffer, title string, data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	buf.WriteString(fmt.Sprintf("<details>\n<summary>%s</summary>\n\n", title))
	buf.WriteString("```json\n")
	buf.WriteString(string(jsonData))
	buf.WriteString("\n```\n")
	buf.WriteString("\n</details>\n\n")
}

// generateMultiRunMarkdownReport creates a markdown report for multiple benchmark runs
func generateMultiRunMarkdownReport(report *MultiRunResult, expected *GoldenExpectation, outputPath string, t *testing.T) error {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("# Multi-Run Benchmark Report: %s\n\n", report.Scenario))
	buf.WriteString(fmt.Sprintf("**Generated:** %s\n\n", report.Timestamp))
	buf.WriteString(fmt.Sprintf("## Summary\n\n%s\n\n", report.Summary))

	// Average Results
	buf.WriteString("## Average Results Across All Runs\n\n")
	buf.WriteString("### Evaluation Results\n\n")
	buf.WriteString(fmt.Sprintf("**Average Score:** %.2f\n", report.AverageBenchmarkResults.AverageScore))
	buf.WriteString(fmt.Sprintf("**Total Expected:** %d\n\n", len(expected.ExpectedRecommendations)))

	buf.WriteString("| Metric | Average Count |\n")
	buf.WriteString("|--------|---------------|\n")
	buf.WriteString(fmt.Sprintf("| Expected Found | %.2f |\n", report.AverageBenchmarkResults.ExpectedMatchesAvg))
	buf.WriteString(fmt.Sprintf("| Missing | %.2f |\n", report.AverageBenchmarkResults.MissingMatchesAvg))
	buf.WriteString(fmt.Sprintf("| Unexpected | %.2f |\n", report.AverageBenchmarkResults.UnexpectedMatchesAvg))
	buf.WriteString("\n")

	// Individual Run Results
	buf.WriteString("## Individual Run Results\n\n")
	buf.WriteString("| Run | Score | Expected Found | Missing | Unexpected |\n")
	buf.WriteString("|-----|-------|----------------|---------|------------|\n")

	for i, run := range report.RunResults {
		buf.WriteString(fmt.Sprintf("| %d | %.2f | %d | %d | %d |\n",
			i+1,
			run.BenchmarkResults.Score,
			len(run.BenchmarkResults.ExpectedMatches),
			len(run.BenchmarkResults.MissingMatches),
			len(run.BenchmarkResults.UnexpectedMatches),
		))
	}
	buf.WriteString("\n")

	// Detailed breakdown for each run
	for i, run := range report.RunResults {
		buf.WriteString(fmt.Sprintf("### Run %d Details\n\n", i+1))

		hasDetails := false

		// Expected Recommendations Details
		if len(run.BenchmarkResults.ExpectedMatches) > 0 {
			addJSONSection(&buf, "Expected Recommendations", run.BenchmarkResults.ExpectedMatches)
			hasDetails = true
		}

		// Missing Recommendations Details
		if len(run.BenchmarkResults.MissingMatches) > 0 {
			addJSONSection(&buf, "Missing Recommendations", run.BenchmarkResults.MissingMatches)
			hasDetails = true
		}

		// Unexpected Recommendations Details
		if len(run.BenchmarkResults.UnexpectedMatches) > 0 {
			addJSONSection(&buf, "Unexpected Recommendations", run.BenchmarkResults.UnexpectedMatches)
			hasDetails = true
		}

		// Actual Recommendations for this run
		if run.RawResult != nil && len(run.RawResult.Recommendations) > 0 {
			addJSONSection(&buf, "Actual Recommendations", run.RawResult.Recommendations)
			hasDetails = true
		}

		// If no details were shown, add a message
		if !hasDetails {
			buf.WriteString("No recommendations were generated for this run.\n\n")
		}

	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// saveResult saves a benchmark result to a timestamped file
func saveResult(result interface{}, scenarioName string, suffix string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s-%s.json", scenarioName, suffix, timestamp)
	outputPath := filepath.Join("_output", filename)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write result: %v", err)
	}

	return outputPath, nil
}

// RunBenchmark runs a benchmark scenario with a specified number of iterations and outputs the results
func RunBenchmark(iterations int, scenarioName string, runFunc func(int) (*CheckResult, error), expected *GoldenExpectation, t *testing.T) {
	// Ensure local output directory exists
	if err := os.MkdirAll("_output", 0755); err != nil {
		t.Fatalf("Failed to create _output directory: %v", err)
	}

	var runResults []RunResult
	var benchmarkScores []float64
	var expectedMatchesCounts []int
	var missingMatchesCounts []int
	var unexpectedMatchesCounts []int

	for i := 0; i < iterations; i++ {
		t.Logf("Running iteration %d of %d for %s", i+1, iterations, scenarioName)

		// Run the benchmark scenario
		result, err := runFunc(i + 1)
		if err != nil {
			t.Fatalf("iteration %d failed: %v", i+1, err)
		}

		// Perform evaluations
		benchmarkResults := evaluateBenchmark(result, expected)

		// Create individual runResult
		runResult := RunResult{
			Scenario:         fmt.Sprintf("%s-run-%d", scenarioName, i+1),
			Summary:          generateSummary(result, expected, benchmarkResults),
			BenchmarkResults: benchmarkResults,
			RawResult:        result,
			Timestamp:        time.Now().Format(time.RFC3339),
		}

		runResults = append(runResults, runResult)

		// Collect metrics for averaging
		benchmarkScores = append(benchmarkScores, benchmarkResults.Score)
		expectedMatchesCounts = append(expectedMatchesCounts, len(benchmarkResults.ExpectedMatches))
		missingMatchesCounts = append(missingMatchesCounts, len(benchmarkResults.MissingMatches))
		unexpectedMatchesCounts = append(unexpectedMatchesCounts, len(benchmarkResults.UnexpectedMatches))

		t.Logf("Iteration %d - Benchmark: %.2f", i+1, benchmarkResults.Score)
	}

	// Calculate averages
	avgBenchmarkResults := AverageBenchmarkResults{
		AverageScore:         calculateAverage(benchmarkScores),
		ExpectedMatchesAvg:   calculateAverage(convertIntToFloat64(expectedMatchesCounts)),
		MissingMatchesAvg:    calculateAverage(convertIntToFloat64(missingMatchesCounts)),
		UnexpectedMatchesAvg: calculateAverage(convertIntToFloat64(unexpectedMatchesCounts)),
	}

	multiRunReport := &MultiRunResult{
		Scenario:                scenarioName,
		Summary:                 expected.Description,
		RunResults:              runResults,
		AverageBenchmarkResults: avgBenchmarkResults,
		Timestamp:               time.Now().Format(time.RFC3339),
	}

	// Save the multi-run evaluation as JSON
	reportPath, err := saveResult(multiRunReport, scenarioName, "multi-run-results")
	if err != nil {
		t.Logf("Warning: Failed to save multi-run evaluation raw results: %v", err)
	} else {
		t.Logf("Multi-run evaluation raw results saved to: %s", reportPath)
	}

	// Generate multi-run markdown report
	timestamp := time.Now().Format("20060102-150405")
	markdownPath := fmt.Sprintf("_output/%s-multi-run-report-%s.md", scenarioName, timestamp)
	if err := generateMultiRunMarkdownReport(multiRunReport, expected, markdownPath, t); err != nil {
		t.Logf("Warning: Failed to generate multi-run markdown report: %v", err)
	} else {
		t.Logf("Multi-run markdown report generated: %s", markdownPath)
	}
}

// calculateAverage calculates the average of a slice of float64 values
func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// convertIntToFloat64 converts a slice of int to float64
func convertIntToFloat64(values []int) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = float64(v)
	}
	return result
}

// generateSummary creates a human-readable summary of the benchmark results
func generateSummary(result *CheckResult, expected *GoldenExpectation, bench BenchmarkResults) string {
	summary := fmt.Sprintf("Scenario: %s\n", expected.Description)
	summary += fmt.Sprintf("Generated %d recommendations (expected %d)\n",
		len(result.Recommendations), len(expected.ExpectedRecommendations))
	summary += fmt.Sprintf("Benchmark score: %.2f", bench.Score)

	if len(bench.MissingMatches) > 0 {
		summary += fmt.Sprintf("\nMissing expected recommendations: %d", len(bench.MissingMatches))
	}
	if len(bench.UnexpectedMatches) > 0 {
		summary += fmt.Sprintf("\nUnexpected recommendations: %d", len(bench.UnexpectedMatches))
	}

	return summary
}
