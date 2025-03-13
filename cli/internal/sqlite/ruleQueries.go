package sqlite

import "database/sql"

type GetDocumentSectionRow struct {
	DocumentID string
	SectionID  string
	Title      string
	RawData    string
}

func GetDocumentSection(document string, section string, systemID string, currentDB *sql.DB, changeDB *sql.DB) (*GetDocumentSectionRow, error) {
	// Check for the section in current
	// TODO

	stmt, err := currentDB.Prepare(`
select
	document_id,
	title,
	raw_data
from
	section
where
  system_id = ?
	AND document_id = ?
	AND title = ?
`)
	if err != nil {
		return nil, err
	}

	var row GetDocumentSectionRow
	err = stmt.QueryRow(systemID, document, section).Scan(&row.DocumentID, &row.Title, &row.RawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}
