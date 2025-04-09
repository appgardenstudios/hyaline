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
	_, err = stmt.Exec(system.ID)
	if err != nil {
		return
	}

	return
}

func GetAllSystem(db *sql.DB) (arr []*System, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID
FROM
  SYSTEM
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row System
		if err := rows.Scan(&row.ID); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

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
	_, err = stmt.Exec(code.ID, code.SystemID, code.Path)
	if err != nil {
		return
	}

	return
}

func GetCode(codeID string, systemID string, db *sql.DB) (*Code, error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, PATH
FROM
  CODE
WHERE
  ID = ? AND SYSTEM_ID = ?
LIMIT 1
`)
	if err != nil {
		return nil, err
	}

	var code Code
	row := stmt.QueryRow(codeID, systemID)
	err = row.Scan(&code.ID, &code.SystemID, &code.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &code, nil
}

func GetAllCode(systemID string, db *sql.DB) (arr []*Code, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, PATH
FROM
  CODE
WHERE
  SYSTEM_ID = ?
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query(systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row Code
		if err := rows.Scan(&row.ID, &row.SystemID, &row.Path); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

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
	_, err = stmt.Exec(file.ID, file.CodeID, file.SystemID, file.Action, file.RawData)
	if err != nil {
		return
	}

	return
}

func GetAllFiles(codeID string, systemID string, db *sql.DB) (arr []*File, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, CODE_ID, SYSTEM_ID, ACTION, RAW_DATA
FROM
  FILE
WHERE
  CODE_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query(codeID, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row File
		if err := rows.Scan(&row.ID, &row.CodeID, &row.SystemID, &row.Action, &row.RawData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

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
	_, err = stmt.Exec(doc.ID, doc.SystemID, doc.Type, doc.Path)
	if err != nil {
		return
	}

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
	_, err = stmt.Exec(doc.ID, doc.DocumentationID, doc.SystemID, doc.Type, doc.Action, doc.RawData, doc.ExtractedData)
	if err != nil {
		return
	}

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
	_, err = stmt.Exec(section.ID, section.DocumentID, section.DocumentationID, section.SystemID, section.Name, section.ParentID, section.PeerOrder, section.ExtractedData)
	if err != nil {
		return
	}

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
	_, err = stmt.Exec(pullRequest.ID, pullRequest.SystemID, pullRequest.Title, pullRequest.Body)
	if err != nil {
		return
	}

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
	_, err = stmt.Exec(issue.ID, issue.SystemID, issue.Title, issue.Body)
	if err != nil {
		return
	}

	return
}
