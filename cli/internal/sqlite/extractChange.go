package sqlite

import (
	"database/sql"
	"log/slog"
)

func CreateChangeSchema(db *sql.DB) (err error) {
	slog.Debug("sqlite.CreateChangeSchema schema creation started")
	_, err = db.Exec(`
create table system(id);
create table code(id, system_id, path);
create table file(id, code_id, system_id, relative_path, action, raw_data);
create table documentation(id, system_id, type, path);
create table document(id, documentation_id, system_id, relative_path, format, action, raw_data, extracted_text);
create table section(id, document_id, documentation_id, system_id, parent_section_id, section_order, title, format, raw_data, extracted_text);
`)

	slog.Debug("sqlite.CreateChangeSchema schema creation complete")
	return
}

type ChangeSystem struct {
	ID string
}

func InsertChangeSystem(row ChangeSystem, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into system
	(id)
values
	(?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID)

	return
}

type ChangeCode struct {
	ID       string
	SystemID string
	Path     string
}

func InsertChangeCode(row ChangeCode, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into code
	(id, system_id, path)
values
	(?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.SystemID, row.Path)

	return
}

type ChangeFile struct {
	ID           string
	CodeID       string
	SystemID     string
	RelativePath string
	Action       string
	RawData      string
}

func InsertChangeFile(row ChangeFile, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into file
	(id, code_id, system_id, relative_path, action, raw_data)
values
	(?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.CodeID, row.SystemID, row.RelativePath, row.Action, row.RawData)

	return
}

type ChangeDocumentation struct {
	ID       string
	SystemID string
	Type     string
	Path     string
}

func InsertChangeDocumentation(row ChangeDocumentation, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into documentation
	(id, system_id, type, path)
values
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.SystemID, row.Type, row.Path)

	return
}

type ChangeDocument struct {
	ID              string
	DocumentationID string
	SystemID        string
	RelativePath    string
	Format          string
	Action          string
	RawData         string
	ExtractedText   string
}

func InsertChangeDocument(row ChangeDocument, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into document
	(id, documentation_id, system_id, relative_path, format, action, raw_data, extracted_text)
values
	(?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.DocumentationID, row.SystemID, row.RelativePath, row.Format, row.Action, row.RawData, row.ExtractedText)

	return
}

type ChangeSection struct {
	ID              string
	DocumentID      string
	DocumentationID string
	SystemID        string
	ParentSectionID string
	Order           int
	Title           string
	Format          string
	RawData         string
	ExtractedText   string
}

func InsertChangeSection(row ChangeSection, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
insert into section
	(id, document_id, documentation_id, system_id, parent_section_id, section_order, title, format, raw_data, extracted_text)
values
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(row.ID, row.DocumentID, row.DocumentationID, row.SystemID, row.ParentSectionID, row.Order, row.Title, row.Format, row.RawData, row.ExtractedText)

	return
}
