---
title: GitHub Actions
purpose: Detail the functionality and usage of Hyaline within GitHub Actions
---
# Overview
Hyaline provides a set of [GitHub Actions](https://github.com/features/actions) that allows you to setup and use Hyaline within a GitHub workflow.

TODO image of overall flow

TODO explanation of image

# Setup
The [setup action](https://github.com/appgardenstudios/hyaline-actions/tree/main/setup) provides an easy way to download and install the Hyaline CLI on your GitHub Actions runner so you can run Hyaline commands.

The default configuration installs a hard-coded version of Hyaline that is updated alongside major Hyaline releases:
```yaml
steps:
  - uses: appgardenstudios/hyaline-actions/setup@v0
```

A specific version of the Hyaline CLI can be installed using:
```yaml
steps:
  - uses: appgardenstudios/hyaline-actions/setup@v0
    with:
      version: "YYYY-MM-DD-HASH"
```

# Check PR
The [check-pr action](https://github.com/appgardenstudios/hyaline-actions/tree/main/check-pr) provides the ability to check a pull request.

TODO image of the steps check PR does (extract current, extract change, check change, update pr)

TODO explanation of image

TODO Note about how Hyaline will update a comment if it already exists

For example, you can configure Hyaline to check your PR using the following workflow:
```yaml
on:
  pull_request:
    types: [opened, reopened, synchronize, ready_for_review]

jobs:
  check-pr:
    runs-on: ubuntu-latest
    # Only run if PR is NOT a draft
    if: ${{ github.event.pull_request.draft == false }}
    permissions:
      pull-requests: write
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Setup Hyaline
        uses: appgardenstudios/hyaline-actions/setup@v0
      - name: Check PR
        uses: appgardenstudios/hyaline-actions/check-pr@v0
        with:
          config: ./hyaline.yml
          system: my-app
          repository: ${{ github.repository }}
          pr_number: ${{ github.event.pull_request.number }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
        env:
          # Set env vars needed by the hyaline CLI when interpolating the hyaline config
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          ANTHROPIC_KEY: ${{ secrets.ANTHROPIC_KEY }}
```

Note that `check-pr` requires the permission `pull-requests: write` to leave a comment on the pull request.
