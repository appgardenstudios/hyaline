package sqlite

import "database/sql"

func GetSystems(db *sql.DB) (*[]System, error) {
	stmt, err := db.Prepare(`
SELECT
	ID
FROM
	SYSTEM
ORDER BY
	ID
	`)
	if err != nil {
		return nil, err
	}
	var rows []System
	rawRows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rawRows.Close()

	for rawRows.Next() {
		var row System
		if err := rawRows.Scan(&row.ID); err != nil {
			return &rows, err
		}
		rows = append(rows, row)
	}
	if err = rawRows.Err(); err != nil {
		return &rows, err
	}

	return &rows, nil
}
