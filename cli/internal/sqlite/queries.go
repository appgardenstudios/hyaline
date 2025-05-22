package sqlite

import (
	"database/sql"
	"log/slog"
)

func CreateSchema(db *sql.DB) (err error) {
	slog.Debug("sqlite.CreateSchema schema creation started")
	_, err = db.Exec(`
CREATE TABLE SYSTEM(ID TEXT PRIMARY KEY);
CREATE TABLE SYSTEM_CODE(ID, SYSTEM_ID, PATH);
CREATE TABLE SYSTEM_FILE(ID, CODE_ID, SYSTEM_ID, ACTION, ORIGINAL_ID, RAW_DATA);
CREATE TABLE SYSTEM_DOCUMENTATION(ID, SYSTEM_ID, TYPE, PATH);
CREATE TABLE SYSTEM_DOCUMENT(ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, ORIGINAL_ID, RAW_DATA, EXTRACTED_DATA);
CREATE TABLE SYSTEM_SECTION(ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA);
CREATE TABLE SYSTEM_CHANGE(ID, SYSTEM_ID, TYPE, TITLE, BODY);
CREATE TABLE SYSTEM_TASK(ID, SYSTEM_ID, TYPE, TITLE, BODY);
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

type SystemCode struct {
	ID       string
	SystemID string
	Path     string
}

func InsertSystemCode(code SystemCode, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_CODE
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

func DeleteSystemCode(codeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_CODE
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

func GetAllSystemCode(systemID string, db *sql.DB) (arr []*SystemCode, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, PATH
FROM
  SYSTEM_CODE
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
		var row SystemCode
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

type SystemFile struct {
	ID         string
	CodeID     string
	SystemID   string
	Action     Action
	OriginalID string
	RawData    string
}

func InsertSystemFile(file SystemFile, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_FILE
	(ID, CODE_ID, SYSTEM_ID, ACTION, ORIGINAL_ID, RAW_DATA)
VALUES
	(?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(file.ID, file.CodeID, file.SystemID, file.Action, file.OriginalID, file.RawData)
	if err != nil {
		return
	}

	return
}

func DeleteSystemFile(codeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_FILE
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

func GetSystemFile(fileID string, codeID string, systemID string, db *sql.DB) (*SystemFile, error) {
	stmt, err := db.Prepare(`
SELECT
  ID, CODE_ID, SYSTEM_ID, ACTION, ORIGINAL_ID, RAW_DATA
FROM
  SYSTEM_FILE
WHERE
  ID = ?
  AND CODE_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return nil, err
	}

	var row SystemFile
	err = stmt.QueryRow(fileID, codeID, systemID).Scan(&row.ID, &row.CodeID, &row.SystemID, &row.Action, &row.OriginalID, &row.RawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func GetAllSystemFiles(codeID string, systemID string, db *sql.DB) (arr []*SystemFile, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, CODE_ID, SYSTEM_ID, ACTION, ORIGINAL_ID, RAW_DATA
FROM
  SYSTEM_FILE
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
		var row SystemFile
		if err := rows.Scan(&row.ID, &row.CodeID, &row.SystemID, &row.Action, &row.OriginalID, &row.RawData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}

func GetAllSystemFilesForSystem(systemID string, db *sql.DB) (arr []*SystemFile, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, CODE_ID, SYSTEM_ID, ACTION, ORIGINAL_ID, RAW_DATA
FROM
  SYSTEM_FILE
WHERE
  SYSTEM_ID = ?
ORDER BY
  CODE_ID, ID
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
		var row SystemFile
		if err := rows.Scan(&row.ID, &row.CodeID, &row.SystemID, &row.Action, &row.OriginalID, &row.RawData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}

type SystemDocumentation struct {
	ID       string
	SystemID string
	Type     string
	Path     string
}

func InsertSystemDocumentation(doc SystemDocumentation, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_DOCUMENTATION
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

func DeleteSystemDocumentation(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_DOCUMENTATION
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

func GetAllSystemDocumentation(systemID string, db *sql.DB) (arr []*SystemDocumentation, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TYPE, PATH
FROM
  SYSTEM_DOCUMENTATION
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
		var row SystemDocumentation
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

type SystemDocument struct {
	ID              string
	DocumentationID string
	SystemID        string
	Type            string
	Action          Action
	OriginalID      string
	RawData         string
	ExtractedData   string
}

func InsertSystemDocument(doc SystemDocument, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_DOCUMENT
	(ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, ORIGINAL_ID, RAW_DATA, EXTRACTED_DATA)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(doc.ID, doc.DocumentationID, doc.SystemID, doc.Type, doc.Action, doc.OriginalID, doc.RawData, doc.ExtractedData)
	if err != nil {
		return
	}

	return
}

func DeleteSystemDocument(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_DOCUMENT
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

func GetSystemDocument(documentID string, documentationID string, systemID string, db *sql.DB) (*SystemDocument, error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, ORIGINAL_ID, RAW_DATA, EXTRACTED_DATA
FROM
  SYSTEM_DOCUMENT
WHERE
  ID = ?
	AND DOCUMENTATION_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return nil, err
	}

	var row SystemDocument
	err = stmt.QueryRow(documentID, documentationID, systemID).Scan(&row.ID, &row.DocumentationID, &row.SystemID, &row.Type, &row.Action, &row.OriginalID, &row.RawData, &row.ExtractedData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func GetAllSystemDocument(documentationID string, systemID string, db *sql.DB) (arr []*SystemDocument, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, ORIGINAL_ID, RAW_DATA, EXTRACTED_DATA
FROM
  SYSTEM_DOCUMENT
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
		var row SystemDocument
		if err := rows.Scan(&row.ID, &row.DocumentationID, &row.SystemID, &row.Type, &row.Action, &row.OriginalID, &row.RawData, &row.ExtractedData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}

func GetAllSystemDocumentsForSystem(systemID string, db *sql.DB) (arr []*SystemDocument, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENTATION_ID, SYSTEM_ID, TYPE, ACTION, ORIGINAL_ID, RAW_DATA, EXTRACTED_DATA
FROM
  SYSTEM_DOCUMENT
WHERE
  SYSTEM_ID = ?
ORDER BY
  DOCUMENTATION_ID, ID
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
		var row SystemDocument
		if err := rows.Scan(&row.ID, &row.DocumentationID, &row.SystemID, &row.Type, &row.Action, &row.OriginalID, &row.RawData, &row.ExtractedData); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}

type SystemSection struct {
	ID              string
	DocumentID      string
	DocumentationID string
	SystemID        string
	Name            string
	ParentID        string
	PeerOrder       int
	ExtractedData   string
}

func InsertSystemSection(section SystemSection, db *sql.DB) (err error) {
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

func DeleteSystemSection(documentationID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_SECTION
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

func GetSystemSection(sectionID string, documentID string, documentationID string, systemID string, db *sql.DB) (*SystemSection, error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA
FROM
  SYSTEM_SECTION
WHERE
  ID = ?
  AND DOCUMENT_ID = ?
	AND DOCUMENTATION_ID = ?
  AND SYSTEM_ID = ?
`)
	if err != nil {
		return nil, err
	}

	var row SystemSection
	err = stmt.QueryRow(sectionID, documentID, documentationID, systemID).Scan(&row.ID, &row.DocumentID, &row.DocumentationID, &row.SystemID, &row.Name, &row.ParentID, &row.PeerOrder, &row.ExtractedData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func GetAllSystemSection(documentationID string, systemID string, db *sql.DB) (arr []*SystemSection, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA
FROM
  SYSTEM_SECTION
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
		var row SystemSection
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

// NOTE that this function MUST return sections in PEER_ORDER as that guarantee is used by the caller
func GetAllSystemSectionsForDocument(documentID string, documentationID string, systemID string, db *sql.DB) (arr []*SystemSection, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, DOCUMENT_ID, DOCUMENTATION_ID, SYSTEM_ID, NAME, PARENT_ID, PEER_ORDER, EXTRACTED_DATA
FROM
  SYSTEM_SECTION
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
		var row SystemSection
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

type SystemChange struct {
	ID       string
	SystemID string
	Type     string
	Title    string
	Body     string
}

func InsertSystemChange(change SystemChange, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_CHANGE
  (ID, SYSTEM_ID, TITLE, BODY)
VALUES
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(change.ID, change.SystemID, change.Type, change.Title, change.Body)
	if err != nil {
		return
	}

	return
}

func DeleteSystemChange(changeID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_CHANGE
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(changeID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSystemChange(systemID string, db *sql.DB) (arr []*SystemChange, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TITLE, BODY
FROM
  SYSTEM_CHANGE
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
		var row SystemChange
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

type SystemTask struct {
	ID       string
	SystemID string
	Type     string
	Title    string
	Body     string
}

func InsertSystemTask(task SystemTask, db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
INSERT INTO SYSTEM_TASK
  (ID, SYSTEM_ID, TYPE, TITLE, BODY)
VALUES
	(?, ?, ?, ?)
`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(task.ID, task.SystemID, task.Type, task.Title, task.Body)
	if err != nil {
		return
	}

	return
}

func DeleteSystemTask(taskID string, systemID string, db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM
  SYSTEM_TASK
WHERE
  ID = ?
	AND SYSTEM_ID = ? 
`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(taskID, systemID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSystemTask(systemID string, db *sql.DB) (arr []*SystemTask, err error) {
	stmt, err := db.Prepare(`
SELECT
  ID, SYSTEM_ID, TYPE, TITLE, BODY
FROM
  SYSTEM_TASK
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
		var row SystemTask
		if err := rows.Scan(&row.ID, &row.SystemID, &row.Type, &row.Title, &row.Body); err != nil {
			return arr, err
		}
		arr = append(arr, &row)
	}
	if err = rows.Err(); err != nil {
		return arr, err
	}

	return
}
