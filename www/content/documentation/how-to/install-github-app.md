---
title: "How To: Install the GitHub App"
description: "Install the Hyaline GitHub App."
purpose: Document how to install the Hyaline GitHub App
---
## Purpose
Install the Hyaline GitHub App into a GitHub organization or personal account.

## Prerequisite(s)
- A GitHub organization or personal account
- One or more documentation sources (e.g. git repo, documentation website, etc...)
- A [supported LLM Provider](../llms/)

Note that the Hyaline GitHub App will trigger workflows in the configuration repository located in your organization or personal account, meaning that you stay in control of your configuration and data. If you wish you can also run your own copy of the Hyaline GitHub app (located in the [configuration repository](https://github.com/appgardenstudios/hyaline-github-app-config)) to prevent any data whatsoever from leaving your organization and being sent to us.

Note also that one of the supported LLM Providers is [GitHub Models](https://github.com/features/models), which offers a free tier and can be used in conjunction with the GitHub Actions free tier.

## Steps

### 1. Create GitHub App Config Repo
All of your configuration for the Hyaline GitHub App will live in a single repository in your organization or personal account. The easiest way to set this up is to go to the [hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) repository, click "Use this template," and create a new repository called `hyaline-github-app-config` in the organization or personal account that you will install the GitHub App into.

Create from template (or otherwise clone/push) the [hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) repository into your organization or personal account. Note that the repository name MUST remain `hyaline-github-app-config` in order to use the hosted version of the Hyaline GitHub App.

Please see GitHub's documentation on [how to create a repository from a template](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template).

Note that when new changes are released to [hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config), you will be able to pull these changes into your repo instance by running the "Update Hyaline" workflow.

### 2. Setup Secrets and Environment Variables
You will need to setup the following secrets and environment variables in your `hyaline-github-app-config` repo instance.

Note: Since Hyaline uses GitHub Personal Access Tokens (PATs) issued by you to act on your behalf, it is recommended to create a dedicated service account to use when issuing PATs. This is so that the comments made by Hyaline on pull requests will use the service account name instead of an individual's name. The service account will need read access to each repository that the GitHub App has access to, and write access to your `hyaline-github-app-config` repo instance.

#### Secrets
The following repository secrets should be created in your `hyaline-github-app-config` repo instance:

**HYALINE_GITHUB_TOKEN** - A GitHub Personal Access Token (PAT) that will be used to extract repo documentation and comment on pull requests (this will be referenced as the value for `github.token` and `extract.crawler.options.auth.password` in the configs via environment substitution). This token should be scoped to the repositories that Hyaline will be extracting documentation from or checking PRs for. It will need to have the following permissions:
- Metadata: Read - Required by GitHub for all PATs
- Contents: Read - Used for extracting documentation from the in-scope repositories
- Pull requests: Read and Write - Used for creating/updating the Pull Request comment containing Hyaline's recommendations

Note that this PAT will include access to public repositories in the organization or personal account as well as any private repositories that were explicitly added to the scope of the PAT.

**HYALINE_CONFIG_GITHUB_TOKEN** - A GitHub Personal Access Token (PAT) that will be used to manage the GitHub App's configuration. This token should be scoped to your `hyaline-github-app-config` repo instance. It will need to have the following permissions:
- Metadata: Read - Required by GitHub for all PATs
- Actions: Read and Write - Used by extract workflows in the config to trigger the merge workflow once extraction is complete
- Contents: Read and Write - Used to clone the configuration in workflows and used by the doctor to push suggested changes to a branch for review
- Pull requests: Read and Write - Used by the doctor to open a pull request with suggested changes
- Workflows: Read and Write - Used by the doctor to push suggested changes to extract and audit workflows to a branch for review

**HYALINE_LLM_TOKEN** - An LLM provider API token used in auditing and checking PRs. This will need to come from the LLM provider and will be referenced as the value for `llm.key` in the configs (using environment substitution).

#### Environment Variables
The following repository variables should be created in your `hyaline-github-app-config` repo instance:

**HYALINE_LLM_PROVIDER** - The LLM provider to be used. This will be referenced as the value for `llm.provider` in the configs (using environment substitution). Please see [configuration reference](../reference/config.md) for supported values.

**HYALINE_LLM_MODEL** - The LLM model to be used. This will be referenced as the value for `llm.provider` in the configs (using environment substitution). Please see [configuration reference](../reference/config.md) for supported values.

### 3. Run Install
To bootstrap the repository in preparation for the Github App installation you will need to run the `Install` workflow and review/edit/merge the generated pull request.

Manually trigger the `Install` workflow in your `hyaline-github-app-config` repo instance and ensure that it completes successfully. It will run `Update Hyaline` to ensure your `hyaline-github-app-config` repo is properly connected and up-to-date. It also runs `Doctor` to generate a pull request with a set of suggested changes and configuration updates based on the repositories in scope of the `HYALINE_GITHUB_TOKEN` generated above.

### 4. Review/Merge the Doctor PR
Review (editing as necessary) and merge the pull request generated by the `Install` workflow to the default (`main`) branch. You can view [how to extract documentation](./extract-documentation.md) and [how to check pull request](./check-pull-request.md) for more detail.

Note that you can explicitly disable unwanted extractions or checks by adding `disabled: true` to the `extract` or `check` section of the configuration (see [configuration reference](../reference/config.md) for more details).

### 5. Run Extract All
Manually trigger the `Extract All` workflow in your `hyaline-github-app-config` repo instance and ensure that it completes successfully. This will run an extract on all configured repositories and merge the documentation together into a single current data set for use in audits and checks.

### 6. Install GitHub App
Install the [Hyaline GitHub App (hyaline.dev)](https://github.com/apps/hyaline-dev) into your organization or personal account. Only grant it access to repositories that you want Hyaline to extract and check pull requests on in addition to your `hyaline-github-app-config` repo instance.

### 7. Verify Installation
You can verify the installation of the GitHub App by opening a non-draft PR in one of the repositories in your organization. Once you do you should see the workflow `_Check PR` kicked off in your `hyaline-github-app-config` repo instance and a comment on the pull request with Hyaline's documentation update recommendations. Then, once the pull request is merged, you should see a corresponding workflow run of `_Extract` followed by a workflow run of `_Merge` in your `hyaline-github-app-config` repo instance.

## Next Steps
Read more about [how to extract documentation from a repository or site](./extract-documentation.md) or visit [an explanation of Hyaline's GitHub App](../explanation/github-app.md).