---
title: "Reference: Audit Results"
description: JSON output format for results produced by the audit documentation command
purpose: Detail the JSON output format for results produced by the audit documentation command
sitemap:
  disable: true
---
## Overview
This documents the results produced by the `audit documentation` command.

## JSON Format
```js
{
  "results": [
    {
      "rule": "content-exists-check",
      "description": "Check that backend documentation exists",
      "pass": true,
      "checks": [
        {
          "source": "backend",
          "document": "CHANGELOG.md",
          "uri": "document://backend/CHANGELOG.md",
          "rule": "content-exists-check",
          "check": "CONTENT_EXISTS",
          "pass": true,
          "message": ""
        }
      ]
    },
    {
      "rule": "content-length-check",
      "description": "Check that README has sufficient content",
      "pass": false,
      "checks": [
        {
          "source": "backend",
          "document": "README.md",
          "uri": "document://backend/README.md",
          "rule": "content-length-check",
          "check": "CONTENT_MIN_LENGTH",
          "pass": false,
          "message": "Content length is 277, minimum required is 10000."
        }
      ]
    }
  ]
}
```

### Fields
A list of fields, their types, and a description of each.

| Field | Type | Description |
|-------|------|-------------|
| results | Array | The array of audit rule results |
| results[n] | Object | An audit rule result |
| results[n].rule | String | The rule ID |
| results[n].description | String | The rule description |
| results[n].pass | Boolean | Whether all checks in the rule passed |
| results[n].checks | Array | The array of individual check results |
| results[n].checks[n] | Object | A check result |
| results[n].checks[n].source | String | The documentation source ID |
| results[n].checks[n].document | String | The document ID |
| results[n].checks[n].section | Array OR undefined | If present, the section path |
| results[n].checks[n].section[n] | String | A section path segment |
| results[n].checks[n].uri | String | The document URI |
| results[n].checks[n].rule | String | The rule ID this check belongs to |
| results[n].checks[n].check | String | The type of check performed |
| results[n].checks[n].pass | Boolean | Whether the check passed |
| results[n].checks[n].message | String | The check message (may be empty) |

### Checks
The list of available checks, their associated config property (under `audit.rules[n]`), and a description of each.

| Check | Config Property | Description |
|-------|-----------------|-------------|
| CONTENT_EXISTS | `checks.content.exists` | Verifies that documentation matching the filter exists |
| CONTENT_MIN_LENGTH | `checks.content.min-length` | Checks if content meets the minimum length requirement |
| CONTENT_MATCHES_REGEX | `checks.content.matches-regex` | Validates content against a regular expression pattern |
| CONTENT_MATCHES_PROMPT | `checks.content.matches-prompt` | Uses an LLM to check if content matches a custom prompt |
| CONTENT_MATCHES_PURPOSE | `checks.content.matches-purpose` | Uses an LLM to verify content aligns with its stated purpose |
| PURPOSE_EXISTS | `checks.purpose.exists` | Checks that a purpose is defined for the document or section |
| TAGS_CONTAINS | `checks.tags.contains` | Verifies required tags are present |

Note: When the `CONTENT_EXISTS` check fails to find matching content, the source, document, section, and uri fields will be empty