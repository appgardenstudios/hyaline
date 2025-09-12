---
title: "Getting Started with Hyaline"
description: "Quick guide to install Hyaline, create configuration, extract documentation, and set up an MCP integration."
purpose: Document how to get started with Hyaline
---
Welcome to Hyaline (pronounced "HIGH-uh-leen"), a documentation tool that helps software development teams keep their documentation current, accurate, and accessible. This guide will walk you through setting up Hyaline.

## Prerequisites
- One or more documentation sources (e.g. git repo, filesystem directory, website)
- A GitHub account OR a supported operating system (64-bit Linux, macOS, or Windows)
- A [supported LLM Provider](./llms/)

> Please note that you are free to install and use Hyaline without a license for testing and evaluation purposes. You will need to [obtain a license](/#pricing) to use Hyaline for any other purpose.

## Try it Out
If you just want to try out Hyaline, please see [How To Try Out Hyaline](./how-to/try-out-hyaline.md). This will guide you through the steps needed to try out Hyaline in a single repository using the free tier of the GitHub Models LLM.

Note that this is not the recommended installation method, as it will only check documentation in a single repository and will not create or use a centralized set of documentation to keep up-to-date. To unlock the full utility of Hyaline please [Install the GitHub App](./how-to/install-github-app.md).

## Install the GitHub App
The [Hyaline GitHub App](https://github.com/apps/hyaline-dev) is the recommended way to use Hyaline. It uses a configuration repository in your organization or personal account to extract, check, and audit documentation. To install the Hyaline GitHub App please see [How To: Install the GitHub App](./how-to/install-github-app.md).

Once installed you can visit other How To guides to learn how to [extract documentation](./how-to/extract-documentation.md), [check pull requests](./how-to/check-pull-request.md), or [run an MCP server](./how-to/run-mcp-server.md). Alternatively you can [learn more about how hyaline works](./explanation/hyaline.md) or visit the [configuration reference](./reference/config.md) to learn how to configure Hyaline.

Note that the Hyaline GitHub App will trigger workflows in the configuration repository located in your organization or personal account, meaning that you stay in control of your configuration and data. If you wish you can also run your own copy of the Hyaline GitHub App (located in the [configuration repository](https://github.com/appgardenstudios/hyaline-github-app-config)) to prevent any data whatsoever from leaving your organization and being sent to us.

## Use Hyaline from Scratch (Advanced)
You can also use Hyaline from scratch to support any number of custom workflows or use cases. For an introduction on how to use Hyaline from scratch please visit [How To: Use Hyaline from Scratch](./how-to/use-hyaline-from-scratch.md).
