---
title: Hyaline config
purpose: Document the configuration options for Hyaline
---
# Overview
TODO

# LLM

# GitHub

# Systems
Stores the configuration for each system known to Hyaline. Defined as a list of system objects.

```yaml
systems:
  - id: system1
    ...
  - id: system2
```

**systems**: A logical grouping of code and documentation that is addressed as a unit. You may define any number of systems in Hyaline.

## System
A system within Hyaline.

```yaml
systems:
  - id: system1
    code:
      - id: codeSource1
        ...
      - id: codeSource2
        ...
    documentation:
      - id: documentationSource1
        ...
      - id: documentationSource2
        ...
```

**id**: Each system has an ID that is used to reference that system, and the ID must be unique across all systems referenced in the documentation. Since IDs are frequently referenced as parameters on the command line, you should only use non-control characters for the ID (preferably `[a-zA-Z0-9_-]+`)

**code**: A list of Code Sources associated with this system. In essence, the set of system code.

**documentation**: A list of Documentation Sources associated with this system. In essence, the set of system documentation.

### Code Source
A set of source code for the system, or system source code.

```yaml
systems:
  - id: system1
    code:
      - id: codeSource1
        extractor: {...}
```

**id**: Each Code Source has an ID used to reference it. The ID must be unique across all Code Sources for a system (note that you can use the same ID for Code Sources in different systems if desired).

**extractor**: Each Code Source has an extractor that is responsible for extracting the source code out and into the data set.

#### Code Source Extractor
An extractor that specifies how code is extracted for this Code Source and placed into a data set.

```yaml
systems:
  - id: system1
    code:
      - id: codeSource1
        extractor:
          type: fs | git
          options: {...} # Dependent on the extractor type
          include: [glob]
          exclude: [glob]
```

**type**: The type of the extractor. For Code Sources there are two extractor types available: `fs` and `git`. For more information see extractor details below.

**options**: The options for the extractor. Note that these are specific to the type of extractor. Please see below for the options available for each extractor.

**include**: The set of globs to include in the set of source code during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

**exclude**: The set of globs to exclude from the set of source code during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

##### Code Source Extractor Options (fs)
Extract source code from a file system path. Note that code sources using this extractor will not be eligible to be included in a change data set. If you have a local git repository use the git extractor with the `path` option.

Note that Include and Exclude globs are relative to the path specified.

Please see the explanation of [Extract Current](../explanation/extract-current.md) for more information.

```yaml
systems:
  - id: system1
    code:
      - id: codeSource1
        extractor:
          type: fs
          options:
            path: path/to/code
```

**path**: The path that code will be extracted from. If the path is not absolute it is joined with the current working directory to turn it into an absolute path. Note that the fs extractor uses [Root](https://pkg.go.dev/os@go1.24.1#Root) when scanning a directory, meaning that while symlinks are followed they must be within the Root.

##### Code Source Extractor Options (git)
Extract source code from a git repository (local or remote). For more information on how extraction works please see the documentation for [Extract Current](../explanation/extract-current.md) and [Extract Change](../explanation/extract-change.md).

Note that Include and Exclude globs are relative to the root of the repository.

```yaml
systems:
  - id: system1
    code:
      - id: codeSource1
        extractor:
          type: git
          options:
            path: path/to/repo
            branch: main
      - id: codeSource2
        extractor:
          type: git
          options:
            path: path/to/repo
            repo: git@github.com:appgardenstudios/hyaline.git
            branch: main
            clone: true
            auth:
              type: ssh
              options:
                user: git
                pem: -----BEGIN OPENSSH...
                password: pem-password...
      - id: codeSource3
        extractor:
          type: git
          options:
            repo: https://github.com/appgardenstudios/hyaline-example.git
            branch: main
            clone: true
            auth:
              type: ssh
              options:
                username: git
                password: github_pat_...
```

**path**: The local path to the repository. If the path is not absolute it is joined with the current working directory to turn it into an absolute path. If `clone` is false the repository at path is opened. If `clone` is true the repository is cloned to the path before being opened. `path` is required if `clone` is false.

**repo**: The remote git repository to use. Can be an ssh or http(s) URL. Only required if `clone` is true.

**branch**: The branch to extract source code from. If not set will default to `main`.

**clone**: Boolean specifying wether or not to clone the repository before opening. If true `repo` is also required. Defaults to false.

**auth**: Authentication information for cloning the repository. Note that if no auth is specified Hyaline will still attempt to clone, and if the repo URL is ssh your local ssh configuration will be used automatically.

**auth.type**: The type of authentication. Can be either `ssh` or `http`. Type should match the type of repo URL supplied (Hyaline does __not__ attempt to auto-detect which auth option to use based on the repo URL)

**auth.options**: Authentication options based on the type specified.

**auth.options.user**: (`ssh`) The ssh user to use when cloning the repository. Defaults to `git`.

**auth.options.pem**: (`ssh`) The contents of the private key to use when cloning the repository.

**auth.options.username**: (`http`) The http username to use when cloning. Defaults to `git`

**auth.options.password**: (`ssh` AND `http`) For `ssh`, the encryption password to use if the PEM contains a password encrypted PEM block. For `http` the password to use when cloning (will usually be a GitHub PAT or equivalent). 

## Documentation Source
TODO

# Common Documents