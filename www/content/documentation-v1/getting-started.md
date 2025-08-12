---
title: "Getting Started with Hyaline"
description: "Quick guide to install Hyaline, create configuration, extract documentation, and set up an MCP integration."
purpose: Document how to get started with Hyaline
sitemap:
  disable: true
---
Welcome to Hyaline (pronounced "HIGH-uh-leen"), a documentation tool that helps software development teams keep their documentation current, accurate, and accessible. This guide will walk you through setting up Hyaline and performing your first documentation extraction and check.

## What You'll Learn

In this guide, you'll learn how to:
- Install the Hyaline CLI
- Create a basic configuration file
- Extract your documentation
- Set up the MCP server for AI integration

## Prerequisites

- **Operating System**: Linux (64-bit), macOS, or Windows (64-bit)
- **Architecture**: AMD64 or ARM64 (ARM64 not available for Windows)
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
   $ unzip hyaline-*-*.zip
   ```

2. **Make the binary executable** (Linux/macOS only):
   ```bash
   $ chmod +x hyaline
   ```

3. **Move to your PATH** (optional but recommended):
   ```bash
   # Linux/macOS
   $ sudo mv hyaline /usr/local/bin/

   # Or add to your local bin directory
   $ mv hyaline ~/bin/
   ```

4. **Verify the installation**:
   ```bash
   $ hyaline version
   ```

For detailed installation instructions, see [How To: Install the CLI](./how-to/install-cli.md).

## Step 2: Create a Configuration File

Create a `hyaline.yml` file in your project root. Here's basic configuration to extract documentation from a typical project that has a locally checked-out git repo and `main` branch:

```yaml
extract:
  source:
    id: my-app
    description: My application documentation
  crawler:
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
  extractors:
    - type: md
      include:
        - "**/*.md"
  metadata:
    - document: README.md
      purpose: Provide an overview of the project, installation instructions, and basic usage examples
    - document: docs/installation.md
      purpose: Detailed installation and setup instructions
```

### Configuration Breakdown

- **extract**: Top-level key for extraction configuration
- **source**: Metadata about the documentation source (id and description)
- **crawler**: Specifies how to crawl the documentation (git repository in this case)
- **extractors**: List of extractors to process documentation files (markdown in this case)
- **metadata**: Define purposes for specific documents

Note that the above configuration is for a simple example, but Hyaline can be configured to extract documentation from multiple systems and sources (e.g. websites, remote git repositories). For complete configuration options and details, see the [Configuration Reference](./reference/config.md).

## Step 3: Extract Documentation

Now let's extract your documentation into a data set (note that you will need to run this from the root of the documentation source):

```bash
$ hyaline extract documentation \
  --config ./hyaline.yml \
  --output ./documentation.db
```

This command will:
- Scan your repository for documentation files based on the configuration
- Extract and process them according to the defined extractors
- Store everything in a SQLite database (`documentation.db`)

For more details on how extraction works, see [Explanation: Extract](./explanation/extract.md).

## Step 4: Set Up MCP Server

To make your documentation available to AI assistants like Claude:

1. **Extract documentation** (if not already done):
   ```bash
   hyaline extract documentation --config ./hyaline.yml --output ./documentation.db
   ```

2. **Start the MCP server**:
   ```bash
   hyaline serve mcp --documentation ./documentation.db
   ```

3. **Configure your AI client** to use the MCP server. The exact steps depend on your client, but you'll typically need to add a server configuration that runs the command above.

For detailed MCP setup instructions, see [How To: Run the MCP Server](./how-to/run-mcp.md) and [MCP Reference](./reference/mcp.md).

## Next Steps
To set up GitHub Actions for automatic PR checks visit [How To: Check a PR](./how-to/check-pr.md), or see the [CLI Reference](./reference/cli.md) for complete command documentation.