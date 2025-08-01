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

## extract documentation
`hyaline extract documentation` extracts documentation from a documentation source. Please see the explanation for [extract](../explanation/extract.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--output` - (required) Path of the data set to create (file must not already exist)

**Example**:
```
$ hyaline extract documentation --config ./hyaline.yml --output ./documentation.db
```
Extract documentation from the system defined in the config file found at `./hyaline.yml` and create a current documentation dataset at `./documentation.db`.

## merge documentation
`hyaline merge documentation` merges 2 or more documentation data sets into a single output database.

**Options**:
* `--input` - (required, multiple allowed) Path of the sqlite databases to merge. At least 2 inputs are required
* `--output` - (required) Path of the sqlite database to create

**Example**:
```
$ hyaline merge documentation --input ./docs1.db --input ./docs2.db --output ./merged.db
```
Merge `./docs1.db` and `./docs2.db` into a single output database `./merged.db`.

**Example**:
```
$ hyaline merge documentation --input ./docs1.db --input ./docs2.db --input ./docs3.db --output ./merged.db
```
Merge multiple documentation databases `./docs1.db`, `./docs2.db`, and `./docs3.db` into a single output database `./merged.db`.

## serve mcp
`hyaline serve mcp` starts an MCP server running locally over stdio and serves up the documentation produced by running `hyaline extract documentation`.

**Options**:
* `--documentation` - (required) Path to the SQLite database containing documentation

**Example**:
```
$ hyaline serve mcp --documentation ./documentation.db
```
Start a local MCP server using the standard I/O transport and have it use the extracted documentation found in `./documentation.db`.