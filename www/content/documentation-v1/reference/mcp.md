---
title: "Reference: MCP"
description: MCP server tools and prompts for making extracted documentation available to LLMs
purpose: Detail the functionality of Hyaline's MCP server
---
## Overview
Hyaline provides a built-in MCP server that can make the [current data set](./data-set.md) extracted by Hyaline's [extract documentation](../explanation/extract-documentation.md) command available to LLMs. For information on how to set up and run the MCP server please see the [how to](../how-to/run-mcp.md) or the [cli reference](./cli.md).

## Tools
Hyaline's MCP server provides the following tools:

### list_documents
List all documents at or under the specified URI path, or all documents if no URI is provided. URIs follow this pattern: `document://<source-id>/<document-id>[?<key>=<value>]`.

**Arguments**
- `document_uri` - The URI path to list documents from (can be partial). Format: `document://<source-id>/<document-id>[?<key>=<value>]`. Query parameters filter results by tags (multiple values fo the same key are comma-separated). If not provided, lists all documents.

**Output**
A list of the documents available for the given URI. If a full URI is not given, documents scoped to the prefix are returned.

### get_documents
Get the contents of documents matching the specified URI, or all documents if no URI provided. Document URIs follow this pattern: `document://<source-id>/<document-id>[?<key>=<value>]`.

**Arguments**
- `document_uri` - The URI specifying which documents to retrieve (can be partial). Format: `document://<source-id>/<document-id>[?<key>=<value>]`. Query parameters filter results by tags. If not provided, retrieves all documents.

**Output**
One or more documents (including the contents of each document).

## Prompts
Hyaline's MCP server provides the following prompts:

### answer_question
Answer a question based on available documentation.

**Arguments**
- `question` - The question to answer