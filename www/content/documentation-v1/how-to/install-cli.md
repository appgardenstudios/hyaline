---
title: "How To: Install the Hyaline CLI"
description: "How to download, install, and set up the Hyaline CLI on Linux, macOS, or Windows."
purpose: Document how to install the Hyaline CLI
sitemap:
  disable: true
---
## Purpose
Install the Hyaline CLI.

## Prerequisite(s)
* (none)

## Steps

### 1. Determine OS and Architecture
Before starting the installation process you need to determine your operating system and architecture. Hyaline supports 64-bit Linux (`linux`), MacOS (`darwin`), and Windows (`windows`) operating systems for either `amd64` or `arm64` architectures (`amd64` only for Windows).

### 2. Download Binary
You can download the appropriate binary from the [Release Page](https://github.com/appgardenstudios/hyaline/releases) on GitHub. Just select the release you would like to use and get the link to the appropriate binary from the assets section.

Alternatively you can use the following URL template: `https://github.com/appgardenstudios/hyaline/releases/download/{RELEASE}/hyaline-{OS}-{ARCH}.zip`.

### 3. Unzip and Make Hyaline Executable
Depending on your operating system you will need to do one or more of the following:

* Unzip the downloaded executable
* Make `hyaline` executable (if applicable)
* Add `hyaline` to your PATH (if desired)

## Next Steps
Visit [How To Run the CLI](./run-cli.md) or the [CLI Reference](../reference/cli.md).