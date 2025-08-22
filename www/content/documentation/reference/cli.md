---
title: "Reference: CLI"
description: Command-line interface reference covering all commands, options, and usage examples
purpose: Detail each command/sub-command of the Hyaline CLI
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

## check diff
`hyaline check diff` checks a diff and outputs a list of recommended documentation updates.

**Options**:
* `--config` - (required) Path to the config file
* `--documentation` - (required) Path to the current documentation data set (output of `hyaline extract documentation`)
* `--path` - (optional) Path to the root of the repository to check. Defaults to `./`
* `--base` - (required if `--base-ref` is not set, mutually exclusive with `--base-ref`) Base branch (where changes will be applied). Tries to resolve to a local branch first, then a remote branch (if there is a single remote), and finally a tag
* `--base-ref` - (required if `--base` is not set, mutually exclusive with `--base`) Base reference (explicit commit hash or fully qualified reference). Passed directly to git resolution
* `--head` - (required if `head-ref` is not set, mutually exclusive with `--head-ref`) Head branch (which changes will be applied). Tries to resolve to a local branch first, then a remote branch (if there is a single remote), and finally a tag
* `--head-ref` - (required if `--head` is not set, mutually exclusive with `--head`) Head reference (explicit commit hash or fully qualified reference). Passed directly to git resolution
* `--pull-request` - (optional) GitHub Pull Request to include in the change (`<owner>/<repo>/<pr_number>`)
* `--issue` - (optional, multiple allowed) GitHub Issue to include in the change (`<owner>/<repo>/<issue_number>`). Accepts multiple issues by setting multiple times
* `--output` - (required) Path of the output file to create (file must not already exist)

**Example**:
```
$ hyaline check diff --config ./hyaline.yml --documentation ./documentation.db --path ./ --base main --head feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3 --output ./recommendations.json
```
Check what documentation in `./documentation.db` should be updated based on the changes between the `main` and `feat-1` branches as well as the configuration in `./hyaline.yml`. It takes into account the contents of the pull request `appgardenstudios/hyaline-example/1` and the issues `appgardenstudios/hyaline-example/2` and `appgardenstudios/hyaline-example/3`. The set of recommendations are output to `./recommendations.json`.

**Example**:
```
$ hyaline check diff --config ./hyaline.yml --documentation ./documentation.db --path ./ --base-ref refs/heads/main --head-ref refs/remotes/origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3 --output ./recommendations.json
```
Check what documentation in `./documentation.db` should be updated based on the changes between the `main` and `feat-1` refs as well as the configuration in `./hyaline.yml`. It takes into account the contents of the pull request `appgardenstudios/hyaline-example/1` and the issues `appgardenstudios/hyaline-example/2` and `appgardenstudios/hyaline-example/3`. The set of recommendations are output to `./recommendations.json`.

## check pr
`hyaline check pr` checks a pull request to see what documentation may need to be updated and adds any recommendations as a comment on the PR.

**Options**:
* `--config` - (required) Path to the config file
* `--documentation` - (required) Path to the current documentation data set
* `--pull-request` - (required) GitHub Pull Request to check (`<owner>/<repo>/<pr_number>`)
* `--issue` - (optional, multiple allowed) GitHub Issue to include in the change (`<owner>/<repo>/<issue_number>`). Accepts multiple issues by setting multiple times
* `--output` - (optional) Path to write the combined (current and previous merged together) recommendations to
* `--output-current` - (optional) Path to write the current recommendations to
* `--output-previous` - (optional) Path to write the previous recommendations to

**Example**:
```
$ hyaline check pr --config ./hyaline.yml --documentation ./documentation.db --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3 --output ./recommendations.md
```
Check what documentation in `./documentation.db` should be updated based on the changes in the pull request `appgardenstudios/hyaline-example/1` as well as the configuration in `./hyaline.yml`. It takes into account the content of the pull request `appgardenstudios/hyaline-example/1` and the issues `appgardenstudios/hyaline-example/2` and `appgardenstudios/hyaline-example/3`. If a comment already exists on the PR, the recommendations from the current run are merged with the recommendations from the previous run, and the comment is updated. Otherwise, a new comment is added with the current recommendations. The set of combined recommendations is output to `./recommendations.json`.

**Example**:
```
$ hyaline check pr --config ./hyaline.yml --documentation ./documentation.db --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3 --output-current ./current-recommendations.md
```
Check what documentation in `./documentation.db` should be updated based on the changes in the pull request `appgardenstudios/hyaline-example/1` as well as the configuration in `./hyaline.yml`. It takes into account the content of the pull request `appgardenstudios/hyaline-example/1` and the issues `appgardenstudios/hyaline-example/2` and `appgardenstudios/hyaline-example/3`. If a comment already exists on the PR, the recommendations from the current run are merged with the recommendations from the previous run, and the comment is updated. Otherwise, a new comment is added with the current recommendations. The set of recommendations from the current run is output to `./current-recommendations.json`

## audit documentation
`hyaline audit documentation` audits documentation against configurable rule checks to ensure compliance with documentation standards.

**Options**:
* `--config` - (required) Path to the config file
* `--documentation` - (required) Path to the documentation database (output of `hyaline extract documentation`)
* `--source` - (optional, multiple allowed) Only audit specific source ID(s). Can be specified multiple times
* `--output` - (required) Path to write the audit results JSON file (file must not already exist)

**Example**:
```
$ hyaline audit documentation --config ./hyaline.yml --documentation ./documentation.db --output ./audit-results.json
```
Audit all documentation in `./documentation.db` against the rules defined in `./hyaline.yml` and output the results to `./audit-results.json`.

**Example**:
```
$ hyaline audit documentation --config ./hyaline.yml --documentation ./documentation.db --source source1 --source source2 --output ./audit-results.json
```
Audit only specific sources (`source1` and `source2`) in `./documentation.db` against the rules defined in `./hyaline.yml` and output the results to `./audit-results.json`.

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

## export documentation
`hyaline export documentation` exports documentation from a documentation data set. Please see the explanation for [export](../explanation/export.md) for more details.

**Options**:
* `--documentation` - (required) Path to the current documentation data set (output of `hyaline extract documentation`).
* `--format` - (required) The format to export the documentation in. Must be one of `fs`, `llmsfulltxt`, `json`, or `sqlite`.
* `--include` - (optional, allows multiple) The documentation to include in the export, specified as a document uri in the form of `document://<source>/<path/to/document>(?tagValue=tagKey)`.
* `--exclude` - (optional, allows multiple) The documentation to exclude from the export, specified as a document uri in the form of `document://<source>/<path/to/document>(?tagValue=tagKey)`.
* `--output` - (required) The path to export the documentation to.

**Example**:
```
$ hyaline export documentation --documentation ./documentation.db --format json --output ./export.json
```
Export the documentation in `./documentation.db` in JSON format and output it to the file `./export.json`.

**Example**:
```
$ hyaline export documentation --documentation ./documentation.db --format json --output ./export.json --include 'document://frontend/**/*' --exclude 'document://*/**/*?type=customer'
```
Export the documentation in `./documentation.db` in JSON format and output it to the file `./export.json`. Only include documentation from the `frontend` source, but exclude any documentation with the tag `type=customer`.

## validate config
`hyaline validate config` validates the configuration and outputs the results of the validation. For more information on the validation output see [Config Validation](./config-validation.md).

**Options**:
* `--config` - (required) Path to the configuration file to validate.
* `--output` - (required) The file to export the results of the validation to.

**Example**:
```
$ hyaline validate config --config ./hyaline.yml --output ./output.json
```
Validate the config in `./hyaline.yml` and output it to the file `./output.json`.