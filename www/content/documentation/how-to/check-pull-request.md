---
title: "How To: Check a GitHub Pull Request"
description: "Use the GitHub App to automatically check for needed documentation updates in pull requests using Hyaline."
purpose: How to check a pull request using the Hyaline GitHub App
---
## Purpose
Configure Hyaline to check a pull request using the Hyaline GitHub App.

## Prerequisite(s)
- [Install GitHub App](./install-github-app.md)
- Have one or more repos that should be checked

## Steps

### 1. Create Configuration
The first step is to create a configuration file for the repo in the `repos/` folder in your `hyaline-github-app-config` repo instance if one does not already exist. Note that the name of the configuration file MUST match the name of the repository.

For example, the configuration to could look something like:

```yml
llm:
  provider: ${HYALINE_LLM_PROVIDER}
  model: ${HYALINE_LLM_MODEL}
  key: ${HYALINE_LLM_TOKEN}

github:
  token: ${HYALINE_GITHUB_TOKEN}

extract:
  source:
    id: <documentation source id>
    description: <documentation source description>
  ...

check:
  code:
    include:
      - "cli/**/*.go"
    exclude:
      - "**/*_test.go"
      - "e2e/**/*"
      - "benchmarks/**/*"
  documentation:
    include:
      - source: "<documentation source id>"
        document: "**/*"
  options:
    detectDocumentationUpdates:
      source: <documentation source id>
```

### 2. Run Doctor
Run the `Doctor` workflow in your `hyaline-github-app-config` repo instance to ensure that the configuration is valid. Merge the resulting PR if needed.

### 3. Create Pull Request
Create a pull request in the repo you created the configuration for. You will see a run of the `_Check PR` workflow being kicked off in your `hyaline-github-app-config` repo instance and a comment with recommendations created on the pull request.

## Next Steps
Read more about [how checking pull requests works](../explanation/check.md) or visit the [configuration reference](../reference/config.md).