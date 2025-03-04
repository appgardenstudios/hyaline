package sqlite

import "database/sql"

func CreateCurrentSchema(db *sql.DB) (err error) {
	_, err = db.Exec(`
create table system(id);
create table code(id, system_id, path);
create table file(id, code_id, system_id, relative_path, raw_data);
create table documentation(id, system_id, type, path);
create table document(id, documentation_id, system_id, relative_path, format, raw_data, extracted_text);
create table section(id, document_id, documentation_id, system_id, parent_section_id, section_order, title, format, raw_data, extracted_text);
`)

	return
}

type CurrentSystem struct {
	ID string
}

func InsertCurrentSystem(row CurrentSystem, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`insert into system values(?)`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID)

	return
}

type CurrentCode struct {
	ID       string
	SystemID string
	Path     string
}

func InsertCurrentCode(row CurrentCode, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`insert into code values(?, ?, ?)`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.SystemID, row.Path)

	return
}

type CurrentFile struct {
	ID           string
	CodeID       string
	SystemID     string
	RelativePath string
	RawData      string
}

func InsertCurrentFile(row CurrentFile, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`insert into file values(?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.CodeID, row.SystemID, row.RelativePath, row.RawData)

	return
}
