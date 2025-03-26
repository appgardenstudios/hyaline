package sqlite

import (
	"database/sql"
	"log/slog"
)

func CreateSchema(db *sql.DB) (err error) {
	slog.Debug("sqlite.CreateSchema schema creation started")
	_, err = db.Exec(`
CREATE TABLE SYSTEM(ID);
CREATE TABLE CODE(ID, SYSTEM_ID, PATH);
CREATE TABLE FILE(ID, CODE_ID, SYSTEM_ID, ACTION, RAW_DATA);
CREATE TABLE DOCUMENTATION(ID, SYSTEM_ID, TYPE, PATH);
CREATE TABLE DOCUMENT(ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, RAW_DATA, EXTRACTED_DATA);
CREATE TABLE SECTION(ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, RAW_DATA, EXTRACTED_DATA);
`)

	slog.Debug("sqlite.CreateSchema schema creation complete")
	return
}

type System struct {
	ID string
}

func InsertSystem(system System, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM
	(ID)
VALUES
	(?)
`)
	if err != nil {
		return
	}
	stmt.Exec(system.ID)

	return
}

type Code struct {
	ID       string
	SystemID string
	Path     string
}

func InsertCode(code Code, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO CODE
	(ID, SYSTEM_ID, PATH)
VALUES
	(?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(code.ID, code.SystemID, code.Path)

	return
}

type File struct {
	ID       string
	CodeID   string
	SystemID string
	Action   string
	RawData  string
}

func InsertFile(file File, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO FILE
	(ID, CODE_ID, SYSTEM_ID, ACTION, RAW_DATA)
values
	(?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(file.ID, file.CodeID, file.SystemID, file.Action, file.RawData)

	return
}
