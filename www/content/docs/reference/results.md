---
title: Hyaline Check Results
purpose: Detail the output results for the check current and check change commands
---
# Overview
This documents the results produced by the `check current` and `check change` commands.

# Check Current
```js
{
  "results": [
    {
      "system": "my-app",
      "documentationSource": "backend",
      "document": "README.md",
      "section": [
        "Running Locally",
      ],
      "check": "COMPLETE",
      "result": "ERROR",
      "message": "This section is not complete. While it does contain an explanation of how to run the app locally, it does not contain an example as required in the stated purpose."
    },
    {
      "system": "my-app",
      "documentationSource": "backend",
      "document": "README.md",
      "section": [
        "Running Locally",
      ],
      "check": "DESIRED_DOCUMENT_EXISTS",
      "result": "PASS",
      "message": ""
    },
    {
      "system": "my-app",
      "documentationSource": "backend",
      "document": "README.md",
      "section": [
        "Running Locally",
      ],
      "check": "MATCHES_PURPOSE",
      "result": "PASS",
      "message": ""
    },
    {
      "system": "my-app",
      "documentationSource": "backend",
      "document": "README.md",
      "section": [
        "Running Locally",
      ],
      "check": "REQUIRED",
      "result": "PASS",
      "message": ""
    },

  ]
}
```

## Fields
A list of fields, their types, and a description of each.

| Field | Type | Description |
|-------|------|-------------|
| results | Array | The array of results |
| results[n] | Object | A result |
| results[n].system | String | The system ID |
| results[n].documentationSource | String | The documentation source ID |
| results[n].document | String | The name of the document |
| results[n].section | Array OR undefined | If present, the section (including parent sections) |
| results[n].section[n] | String | The section name |
| results[n].check | String | The check being run |
| results[n].result | String | The result of the check |
| results[n].message | String | The message (may be an empty string) |

## Checks
The list of available checks, what cli option is required to perform the chck (if any), and a description of the check

| Check | CLI Option | Description |
|-------|------------|-------------|
| COMPLETE | --check-completeness | If the document or section contents are complete |
| DESIRED_DOCUMENT_EXISTS | (none) | If there is a corresponding document or section in the configuration |
| MATCHES_PURPOSE | --check-purpose | If the document or section contents match the stated purpose in the config |
| REQUIRED | (none) | If the document or section is present if marked as required in the config |

## Results
The list of possible results and a description of the result

| Check | Description |
|-------|-------------|
| ERROR | The check did not pass |
| PASS | The check passed |
| SKIPPED | The check was skipped due to the document or section being ignored |
| WARN | The document or section does not have a stated purpose |

# Check Change

```js
{
  "recommendations": [
    {
      "system": "my-app",
      "documentationSource": "backend",
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

## Fields
A list of fields, their types, and a description of each.

| Field | Type | Description |
|-------|------|-------------|
| recommendations | Array | The array of recommendations |
| recommendations[n] | Object | A result |
| recommendations[n].system | String | The system ID |
| recommendations[n].documentationSource | String | The documentation source ID |
| recommendations[n].document | String | The name of the document |
| recommendations[n].section | Array OR undefined | If present, the section (including parent sections) |
| recommendations[n].section[n] | String | The section name |
| recommendations[n].recommendation | String | The recommendation |
| recommendations[n].reasons | Array | A list of reasons |
| recommendations[n].reasons[n] | String | A reason |
| recommendations[n].changed | Boolean | If the document or section was marked as changed in the change data set |
