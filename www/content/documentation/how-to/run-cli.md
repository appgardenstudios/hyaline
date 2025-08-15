---
title: "How To: Run the Hyaline CLI"
description: "Run the Hyaline CLI."
purpose: Document how to run the Hyaline CLI
---
## Purpose
Run the hyaline cli locally or on a remote machine.

## Prerequisite(s)
* [Install the CLI](./install-cli.md)

## Steps

### 1. Ensure CLI is Installed
Run `hyaline version` to ensure that the Hyaline CLI is installed and working properly.

### 2. Get Path to Config
Most commands require a hyaline configuration file. You will need to pass the path of the file in as the value of the `--config` commandline parameter. Either locate or create that configuration file now. For more information on the format of the configuration file please see the [Hyaline Config Reference](../reference/config.md).

### 3. Export Env Vars (optional)
Based on the contents of the configuration file above there may be one or more environment variables that need to be set before running Hyaline. You can find these in the config by looking for references like `${MY_VAR}`, which specifies that Hyaline should use the value of the environment variable `MY_VAR` for that key in the config.

### 4. Execute the Command
Now you can execute the hyaline command, passing in all of the required and (optionally) optional parameters. Please see [CLI Reference](../reference/cli.md) for a full list of commands and their associated options.

## Next Steps
Visit [CLI Reference](../reference/cli.md) or [How to build the CLI](./build-cli.md).