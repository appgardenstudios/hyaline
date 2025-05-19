---
title: Hyaline CLI
purpose: Detail each command/sub-command of the Hyaline CLI
---
# Overview
This documents the commandline options for the Hyaline Command Line Interface (CLI)

# Commands
The following commands and sub-commands are available within hyaline.

**Common Options**:
* `debug` - (optional) Enables debug output

**Example**:
```
$ hyaline --debug extract current --config ./hyaline.yml --system app --output ./current.db
```

# help
`hyaline help` prints out usage information to stdout.

**Options**:
* (none)

**Example**:
```
$ hyaline help
```

## version
`hyaline version` prints out the currently install version to stdout.

**Options**:
* (none)

**Example**:
```
$ hyaline version
```

## extract current
`hyaline extract current` extracts current code and documentation for a system.

**Options**:
* `config` - (required) Path to the config file
* `system` - (required) ID of the system to extract
* `output` - (required) Path of the data set to create (file must not already exist)

**Example**:
```
$ hyaline extract current --config ./hyaline.yml --system app --output ./current.db
```
Extract code and documentation from the system `app` using the config file found at `./hyaline.yml` and create a current dataset at `./current.db`.

## extract change
`hyaline extract change` extracts changed code, documentation, and metadata for a system based on a change.

**Options**:
* `config` - (required) Path to the config file
* `system` - (required) ID of the system to extract
* `base` - (required) Base branch (where changes will be applied)
* `head` - (required) Head branch (which changes will be applied)
* `code-id` - (optional, multiple allowed) ID of the code source(s) that will be extracted
* `documentation-id` - (optional, multiple allowed) ID of the documentation source(s) that will be extracted
* `pull-request` - (optional) GitHub Pull Request to include in the change (OWNER/REPO/PR_NUMBER)
* `issue` - (optional, multiple allowed) GitHub Issue to include in the change (OWNER/REPO/PR_NUMBER)
* `output` - (required) Path of the data set to create (file must not already exist)

**Example**:
```
$ hyaline extract change --config ./hyaline.yml --system app --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./change.db
```
Extract code and documentation from the system `app` using the config file found at `./hyaline.yml` and create a change dataset at `./change.db`. This change set will contain the code and documentation diffs between the `main` and `origin/feat-1` branches, as well as the pull request `appgardenstudios/hyaline-example/1` and issues `appgardenstudios/hyaline-example/2` and `appgardenstudios/hyaline-example/3`.

**Example**:
```
$ hyaline extract change --config ./hyaline.yml --system app --base main --head origin/feat-1 --code-id backend --documentation-id backend  --output ./change.db
```
Extract code and documentation from the system `app` using the config file found at `./hyaline.yml` and create a change dataset at `./change.db`. It will only extract changes for the code source `backend` and documentation source `backend`.

## check current

## check change

## generate config

## merge
