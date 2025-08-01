---
title: "Reference: Recommendations"
description: JSON output format for recommendations produced by the check diff command
purpose: Detail the JSON output format for recommendations produced by the check diff command
---
## Overview
This documents the recommendations produced by the `check diff` command.

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
        "Update this section if any files matching package.json were modified"
      ],
      "changed": true
    }
  ]
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
| recommendations[n].section | Array OR undefined | If present, the section (including parent sections) |
| recommendations[n].section[n] | String | The section name |
| recommendations[n].recommendation | String | The recommendation |
| recommendations[n].reasons | Array | A list of reasons |
| recommendations[n].reasons[n] | String | A reason |
| recommendations[n].changed | Boolean | If the document or section was changed in the diff |