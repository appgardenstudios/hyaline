---
title: "Reference: Hyaline Config"
description: Configuration file schema including extract, check, and audit options
purpose: Document the configuration options for Hyaline
---
## Overview
This documents the configuration options and format present in the Hyaline configuration file.

Note that sections prefixed with **(Common)** are used in multiple places and are defined at the bottom.

## Secrets
Hyaline has the ability to pull configuration values from environment variables. To use this functionality set the value of a key to `${ENV_VAR_NAME}` to use the value of the environment variable called `ENV_VAR_NAME`.

```yaml
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  key: ${HYALINE_ANTHROPIC_KEY}

github:
  token: ${HYALINE_GITHUB_PAT}
```

In the configuration example above `llm.key` will be set to the value of the environment variable `HYALINE_ANTHROPIC_KEY`, and `github.token` will be set to the value of the environment variable `HYALINE_GITHUB_PAT`

## LLM
The connection information to use when calling out to an LLM.

```yaml
llm:
  provider: anthropic | openai | github-models | testing
  model: model-identifier
  key: ${LLM_API_KEY}
  endpoint: <custom-provider-url>
```

**provider**: The provider to use when calling out to an LLM. See [LLMs](../llms/) to learn more about the supported providers.

**model**: The LLM model to use. See each provider's documentation for a list of possible values.

**key**: The API key to use in requests. Note that this should be pulled from the environment and not hard-coded in the configuration file itself (see Secrets above)

**endpoint**: An optional custom provider URL. Specify this if your LLM provider is hosted from a non-standard URL. For `github-models`, this defaults to `https://models.github.ai/inference`.

## GitHub
The configuration for calling out to GitHub (not used for extraction, just for PR and issue retrieval during checks)

```yaml
github:
  token: ${GITHUB_PAT}
```

**token**: The GitHub token. Should be able to read pull requests and issues from relevant repositories when using `check diff`. Should be able to read pull requests, read issues, read/write issue comments, and read repo files when using `check pr`.

## Extract
Stores the configuration to use when extracting documentation.

```yaml
extract:
  disabled: false
  source:
  crawler:
  extractors:
  metadata:
```

**disabled**: If extract is disabled or not. Defaults to false.

**source**: Metadata assigned to the source being extracted.

**crawler**: The crawler to use to extract documentation.

**extractors**: A list of extractors to use when extracting documentation.

**metadata**: Metadata to add to the extracted documents and sections.

### Extract Source
Metadata about the source being extracted.

```yaml
extract:
  source:
    id: Source1
    description: A description of this source
    root: git@github.com:appgardenstudios/hyaline.git
```

**id**: Each documentation source is assigned an ID. This ID must be unique across all documentation sources used within an organization. The ID must match the regex `/^[A-z0-9][A-z0-9_-]{0,63}$/`.

**description**: A description of this documentation source.

**root**: An optional override for the root of this documentation source. Normally this is calculated based on the crawler used, but this property can be used to override the derived root. See **Extract Crawler** for more information on how the source root is calculated.

### Extract Crawler
Crawler configuration for the documentation source being extracted.

A note about the source root: it is calculated using the following algorithm based on the crawler that is configured

- If `extract.source.root` is set then that value is used.
- Else if the crawler type is `fs` then the value of `crawler.options.path` is used.
- Else if the crawler type is `git`:
  - If `crawler.options.repo` is set then that value is used.
  - Else the value of `crawler.options.path` is used.
- Else if the crawler type is `http` then the scheme and host from `crawler.options.baseUrl` is used (e.g. `https://example.com`)

```yaml
extract:
  crawler:
    type: fs | git | http
    options: {...} # Dependent on the crawler type
    include: ["**/*.md"]
    exclude: ["LICENSE.md"]
```

**type**: The type of the crawler. For Documentation Sources there are three crawler types available: `fs`, `git`, and `http`. For more information see crawler details below.

**options**: The options for the crawler. Note that these are specific to the type of crawler. Please see below for the options available for each crawler.

**include**: The set of globs to include in the set of documentation during the crawling process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See crawl option details below for how path comparisons are made and how relative glob paths work.

**exclude**: The set of globs to exclude from the set of documentation during the crawling process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. See crawl option below for how path comparisons are made and how relative glob paths work.

#### Extract Crawler Options (fs)
Crawl a file system path. If you have a local git repository use the git crawler with the `path` option instead.

Note that Include and Exclude globs are relative to the path specified.

Please see the explanation of [Extract](../explanation/extract.md) for more information.

```yaml
extract:
  crawler:
    type: fs
    options:
      path: path/to/documentation
```

**path**: The path that documentation will be crawled. If the path is not absolute it is joined with the current working directory to turn it into an absolute path. Note that the fs crawler uses [Root](https://pkg.go.dev/os@go1.24.1#Root) when scanning a directory, meaning that while symlinks are followed they must be within the Root to be crawled.

#### Extract Crawler Options (git)
Crawl a local or remote git repository. The behavior of the crawler when resolving a git repo is as follows:

- If `clone` is set:
  - If `path` is set, then the remote repo is cloned to path specified on disk
  - Else the remote repo is cloned into an in-memory file system
- Else:
  - If `path` is not set, error
  - Else open the local repository specified by `path` on disk

When cloning, authorization is handled as follows:

- If `auth.type` is `http`:
  - If `auth.username` is set, then basic auth uses that for the username
  - Else basic auth uses the value `git` for the username
  - Finally basic auth uses `auth.password` for the password
- If `auth.type` is `ssh`:
  - If `auth.user` is set, then ssh auth uses that for the user
  - Else ssh auth uses the value `git` for the user
  - Finally ssh auth uses the PEM key specified by `auth.pem` as the ssh auth key. Note that if `auth.password` is set it is used as the password when decoding in the PEM key.

For more information on how extraction works please see the documentation for [Extract](../explanation/extract.md).

Note that Include and Exclude globs are relative to the root of the repository.

```yaml
extract:
  crawler:
    type: git
    options:
      path: path/to/repo
      branch: main
```

```yaml
extract:
  crawler:
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
          pem: -----BEGIN OPENSSH... # Or use an env var like ${HYALINE_SSH_PEM}
          password: pem-password... # Or use an env var like ${HYALINE_SSH_PWD}
```

```yaml
extract:
  crawler:
    type: git
    options:
      repo: https://github.com/appgardenstudios/hyaline-example.git
      branch: main
      clone: true
      auth:
        type: http
        options:
          username: git
          password: github_pat_... # Or use an env var like ${HYALINE_GITHUB_PAT}
```

**path**: The local path to the repository. If the path is not absolute it is joined with the current working directory to turn it into an absolute path. If `clone` is false the repository at path is opened. If `clone` is true the repository is cloned to the path before being opened. `path` is required if `clone` is false.

**repo**: The remote git repository to use. Can be an ssh or http(s) URL. Only required if `clone` is true.

**branch**: The branch to crawl. If not set will default to `main`. Tries to resolve to a local branch first, then a remote branch (if there is a single remote), and finally a tag.

**clone**: Boolean specifying wether or not to clone the repository before opening. If true `repo` is also required. Defaults to false.

**auth**: Authentication information for cloning the repository. Note that if no auth is specified Hyaline will still attempt to clone, and if the repo URL is ssh your local ssh configuration will be used automatically.

**auth.type**: The type of authentication. Can be either `ssh` or `http`. Type should match the type of repo URL supplied (Hyaline does __not__ attempt to auto-detect which auth option to use based on the repo URL).

**auth.options**: Authentication options based on the type specified.

**auth.options.user**: (`ssh`) The ssh user to use when cloning the repository. Defaults to `git`.

**auth.options.pem**: (`ssh`) The contents of the private key to use when cloning the repository. Note that the encoded pem must contain the standard newlines, so use double quotes a la `"-----BEGIN OPENSSH PRIVATE KEY-----\n..." when exporting it to the relevant environment variable.

**auth.options.username**: (`http`) The http username to use when cloning. Defaults to `git`.

**auth.options.password**: (`ssh` AND `http`) For `ssh`, the encryption password to use if the PEM contains a password encrypted PEM block. For `http` the password to use when cloning (will usually be a GitHub PAT or equivalent). 

#### Extract Crawler Options (http)
Crawl a local or remote http or https website.

Note that Include and Exclude globs are relative to the baseURL.

Please see the explanation of [Extract](../explanation/extract.md) for more information.

```yaml
extract:
  crawler:
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

### Extract Extractors
Extractor configuration for the documentation source being extracted.

A note about the extractor being used. The first extractor that matches the document path (relative to the root of the crawler) is used. If there is no extractor configured to handle the document an error is returned.

Please see the explanation of [Extract](../explanation/extract.md) for more information.

```yaml
extract:
  extractors:
    - type: md | html
      options: # Dependent on the extractor type
      include: ["**/*.md"]
      exclude: []
```

**type**: The type of documentation extractor. `md` and `html` are the currently supported types.

**options**: Options used when extracting documentation and converting it into markdown (if applicable).

**include**: The set of globs to match against during the extraction process. Crawled documents must match at least one glob in order to be extracted using the extractor. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths.

**exclude**: The set of globs to exclude from the set of documentation during the extraction process. Crawled documents must match at none of these globs in order to be extracted using the extractor. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths.

#### Extract Extractors Options
Extractor options based on the type of the extractor.

```yaml
extract:
  extractors:
    - type: md
      options: # There are no options for the md extractor
```

```yaml
extract:
  extractors:
    - type: html
      options:
        selector: main
```

**selector**: A css-style selector used to extract documentation when the type of documentation is html. Only documentation that is a child of this selector will be extracted. Uses [Cascadia](https://pkg.go.dev/github.com/andybalholm/cascadia). Please see the explanation of [Extract](../explanation/extract.md) for more information.

### Extract Metadata
Configuration about what metadata to add to the extracted documentation.

Note that the specified metadata (`purpose` and/or `tags`) is added to each document or section that matches. If only `document` is specified only matching documents have the metadata applied. If both `document` and `section` are specified only sections matching both `document` and `section` have the metadata applied.

Note that metadata is applied sequentially, meaning that any overlapping documents or sections will have their purpose(s) overwritten and the tags added to.

```yaml
extract:
  metadata:
    - document: README.md
      section: About
      purpose: My document or section purpose
      tags:
        - key: system
          value: my-app
```

**document**: A glob to match a set of documents. Documents must match this extractor to have metadata applied. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. `document` is required.

**section**: A glob to match a set of sections. Sections must match this extractor to have metadata applied. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. `section` is optional.

**purpose**: The purpose to associate with the specified document or section.

**tags**: The set of tags to associate with the specified document or section.

**tags[].key**: The key of the tag to add. Must match the regex `/^[A-z0-9][A-z0-9_-]{0,63}$/`.

**tags[].value**: The value of the tag to add. Must match the regex `/^[A-z0-9][A-z0-9_-]{0,63}$/`.

## Check
Stores the configuration to use when checking documentation.

```yaml
check:
  disabled: false
  code:
  documentation:
  options:
```

**disabled**: If check is disabled or not. Defaults to false.

**code**: The set of code to evaluate when checking for recommended updates.

**documentation**: The set of documentation to include when evaluating which documents/sections need to be updated.

**options**: Options used to configure how the check process runs.

### Check Code
Determine what code is included when checking for recommended updates. Only code that is included is used when evaluating what documentation should be updated, so only include code that affects documentation (i.e. source code and not tests or tool configuration files).

```yaml
check:
  code:
    include:
      - "**/*.js"
      - "package.json"
    exclude:
      - "old/**/*"
      - "**/*.test.js"
```

**include**: The set of globs dictating what code files to include and consider during the check process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. Each glob is relative to the root of the repository.

**exclude**: The set of globs dictating what code files to exclude and not consider during the check process. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. Each glob is relative to the root of the repository.

### Check Documentation
Determine what documentation should be included in the set of documentation considered. Note that documents and sections must be included and not excluded to be considered when recommending what documentation to update.

```yaml
check:
  documentation:
    include:
      - source: "my-app"
        document: "**/*"
      - source: "**/*"
        tags:
          - key: system
            value: my-app
      - uri: document://product-docs/**/*
    exclude:
      - source: my-app
        document: README.md
        section: License
```

**include**: A set of Documentation Filters (see below) dictating what documentation is in scope of this check.

**exclude**: A set of Documentation Filters (see below) dictating what documentation is not in scope of this check.

### Check Options
Various options used when checking what documentation needs to be updated based on a code change.

```yaml
check:
  options:
    detectDocumentationUpdates:
    updateIf:
```

**detectDocumentationUpdates**: Option to detect documentation updates and mark recommendations as changed.

**updateIf**: Options to link code and documents so that code changes will generate documentation update recommendations based on the configuration.

#### Check Options DetectDocumentationUpdates
Detect documentation updates and mark recommendations as changed.

```yaml
check:
  options:
    detectDocumentationUpdates:
      source: my-app
```

**source**: If set, Hyaline will mark documents and sections as changed if they 1) have the same source and 2) the document was touched as a part of the change being examined (i.e. the document was changed in the diff or the pull request)

#### Check Options UpdateIf
Configure Hyaline to recommend that documentation be updated if a corresponding file change occurs.

```yaml
check:
  options:
    updateIf:
      touched: [...]
      added: [...]
      modified: [...]
      deleted: [...]
      renamed: [...]
```

**touched**: A list of UpdateIf Entries (see UpdateIf Entry below) detailing that this document should be updated if any matching files are touched (e.g. added, modified, deleted, or renamed).

**added**: A list of UpdateIf Entries (see UpdateIf Entry below) detailing that this document should be updated if any matching files are added (e.g. created or inserted).

**modified**: A list of UpdateIf Entries (see UpdateIf Entry below) detailing that this document should be updated if any matching files are modified (e.g. changed).

**deleted**: A list of UpdateIf Entries (see UpdateIf Entry below) detailing that this document should be updated if any matching files are deleted (e.g. removed).

**renamed**: A list of UpdateIf Entries (see UpdateIf Entry below) detailing that this document should be updated if any matching files are renamed (e.g. moved).

##### **Check Options UpdateIf Entry**
An entry that specifies that matching documentation should be updated if matching code was changed.

```yaml
check:
  options:
    updateIf:
      touched: # A list of UpdateIf Entries
        - code:
            path: "src/routes.js"
          documentation: # A Documentation Filter
            source: "my-app"
            document: "docs/routes.md"
```

**code**: The code that triggers the update.

**code.path**: A glob dictating what code files to match. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. The glob is relative to the root of the repository.

**documentation**: The Documentation Filter (see below) that determines which documentation to match.

## Audit
Stores the configuration to use when auditing documentation.

```yaml
audit:
  disabled: false
  rules:
```

**disabled**: If audit is disabled or not. Defaults to false.

**rules**: A list of audit rules to evaluate against the documentation.

### Audit Rules
Define specific checks to run against your documentation.

```yaml
audit:
  rules:
    - id: "content-exists-check"
      description: "Check that critical documentation exists"
      documentation:
        - source: "**/*"
          document: "README.md"
      ignore:
        - source: "internal"
      checks:
        content:
        purpose:
        tags:
```

**id**: A unique identifier for the rule. Must match the regex `/^[A-z0-9][A-z0-9_-]{0,63}$/`. If not provided, an auto-generated ID will be assigned (e.g., `_0`, `_1`).

**description**: A human-readable description of what this rule audits.

**documentation**: A list of Documentation Filters (see below) dictating what documentation this rule applies to. At least one filter is required.

**ignore**: A list of Documentation Filters (see below) dictating what documentation to exclude from this rule's evaluation.

**checks**: The validation checks to perform on matching documentation.

#### Audit Rules Checks Content
Validate the content of documentation.

```yaml
audit:
  rules:
    - checks:
        content:
          exists: true
          min-length: 100
          matches-regex: "(?i)installation"
          matches-prompt: "Does this document contain deployment instructions?"
          matches-purpose: true
```

**exists**: Boolean indicating whether the content must exist. When true, the check passes if at least one document or section matches the documentation filters. When false or not set, this check is disabled.

**min-length**: Minimum required content length in characters. The check fails if the content is shorter than this value. When not set, this check is disabled.

**matches-regex**: A regular expression pattern that the content must match. When not set, this check is disabled.

**matches-prompt**: An LLM prompt to evaluate against the content. The LLM will assess whether the content satisfies the prompt's criteria. When not set, this check is disabled.

**matches-purpose**: Boolean indicating whether the content should match its defined purpose. Requires the document or section to have a purpose defined. When false or not set, this check is disabled.

#### Audit Rules Checks Purpose
Validate that documentation has defined purposes.

```yaml
audit:
  rules:
    - checks:
        purpose:
          exists: true
```

**exists**: Boolean indicating whether a purpose must be defined for the documentation. When false or not set, this check is disabled.

#### Audit Rules Checks Tags
Validate the presence of specific tags on documentation.

```yaml
audit:
  rules:
    - checks:
        tags:
          contains:
            - key: "type"
              value: ".*guide"
            - key: "level"
              value: "beginner"
```

**contains**: A list of tags that must be present on the documentation. When set, the check passes if documentation has all of the listed tags. When not set, this check is disabled.

**contains[n].key**: The tag key to check for. Uses regex pattern matching.

**contains[n].value**: The tag value to check for. Uses regex pattern matching.

## (Common) Documentation Filter
A filter to use to select a subset of documentation.

```yaml
check:
  documentation:
    include: # An array of Documentation Filters
      - source: "my-app"
      - source: "api"
        document: "my-app/**/*"
      - source: "security"
        document: "frontend.md"
        section: "my-app"
      - source: "**/*"
        tags:
          - key: system
            value: my-app
      - uri: document://product-docs/**/*
```

**source**: A glob that matches against a document or section's source ID. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths. Must be set if `uri` is not set.

**document**: A glob that matches against a document or section's document ID. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths.

**section**: A glob that matches against a section's section ID. This uses the [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) package to match paths.

**tags**: A set of tags to match the document or section against.

**tags[n].key**: A tag key. Must match `/^[A-z0-9][A-z0-9_-]{0,63}$/`

**tags[n].value**: A tag value. Must match `/^[A-z0-9][A-z0-9_-]{0,63}$/`

**uri**: An encoded document URI in the format of `document://<source-id>/<path/of/document.md>#<path/of/section>`. Must start with `document://` and contain at least a source and document glob. Each section (source, document, section) must be a valid [doublestar](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) glob. Must be set if `source` is not set.