---
title: "Reference: Recommendations"
description: JSON output format for recommendations produced by the check diff and check pr commands
purpose: Detail the JSON output format for recommendations produced by the check diff and check pr commands
---
## Overview
This documents the recommendations produced by the `check diff` and `check pr` commands.

## JSON Format
```js
{
  "recommendations": [
    {
      "documentationSource": "my-app",
      "document": "README.md",
      "section": [
        "Running Locally"
      ],
      "recommendation": "Consider reviewing and updating this documentation",
      "reasons": [
        {
          "reason": "Update this section if any files matching package.json were modified",
          "outdated": false,
          "check": {
            "type": "UPDATE_IF_MODIFIED",
            "file": "package.json",
            "contextHash": "a1b2c3d4"
          }
        }
      ],
      "changed": true,
      "checked": true,
      "outdated": false
    }
  ],
  "head": "b4c5c736fd31d30a04067af9c0929d7dc42f049e",
  "base": "b564300250288b332d50e2925dbd25e98831adbd"
}
```

### Fields
A list of fields, their types, and a description of each.

| Field | Type | Description |
|-------|------|-------------|
| recommendations | Array | The array of recommendations |
| recommendations[n] | Object | A recommendation |
| recommendations[n].documentationSource | String | The documentation source ID |
| recommendations[n].document | String | The name of the document |
| recommendations[n].section | Array | If not empty, the section (including parent sections) |
| recommendations[n].section[n] | String | The section name |
| recommendations[n].recommendation | String | The recommendation |
| recommendations[n].reasons | Array | A list of reasons |
| recommendations[n].reasons[n] | Object | A reason |
| recommendations[n].reasons[n].reason | String | The human-readable reason |
| recommendations[n].reasons[n].outdated | Boolean | Whether this reason is outdated (e.g. due to subsequent code changes) |
| recommendations[n].reasons[n].check | Object | Information about the check that generated this reason |
| recommendations[n].reasons[n].check.type | String | The type of check performed |
| recommendations[n].reasons[n].check.file | String | The file that triggered this check |
| recommendations[n].reasons[n].check.contextHash | String | A hash representing the context when this check was performed |
| recommendations[n].changed | Boolean | If the document or section was changed in the diff |
| recommendations[n].checked | Boolean | If the recommendation has been checked (such as by updating the recommended document or section) |
| recommendations[n].outdated | Boolean | Whether this entire recommendation is outdated (true if all reasons are outdated) |
| head | String | The commit hash used as the head reference in the diff |
| base | String | The commit hash used as the base reference in the diff |

### Check Types
The list of available check types, their associated config property, what goes into the context hash, and a description of each.

| Check Type | Config Property | Context Hash | Description |
|------------|-----------------|--------------|-------------|
| `LLM` | N/A (automatic) | Hash of the prompt | Uses an LLM to analyze code changes and determine if documentation should be updated |
| `UPDATE_IF_TOUCHED` | `check.options.updateIf.touched` | Hash of the check type | Triggers when a file matching the configured pattern is modified in any way (added, modified, deleted, or renamed) |
| `UPDATE_IF_ADDED` | `check.options.updateIf.added` | Hash of the check type | Triggers when a file matching the configured pattern is added |
| `UPDATE_IF_MODIFIED` | `check.options.updateIf.modified` | Hash of the check type | Triggers when a file matching the configured pattern is modified |
| `UPDATE_IF_DELETED` | `check.options.updateIf.deleted` | Hash of the check type | Triggers when a file matching the configured pattern is deleted |
| `UPDATE_IF_RENAMED` | `check.options.updateIf.renamed` | Hash of the check type | Triggers when a file matching the configured pattern is renamed |