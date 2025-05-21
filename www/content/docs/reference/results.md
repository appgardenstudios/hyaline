---
title: Hyaline Check Results
purpose: Detail the output results for the check current and check change commands
---
# Overview
This documents the results produced by the `check current` and `check change` commands.

## Check Current

```js
{
  "results": [
    {
      "system": "check-current",
      "documentationSource": "app",
      "document": "README.md",
      "check": "COMPLETE",
      "result": "ERROR",
      "message": "This document is not complete. "
    },
    {
      "system": "check-current",
      "documentationSource": "app",
      "document": "README.md",
      "check": "DESIRED_DOCUMENT_EXISTS",
      "result": "PASS",
      "message": ""
    },
    {
      "system": "check-current",
      "documentationSource": "app",
      "document": "README.md",
      "check": "MATCHES_PURPOSE",
      "result": "ERROR",
      "message": "This document does not match it's purpose. "
    },
    {
      "system": "check-current",
      "documentationSource": "app",
      "document": "README.md",
      "check": "REQUIRED",
      "result": "PASS",
      "message": ""
    },

  ]
}
```

TODO discuss / show example with section


## Check Change

```js
{
  "recommendations": [
    {
      "system": "check-change",
      "documentationSource": "app",
      "document": "README.md",
      "section": [
        "Running Locally"
      ],
      "recommendation": "Consider reviewing and updating this documentation",
      "reasons": [
        "Update this section if any files matching package.json were modified"
      ],
      "changed": true
    }
  ]
}
```

TODO discuss / show example with section