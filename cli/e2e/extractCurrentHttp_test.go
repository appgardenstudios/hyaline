package e2e

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExtractCurrentHttp(t *testing.T) {
	// Start server on 8081 instead of an ephemeral port
	l, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Fatal(err)
	}
	fileServer := http.FileServer(http.Dir("./_input/extract-current-http/"))
	server := httptest.NewUnstartedServer(fileServer)
	server.Listener.Close()
	server.Listener = l
	server.Start()
	defer server.Close()

	goldenPath := "./_golden/extract-current-http.sqlite"
	outputPath := fmt.Sprintf("./_output/extract-current-http-%d.db", time.Now().UnixMilli())
	args := []string{
		"--debug",
		"extract", "current",
		"--config", "./_input/extract-current-http/config.yml",
		"--system", "my-app",
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
