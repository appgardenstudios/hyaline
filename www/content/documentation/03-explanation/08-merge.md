---
title: "Explanation: Merge"
linkTitle: Merge
purpose: Explain how Hyaline merge works
url: documentation/explanation/merge
---
## Overview
Hyaline has the ability to merge data sets together. This can be used to create a single data set containing all current documentation, or to merge newly extracted data into an existing data set (for example when extracting to prepare for checking a change).

![Overview](/documentation/03-explanation/_img/merge-overview.svg)

Merging happens at the system level. If a system being merged in does not yet exist in the data set the system is pulled in wholesale. IF a system exists, code sources, documentation sources, changes, and tasks are added if they do not exists or overwritten if they do.

## Example
![Example](/documentation/03-explanation/_img/merge-example.svg)

For example, if we have 2 input data sets that have the same system, we merge the systems. In this example The system in Input 1 has two Code Sources (1 and 2) and the system in Input 2 has only Code Source 2. When merged, the system will contain both Code Sources. Since Code Source 1 is not present in input 2, Code Source 1 and its files are pulled into the output directly from Code Source 1. And since Code Source 2 is present in both inputs, Hyaline pulls Code Source 2 and its associated files into the output from Input 2 (and nothing from Code Source 2 in Input 1 is pulled into the output).