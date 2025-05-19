---
title: How To install the Hyaline CLI
purpose: Document how to install the Hyaline CLI
---
# Purpose
Install the Hyaline CLI

# Prerequisite(s)
* (none)

# Steps

## 1. Determine OS and Architecture
Before starting the installation process you need to determine your operating system and architecture. Hyaline supports 64-bit Linux (`linux`), MacOS (`darwin`), and Windows (`windows`) operating systems for either `amd64` or `arm64` architectures (`amd64` only for Windows).

## 2. Download Binary
You can download the appropriate binary from the [Release Page](https://github.com/appgardenstudios/hyaline/releases) on GitHub. Just select the release you would like to use and get the link to the binary from the assets version.

Alternatively you can use the following URL template: `https://github.com/appgardenstudios/hyaline/releases/download/{RELEASE}/hyaline-{OS}-{ARCH}` (windows has an `.exe` postfix).

## 3. Make Hyaline Executable
Depending on your operating system you will need to do one or more of the following:

* Rename the downloaded executable to `hyaline`
* Make `hyaline` executable
* Add `hyaline` to your PATH

# Next Steps
Visit [How To Run the CLI](./run-cli.md) or the [CLI Reference](../reference/cli.md).