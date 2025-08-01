CREATE TABLE SOURCE (
  ID TEXT NOT NULL,
  DESCRIPTION TEXT NOT NULL,
  CRAWLER TEXT NOT NULL,
  ROOT TEXT NOT NULL,
  PRIMARY KEY(ID)
);

CREATE TABLE DOCUMENT (
  ID TEXT NOT NULL,
  SOURCE_ID TEXT NOT NULL,
  TYPE TEXT NOT NULL,
  PURPOSE TEXT NOT NULL,
  RAW_DATA TEXT NOT NULL,
  EXTRACTED_DATA TEXT NOT NULL,
  PRIMARY KEY(ID, SOURCE_ID)
);

CREATE TABLE DOCUMENT_TAG (
  SOURCE_ID TEXT NOT NULL,
  DOCUMENT_ID TEXT NOT NULL,
  TAG_KEY TEXT NOT NULL,
  TAG_VALUE TEXT NOT NULL,
  PRIMARY KEY(SOURCE_ID, DOCUMENT_ID, TAG_KEY, TAG_VALUE)
);

CREATE INDEX DOCUMENT_TAG_KEY_VALUE ON DOCUMENT_TAG(TAG_KEY, TAG_VALUE);

CREATE TABLE SECTION (
  ID TEXT NOT NULL,
  DOCUMENT_ID TEXT NOT NULL,
  SOURCE_ID TEXT NOT NULL,
  PARENT_ID TEXT NOT NULL,
  PEER_ORDER NUM NOT NULL,
  NAME text NOT NULL,
  PURPOSE TEXT NOT NULL,
  EXTRACTED_DATA text NOT NULL,
  PRIMARY KEY(ID, DOCUMENT_ID, SOURCE_ID)
);

CREATE TABLE SECTION_TAG (
  SOURCE_ID TEXT NOT NULL,
  DOCUMENT_ID TEXT NOT NULL,
  SECTION_ID TEXT NOT NULL,
  TAG_KEY TEXT NOT NULL,
  TAG_VALUE TEXT NOT NULL,
  PRIMARY KEY(SOURCE_ID, DOCUMENT_ID, SECTION_ID, TAG_KEY, TAG_VALUE)
);

CREATE INDEX SECTION_TAG_KEY_VALUE ON SECTION_TAG(TAG_KEY, TAG_VALUE);