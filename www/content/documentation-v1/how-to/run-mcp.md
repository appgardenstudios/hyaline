---
title: "How To: Run the Hyaline MCP Server"
description: "Configure and run the Hyaline MCP server to make documentation available to AI assistants."
purpose: Document how to run the Hyaline MCP Server
---
## Purpose
Run the hyaline mcp server locally.

## Prerequisite(s)
* [Install the CLI](./install-cli.md)

## Steps

### 1. Ensure CLI is Installed
Run `hyaline version` to ensure that the Hyaline CLI is installed and working properly.

### 2. Extract Documentation
Run `hyaline extract documentation` for each system you want to include. If you have multiple systems you can combine them into a single data set using `hyaline merge documentation`.

Place the resulting data set in a convenient location.

### 3. Add MCP Server to Client
This will vary by client, but the gist of it is you want to have the client run the command `hyaline serve mcp --documentation ./path/to/documentation-data-set.db` to start a local MCP server listening over stdio.

Please see your client documentation for specific steps.

## Next Steps
Visit the [CLI Reference](../reference/cli.md).