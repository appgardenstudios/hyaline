---
title: "Getting Started with Hyaline"
description: "Quick guide to install Hyaline, create configuration, extract documentation, and set up an MCP integration."
purpose: Document how to get started with Hyaline
---
Welcome to Hyaline (pronounced "HIGH-uh-leen"), a documentation tool that helps software development teams keep their documentation current, accurate, and accessible. This guide will walk you through setting up Hyaline.

## Prerequisites
- One or more documentation sources (e.g. git repo, filesystem directory, website)
- A GitHub account OR a supported operating system (64-bit Linux, macOS, or Windows)
- An LLM Provider (TODO)

> Please note that while you are free to install and use Hyaline for testing and evaluation purposes without a license, you still need to [obtain a license](/#pricing) to use Hyaline for any and all other purposes.

## Installation
Hyaline offers two means of installation depending on your needs and setup. The first is to install Hyaline as a GitHub App for a GitHub organization or personal account. This is the most straightforward installation method, and is recommended. The second is to install and use Hyaline directly within your CI/CD workflows or locally. This offers more control and customization, but comes at the cost of more manual configuration and management. Both options are supported, so it's up to you to decide.

If you would like to install the Hyaline GitHub App in your organization or personal account please see [How To: Install the GitHub App](./how-to/install-github-app.md).

If you would like to use Hyaline solo and set everything up directly please see [How To: Use Hyaline Solo (Advanced)](./how-to/use-hyaline-solo.md).

## Using Hyaline
Depending on the installation method you chose above, you can visit other How To guides to learn how to [extract documentation](./how-to/extract-documentation.md), [check pull requests](./how-to/check-pull-request.md), or [run an MCP server](./how-to/run-mcp-server.md). Alternatively you can [learn more about how hyaline works](./explanation/hyaline.md) or visit the [configuration reference](./reference/config.md) to learn how to configure Hyaline.
