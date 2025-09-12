---
title: "How To: Try out Hyaline"
description: "How to try out Hyaline without installing the GitHub App."
purpose: Document how to try out Hyaline without installing the GitHub App
---
## Purpose
Try out Hyaline without installing the GitHub App.

## Prerequisite(s)
- Have a GitHub Repository with some code and documentation in it.

**Note**: This installation method is not the recommended way to install or use Hyaline, and is only intended to allow you to quickly setup an isolated test of Hyaline for evaluation purposes. To install Hyaline for an entire organization so that all of your documentation is checked and kept up-to-date please visit [How To Install the GitHub App](../how-to/install-github-app.md).

## Steps

### 1. Create Hyaline Configuration
Create a [Hyaline Configuration File](../reference/config.md) named `hyaline.yml` at the root of your repository. You can use the example file below for reference:

```yml
llm:
  provider: github-models
  model: openai/gpt-5
  key: ${GITHUB_TOKEN}

github:
  token: ${GITHUB_TOKEN}

extract:
  source:
    id: try-out-hyaline # This will typically be the name of your repository
  crawler:
    type: fs
    options:
      path: ./
    include:
      - "**/*.md" # Crawl all markdown files
  extractors:
    - type: md
      include:
        - "**/*.md" # Extract all markdown files

check:
  code:
    include: # Modify this section to include source files you want checked
      - "**/*.go"
      - "go.mod"
      - "**/*.js"
      - "**/*.ts"
      - "**/*.tsx"
      - "package.json"
      - "**/*.py"
      - "requirements.txt"
      - "pyproject.toml"
      - "**/*.java"
      - "pom.xml"
      - "*.gradle"
    exclude: # Modify this section to exclude files you do not want checked (tests, fixtures, etc...)
      - "**/*_test.go"
      - "**/*.test.js"
      - "**/*.spec.js"
      - "**/*.test.ts"
      - "**/*.spec.ts"
      - "**/*.test.tsx"
      - "**/*.spec.tsx"
      - "**/test_*.py"
      - "**/*_test.py"
      - "**/*Test.java"
      - "**/Test*.java"
      - "**/test/**"
      - "**/tests/**"
  documentation:
    include:
      - source: "*" # Check against all extracted documentation
  options:
    detectDocumentationUpdates:
      source: try-out-hyaline # This will typically be the name of your repository
```

### 2. Create Workflow
Create a GitHub Action Workflow at `.github/workflows/try-out-hyaline.yml`. You can use the example workflow below:

```yml
name: Try out Hyaline

on:
  pull_request:
    types: [opened, reopened, synchronize, ready_for_review]

jobs:
  check-pr:
    runs-on: ubuntu-latest
    # Only run if PR is NOT a draft
    if: ${{ github.event.pull_request.draft == false }}
    permissions:
      pull-requests: write # To allow Hyaline to post recommendations to the PR
      models: read # To allow Hyaline to use GitHub Models
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Setup Hyaline
        uses: appgardenstudios/hyaline-actions/setup@v1
      - name: Check PR
        uses: appgardenstudios/hyaline-actions/check-pr@v1
        with:
          config: ./hyaline.yml
          repository: ${{ github.repository }}
          pr_number: ${{ github.event.pull_request.number }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 3. Open a Pull Request
Once the Hyaline configuration file and workflow have been merged to your default branch, open a non-draft pull request and you should see Hyaline post a comment to your PR with documentation recommendations. Feel free to play around with the configuration as needed.

## Next Steps
Install the [Hyaline GitHub App](../how-to/install-github-app.md) or read more about [how Hyaline works](../explanation/hyaline.md).