---
title: "Reference: Hyaline CLI"
description: Command-line interface reference covering all commands, options, and usage examples
purpose: Detail each command/sub-command of the Hyaline CLI
sitemap:
  disable: true
---
## Overview
This documents the commandline options for the Hyaline Command Line Interface (CLI).

## Commands
The following commands and sub-commands are available within hyaline.

**Common Options**:
* `--debug` - (optional) Enables debug output

## help
`hyaline help` prints out usage information.

**Options**:
* (none)

**Example**:
```
$ hyaline help
```

## version
`hyaline version` prints out the currently installed version.

**Options**:
* (none)

**Example**:
```
$ hyaline version
```

## serve mcp
`hyaline serve mcp` starts an MCP server running locally over stdio and serves up the documentation produced by running `hyaline extract documentation`.

**Options**:
* `--documentation` - (required) Path to the SQLite database containing documentation

**Example**:
```
$ hyaline serve mcp --documentation ./documentation.db
```
Start a local MCP server using the standard I/O transport and have it use the extracted documentation found in `./documentation.db`.