-- name: InsertSource :exec
INSERT INTO SOURCE (
  ID, DESCRIPTION, CRAWLER, ROOT
) VALUES (
  ?, ?, ?, ?
);

-- name: InsertDocument :exec
INSERT INTO DOCUMENT (
  ID, SOURCE_ID, TYPE, PURPOSE, RAW_DATA, EXTRACTED_DATA
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: GetDocumentIDsForSource :many
SELECT
  ID
FROM
  DOCUMENT
WHERE
  SOURCE_ID = ?;

-- name: UpdateDocumentPurpose :exec
UPDATE DOCUMENT 
SET PURPOSE = ?
WHERE ID = ? AND SOURCE_ID = ?;

-- name: UpsertDocumentTag :exec
INSERT INTO DOCUMENT_TAG (
  SOURCE_ID, DOCUMENT_ID, TAG_KEY, TAG_VALUE
) VALUES (
  ?, ?, ?, ?
) ON CONFLICT DO NOTHING;

-- name: InsertSection :exec
INSERT INTO SECTION (
  ID, DOCUMENT_ID, SOURCE_ID, PARENT_ID, PEER_ORDER, NAME, PURPOSE, EXTRACTED_DATA
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSectionIDsForSource :many
SELECT
  ID, DOCUMENT_ID
FROM
  SECTION
WHERE
  SOURCE_ID = ?;

-- name: UpdateSectionPurpose :exec
UPDATE SECTION 
SET PURPOSE = ?
WHERE ID = ? AND DOCUMENT_ID = ? AND SOURCE_ID = ?;

-- name: UpsertSectionTag :exec
INSERT INTO SECTION_TAG (
  SOURCE_ID, DOCUMENT_ID, SECTION_ID, TAG_KEY, TAG_VALUE
) VALUES (
  ?, ?, ?, ?, ?
) ON CONFLICT DO NOTHING;