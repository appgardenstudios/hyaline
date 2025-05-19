---
title: Extract Change
purpose: Explain how Hyaline extracts changed documentation, code, and other metadata
---
# Overview
Hyaline has the ability to extract the set of changed code and documentation, along with other change metadata, into a change data set that can be used to check for needed documentation updates based on the changes made.

TODO image of overall flow

TODO description of image above

The main unit of organization within Hyaline is the system. A system can contain any number of code and/or documentation sources. When using Hyaline, it is helpful to create multiple focused, single purpose systems rather than a single system with everything in it. Also note that there is no restriction on where the code and/or documentation of a system comes from, meaning that you can break up a mono-repo into multiple, smaller systems or piece together a system from code and documentation spread across a large number of repositories and sites.

TODO talk about when extracting change(s), how to target a single system, or even target just a single or the small set of code/documentation sources that actually changed. The goal is to extract only what changed into an addressable unit with all applicable data and metadata about the change.

# Extracting Changed Code
System source code that changed is extracted for each targeted code source in the system. Note that the code source must be configured to use the `git` extractor for change extraction to work, as Hyaline compares two branches to extract the diffs used when extracting the change.

The extraction process uses the same configuration as the extract current process does, so if you haven't read up on how [extract current](./extract-current.md) works it would be helpful to do so now.

TODO image of change extraction at the file level, based on changed, include, and exclude

TODO explanation of the image

TODO link to data set documentation

# Extracting Changed Documentation
System documentation that changed is extracted for each targeted documentation source in the system. Note that the documentation source must be configured to use the `git` extractor for change extraction to work, as Hyaline compares two branches to extract the diffs used when extracting the change.

The extraction process uses the same configuration as the extract current process does, so if you haven't read up on how [extract current](./extract-current.md) works it would be helpful to do so now.

TODO image of change extraction at the file level, based on changed, include, and exclude

TODO explanation of the image

TODO link to data set documentation

# Extracting Metadata
Hyaline also supports extracting additional metadata and context about the change, such as any pull request or issue information available.

## Pull Request
Hyaline supports extracting the title and contents of a GitHub pull request and the inclusion of that information in the change data set.

TODO image of pull request, extraction, and placement of that information within the data set associated with the system

TODO explanation of the image above.

## Issues
Hyaline supports extracting the title and contents of one or more issues and that inclusion of that information in the change data set.

TODO image of issue(s), extraction, and placement of that info within the data set (linked to the system)

# Next Steps
You can continue on to see how Hyaline checks [current](./check-current.md) or [change](./check-change.md) data sets.