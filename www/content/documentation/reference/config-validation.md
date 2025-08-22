---
title: "Reference: Config Validation"
description: JSON output format for output produced by the validate config command
purpose: Detail the JSON output format for output produced by the validate config command
---
## Overview
This documents the output produced by the `validate config` command.

## JSON Format
```js
{
  "valid": true,
  "error": "",
  "detail": {
    "llm": {
      "present": true,
      "valid": true,
      "error": ""
    },
    "github": {
      "present": true,
      "valid": true,
      "error": ""
    },
    "extract": {
      "present": true,
      "disabled": true,
      "valid": true,
      "error": ""
    },
    "check": {
      "present": true,
      "disabled": true,
      "valid": true,
      "error": ""
    },
    "audit": {
      "present": true,
      "disabled": true,
      "valid": true,
      "error": ""
    }
  }
}
```

### Fields
A list of fields, their types, and a description of each.

| Field | Type | Description |
|-------|------|-------------|
| valid | Boolean | True if the entire configuration is valid |
| error | String | The validation error, if any |
| detail | Object | Details about each section of the config |
| detail.llm | Object | Details about the llm section |
| detail.llm.present | Boolean | True if the llm.provider is set |
| detail.llm.valid | Boolean | True if the llm section is valid |
| detail.llm.error | String | The validation error for the llm section, if any |
| detail.github | Object | Details about the github section |
| detail.github.present | Boolean | True if the github.token is set |
| detail.github.valid | Boolean | True if the github section is valid |
| detail.github.error | String | The validation error for the github section, if any |
| detail.extract | Object | Details about the extract section |
| detail.extract.present | Boolean | True if the extract section is present |
| detail.extract.disabled | Boolean | True if extract.disabled is true |
| detail.extract.valid | Boolean | True if the extract section is valid |
| detail.extract.error | String | The validation error for the extract section, if any |
| detail.check | Object | Details about the check section |
| detail.check.present | Boolean | True if the check section is present |
| detail.check.disabled | Boolean | True if check.disabled is true |
| detail.check.valid | Boolean | True if the check section is valid |
| detail.check.error | String | The validation error for the check section, if any |
| detail.audit | Object | Details about the audit section |
| detail.audit.present | Boolean | True if the audit section is present |
| detail.audit.disabled | Boolean | True if audit.disabled is true |
| detail.audit.valid | Boolean | True if the audit section is valid |
| detail.audit.error | String | The validation error for the audit section, if any |