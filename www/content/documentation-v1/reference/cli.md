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