---
title: "Reference: Hyaline Config"
description: Configuration file schema including extract, check, and audit options
purpose: Document the configuration options for Hyaline
sitemap:
  disable: true
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
The configuration for calling out to GitHub (not used for extraction, just for PR and issue retrieval during checks)

```yaml
github:
  token: ${GITHUB_PAT}
```

**token**: The GitHub token. Should be able to read pull requests and issues from relevant repositories.

## Extract
Stores the configuration to use when extracting documentation.

```yaml
extract:
  source:
  crawler:
  extractors:
  metadata:
```

**source**: TODO

**crawler**: TODO

**extractors**: TODO

**metadata**: TODO

### Source
TODO

### Crawler
TODO

### Extractors
TODO

### Metadata
TODO