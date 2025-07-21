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

-- name: InsertDocumentTag :exec
INSERT INTO DOCUMENT_TAG (
  SOURCE_ID, DOCUMENT_ID, TAG_KEY, TAG_VALUE
) VALUES (
  ?, ?, ?, ?
);

-- name: InsertSection :exec
INSERT INTO SECTION (
  ID, DOCUMENT_ID, SOURCE_ID, PARENT_ID, PEER_ORDER, NAME, EXTRACTED_DATA
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: InsertSectionTag :exec
INSERT INTO SECTION_TAG (
  SOURCE_ID, DOCUMENT_ID, SECTION_ID, TAG_KEY, TAG_VALUE
) VALUES (
  ?, ?, ?, ?, ?
);