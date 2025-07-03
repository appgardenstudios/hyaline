---
title: "How To: Check a PR using the Hyaline CLI"
linkTitle: Check a PR
purpose: How to check a PR using the Hyaline CLI
url: documentation/how-to/check-pr
---
## Purpose
Use hyaline to check a pull request in GitHub.

## Prerequisite(s)
* A repository on GitHub with [GitHub Actions](https://github.com/features/actions) available.

## Steps

### 1. Check In Hyaline Configuration
Create a Hyaline Configuration File and check it into your GitHub repository.

Make sure you don't check in any secrets, like the LLM Provider Key or the GitHub Token. Instead, setup the configuration to pull them from the environment.

See [config reference](../05-reference/01-config.md) for more information on creating a configuration file and referencing secrets.

### 2. Set Up Secrets
For each environment variable used in the Hyaline configuration that references a secret, set that variable up to be pulled in as a [secret in GitHub](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions).

### 3. Create Workflow File
Create a [GitHub Workflow](https://docs.github.com/en/actions/writing-workflows/quickstart) file to run when a pull request is updated. You can see an example file in the [GitHub Actions reference](../05-reference/05-github-actions.md).

Alternatively you could set Hyaline up to be run [manually on dispatch](https://docs.github.com/en/actions/managing-workflow-runs-and-deployments/managing-workflow-runs/manually-running-a-workflow).

## Next Steps
Visit [Config Reference](../05-reference/01-config.md) or [GitHub Actions Reference](../05-reference/05-github-actions.md).