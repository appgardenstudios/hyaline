---
title: "Reference: Hyaline Config"
linkTitle: Config
purpose: Document the configuration options for Hyaline
url: documentation/reference/config
---
## Overview
This documents the configuration options and format present in the Hyaline configuration file.

## Secrets
Hyaline has the ability to pull configuration values from environment variables. To use this functionality set the value of a key to `${ENV_VAR_NAME}` to use the value of the environment variable called `ENV_VAR_NAME`.

```yaml
llm:
  provider: anthropic
  model: claude-3-5-sonnet-20241022
  key: ${HYALINE_ANTHROPIC_KEY}

github:
  token: ${HYALINE_GITHUB_PAT}
```

In the configuration example above `llm.key` will be set to the value of the environment variable `HYALINE_ANTHROPIC_KEY`, and `github.token` will be set to the value of the environment variable `HYALINE_GITHUB_PAT`

## LLM
The connection information to use when calling out to an LLM.

```yaml
llm:
  provider: anthropic | testing
  model: model-identifier
  key: ${LLM_API_KEY}
```

**provider**: The provider to use when calling out to an LLM. possible values are `anthropic` and  `testing`.

**model**: The LLM model to use. See each provider's documentation for a list of possible values.

**key**: The API key to use in requests. Note that this should be pulled from the environment and not hard-coded in the configuration file itself (see Secrets above)

## GitHub
The configuration for calling out to GitHub (not used for extraction, just for PR and issue retrieval)

```yaml
github:
  token: ${GITHUB_PAT}
```

**token**: The GitHub token. Should be able to read pull requests and issues from relevant repositories.

## Systems
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
          include: [package.json, "**/*.js"]
          exclude: ["**/*.test.js"]
```

**type**: The type of the extractor. For Code Sources there are two extractor types available: `fs` and `git`. For more information see extractor details below.

**options**: The options for the extractor. Note that these are specific to the type of extractor. Please see below for the options available for each extractor.

**include**: The set of globs to include in the set of source code during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

**exclude**: The set of globs to exclude from the set of source code during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

##### Code Source Extractor Options (fs)
Extract source code from a file system path. Note that code sources using this extractor will not be eligible to be included in a change data set. If you have a local git repository use the git extractor with the `path` option.

Note that Include and Exclude globs are relative to the path specified.

Please see the explanation of [Extract Current](../03-explanation/02-extract-current.md) for more information.

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
Extract source code from a git repository (local or remote). For more information on how extraction works please see the documentation for [Extract Current](../03-explanation/02-extract-current.md) and [Extract Change](../03-explanation/03-extract-change.md).

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

**branch**: The branch to extract source code from. If not set will default to `main`. If the specified branch cannot be resolved locally and doesn't contain a `/`, the system will attempt to resolve `origin/<branch>` as a fallback.

**clone**: Boolean specifying wether or not to clone the repository before opening. If true `repo` is also required. Defaults to false.

**auth**: Authentication information for cloning the repository. Note that if no auth is specified Hyaline will still attempt to clone, and if the repo URL is ssh your local ssh configuration will be used automatically.

**auth.type**: The type of authentication. Can be either `ssh` or `http`. Type should match the type of repo URL supplied (Hyaline does __not__ attempt to auto-detect which auth option to use based on the repo URL)

**auth.options**: Authentication options based on the type specified.

**auth.options.user**: (`ssh`) The ssh user to use when cloning the repository. Defaults to `git`.

**auth.options.pem**: (`ssh`) The contents of the private key to use when cloning the repository.

**auth.options.username**: (`http`) The http username to use when cloning. Defaults to `git`.

**auth.options.password**: (`ssh` AND `http`) For `ssh`, the encryption password to use if the PEM contains a password encrypted PEM block. For `http` the password to use when cloning (will usually be a GitHub PAT or equivalent). 

### Documentation Source
A set of documentation for the system, or system documentation.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        type: md | html
        options: {...}
        extractor: {...}
        documents:
          - path: path/to/document.md
            ...
        includeDocuments: [commonDocumentID, ...]
```

**id**: Each Documentation Source has an ID used to reference it. The ID must be unique across all Documentation Sources for a system (note that you can use the same ID for Documentation Sources in different systems if desired).

**type**: The type of documentation. `md` and `html` are the currently supported documentation types.

**options**: Documentation source options used when converting documentation into markdown.

**extractor**: Each Documentation Source has an extractor that is responsible for extracting the documentation out and into the data set.

**documents**: A set of documents specifying the desired structure and contents of the system documentation.

**includeDocuments**: A set of IDs referencing one or more sets of Common Documents

#### Documentation Source Options
Documentation source options.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        type: md | html
        options:
          selector: "#main"
```

**selector**: A css-style selector used to extract system documentation when the type of documentation is html. Uses [Cascadia](https://pkg.go.dev/github.com/andybalholm/cascadia). See [Extract Current](../03-explanation/02-extract-current.md) and [Extract Change](../03-explanation/03-extract-change.md) for more information.

#### Documentation Source Extractor
An extractor that specifies how documentation is extracted for this Documentation Source and placed into a data set.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        extractor:
          type: fs | git | http
          options: {...} # Dependent on the extractor type
          include: ["**/*.md"]
          exclude: ["LICENSE.md"]
```

**type**: The type of the extractor. For Documentation Sources there are three extractor types available: `fs`, `git`, and `http`. For more information see extractor details below.

**options**: The options for the extractor. Note that these are specific to the type of extractor. Please see below for the options available for each extractor.

**include**: The set of globs to include in the set of documentation during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

**exclude**: The set of globs to exclude from the set of documentation during the extraction process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See extractor details below for how path comparisons are made and how relative glob paths work.

##### Documentation Source Extractor Options (fs)
Extract documentation from a file system path. Note that documentation sources using this extractor will not be eligible to be included in a change data set. If you have a local git repository use the git extractor with the `path` option.

Note that Include and Exclude globs are relative to the path specified.

Please see the explanation of [Extract Current](../03-explanation/02-extract-current.md) for more information.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        extractor:
          type: fs
          options:
            path: path/to/documentation
```

**path**: The path that documentation will be extracted from. If the path is not absolute it is joined with the current working directory to turn it into an absolute path. Note that the fs extractor uses [Root](https://pkg.go.dev/os@go1.24.1#Root) when scanning a directory, meaning that while symlinks are followed they must be within the Root.

##### Documentation Source Extractor Options (git)
Extract documentation from a git repository (local or remote). For more information on how extraction works please see the documentation for [Extract Current](../03-explanation/02-extract-current.md) and [Extract Change](../03-explanation/03-extract-change.md).

Note that Include and Exclude globs are relative to the root of the repository.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        extractor:
          type: git
          options:
            path: path/to/repo
            branch: main
      - id: documentationSource2
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
      - id: documentationSource3
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

**branch**: The branch to extract documentation from. If not set will default to `main`. If the specified branch cannot be resolved locally and doesn't contain a `/`, the system will attempt to resolve `origin/<branch>` as a fallback.

**clone**: Boolean specifying wether or not to clone the repository before opening. If true `repo` is also required. Defaults to false.

**auth**: Authentication information for cloning the repository. Note that if no auth is specified Hyaline will still attempt to clone, and if the repo URL is ssh your local ssh configuration will be used automatically.

**auth.type**: The type of authentication. Can be either `ssh` or `http`. Type should match the type of repo URL supplied (Hyaline does __not__ attempt to auto-detect which auth option to use based on the repo URL)

**auth.options**: Authentication options based on the type specified.

**auth.options.user**: (`ssh`) The ssh user to use when cloning the repository. Defaults to `git`.

**auth.options.pem**: (`ssh`) The contents of the private key to use when cloning the repository.

**auth.options.username**: (`http`) The http username to use when cloning. Defaults to `git`.

**auth.options.password**: (`ssh` AND `http`) For `ssh`, the encryption password to use if the PEM contains a password encrypted PEM block. For `http` the password to use when cloning (will usually be a GitHub PAT or equivalent). 

##### Documentation Source Extractor Options (http)
Extract documentation from an http(s) source via crawling. Note that documentation sources using this extractor will not be eligible to be included in a change data set.

Note that Include and Exclude globs are relative to the baseURL.

Please see the explanation of [Extract Current](../03-explanation/02-extract-current.md) for more information.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        extractor:
          type: http
          options:
            baseUrl: https://www.hyaline.dev/
            start: ./documentation
            headers:
              custom-header: My Header Value
```

**baseUrl**: The base URL to start with. The baseUrl will be the starting URL if `start` is not defined. Also note that the crawler is limited to the same domain as that on the baseUrl.

**start**: An (optional) starting path relative to the baseURL. If set the crawler will start on the `baseUrl` joined with `start` path.

**headers**: A set of (optional) headers to include with each request.

#### Documentation Source Documents
Documents describing the desired layout and contents of each piece of documentation for this documentation source.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        documents:
          - name: README.md
            purpose: The purpose of this document is...
            required: true
            ignore: false
            updateIf: {...}
            sections:
              - name: Running Locally
                purpose: The purpose of this section is...
                ...
```

**name**: The name of the document, including the relative path (uses the same relative path logic as Include/Exclude).

**purpose**: The purpose of the document.

**required**: True if this document is required (default false).

**ignore**: True if this document should be ignored (default false).

**updateIf**: Configuration on when to update this document based on code changes that may occur.

**sections**: A list of document sections.

##### Documentation Source Document Update If
Conditions in which documentation should be reviewed and/or updated based on code changes that may occur.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        documents:
          - name: README.md
            updateIf:
              touched: [...]
              added: [...]
              modified: [...]
              deleted: [...]
              renamed: [...]
```

**touched**: A list of entries detailing that this document should be updated if any matching files are touched.

**added**: A list of entries detailing that this document should be updated if any matching files are added.

**modified**: A list of entries detailing that this document should be updated if any matching files are modified.

**deleted**: A list of entries detailing that this document should be updated if any matching files are deleted.

**renamed**: A list of entries detailing that this document should be updated if any matching files are renamed.


##### Documentation Source Document Update If Entry
Specifies which files should be looked at to see if a document or section should be updated.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        documents:
          - name: README.md
            updateIf:
              touched:
                - codeID: "*"
                  glob: "subDir1/**/*.js"
                - codeID: "codeSource1"
                  glob: "subDir2/**/*.js"
                - glob: "subDir3/**/*.js"
```

**codeID**: The ID of the code source in the system that should be matched against. If blank or `*` it will match against any code source.

**glob**: The glob to use for matching files. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. File paths are matched using the same logic as that the file's code source extractor configuration uses.

##### Documentation Source Document Sections
Sections describing the desired layout and contents of each section of documentation for this document in this documentation source.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        documents:
          - name: README.md
            purpose: The purpose of this document is...
            required: true
            ignore: false
            updateIf: {...}
            sections:
              - name: Running Locally
                purpose: The purpose of this section is...
                required: true
                ignore: false
                updateIf: {...}
                sections: [...]
```

**name**: The name of the section. Note that the section name cannot contain `#`.

**purpose**: The purpose of the section.

**required**: True if this section is required (default false).

**ignore**: True if this section should be ignored (default false).

**updateIf**: Configuration on when to update this section based on code changes that may occur. These operate the same way as the updateIfs at the document level (see above for documentation).

**sections**: A list of child sections matching this same schema.

#### Documentation Source Include Documents
Include common documents describing the desired layout and contents of each piece of documentation for this documentation source. These are added to the list of documents in this documentation source. Note that documents in the documentation source take priority over any included documents.

```yaml
systems:
  - id: system1
    documentation:
      - id: documentationSource1
        includeDocuments: [commonDocumentID, ...]

commonDocuments:
  - id: commonDocumentID
    documents:
      - name: README.md
        ...
```

## Common Documents
A set of common document that can be referenced by documentation sources to reduce duplicate configuration and enforce consistency across systems (e.g. common formats for a README.md document, for example)

```yaml
commonDocuments:
  - id: commonDocumentID
    documents:
      - name: README.md
        ...
```

**id**: The ID of the document set.

**documents**: The set of documents. This has the exact same format as the documents block in the documentation source. See details above.
