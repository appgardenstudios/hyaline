---
title: Check Change
purpose: Explain how Hyaline checks a change for what documentation needs to be updated
---
# Overview
Hyaline has the ability to check the current set of documentation against a changed set of code and documentation to determine which pieces of documentation (if any) need to be updated. The goal is to identify each piece of documentation that could be affected and present that list to the human making the change. Hyaline also has the ability to call out to an LLM to generate a suggested change to each piece of documentation if desired.

TODO image of flow, show extract current, extract change, maybe merge?

Explanation of image

TODO Link to Extract Current and Extract Change

TODO link to output results reference

# Recommendations
Hyaline loops through each code change and conceptually examines the links between the code and documentation and asks an LLM "based on this change and the associated metadata, what documentation should be updated and why?". Hyaline then compiles the set of documentation that needs to be updated into a unified list of recommendations and presents those results to the human(s) that made the change.

## Directly via updateIf
TODO image of updateIf and example glob match

Hyaline's configuration has the ability to express direct relationships between documentation and code via a set of updateIf statements. These statements direct Hyaline to mark documentation as needing an update if any code matching a glob is updated in a certain way in the change. For example:

For a full list of what configuration options are available please visit the [configuration reference](../reference/config.md).

## Indirectly via LLM
TODO image of indirectly via llm

Hyaline will call out to an LLM to determine which documentation (if any) should be updated for each specific change. To do that Hyaline formats the following information and includes it as context in the LLM call:

* The diff of the changed file
* A list of system documentation (including each document and section's purpose)
* The pull request (if available)
* The list of related issues (if available)

The LLM then responds with an indication of the set of documents and/or sections that should be updated or with an indication that no updates are needed based on the supplied information

## Compilation of results
TODO image of compilation of results

The results of the updateIfs and LLM calls are then compiled into a single list of documents and sections needing to be looked at. If a document or section was identified as needing to be updated for more than one reason, a list of reasons for that update are returned.

# Suggestions
TODO image of suggestions

If configured to do so, Hyaline will take the list of recommendation results and ask an LLM what updates should be made to each document or section. To do that Hyaline formats the following information and includes it as context in the LLM call:

* The set of diffs that were identified by the LLM as relating to this update
* The pull request (if available)
* The list of related issues (if available)
* The contents of the document or section that needs to be updated.

The LLM responds with the suggested update to the document or section, and that suggestion is then added to the list of recommendation results that are returned.

# Next Steps
You can continue on to see how Hyaline can [generate a configuration](./generate-config.md) or [merge data sets](./merge.md), or see how Hyaline extracts [current](./extract-current.md) or [change](./extract-change.md) data sets.
