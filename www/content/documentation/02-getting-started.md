---
title: "Getting Started with Hyaline"
linkTitle: Getting Started
purpose: Document how to get started with Hyaline
url: documentation/getting-started
---
Welcome to Hyaline (pronounced "HIGH-uh-leen"), a documentation tool that helps software development teams keep their documentation current, accurate, and accessible. This guide will walk you through setting up Hyaline and performing your first documentation extraction and check.

## What You'll Learn

In this guide, you'll learn how to:
- Install the Hyaline CLI
- Create a basic configuration file
- Extract your current documentation
- Set up the MCP server for AI integration

## Prerequisites

- **Operating System**: Linux (64-bit), macOS, or Windows (64-bit)
- **Architecture**: AMD64 or ARM64 (ARM64 not available for Windows)
- **LLM Access**: API key for a supported LLM provider (currently Anthropic Claude)
- **Documentation Source**: A documentation source (e.g. git repo, filesystem directory, website)

## Step 1: Install the Hyaline CLI

### Download the Latest Release

1. Visit the [Hyaline releases page](https://github.com/appgardenstudios/hyaline/releases)
2. Download the appropriate binary for your system:
   - **Linux AMD64**: `hyaline-linux-amd64.zip`
   - **Linux ARM64**: `hyaline-linux-arm64.zip`
   - **macOS AMD64**: `hyaline-darwin-amd64.zip`
   - **macOS ARM64**: `hyaline-darwin-arm64.zip`
   - **Windows AMD64**: `hyaline-windows-amd64.zip`

### Install the Binary

1. **Unzip the downloaded file**:
   ```bash
   unzip hyaline-*-*.zip
   ```

2. **Make the binary executable** (Linux/macOS only):
   ```bash
   chmod +x hyaline
   ```

3. **Move to your PATH** (optional but recommended):
   ```bash
   # Linux/macOS
   sudo mv hyaline /usr/local/bin/

   # Or add to your local bin directory
   mv hyaline ~/bin/
   ```

4. **Verify the installation**:
   ```bash
   hyaline version
   ```

For detailed installation instructions, see [How To: Install the CLI](./03-how-to/01-install-cli.md).

## Step 2: Create a Configuration File

Create a `hyaline.yml` file in your project root. Here's a basic configuration for a typical project that has a locally checked-out git repo and `main` branch:

```yaml
llm:
  provider: anthropic
  model: claude-3-5-sonnet-20241022
  key: ${ANTHROPIC_KEY}

github:
  token: ${GITHUB_TOKEN}

systems:
  - id: my-app
    documentation:
      - id: main-docs
        type: md
        extractor:
          type: git
          options:
            path: .
            branch: main
          include: 
            - "README.md"
            - "docs/**/*.md"
            - "**/*.md"
          exclude:
            - "node_modules/**/*"
            - ".git/**/*"
        documents:
          - name: README.md
            purpose: Provide an overview of the project, installation instructions, and basic usage examples
            required: true
          - name: docs/installation.md
            purpose: Detailed installation and setup instructions
            required: false
```

### Configuration Breakdown

- **llm**: Configuration for the AI provider
- **systems**: Define your projects/systems
- **documentation**: Specify where your docs are located
- **documents**: Define expected documentation structure and purposes

For complete configuration options and details, see the [Configuration Reference](./05-reference/01-config.md).

## Step 3: Set Up Environment Variables

Based on the configuration above, you'll need to set the environment variables referenced in your config file. For detailed instructions on setting up environment variables, see [How To: Run the CLI](./03-how-to/02-run-cli.md) and [Configuration Reference](./05-reference/01-config.md).

## Step 4: Extract Current Documentation

Now let's extract your current documentation into a data set:

```bash
hyaline extract current \
  --config ./hyaline.yml \
  --system my-app \
  --output ./current.db
```

This command will:
- Scan your repository for documentation files
- Extract and organize them by system
- Store everything in a SQLite database (`current.db`)

For more details on how extraction works, see [Explanation: Extract Current](./04-explanation/02-extract-current.md).

## Step 5: Set Up MCP Server

To make your documentation available to AI assistants like Claude:

1. **Extract current documentation** (if not already done):
   ```bash
   hyaline extract current --config ./hyaline.yml --system my-app --output ./current.db
   ```

2. **Start the MCP server**:
   ```bash
   hyaline mcp stdio --current ./current.db
   ```

3. **Configure your AI client** to use the MCP server. The exact steps depend on your client, but you'll typically need to add a server configuration that runs the command above.

For detailed MCP setup instructions, see [How To: Run the MCP Server](./03-how-to/04-run-mcp.md) and [MCP Reference](./05-reference/06-mcp.md).

## Next Steps
To set up GitHub Actions for automatic PR checks visit [How To: Check a PR](./03-how-to/03-check-pr.md), or see the [CLI Reference](./05-reference/02-cli.md) for complete command documentation.