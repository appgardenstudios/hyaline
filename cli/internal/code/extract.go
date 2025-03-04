package code

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"path/filepath"
)

func ExtractCurrent(system string, cfg *config.Config, db *sql.DB) (err error) {
	// Find our target system (error if not found)
	var targetSystem *config.System
	for _, s := range cfg.Systems {
		if s.ID == system {
			targetSystem = &s
		}
	}
	if targetSystem == nil {
		// TODO better error message here
		return errors.New("system not found")
	}

	// Insert each code source
	for _, c := range targetSystem.Code {
		err = sqlite.InsertCurrentCode(sqlite.CurrentCode{
			ID:       c.ID,
			SystemID: targetSystem.ID,
			Path:     c.Path,
		}, db)

		absPath, err := filepath.Abs(c.Path)
		if err != nil {
			return err
		}
		fmt.Println(absPath)
		// TODO get code from path
		// TODO based on preset, get package.json and *.js files
		// TODO look for and use a glob library
	}

	return
}
