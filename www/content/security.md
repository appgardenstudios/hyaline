---
title: Security
purpose: Provide an overview of the security practices and policies of Hyaline
---

# Security

Hyaline is designed with security and privacy as core principles.

## Deployment Models

Hyaline can be used in two ways:

1. **GitHub App**: The [Hyaline GitHub App](/documentation/explanation/github-app/) automates most of the setup and running of Hyaline.
2. **From Scratch**: The [from scratch](/documentation/how-to/use-hyaline-from-scratch/) model allows you to fully customize your usage of Hyaline.

Each model has different security considerations outlined below.

## GitHub App Model

### Data Privacy
<!-- purpose: Explain how Hyaline handles data privacy in the GitHub App Model -->
When using the hosted Hyaline GitHub App, our service receives GitHub webhook payloads for pull request events. These payloads contain code diffs and metadata as documented in [GitHub's webhook payload reference](https://docs.github.com/en/webhooks/webhook-events-and-payloads#pull_request). The GitHub App **ignores and does not process** the code content in these payloads. It only uses the webhook events to trigger workflows in your `hyaline-github-app-configuration` repository.

All workflows, configuration, and data processing happen in your own GitHub repository (`hyaline-github-app-config`). You maintain full control over your configuration and extracted documentation data.

You can optionally [host your own copy of the GitHub App](https://github.com/appgardenstudios/hyaline-github-app-config/tree/main/.github/apps/_hyaline) to prevent any data from being sent to our servers.

### Operation Environment
<!-- purpose: Explain the operations that Hyaline performs in the GitHub App Model -->
The GitHub App triggers workflows that run in your GitHub Actions environment:

- Workflows execute in GitHub's runner infrastructure (or your self-hosted runners)
- All processing happens within your GitHub organization/account
- Documentation data is stored as artifacts in your configuration repository

### Network Communications
<!-- purpose: Document the network calls that Hyaline makes in the GitHub App Model -->
Network requests are made by GitHub Actions workflows in your repository to:

- **GitHub API**: Using your provided Personal Access Tokens to read repository content and comment on pull requests
- **LLM API calls**: Using your provided API keys for documentation analysis and recommendations
- **Source repositories**: If configured to extract from remote repositories
- **Documentation websites**: If configured to crawl external documentation

### Authentication & Secrets
<!-- purpose: Explain the types of authentication Hyaline uses and how secrets should be managed in the GitHub App Model -->
Authentication is managed through GitHub repository secrets:

- **GitHub Personal Access Tokens**: For accessing repositories and commenting on pull requests
- **LLM API keys**: For documentation analysis and recommendations
- **Environment variable substitution**: Secrets are referenced in configuration files and substituted at runtime

See: [How To: Install the GitHub App](/documentation/how-to/install-github-app/) to learn more about the specific permissions the GitHub App requires.

All secrets are stored in your GitHub repository's secret management system and never exposed in logs or artifacts.

When working with remote repositories, Hyaline supports either ssh key or http authentication.

## From Scratch Model

### Data Privacy
<!-- purpose: Explain how Hyaline handles data privacy in the From Scratch Model -->
When using Hyaline from scratch, Hyaline runs entirely on your infrastructure, whether that's your local machine, CI environment, or your own servers. Your source code and documentation never leave your control. Hyaline does not send usage analytics, telemetry, or any other data to external servers. We have no visibility into how you use the tool or what content you're processing.

### Operation Environment
<!-- purpose: Explain the operations that Hyaline performs in the From Scratch Model -->
Hyaline operates as a standalone CLI tool that:

- Reads your source code and documentation from locations you specify
- Writes data to a SQLite database at a path you control

### Network Communications
<!-- purpose: Document the network calls that Hyaline makes in the From Scratch Model -->
Hyaline only makes network requests when you explicitly configure it to:

- **Source code repositories**: If you configure git extractors to clone remote repositories
- **Documentation sources**: If you configure HTTP extractors to crawl documentation websites
- **LLM API calls**: If you provide an API key for LLM services (like Anthropic's Claude)
- **GitHub API calls**: If you provide an API key for PR commenting

All connections are made directly from your environment to the configured services - Hyaline does not proxy or intercept this traffic.

### File System Access
<!-- purpose: Document Hyaline's file system access requirements in the From Scratch Model -->
Hyaline requires:

- **Read access** to your source code and documentation files as specified in your configuration
- **Write access** to create and update the SQLite database you specify

The tool respects standard file system permissions and only accesses files within the paths you configure.

### Authentication & Secrets
<!-- purpose: Explain the types of authentication Hyaline uses and how secrets should be managed in the GitHub App Model -->
When working with remote repositories or LLM APIs, Hyaline supports:

- Environment variable substitution for API keys and tokens
- SSH key authentication for git repositories
- HTTP authentication for private repositories

**Best practice**: Store all sensitive credentials as environment variables rather than hard-coding them in configuration files.

## MCP Server
<!-- purpose: Document the security aspects of the MCP server -->
The MCP server runs locally and provides your documentation to AI tools. Currently:

- No authentication or access controls are implemented
- The server only runs on your local machine
- Access is limited to applications you explicitly connect to the server

## Reporting Security Issues
<!-- purpose: Provide contact information for reporting security vulnerabilities -->
If you discover a security vulnerability in Hyaline, please report it to us at:

**Email**: [support@hyaline.dev](mailto:support@hyaline.dev)

We take security issues seriously and will respond promptly to any reports.