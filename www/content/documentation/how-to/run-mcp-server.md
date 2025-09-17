---
title: "How To: Run the Hyaline MCP Server"
description: "Configure and run the Hyaline MCP server to make documentation available to AI assistants."
purpose: Document how to run the Hyaline MCP Server
---
## Purpose
Run the hyaline mcp server locally.

## Prerequisite(s)
- [Install GitHub App](./install-github-app.md)
- [Install the CLI Locally](./install-cli-locally.md)
- Have at least some documentation that has been extracted

## Steps

### 1. Download Current Documentation
Go to the latest `_Merge` workflow run in your `hyaline-github-app-config` repo instance and download the artifact `_current-documentation`. Once downloaded extract the folder and note the location of the extracted `documentation.db` file for later use.

### 2. Add MCP Server to Client
This will vary by client, but the gist of it is you want to have the client run the command `hyaline serve mcp --documentation ./path/to/documentation.db` to start a local MCP server listening over stdio.

Please see your client documentation for specific steps.

## Next Steps
Read more about [Hyaline's MCP server](../explanation/mcp.md) or visit the [CLI reference](../reference/cli.md).