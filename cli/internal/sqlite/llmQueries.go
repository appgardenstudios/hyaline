package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

type GetCurrentCodeFileRow struct {
	ID      string
	RawData string
}

func GetCurrentCodeFiles(files []string, systemID string, db *sql.DB) (*[]GetCurrentCodeFileRow, error) {
	// Guard against an empty array of files and short-circuit
	if len(files) == 0 {
		return nil, nil
	}

	// Construct our list of ids and the corresponding placeholder string (for in)
	ids := []any{systemID}
	for _, file := range files {
		ids = append(ids, file)
	}
	placeholders := strings.Join(strings.Split(strings.Repeat("?", len(files)), ""), ",")

	// Construct our statement
	stmt, err := db.Prepare(fmt.Sprintf(`
select
	id,
	raw_data
from
	file
where
  system_id = ?
	AND id IN (%s)
`, placeholders))
	if err != nil {
		return nil, err
	}

	// Get rows
	var rows []GetCurrentCodeFileRow
	slog.Debug("sqlite.GetCodeFile executing query", "query", fmt.Sprintf("%v", stmt), "ids", ids)
	rawRows, err := stmt.Query(ids...)
	if err != nil {
		return nil, err
	}
	defer rawRows.Close()

	for rawRows.Next() {
		var row GetCurrentCodeFileRow
		if err := rawRows.Scan(&row.ID, &row.RawData); err != nil {
			return &rows, err
		}
		rows = append(rows, row)
	}
	if err = rawRows.Err(); err != nil {
		return &rows, err
	}

	return &rows, nil
}

type GetChangeCodeFileRow struct {
	ID      string
	Action  string
	RawData string
}

func GetChangeCodeFiles(files []string, systemID string, db *sql.DB) (*[]GetChangeCodeFileRow, error) {
	// Guard against an empty array of files and short-circuit
	if len(files) == 0 {
		return nil, nil
	}

	// Construct our list of ids and the corresponding placeholder string (for in)
	ids := []any{systemID}
	for _, file := range files {
		ids = append(ids, file)
	}
	placeholders := strings.Join(strings.Split(strings.Repeat("?", len(files)), ""), ",")

	// Construct our statement
	stmt, err := db.Prepare(fmt.Sprintf(`
select
	id,
	action,
	raw_data
from
	file
where
  system_id = ?
	AND id IN (%s)
`, placeholders))
	if err != nil {
		return nil, err
	}

	// Get rows
	var rows []GetChangeCodeFileRow
	slog.Debug("sqlite.GetCodeFile executing query", "query", fmt.Sprintf("%v", stmt), "ids", ids)
	rawRows, err := stmt.Query(ids...)
	if err != nil {
		return nil, err
	}
	defer rawRows.Close()

	for rawRows.Next() {
		var row GetChangeCodeFileRow
		if err := rawRows.Scan(&row.ID, &row.Action, &row.RawData); err != nil {
			return &rows, err
		}
		rows = append(rows, row)
	}
	if err = rawRows.Err(); err != nil {
		return &rows, err
	}

	return &rows, nil
}
