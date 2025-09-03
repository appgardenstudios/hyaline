---
title: "Getting Started with Hyaline"
description: "Quick guide to install Hyaline, create configuration, extract documentation, and set up an MCP integration."
purpose: Document how to get started with Hyaline
---
Welcome to Hyaline (pronounced "HIGH-uh-leen"), a documentation tool that helps software development teams keep their documentation current, accurate, and accessible. This guide will walk you through setting up Hyaline.

## Prerequisites
- One or more documentation sources (e.g. git repo, filesystem directory, website)
- A GitHub account OR a supported operating system (64-bit Linux, macOS, or Windows)
- A [supported LLM Provider](./reference/config.md)

> Please note that while you are free to install and use Hyaline for testing and evaluation purposes without a license, you still need to [obtain a license](/#pricing) to use Hyaline for any and all other purposes.

## GitHub App
The [Hyaline GitHub App](https://github.com/apps/hyaline-dev) is the recommended way of using Hyaline. It uses a configuration repository in your organization or personal account to extract, check, and audit documentation. To install the Hyaline GitHub App please see [How To: Install the GitHub App](./how-to/install-github-app.md).

Once installed you can visit other How To guides to learn how to [extract documentation](./how-to/extract-documentation.md), [check pull requests](./how-to/check-pull-request.md), or [run an MCP server](./how-to/run-mcp-server.md). Alternatively you can [learn more about how hyaline works](./explanation/hyaline.md) or visit the [configuration reference](./reference/config.md) to learn how to configure Hyaline.

## From Scratch (Advanced)
You can also use Hyaline from scratch to support any number of custom workflows or use cases. For an introduction on how to use Hyaline from scratch please see [How To: Use Hyaline From Scratch](./how-to/use-hyaline-from-scratch.md).
