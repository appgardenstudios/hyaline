---
title: "Reference: MCP"
description: MCP server tools and prompts for making extracted documentation available to LLMs
purpose: Detail the functionality of Hyaline's MCP server
---
## Overview
Hyaline provides a built-in MCP server that can make the [current data set](./data-set.md) extracted by Hyaline's [extract documentation](../explanation/extract.md) command available to LLMs. For information on how to set up and run the MCP server please see the [how to](../how-to/run-mcp-server.md) or the [cli reference](./cli.md).

## Tools
<!-- purpose: Document all the available tools the Hyaline MCP server provides -->
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

### reload_documentation
Reload the documentation dataset. When running in GitHub Artifacts mode (with `--github-repo`), this downloads the latest artifact from the configured repository. When running in local filesystem mode (with `--documentation`), this reloads the documentation from the local database file.

**Arguments**
None.

**Output**
A success message indicating the documentation was reloaded successfully.

## Resources
<!-- purpose: Document all the resource features the Hyaline MCP server provides -->
Hyaline's MCP server exposes all extracted documents as MCP resources. Access to resources varies by client. For example, Claude Code and many other clients use an `@` symbol (e.g. `@hyaline:document://hyaline-example/docs/index.html`). Refer to your client documentation for specific usage guidance.

### List Resources
List all available documentation resources. Each resource represents a single document in the documentation data set.

**Output**
A list of all available resources, where each resource includes:
- `uri` - The document URI in the format `document://<source-id>/<document-id>`
- `name` - The name of the document, including its file path in the documentation source (e.g. `docs/index.html`)
- `description` - The document's purpose (if available)
- `mimeType` - Always `text/markdown` (extracted documentation is in markdown format)

### Read Resource
Read the contents of a specific documentation resource.

**Arguments**
- `uri` - The document URI in the format `document://<source-id>/<document-id>`

**Output**
The complete contents of the requested document.

## Prompts
<!-- purpose: Document all the available prompts the Hyaline MCP server provides -->
Hyaline's MCP server provides the following prompts:

### answer_question
Answer a question based on available documentation.

**Arguments**
- `question` - The question to answer