---
title: "Explanation: Audit"
description: Learn how Hyaline audits documentation for compliance with rules and checks
purpose: Explain how Hyaline audits documentation against configurable rules to ensure compliance
---
## Overview

<div class="portrait">

![Overview](./_img/hyaline-audit.svg)

Hyaline has the ability to audit your current set of documentation based on a set of rules and checks. Using Hyaline, you can ensure that your documentation contains the information necessary to comply with industry regulations or internal compliance rules. You can also check for consistency in documentation across all of your products and systems.

In this example, the current documentation (optionally filtered by specific sources) and the audit rules from your configuration file are provided as inputs to the audit command. The audit process evaluates each rule against the matching documentation. Hyaline uses an AI/LLM to perform qualitative checks, such as whether the documentation matches its intended purpose or a custom provided prompt.

The results of this audit, including the results of all checks that were preformed, are then output as a JSON file.

</div>

## Rules
Hyaline uses audit rules defined in the configuration file to determine what documentation to audit and what checks to perform. You can define multiple rules, and each rule can target specific documentation using filters and apply multiple checks.

<div class="code-example">

```yml
audit:
  rules:
    - id: "content-exists-check"
      description: "Check that content exists and has a purpose"
      documentation:
        - source: "**/*"
          document: "README.md"
      ignore:
        - source: "internal"
      checks:
        content:
          exists: true
          min-length: 100
        purpose:
          exists: true
```
</div>

In the example above, the audit rule `content-exists-check` targets the `README.md` document for all sources except for `internal` and performs three checks: verify the content exists, ensure it meets a minimum length of 100 characters, and confirm a purpose is defined.

## Checks
Hyaline supports several types of checks that can be applied to documentation. These checks fall into three categories: "content checks" for validating documentation structure and content, "purpose checks" for ensuring documentation purposes are correct, and "tag checks" for verifying documents and sections have the right tags.

### Content Checks
Content checks validate the existence, structure, and content of documentation.

#### Content Exists
Verifies that documentation matching the specified filters exists. This check is useful for ensuring required documentation is present.

<div class="code-example">

```yml
audit:
  rules:
    - id: "required-docs"
      description: "Ensure critical docs exist"
      documentation:
        - source: "**/*"
          document: "README.md"
      checks:
        content:
          exists: true
```

</div>

This example checks that a README.md file exists in all sources. The check passes if at least one document matches the filters.

#### Content Min Length
Validates that documentation content meets a minimum length requirement in characters. This helps ensure documentation provides sufficient detail.

<div class="code-example">

```yml
audit:
  rules:
    - id: "sufficient-content"
      description: "Ensure docs are detailed enough"
      documentation:
        - source: "api"
          document: "**/*.md"
      checks:
        content:
          min-length: 500
```

</div>

This example ensures all markdown files in the "api" source contain at least 500 characters. Documentation shorter than this will fail the check.

#### Content Matches Regex
Validates documentation content against a regular expression pattern. This is useful for ensuring specific information or formatting is present.

<div class="code-example">

```yml
audit:
  rules:
    - id: "installation-instructions"
      description: "Ensure install steps are present"
      documentation:
        - source: "**/*"
          document: "README.md"
      checks:
        content:
          matches-regex: "(?i)(install|setup|getting.started)"
```

</div>

This example verifies that README files contain installation-related keywords (case-insensitive). Documentation that doesn't match the regex will fail the check.

#### Content Matches Prompt
Uses an LLM to validate content against a custom prompt or criteria. This provides flexible validation for complex requirements that cannot be expressed with simple rules.

<div class="code-example">

```yml
audit:
  rules:
    - id: "security-guidelines"
      description: "Ensure security best practices"
      documentation:
        - source: "security"
          document: "**/*.md"
      checks:
        content:
          matches-prompt: "Does this document contain specific security guidelines for handling user data?"
```

</div>

This example uses an LLM to evaluate whether security documentation adequately covers user data handling guidelines. The LLM provides a reason for its pass/fail decision.

#### Content Matches Purpose
Uses an LLM to verify that documentation content aligns with its stated purpose. This ensures documentation actually serves its intended function.

<div class="code-example">

```yml
audit:
  rules:
    - id: "purpose-alignment"
      description: "Ensure content matches purpose"
      documentation:
        - source: "**/*"
          document: "**/*.md"
      checks:
        content:
          matches-purpose: true
```

</div>

This example uses an LLM to verify that the actual content of each document aligns with its stated purpose. The LLM provides a reason for its pass/fail decision.

### Purpose Checks
Purpose checks ensure that documentation has a defined purpose. Purposes help ensure documentation serves a clear function and can be maintained effectively.

#### Purpose Exists
Validates that documentation has a defined purpose. 

<div class="code-example">

```yml
audit:
  rules:
    - id: "documented-purposes"
      description: "Ensure all docs have purposes"
      documentation:
        - source: "**/*"
          document: "**/*.md"
      checks:
        purpose:
          exists: true
```

</div>

This example ensures all markdown files have a defined purpose. Documentation without a purpose statement will fail this check.

### Tag Checks
Tag checks validate the presence and values of metadata tags on documentation. Tags are used for categorization, compliance tracking, and filtering, making these checks essential for maintaining organized and compliant documentation.

#### Tags Contains
Validates that required tags are present on documentation. Tag keys and values can be specified as a regex pattern.

<div class="code-example">

```yml
audit:
  rules:
    - id: "compliance-tags"
      description: "Ensure compliance metadata"
      documentation:
        - source: "**/*"
          document: "**/*.md"
      checks:
        tags:
          contains:
            - key: "compliance"
              value: "required"
            - key: "reviewed"
              value: "true"
```

</div>

This example verifies that documents have both a "compliance: required" tag and a "reviewed: true" tag. Both tags must be present for the check to pass.

## Results

Once Hyaline completes the audit, it generates a JSON file containing detailed results for each rule and check. The results provide information about what passed, what failed, and why.

For detailed information about the results schema, see the [Audit Results Reference](../reference/audit-results.md).

## Next Steps
Read more about [merging documentation](./merge.md) or visit the [configuration reference documentation](../reference/config.md).