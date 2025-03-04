package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/code"
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
	cfg, err := config.Load(args.Config)
	if err != nil {
		return err
	}

	// Create/Scaffold SQLite
	absPath, err := filepath.Abs(args.Output)
	if err != nil {
		return err
	}
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

	// Insert System
	sqlite.InsertCurrentSystem(sqlite.CurrentSystem{
		ID: args.System,
	}, db)

	// Extract/Insert Code
	err = code.ExtractCurrent(args.System, cfg, db)
	if err != nil {
		return err
	}

	// Extract/Insert Docs
	// TODO

	return nil
}
