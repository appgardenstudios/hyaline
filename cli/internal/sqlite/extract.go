package sqlite

import (
	"database/sql"
	"log/slog"
)

func CreateSchema(db *sql.DB) (err error) {
	slog.Debug("sqlite.CreateSchema schema creation started")
	_, err = db.Exec(`
CREATE TABLE SYSTEM(ID TEXT PRIMARY KEY);
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

func UpsertSystem(system System, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM
	(ID)
VALUES
	(?)
ON CONFLICT DO NOTHING
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

func DeleteCode(codeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  CODE
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(codeID, systemID)
	if err != nil {
		return err
	}

	return nil
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

func DeleteFile(codeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  FILE
WHERE
  CODE_ID = ?
	AND SYSTEM_ID = ?
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(codeID, systemID)
	if err != nil {
		return err
	}

	return nil
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

func DeleteDocumentation(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  DOCUMENTATION
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(documentationID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllDocumentation(systemID string, db *sql.DB) (arr []*Documentation, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TYPE, PATH
FROM
  DOCUMENTATION
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
		var row Documentation
		if err := rows.Scan(&row.ID, &row.SystemID, &row.Type, &row.Path); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
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

func DeleteDocument(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  DOCUMENT
WHERE
  DOCUMENTATION_ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(documentationID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllDocument(documentationID string, systemID string, db *sql.DB) (arr []*Document, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, RAW_DATA, EXTRACTED_DATA
FROM
  DOCUMENT
WHERE
  DOCUMENTATION_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query(documentationID, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row Document
		if err := rows.Scan(&row.ID, &row.DocumentationID, &row.SystemID, &row.Type, &row.Action, &row.RawData, &row.ExtractedData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
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

func DeleteSection(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SECTION
WHERE
  DOCUMENTATION_ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(documentationID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSection(documentationID string, systemID string, db *sql.DB) (arr []*Section, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA
FROM
  SECTION
WHERE
  DOCUMENTATION_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query(documentationID, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row Section
		if err := rows.Scan(&row.ID, &row.DocumentID, &row.DocumentationID, &row.SystemID, &row.Name, &row.ParentID, &row.PeerOrder, &row.ExtractedData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}

func GetAllSectionsForDocument(documentID string, documentationID string, systemID string, db *sql.DB) (arr []*Section, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA
FROM
  SECTION
WHERE
  DOCUMENT_ID = ?
  AND DOCUMENTATION_ID = ?
  AND SYSTEM_ID = ?
ORDER BY
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, PARENT_ID, PEER_ORDER
`)
	if err != nil {
		return
	}

	rows, err := stmt.Query(documentID, documentationID, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row Section
		if err := rows.Scan(&row.ID, &row.DocumentID, &row.DocumentationID, &row.SystemID, &row.Name, &row.ParentID, &row.PeerOrder, &row.ExtractedData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
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

func DeletePullRequest(pullRequestID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  PULL_REQUEST
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pullRequestID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllPullRequest(systemID string, db *sql.DB) (arr []*PullRequest, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TITLE, BODY
FROM
  PULL_REQUEST
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
		var row PullRequest
		if err := rows.Scan(&row.ID, &row.SystemID, &row.Title, &row.Body); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
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

func DeleteIssue(issueID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  ISSUE
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(issueID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllIssue(systemID string, db *sql.DB) (arr []*Issue, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TITLE, BODY
FROM
  PULL_REQUEST
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
		var row Issue
		if err := rows.Scan(&row.ID, &row.SystemID, &row.Title, &row.Body); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}
