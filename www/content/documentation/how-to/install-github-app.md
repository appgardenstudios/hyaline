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

## Steps

### 1. Create GitHub App Config Repo
All of your configuration for the Hyaline GitHub App will live in a single repository in your organization or personal account. The easiest way to set this up is to fork the 
[hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) repository into the organization or personal account that you will install the GitHub App into.

#### 1.1 Fork hyaline-github-app-config
Fork (or otherwise clone/push) the [hyaline-github-app-config](https://github.com/appgardenstudios/hyaline-github-app-config) into your organization or personal account. Note that the repository name MUST remain `hyaline-github-app-config` in order to use the hosted version of the Hyaline GitHub App.

Please see GitHub's documentation on [how to fork a repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository).

Note that we ask that you fork the configuration repository to make it easy to pull in updated documentation, new features, and bug fixes from the source repository using [GitHub's sync functionality](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork).

### 2. Setup Secrets and Environment Variables
You will need to setup the following secrets and environment variables in the forked `hyaline-github-app-config` repository.

#### 2.1 Secrets
HYALINE_GITHUB_TOKEN
HYALINE_CONFIG_GITHUB_TOKEN
HYALINE_LLM_TOKEN

#### 2.2 Environment Variables
HYALINE_LLM_PROVIDER
HYALINE_LLM_MODEL

### 3. Run Doctor

#### 3.1 Trigger Workflow

#### 3.2 Review and Merge PR

### 4. Run Extract All

### 5. Install GitHub App

## Next Steps
TODO