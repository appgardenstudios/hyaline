---
title: Hyaline CLI
purpose: Detail each command/sub-command of the Hyaline CLI
---
# Overview
This documents the commandline options for the Hyaline Command Line Interface (CLI).

# Commands
The following commands and sub-commands are available within hyaline.

**Common Options**:
* `--debug` - (optional) Enables debug output

**Example**:
```
$ hyaline --debug extract current --config ./hyaline.yml --system app --output ./current.db
```

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

## extract current
`hyaline extract current` extracts current code and documentation for a system. Please visit the explanation documentation for [extract current](../explanation/extract-current.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--system` - (required) ID of the system to extract
* `--output` - (required) Path of the data set to create (file must not already exist)

**Example**:
```
$ hyaline extract current --config ./hyaline.yml --system app --output ./current.db
```
Extract code and documentation from the system `app` using the config file found at `./hyaline.yml` and create a current dataset at `./current.db`.

## extract change
`hyaline extract change` extracts changed code, documentation, and metadata for a system based on a change. Please visit the explanation documentation for [extract change](../explanation/extract-change.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--system` - (required) ID of the system to extract
* `--base` - (required) Base branch (where changes will be applied)
* `--head` - (required) Head branch (which changes will be applied)
* `--code-id` - (optional, multiple allowed) ID of the code source(s) that will be extracted
* `--documentation-id` - (optional, multiple allowed) ID of the documentation source(s) that will be extracted
* `--pull-request` - (optional) GitHub Pull Request to include in the change (OWNER/REPO/PR_NUMBER)
* `--issue` - (optional, multiple allowed) GitHub Issue to include in the change (OWNER/REPO/PR_NUMBER)
* `--output` - (required) Path of the data set to create (file must not already exist)

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
`hyaline check current` checks current code and documentation for a system. Please visit the explanation documentation for [check current](../explanation/check-current.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--current` - (required) Path to the current data set to check (output of `hyaline extract current`)
* `--system` - (required) ID of the system to extract
* `--output` - (required) Path of the output file to create (file must not already exist)
* `--check-purpose` - (optional, boolean) Call the llm to check that the purpose of each document/section matches the content
* `--check-completeness` - (optional, boolean) Call the llm to check that each document/section is complete

**Example**:
```
$ hyaline check current --config ./hyaline.yml --current ./current.db --system app --output ./results.json
```
Check the documentation for the system `app` from `./current.db` using the config file found at `./hyaline.yml` and writing the results to `./results.json`.

**Example**:
```
$ hyaline check current --config ./hyaline.yml --current ./current.db --system app --output ./results.json --check-purpose --check-completeness
```
Check the documentation for the system `app` from `./current.db` using the config file found at `./hyaline.yml` and writing the results to `./results.json`. Also ensure documentation is complete and that it matches its stated purpose while checking.

## check change
`hyaline check change` checks changed code and documentation for a system. Please visit the explanation documentation for [check change](../explanation/check-change.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--current` - (required) Path to the current data set to check (output of `hyaline extract current`)
* `--change` - (required) Path to the change data set to check (output of `hyaline extract change`)
* `--system` - (required) ID of the system to extract
* `--output` - (required) Path of the output file to create (file must not already exist)
* `--suggest` - (optional, boolean) Call the llm to suggest what edits to make to the documentation for each recommended update

**Example**:
```
$ hyaline check change --config ./hyaline.yml --current ./current.db --change ./change.db --system app --output ./results.json
```
Check which documentation should be updated for the system `app` using the data sets `./current.db` and `./change.db`, using the config file found at `./hyaline.yml`, and writing the results to `./results.json`.

**Example**:
```
$ hyaline check change --config ./hyaline.yml --current ./current.db --change ./change.db --system app --output ./results.json --suggest
```
Check which documentation should be updated for the system `app` using the data sets `./current.db` and `./change.db`, using the config file found at `./hyaline.yml`, and writing the results to `./results.json`. Also generate suggested edits for the documentation that should be updated.

## generate config
`hyaline generate config` generates hyaline configuration for a current data set system. Please visit the explanation documentation for [generate config](../explanation/generate-config.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--current` - (required) Path to the current data set to check (output of `hyaline extract current`)
* `--system` - (required) ID of the system to extract
* `--output` - (required) Path of the output file to create (file must not already exist)
* `--include-purpose` - (optional, boolean) Call the llm to generate the document/section purpose

**Example**:
```
$ hyaline generate config --config ./hyaline.yml --current ./current.db --system app --output ./config-additions.yml
```
Generate a configuration for the system `app` from `./current.db` using the config file found at `./hyaline.yml` and writing the suggested updates to `./config-additions.yml`.

**Example**:
```
$ hyaline generate config --config ./hyaline.yml --current ./current.db --system app --output ./config-additions.yml --include-purpose
```
Generate a configuration for the system `app` from `./current.db` using the config file found at `./hyaline.yml` and writing the suggested updates to `./config-additions.yml`. Also include the purpose of each document/section.

## merge
`hyaline merge` merge 2 or more data sets into a single data set. Please visit the explanation documentation for [merge](../explanation/merge.md) for more details.

**Options**:
* `--input` - (required, multiple allowed) Path to a data set
* `--output` - (required) Path of the merged data set (file must not already exist)

**Example**:
```
$ hyaline merge --input ./current.db --input ./new.db --output ./combined.db
```
Merge `./new.db` into `./current.db` and output the result to `./combined.db`

**Example**:
```
$ hyaline merge --input ./current.db --input ./new1.db --input ./new2.db --output ./combined.db
```
Merge `./new1.db` into `./current.db` followed by `./new2.db` and output the result to `./combined.db`