---
title: "How To: Install the GitHub App"
description: "Install the Hyaline GitHub App."
purpose: Document how to install the Hyaline GitHub App
---
## Purpose
Install the Hyaline GitHub App into a GitHub organization or personal account

## Prerequisite(s)
- A GitHub organization or personal account
- One or more documentation sources (e.g. git repo, documentation website, etc...)
- A [supported LLM Provider](../reference/config.md)

## Steps

### 1. Create GitHub App Config Repo
All of your configuration for the Hyaline GitHub App will live in a single repository in your organization or personal account. The easiest way to set this up is to fork the 
[hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) repository into the organization or personal account that you will install the GitHub App into.

#### Fork hyaline-github-app-config

Fork (or otherwise clone/push) the [hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) into your organization or personal account. Note that the repository name MUST remain `hyaline-github-app-config` in order to use the hosted version of the Hyaline GitHub App.

Please see GitHub's documentation on [how to fork a repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository).

Note that we ask that you fork the configuration repository to make it easy to pull in updated documentation, new features, and bug fixes from the source repository using [GitHub's sync functionality](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork).

### 2. Setup Secrets and Environment Variables
You will need to setup the following secrets and environment variables in the forked `hyaline-github-app-config` repository.

TODO note about since Hyaline will be using a PAT to act within the org and whomever owns that name will be the name on the comments/PRs, we suggest using a service account (link to GitHub Doc)
TODO A good way to do this is using a Service Account rather than creating a PAT from a user in the org. That service account will only need read access to the repos (since it )

#### Secrets
The following repository secrets should be created in the forked `hyaline-github-app-config` repo:

**HYALINE_GITHUB_TOKEN** - A GitHub Personal Access Token (PAT) that will be used to extract repo documentation and comment on pull requests (this will be referenced as the value for `github.token` and `extract.crawler.options.auth.password` in the configs via environment substitution). This token should be scoped to the repositories that Hyaline will be extracting documentation from or checking PRs for. It will need to have the following permissions:
- Metadata: Read - Required by GitHub for all PATs
- Contents: Read - Used for extracting documentation from the in-scope repositories
- Pull requests: Read and Write - Used for creating/updating the Pull Request comment containing Hyaline's recommendations

Note that this PAT will include access to public repositories in the organization or personal account as well as any private repositories that were explicitly added to the scope of the PAT.

**HYALINE_CONFIG_GITHUB_TOKEN** - A GitHub Personal Access Token (PAT) that will be used to manage the GitHub App's configuration and use. This token should be scoped to the forked `hyaline-github-app-config` repository. It will need to have the following permissions:
- Metadata: Read - Required by GitHub for all PATs
- Actions: Read and Write - Used by extract workflows in the config to trigger the merge workflow once extraction is complete
- Contents: Read and Write - Used to clone the configuration in workflows and used by the doctor to push suggested changes to a branch for review
- Pull requests: Read and Write - Used by the doctor to open a pull request with suggested changes
- Workflows: Read and Write - Used by the doctor to push suggested changes to extract and audit workflows to a branch for review

**HYALINE_LLM_TOKEN** - A LLM provider API token used in auditing and checking PRs. This will need to come from the LLM provider and will be referenced as the value for `llm.key` in the configs (using environment substitution)

#### Environment Variables
The following repository variables  should be created in the forked `hyaline-github-app-config` repo:

**HYALINE_LLM_PROVIDER** - The LLM provider to be used. This will be referenced as the value for `llm.provider` in the configs (using environment substitution). Please see [configuration reference](../reference/config.md) for supported values.

**HYALINE_LLM_MODEL** - The LLM model to be used. This will be referenced as the value for `llm.provider` in the configs (using environment substitution). Please see [configuration reference](../reference/config.md) for supported values.

### 3. Run Doctor
To bootstrap the repository in preparation for the Github App installation you will need to run the `Doctor` workflow and review/edit/merge the generated pull request.

Manually trigger the `Doctor` workflow in the forked `hyaline-github-app-config` repository and ensure that it completes successfully. It should generate a pull request with a set of suggested changes and configuration updates based on the repositories in scope of the `HYALINE_GITHUB_TOKEN` generated above.

### 4. Review/Merge Doctor PR
Review (editing as necessary) and merge the pull request generated by the doctor to the default (`main`) branch. You can view [how to extract documentation](./extract-documentation.md) and [how to check pull request](./check-pull-request.md) for more detail.

Note that you can explicitly disable unwanted extractions or checks by adding `disabled: true` to the `extract` or `check` section of the configuration (see [configuration reference](../reference/config.md) for more details).

### 5. Run Extract All
Manually trigger the `Manual - Extract All` workflow in the forked `hyaline-github-app-config` repository and ensure that it completes successfully. This will run an extract on all configured repositories and merge the documentation together into a single current dataset for use in audits and checks.

### 6. Install GitHub App
Install the [Hyaline GitHub App (hyaline.dev)](https://github.com/apps/hyaline-dev) into your organization or personal account. Only grant it access to repositories that you want Hyaline to extract and check pull requests on.

### 7. Verify Installation
You can verify the installation of the GitHub App by opening a non-draft PR in one of the repositories in your organization. Once you do you should see the workflow `Internal - Check PR` kicked off in the forked `hyaline-github-app-config` repository and a comment on the pull request with Hyaline's documentation update recommendations. Then, once the pull request is merged, you should see a corresponding workflow run of `Internal - Extract` followed by a workflow run of `Internal - Merge` also in the forked `hyaline-github-app-config` repository.

## Next Steps
To read more about the Hyaline GitHub App please see [an explanation of Hyaline's GitHub App](../explanation/github-app.md) or read on to see how to [extract documentation from a repository or site](./extract-documentation.md).