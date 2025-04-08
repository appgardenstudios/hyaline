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
CREATE TABLE SECTION(ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA);
CREATE TABLE PULL_REQUEST(ID, SYSTEM_ID, TITLE, BODY);
CREATE TABLE ISSUE(ID, SYSTEM_ID, TITLE, BODY);
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
VALUES
	(?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(file.ID, file.CodeID, file.SystemID, file.Action, file.RawData)

	return
}

type Documentation struct {
	ID       string
	SystemID string
	Type     string
	Path     string
}

func InsertDocumentation(doc Documentation, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO DOCUMENTATION
	(ID, SYSTEM_ID, TYPE, PATH)
VALUES
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(doc.ID, doc.SystemID, doc.Type, doc.Path)

	return
}

type Document struct {
	ID              string
	DocumentationID string
	SystemID        string
	Type            string
	Action          string
	RawData         string
	ExtractedData   string
}

func InsertDocument(doc Document, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO DOCUMENT
	(ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, RAW_DATA, EXTRACTED_DATA)
VALUES
	(?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(doc.ID, doc.DocumentationID, doc.SystemID, doc.Type, doc.Action, doc.RawData, doc.ExtractedData)

	return
}

type Section struct {
	ID              string
	DocumentID      string
	DocumentationID string
	SystemID        string
	Name            string
	ParentID        string
	PeerOrder       int
	ExtractedData   string
}

func InsertSection(section Section, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SECTION
	(ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(section.ID, section.DocumentID, section.DocumentationID, section.SystemID, section.Name, section.ParentID, section.PeerOrder, section.ExtractedData)

	return
}

type PullRequest struct {
	ID       string
	SystemID string
	Title    string
	Body     string
}

func InsertPullRequest(pullRequest PullRequest, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO PULL_REQUEST
  (ID, SYSTEM_ID, TITLE, BODY)
VALUES
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(pullRequest.ID, pullRequest.SystemID, pullRequest.Title, pullRequest.Body)

	return
}

type Issue struct {
	ID       string
	SystemID string
	Title    string
	Body     string
}

func InsertIssue(issue Issue, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO ISSUE
  (ID, SYSTEM_ID, TITLE, BODY)
VALUES
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	stmt.Exec(issue.ID, issue.SystemID, issue.Title, issue.Body)

	return
}
