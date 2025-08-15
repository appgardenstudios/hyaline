package io

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

// InitOutput initializes an output file, ensuring it doesn't already exist
// Returns the file handle and any error
func InitOutput(outputPath string) (file *os.File, err error) {
	// Get absolute path
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		slog.Debug("io.InitOutput could not get an absolute path for output", "output", outputPath, "error", err)
		return
	}

	// Ensure output path does not exist
	_, err = os.Stat(absPath)
	if err == nil {
		err = errors.New("output file already exists")
		slog.Debug("io.InitOutput detected that output already exists", "absPath", absPath)
		return
	}

	// Create and open the file
	file, err = os.Create(absPath)
	if err != nil {
		slog.Debug("io.InitOutput could not create output file", "absPath", absPath, "error", err)
		return
	}

	return
}
