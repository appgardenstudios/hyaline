---
title: "Reference: Hyaline CLI"
linkTitle: CLI
purpose: Detail each command/sub-command of the Hyaline CLI
url: documentation/reference/cli
---
## Overview
This documents the commandline options for the Hyaline Command Line Interface (CLI).

## Commands
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
`hyaline extract current` extracts current code and documentation for a system. Please visit the explanation documentation for [extract current](../04-explanation/02-extract-current.md) for more details.

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
`hyaline extract change` extracts changed code, documentation, and metadata for a system based on a change. Please visit the explanation documentation for [extract change](../04-explanation/03-extract-change.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--system` - (required) ID of the system to extract
* `--base` - (required if `--base-ref` is not set, mutually exclusive with `--base-ref`) Base branch (where changes will be applied). Tries to resolve to a local branch first, then a remote branch (if there is a single remote), and finally a tag.
* `--head` - (required if `head-ref` is not set, mutually exclusive with `--head-ref`) Head branch (which changes will be applied). Tries to resolve to a local branch first, then a remote branch (if there is a single remote), and finally a tag.
* `--base-ref` - (required if `--base` is not set, mutually exclusive with `--base`) Base reference (explicit commit hash or fully qualified reference). Passed directly to git resolution.
* `--head-ref` - (required if `--head` is not set, mutually exclusive with `--head`) Head reference (explicit commit hash or fully qualified reference). Passed directly to git resolution.
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

**Example** (using explicit references):
```
$ hyaline extract change --config ./hyaline.yml --system app --base-ref abc123def --head-ref refs/heads/feature-branch --output ./change.db
```
Extract code and documentation from the system `app` using explicit references. The base is specified as commit hash `abc123def` and head as the fully qualified reference `refs/heads/feature-branch`.

## check current
`hyaline check current` checks current code and documentation for a system. Please visit the explanation documentation for [check current](../04-explanation/04-check-current.md) for more details.

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
`hyaline check change` checks changed code and documentation for a system. Please visit the explanation documentation for [check change](../04-explanation/05-check-change.md) for more details.

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
`hyaline generate config` generates hyaline configuration for a current data set system. Please visit the explanation documentation for [generate config](../04-explanation/07-generate-config.md) for more details.

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
`hyaline merge` merge 2 or more data sets into a single data set. Please visit the explanation documentation for [merge](../04-explanation/08-merge.md) for more details.

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

## update pr
`hyaline update pr` updates a GitHub pull request by adding or updating a comment with the recommendations output by `hyaline check change`. Please visit the explanation documentation for [update pr](../04-explanation/06-update-pr.md) for more details.

**Options**:
* `--config` - (required) Path to the config file
* `--pull-request` - (required) GitHub Pull Request to update (OWNER/REPO/PR_NUMBER)
* `--comment` - (optional) GitHub Pull Request comment to update (OWNER/REPO/COMMENT_NUMBER)
* `--sha` - (optional) Git sha to add to the comment
* `--recommendations` - (required) Recommendations to use (output of check change)
* `--output` - (required) Path of the comment metadata (file must not already exist)

**Example**:
```
$ hyaline update pr --config ./hyaline.yml --pull-request appgardenstudios/hyaline-example/1 --recommendations ./results.json --output ./comment-metadata.json
```
Update the PR `appgardenstudios/hyaline-example/1` by adding a comment with the recommendations in `./results.json` and output the resulting comment metadata to `./comment-metadata.json`

**Example**:
```
$ hyaline update pr --config ./hyaline.yml --pull-request appgardenstudios/hyaline-example/1 --comment appgardenstudios/hyaline-example/2916854796 --sha b4c5c73 --recommendations ./results.json --output ./comment-metadata.json
```
Update the PR `appgardenstudios/hyaline-example/1` by updating the comment identified by `appgardenstudios/hyaline-example/2916854796` with the recommendations in `./results.json` (including the sha `b4c5c73`) and output the resulting comment metadata to `./comment-metadata.json`

## mcp stdio
`hyaline mcp stdio` starts an MCP server running locally over stdio and serves up the documentation produced by running `hyaline extract current`. Please visit the explanation documentation for [mcp](../04-explanation/09-mcp.md) for more details.

**Options**:
* `--current` - (required) Path to the current data set to serve (output of `hyaline extract current`)

**Example**:
```
$ hyaline mcp stdio --current ./current.db
```
Start a local MCP server using the standard I/O transport and have it use the extracted documentation found in `./current.db`.