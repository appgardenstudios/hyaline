---
title: Extract Current
purpose: Explain how Hyaline extracts current documentation, code, and other metadata
---
# Overview
Hyaline has the ability to extract code, documentation, and other metadata into a current data set that can be used to build systems and products as well as verify that existing documentation is accurate and complete.

TODO image of overall extraction process/flow

TODO explanation of image

TODO Explanation of system and how you have multiple sources for system code and documentation

# Extracting Code
System source code is extracted for each defined code source in the configuration (TODO link to config). Code can be extracted using one of two available extractors: `fs` and `git`.

TODO talk about what to extract (source code, not tests and other items)

TODO image of conceptual relationships between system, code, and files

The code that is extracted is placed into a data set that is stored in sqlite. TODO link to reference

## Extracting Code - fs
The `fs` extractor extracts source code from the local file system.

TODO image of extraction from filesystem

TODO explanation of image

## Extracting Code - git
The `git` extractor extracts source code from a local or remote git repository. It supports several different setups that are detailed below.

Note that Hyaline extracts code from a specific branch as specified in the configuration. It does this extraction via the git metadata itself, rather than requiring the repository to be in a specific state. In other words, you don't need to check out the main branch to extract code from it. Hyaline will use the internal git structure to scan and extract the code. 

### Local Repo
TODO image of local repo

In this scenario a local repository already exists on the local file system, and Hyaline uses that repository to extract the code.

### Remote Repo, Cloned Locally
TODO image of remote repo being cloned to the fs

In this scenario Hyaline clones a remote repository down to the local file system, and then uses that local repo to extract code from. 

### Remote Repo, Cloned In Memory
TODO image of remote repo being cloned to an in-memory fs

In this scenario Hyaline clones a remote repository into a local in-memory filesystem, and then uses that in-memory repository to extract code from.

# Extracting Documentation
System documentation is extracted for each defined documentation source in the configuration (TODO link to config). Documentation can be extracted using one of three available extractors: `fs`, `git`, and `http`.

TODO discuss non-markdown to markdown conversion and selector for html sources

TODO talk about what to extract (documentation, can ignore items as needed)

TODO image of conceptual relationships between system, documentation, document, and section.

TODO talk about how sections are extracted.

The documentation that is extracted is placed into a data set that is stored in sqlite. TODO link to reference

## Extracting Documentation - fs
The `fs` extractor extracts documentation from the local file system. It operates the same way as the Code `fs` extractor (See Extracting Code - fs).

TODO image 

## Extracting Documentation - git
The `git` extractor extracts documentation from a local or remote git repository. It operates the same way as the Code `git` extractor and supports the same set of setups (See Extracting Code - git).

TODO image 

## Extracting Documentation - http
The `http` extractor extracts documentation from an http(s) server via crawling.

TODO image of crawling

TODO explanation of image

# Extracting Metadata
Hyaline will be extended to extract additional organizational metadata in the future. As of now, Hyaline only supports extracting pull request and issue information when extracting changes (TODO link to extract-change)

# Next Steps
You can continue on to see how [Hyaline extracts change information](./extract-change.md), or see how Hyaline can [merge together data sets](./merge.md).