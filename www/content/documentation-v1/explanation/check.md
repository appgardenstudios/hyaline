---
title: "Explanation: Check"
description: Learn how Hyaline analyzes code changes to identify documentation that needs updating
purpose: Explain how Hyaline checks for what documentation needs to be updated
sitemap:
  disable: true
---
## Overview

<div class="portrait">

![Overview](./_img/check-overview.svg)

TODO portrait image of 2 branches (to get diff), left input of documentation, right input of pr and issues, input of config on lower left, input of llm on lower right, output of recommendations bottom middle. Check process in middle

</div>

## Inputs
The following are all inputs into the Check Process...

### Code Diff
TODO side-by-side config (left) and changed files (right) selected

<div class="side-by-side">

```yml
check:
  code:
    include:
      - "**/*.js"
      - "package.json"
    exclude:
      - "old/**/*"
      - "**/*.test.js"
  ...
```

![Code Diff](./_img/check-code-diff.svg)

</div>

### Documentation
Selectors to get subset of documents

TODO side-by-side config (left) and selected documentation (right)

<div class="side-by-side">

```yml
check:
  ...
  documentation:
    include:
      - source: "**/*"
    exclude:
      - source: my-app
        document: CHANGELOG.md
  ...
```

![Documentation](./_img/check-documentation.svg)

</div>

### PR and Issue Context
TODO side-by-side command (left) and pr/issue info selected (right)

<div class="side-by-side">

```bash
$ hyaline check diff /
  ...
  --pull-request /
    appgardenstudios/hyaline-example/7 /
  --issue /
    appgardenstudios/hyaline-example/2 /
  --issue /
    appgardenstudios/hyaline-example/3 /
  ...
```

![PR and Issue Context](./_img/check-pr-and-issues.svg)

</div>

### UpdateIf Rules
TODO side-by-side config (left) and matched code/docs on right

<div class="side-by-side">

```yaml
check:
  options:
    updateIf:
      touched:
        - code:
            path: "src/routes.js"
          documentation:
            source: "my-app"
            document: "docs/routes.md"
```

![UpdateIf Rules](./_img/check-updateif.svg)

</div>

### LLM Prompt
TODO side-by-side code diff (text diff) on left, input to prompt on right

<div class="side-by-side">

```diff
--- src/server.js
+++ src/server.js
@@ -15,6 +15,9 @@ function serve() {
 }
 
 function isValidUrl(string) {
+  if (!string) {
+    return false;
+  }
   try {
     new URL(string);
     return true;
```

![LLM Prompt](./_img/check-prompt.svg)

</div>


## Check Process

<div class="portrait">

![Check Process](./_img/check-process.svg)

TODO talk about process of looping through list

TODO portrait (square or portrait) image of check process (zoom in on) of looping through diffs to collect llm recommendations and updateIf matches, then merging everything together, then marking what documentation has already been updated.

</div>


## Recommendations

<div class="portrait">

![Recommendations](./_img/check-recommendations.svg)

output to json file or update PR

TODO square image of recommendations output from process to either json file or pr comment

</div>