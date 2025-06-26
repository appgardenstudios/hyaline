package e2e

import (
	"fmt"
	"testing"
	"time"
)

func TestExportLlmsTxt(t *testing.T) {
	goldenPath := "./_golden/export-llms-txt.txt"
	outputPath := fmt.Sprintf("./_output/export-llms-txt-%d.txt", time.Now().UnixMilli())
	args := []string{
		"export", "llms-txt",
		"--current", "./_input/check-current/current.sqlite",
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

func TestExportLlmsTxtFull(t *testing.T) {
	goldenPath := "./_golden/export-llms-txt-full.txt"
	outputPath := fmt.Sprintf("./_output/export-llms-txt-full-%d.txt", time.Now().UnixMilli())
	args := []string{
		"export", "llms-txt",
		"--current", "./_input/check-current/current.sqlite",
		"--output", outputPath,
		"--full",
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

func TestExportLlmsTxtFiltered(t *testing.T) {
	goldenPath := "./_golden/export-llms-txt-filtered.txt"
	outputPath := fmt.Sprintf("./_output/export-llms-txt-filtered-%d.txt", time.Now().UnixMilli())
	args := []string{
		"export", "llms-txt",
		"--current", "./_input/check-current/current.sqlite",
		"--output", outputPath,
		"--document-uri", "docs",
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