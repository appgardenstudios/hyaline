package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var binaryName = "hyaline-coverage"

var binaryPath = ""

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get current dir: %v", err)
	}

	binaryPath = filepath.Join(dir, binaryName)

	os.Exit(m.Run())
}

func runBinary(args []string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}

func TestExtractCurrent(t *testing.T) {
	output := fmt.Sprintf("./e2e/_output/extract-current-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "current",
		"--config", "./e2e/_input/extract-current/config.yml",
		"--system", "my-app",
		"--output", output,
	}

	stdOutStdErr, err := runBinary(args)
	if err != nil {
		t.Log(string(stdOutStdErr))
		t.Fatal(err)
	}

	// TODO ensure golden file and output are the same
}

func compareDBs(db1 string, db2 string) {
	// TODO
}
