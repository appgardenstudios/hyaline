package io

import (
	"encoding/json"
	"log/slog"
	"os"
)

// WriteJSON writes JSON data to a file
func WriteJSON(file *os.File, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Debug("io.WriteJSON could not marshal json", "error", err)
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		slog.Debug("io.WriteJSON could not write to file", "error", err)
		return err
	}

	return nil
}
