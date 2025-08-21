---
title: Security
purpose: Provide an overview of the security practices and policies of Hyaline
---

# Security

Hyaline is designed with security and privacy as core principles. This page outlines our security practices and what you can expect when using Hyaline.

## Data Privacy

**We don't see your code or documentation.** Hyaline runs entirely on your infrastructure - whether that's your local machine, CI environment, or your own servers. Your source code and documentation never leave your control.

**No data collection.** Hyaline does not send usage analytics, telemetry, or any other data to our servers. We have no visibility into how you use the tool or what content you're processing.

## Local Operation

Hyaline operates as a standalone CLI tool that:

- Reads your source code and documentation from locations you specify
- Writes data to a SQLite database at a path you control
- Runs an MCP server locally for integration with AI tools

## Network Communications

Hyaline only makes network requests when you explicitly configure it to:

- **Source code repositories**: If you configure git extractors to clone remote repositories
- **Documentation sources**: If you configure HTTP extractors to crawl documentation websites
- **LLM API calls**: If you provide an API key for LLM services (like Anthropic's Claude)

These connections are made directly from your environment to the configured services - Hyaline does not proxy or intercept this traffic.

## File System Access

Hyaline requires:

- **Read access** to your source code and documentation files as specified in your configuration
- **Write access** to create and update the SQLite database you specify

The tool respects standard file system permissions and only accesses files within the paths you configure.

## Authentication & Secrets

When working with remote repositories or LLM APIs, Hyaline supports:

- Environment variable substitution for API keys and tokens
- SSH key authentication for git repositories
- HTTP authentication for private repositories

**Best practice**: Store all sensitive credentials as environment variables rather than hard-coding them in configuration files.

## MCP Server

The MCP server runs locally and provides your documentation to AI tools. Currently:

- No authentication or access controls are implemented
- The server only runs on your local machine
- Access is limited to applications you explicitly connect to the server

## Reporting Security Issues

If you discover a security vulnerability in Hyaline, please report it to us at:

**Email**: [support@hyaline.dev](mailto:support@hyaline.dev)

We take security issues seriously and will respond promptly to any reports.