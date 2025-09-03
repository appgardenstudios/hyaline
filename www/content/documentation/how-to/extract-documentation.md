---
title: "How To: Extract Documentation"
description: "Extract documentation using the Hyaline GitHub App."
purpose: Document how to extract documentation using the Hyaline GitHub App
---
## Purpose
Configure Hyaline to extract documentation using the Hyaline GitHub App

## Prerequisite(s)
- [Install GitHub App](./install-github-app.md)
- Have one or more documentation sources to be extracted (e.g. git repo, documentation website, etc...)

## Steps

### 1. Create Configuration
The first step is to create a configuration file for the documentation source in the appropriate folder in the forked `hyaline-github-app-config` repo.

For example, the configuration to extract documentation from a repository should be places in the `repos/` folder named `<repo-name>.yml` and should look something like:

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
  crawler:
    type: git
    options:
      repo: https://github.com/<owner>/<repo>.git
      branch: main
      clone: true
      auth:
        type: http
        options:
          username: git
          password: ${HYALINE_GITHUB_TOKEN}
    include:
      - "**/*.md"
  extractors:
    - type: md
      include:
        - "**/*.md"
```

Configuration to extract documentation from a documentation site should be placed in `sites/` and the crawler/extractors should be configured as needed (see the [configuration reference](../reference/config.md) for more information).

### 2. Run Doctor
Run the `Doctor` workflow in the forked `hyaline-github-app-config` repo to 1) ensure that the configuration is valid and 2) to add the repository or site to the list of available extraction targets. Merge the resulting PR if needed.

### 3. Run Extract
Run the `Manual - Extract` workflow in the forked `hyaline-github-app-config` repo to trigger an extraction. Note that you can trigger a merge of this documentation into the current documentation data set by leaving the `Trigger Merge Workflow` option enabled.

## Next Steps
Read more about [how extraction works](../explanation/extract.md) or more about Hyaline's [configuration](../reference/config.md).