package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

type GetCodeFileRow struct {
	ID      string
	RawData string
}

func GetCodeFile(files []string, systemID string, db *sql.DB) (*[]GetCodeFileRow, error) {
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
	var rows []GetCodeFileRow
	slog.Debug("sqlite.GetCodeFile executing query", "query", fmt.Sprintf("%v", stmt), "ids", ids)
	rawRows, err := stmt.Query(ids...)
	if err != nil {
		return nil, err
	}
	defer rawRows.Close()

	for rawRows.Next() {
		var row GetCodeFileRow
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
