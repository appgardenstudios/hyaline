---
title: "How To: Run the Hyaline MCP Server"
linkTitle: Run the MCP Server
description: "Configure and run the Hyaline MCP server to make documentation available to AI assistants."
purpose: Document how to run the Hyaline MCP Server
url: documentation/how-to/run-mcp
---
## Purpose
Run the hyaline mcp server locally.

## Prerequisite(s)
* [Install the CLI](./01-install-cli.md)

## Steps

### 1. Ensure CLI is Installed
Run `hyaline version` to ensure that the Hyaline CLI is installed and working properly.

### 2. Extract Current Documentation
Run `hyaline extract current` for each system you want to include. If you have multiple systems you can combine them into a single data set using `hyaline merge`.

Place the resulting data set in a convenient location.

### 3. Add MCP Server to Client
This will vary by client, but the gist of it is you want to have the client run the command `hyaline mcp stdio --current ./path/to/current-data-set.db` to start a local MCP server listening over stdio.

Please see your client documentation for specific steps.

## Next Steps
Visit the [CLI Reference](../05-reference/02-cli.md) or [MCP reference](../05-reference/06-mcp.md).
