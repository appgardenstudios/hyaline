package sqlite

import "database/sql"

type GetCurrentDocumentSectionRow struct {
	DocumentID string
	SectionID  string
	Title      string
	RawData    string
}

func GetCurrentDocumentSection(document string, section string, systemID string, db *sql.DB) (*GetCurrentDocumentSectionRow, error) {
	stmt, err := db.Prepare(`
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

	var row GetCurrentDocumentSectionRow
	err = stmt.QueryRow(systemID, document, section).Scan(&row.DocumentID, &row.Title, &row.RawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

type GetChangeDocumentRow struct {
	ID      string
	Action  string
	RawData string
}

func GetChangeDocument(document string, systemID string, db *sql.DB) (*GetChangeDocumentRow, error) {
	stmt, err := db.Prepare(`
select
	id,
	action,
	raw_data
from
	document
where
  system_id = ?
	AND id = ?
`)
	if err != nil {
		return nil, err
	}

	var row GetChangeDocumentRow
	err = stmt.QueryRow(systemID, document).Scan(&row.ID, &row.Action, &row.RawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

type GetChangeDocumentSectionRow struct {
	DocumentID string
	SectionID  string
	Title      string
	RawData    string
}

func GetChangeDocumentSection(document string, section string, systemID string, db *sql.DB) (*GetChangeDocumentSectionRow, error) {
	stmt, err := db.Prepare(`
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

	var row GetChangeDocumentSectionRow
	err = stmt.QueryRow(systemID, document, section).Scan(&row.DocumentID, &row.Title, &row.RawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}
