---
title: "How To: Use Hyaline Solo (Advanced)"
description: "Configure and run Hyaline solo without using the GitHub App."
purpose: Document how to run Hyaline solo without using the GitHub App
---
## Purpose
Run Hyaline solo without using the GitHub App

## Prerequisite(s)
- [Go Toolchain Installed](https://go.dev/) (version 1.24+)

## Steps

### 1. Install CLI
You can either download a pre-build version from GitHub (**1.1**) or build the CLI yourself (**1.2**)

#### 1.1 Download CLI
Before starting the installation process you need to determine your operating system and architecture. Hyaline supports 64-bit Linux (`linux`), MacOS (`darwin`), and Windows (`windows`) operating systems for either `amd64` or `arm64` architectures (`amd64` only for Windows).

You can download the appropriate binary from the [Release Page](https://github.com/appgardenstudios/hyaline/releases) on GitHub. Just select the release you would like to use and get the link to the appropriate binary from the assets section.

Alternatively you can use the following URL template: `https://github.com/appgardenstudios/hyaline/releases/download/{RELEASE}/hyaline-{OS}-{ARCH}.zip`.

Depending on your operating system you will need to do one or more of the following:

* Unzip the downloaded executable
* Make `hyaline` executable (if applicable)
* Add `hyaline` to your PATH (if desired)

#### 1.2 Build CLI
Ensure that the [Hyaline Repository](https://github.com/appgardenstudios/hyaline) is cloned and checked out to the version you wish to build. Once the version is checked out you can build hyaline using the following command (from the root of the repo):

```bash
$ go build -o ./dist/hyaline -ldflags="-X 'main.Version=$VERSION'" ./cmd/hyaline.go
```

Note: Hyaline uses the pattern `v1YYYY-MM-DD-HASH` for release versions, and defaults to the version `unknown` when `main.Version` is not set via flag

You can specify the OS and architecture to use by setting the appropriate [GOOS/GOARCH environment variables](https://go.dev/doc/install/source#environment). For example, to build for the 64bit ARM version of MacOS:

```bash
$ GOOS=darwin GOARCH=arm64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
```

### 2. Run the CLI
First run `hyaline version` to ensure that the Hyaline CLI is installed and working properly.

Most commands require a hyaline configuration file. You will need to pass the path of the file in as the value of the `--config` commandline parameter. Either locate or create that configuration file now. For more information on the format of the configuration file please see the [Hyaline Config Reference](../reference/config.md).

Based on the contents of the configuration file above there may be one or more environment variables that need to be set before running Hyaline. You can find these in the config by looking for references like `${MY_VAR}`, which specifies that Hyaline should use the value of the environment variable `MY_VAR` for that key in the config.

Now you can execute the hyaline command, passing in all of the required and (optionally) optional parameters. Please see [CLI Reference](../reference/cli.md) for a full list of commands and their associated options.

### 3. Extract Documentation
TODO

### 3.1 Create a Configuration File

Create a `hyaline.yml` file in your project root. Here's basic configuration to extract documentation from a typical project that has a locally checked-out git repo and `main` branch:

```yaml
extract:
  source:
    id: my-app
    description: My application documentation
  crawler:
    type: git
    options:
      path: .
      branch: main
    include:
      - "README.md"
      - "docs/**/*.md"
      - "**/*.md"
    exclude:
      - "node_modules/**/*"
      - ".git/**/*"
  extractors:
    - type: md
      include:
        - "**/*.md"
  metadata:
    - document: README.md
      purpose: Provide an overview of the project, installation instructions, and basic usage examples
    - document: docs/installation.md
      purpose: Detailed installation and setup instructions
```

### Configuration Breakdown

- **extract**: Top-level key for extraction configuration
- **source**: Metadata about the documentation source (id and description)
- **crawler**: Specifies how to crawl the documentation (git repository in this case)
- **extractors**: List of extractors to process documentation files (markdown in this case)
- **metadata**: Define purposes for specific documents

Note that the above configuration is for a simple example, but Hyaline can be configured to extract documentation from multiple systems and sources (e.g. websites, remote git repositories). For complete configuration options and details, see the [Configuration Reference](../reference/config.md).

## Step 3.2: Extract Documentation

Now let's extract your documentation into a data set (note that you will need to run this from the root of the documentation source):

```bash
$ hyaline extract documentation \
  --config ./hyaline.yml \
  --output ./documentation.db
```

This command will:
- Scan your repository for documentation files based on the configuration
- Extract and process them according to the defined extractors
- Store everything in a SQLite database (`documentation.db`)

For more details on how extraction works, see [Explanation: Extract](../explanation/extract.md).

### 4. Check a PR
To check a GitHub pull request you need to check in your configuration file, setup the required secrets, and then create a GitHub Workflow file to run a GitHub Action to check a pull request.

#### 4.1 Check in Configuration
Create a Hyaline Configuration File with [extract](../reference/config.md#extract) and [check](../reference/config.md#check) configured, and check it into your GitHub repository.

Make sure you don't check in any secrets, like the LLM Provider Key or the GitHub Token. Instead, setup the configuration to pull them from the environment.

See [config reference](../reference/config.md) for more information on creating a configuration file and referencing secrets.

#### 4.2. Set Up Secrets
For each environment variable used in the Hyaline configuration that references a secret, set that variable up to be pulled in as a [secret in GitHub](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions).

#### 4.3. Create Workflow File
Create a [GitHub Workflow](https://docs.github.com/en/actions/writing-workflows/quickstart) file to run when a pull request is updated. You can see an example file in the [GitHub Actions reference](../reference/github-actions.md).

Alternatively you could set Hyaline up to be run [manually on dispatch](https://docs.github.com/en/actions/managing-workflow-runs-and-deployments/managing-workflow-runs/manually-running-a-workflow).

### 5. Run MCP Server
Ensure that you have extracted your documentation and placed the resulting data set in a convenient location.

Running the server will vary by client, but the gist of it is you want to have the client run the command `hyaline serve mcp --documentation ./path/to/documentation-data-set.db` to start a local MCP server listening over stdio.

Please see your client documentation for specific steps.

## Next Steps
Visit the [CLI Reference](../reference/cli.md).