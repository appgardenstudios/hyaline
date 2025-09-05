---
title: "How To: Audit Documentation"
description: "Audit documentation using the Hyaline GitHub App."
purpose: Document how to audit documentation using the Hyaline GitHub App
---
## Purpose
Configure Hyaline to audit documentation using the Hyaline GitHub App.

## Prerequisite(s)
- [Install GitHub App](./install-github-app.md)
- Have documentation extracted that can be audited

## Steps

### 1. Create Configuration
The first step is to create a configuration file for the audit in the `audits/` folder in the forked `hyaline-github-app-config` repo.

For example, the configuration could look something like:

```yml
llm:
  provider: ${HYALINE_LLM_PROVIDER}
  model: ${HYALINE_LLM_MODEL}
  key: ${HYALINE_LLM_TOKEN}

audit:
  rules:
    - id: "content-length-check"
      description: "Check that README has sufficient content"
      documentation:
        - source: "**/*"
          document: "README.md"
      checks:
        content:
          min-length: 100
```

For more information on how to configure an audit and what checks are available please see the [audit explanation](../explanation/audit.md) or the [configuration reference](../reference/config.md).

### 2. Run Doctor
Run the `Doctor` workflow in the forked `hyaline-github-app-config` repo to 1) ensure that the configuration is valid and 2) to add the audit to the list of available audits. Merge the resulting PR if needed.

### 3. Run Audit
Run the `Run Audit` workflow in the forked `hyaline-github-app-config` repo to trigger an audit of the current extracted documentation. Audit results will be attached to the workflow run once it completes.

## Next Steps
Read more about [how auditing documentation works](../explanation/audit.md) or more about Hyaline's [configuration](../reference/config.md).