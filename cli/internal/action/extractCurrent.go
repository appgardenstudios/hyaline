package action

import (
	"fmt"
	"hyaline/internal/config"
)

type ExtractCurrentArgs struct {
	Config string
	System string
	Output string
}

func ExtractCurrent(args *ExtractCurrentArgs) error {
	fmt.Println("Extracting Current", args)

	// Load Config
	_, err := config.Load(args.Config)
	if err != nil {
		return err
	}

	// Create/Scaffold SQLite
	// TODO

	// Extract Code
	// TODO

	// Extract Docs
	// TODO

	// Cleanup
	// TODO close sqlite

	return nil
}
