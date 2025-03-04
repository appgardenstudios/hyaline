package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
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
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		return err
	}
	fmt.Println(absPath)
	// Error if file exists as we want to ensure we start from an empty DB
	// NOTE this is fragile as the file could be created by another program between here and
	// when we open the DB
	_, err = os.Stat(absPath)
	if err == nil {
		return errors.New("output file already exists")
	}
	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		return err
	}
	defer db.Close()
	err = sqlite.CreateCurrentSchema(db)
	if err != nil {
		return err
	}

	// Extract Code
	// TODO

	// Extract Docs
	// TODO

	// Cleanup
	// TODO close sqlite

	return nil
}
