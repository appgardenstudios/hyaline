---
title: "Reference: Recommendations"
description: JSON output format for recommendations produced by the check diff and check pr commands
purpose: Detail the JSON output format for recommendations produced by the check diff and check pr commands
sitemap:
  disable: true
---
## Overview
This documents the recommendations produced by the `check diff` and `check pr` commands.

## JSON Format
```js
{
  "recommendations": [
    {
      "source": "my-app",
      "document": "README.md",
      "section": [
        "Running Locally"
      ],
      "recommendation": "Consider reviewing and updating this documentation",
      "reasons": [
        {
          "reason": "Update this section if any files matching package.json were modified"
        }
      ],
      "changed": true,
      "checked" true
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
| recommendations[n] | Object | A result |
| recommendations[n].source | String | The documentation source ID |
| recommendations[n].document | String | The name of the document |
| recommendations[n].section | Array | If not empty, the section (including parent sections) |
| recommendations[n].section[n] | String | The section name |
| recommendations[n].recommendation | String | The recommendation |
| recommendations[n].reasons | Array | A list of reasons |
| recommendations[n].reasons[n] | Object | A reason |
| recommendations[n].reasons[n].reason | String | The human-readable reason |
| recommendations[n].changed | Boolean | If the document or section was changed in the diff |
| recommendations[n].checked | Boolean | If the recommendation has been checked (such as by updating the recommended document or section) |
| head | String | The commit hash used as the head reference in the diff |
| base | String | The commit hash used as the base reference in the diff |