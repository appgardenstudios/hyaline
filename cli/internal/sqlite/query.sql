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

-- name: GetAllSources :many
SELECT
  ID, DESCRIPTION, CRAWLER, ROOT
FROM
  SOURCE
ORDER BY
  ID;

-- name: GetDocumentsForSource :many
SELECT
  ID, SOURCE_ID, TYPE, PURPOSE, RAW_DATA, EXTRACTED_DATA
FROM
  DOCUMENT
WHERE
  SOURCE_ID = ?
ORDER BY
  ID;

-- name: GetSectionsForDocument :many
SELECT
  ID, DOCUMENT_ID, SOURCE_ID, PARENT_ID, PEER_ORDER, NAME, PURPOSE, EXTRACTED_DATA
FROM
  SECTION
WHERE
  SOURCE_ID = ? AND DOCUMENT_ID = ?
ORDER BY
  PEER_ORDER, ID;

-- name: GetDocumentTags :many
SELECT
  TAG_KEY, TAG_VALUE
FROM
  DOCUMENT_TAG
WHERE
  SOURCE_ID = ? AND DOCUMENT_ID = ?
ORDER BY
  TAG_KEY, TAG_VALUE;

-- name: GetSectionTags :many
SELECT
  TAG_KEY, TAG_VALUE
FROM
  SECTION_TAG
WHERE
  SOURCE_ID = ? AND DOCUMENT_ID = ? AND SECTION_ID = ?
ORDER BY
  TAG_KEY, TAG_VALUE;

-- name: DeleteSource :exec
DELETE FROM SOURCE WHERE ID = ?;

-- name: DeleteDocumentsForSource :exec
DELETE FROM DOCUMENT WHERE SOURCE_ID = ?;

-- name: DeleteDocumentTagsForSource :exec
DELETE FROM DOCUMENT_TAG WHERE SOURCE_ID = ?;

-- name: DeleteSectionsForSource :exec
DELETE FROM SECTION WHERE SOURCE_ID = ?;

-- name: DeleteSectionTagsForSource :exec
DELETE FROM SECTION_TAG WHERE SOURCE_ID = ?;