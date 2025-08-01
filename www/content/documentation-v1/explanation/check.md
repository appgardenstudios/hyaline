---
title: "Explanation: Check"
description: Learn how Hyaline analyzes code changes to identify documentation that needs updating
purpose: Explain how Hyaline checks for what documentation needs to be updated
---
## Overview
TODO

TODO portrait image of 2 branches (to get diff), left input of documentation, right input of pr and issues, input of config on lower left, input of llm on lower right, output of recommendations bottom middle. Check process in middle

## Inputs
The following are all inputs into the Check Process...

### Code Diff

TODO side-by-side config (left) and changed files (right) selected

### Documentation
Selectors to get subset of documents

TODO side-by-side config (left) and selected documentation (right)

### PR and Issue Context

TODO side-by-side command (left) and pr/issue info selected (right)

### UpdateIf Rules

TODO side-by-side config (left) and matched code/docs on right

### LLM Prompt

TODO side-by-side code diff (text diff) on left, input to prompt on right

## Check Process
TODO talk about process of looping through list

TODO portrait (square or portrait) image of check process (zoom in on) of looping through diffs to collect llm recommendations and updateIf matches, then merging everything together, then marking what documentation has already been updated.

## Recommendations

output to json file or update PR

TODO square image of recommendations output from process to either json file or pr comment