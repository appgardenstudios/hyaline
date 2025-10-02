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

The Hyaline MCP server can be run in two modes: **GitHub Artifacts mode** (recommended) or **Local Filesystem mode**.

### Option 1: GitHub Artifacts Mode (Recommended)

This mode automatically downloads the latest documentation from your `hyaline-github-app-config` repo instance, making it easier to keep your documentation up to date.

#### 1. Create a Personal Access Token
Create a GitHub Personal Access Token (classic) with access to read action artifacts from your `hyaline-github-app-config` repo instance.

#### 2. Add MCP Server to Client
Configure your MCP client to run the Hyaline Docker image with GitHub repository access. Here is example configuration for Claude Code:

```json
{
  "mcpServers": {
    "hyaline": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "ghcr.io/appgardenstudios/hyaline:latest",
        "serve",
        "mcp",
        "--github-repo",
        "<your_github_account>/hyaline-github-app-config",
        "--github-token",
        "<ghp_yourpersonalaccesstoken>"
      ]
    }
  }
}
```

Replace `<your_github_account>/hyaline-github-app-config` with your actual repository path an `<ghp_yourpersonalaccesstoken>` with your personal access token.

### Option 2: Local Filesystem Mode

This mode uses a locally downloaded documentation database file.

#### 1. Download Current Documentation
Go to the latest `_Merge` workflow run in your `hyaline-github-app-config` repo instance and download the artifact `_current-documentation`. Once downloaded extract the folder and note the location of the extracted `documentation.db` file for later use.

#### 2. Add MCP Server to Client
Configure your MCP client to run the Hyaline Docker image with the documentation mounted as a volume. Here is example configuration for Claude Code (substituting `<path-to-documentation.db>` with the path to the downloaded `documentation.db` on your local machine):

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

Please see your MCP client documentation for specific configuration steps.

## Next Steps
Read more about [Hyaline's MCP server](../explanation/mcp.md) or visit the [CLI reference](../reference/cli.md).