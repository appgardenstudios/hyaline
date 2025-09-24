---
title: "How To: Run the Hyaline MCP Server"
description: "Configure and run the Hyaline MCP server to make documentation available to AI assistants."
purpose: Document how to run the Hyaline MCP Server
---
## Purpose
Run the hyaline mcp server locally.

## Prerequisite(s)
- [Install GitHub App](./install-github-app.md)
- Docker (v28.4+)
- Have at least some documentation that has been extracted

## Steps

### 1. Download Current Documentation
Go to the latest `_Merge` workflow run in your `hyaline-github-app-config` repo instance and download the artifact `_current-documentation`. Once downloaded extract the folder and note the location of the extracted `documentation.db` file for later use.

### 2. Add MCP Server to Client
This will vary by client, but the gist is to configure the client to run the[ Hyaline Docker image](https://github.com/appgardenstudios/hyaline/pkgs/container/hyaline), mount the documentation as a volume, and start the MCP server. Here is example configuration for configuring the Hyaline MCP server for Claude Code (substituting `<path-to-documentation.db>` with the path to the downloaded `documentation.db` on your local machine):

```json
{
  "mcpServers": {
    "hyaline": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-v",
        "<path-to-documentation.db>:/home/appuser/documentation.db:ro",
        "ghcr.io/appgardenstudios/hyaline:latest",
        "serve",
        "mcp",
        "--documentation",
        "/home/appuser/documentation.db"
      ]
    }
  }
}
```

Please see your client documentation for specific steps.

## Next Steps
Read more about [Hyaline's MCP server](../explanation/mcp.md) or visit the [CLI reference](../reference/cli.md).