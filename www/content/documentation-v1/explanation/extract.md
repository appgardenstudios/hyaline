---
title: "Explanation: Extract Documentation"
description: Learn how Hyaline extracts documentation from various sources into a unified current data set
purpose: Explain how Hyaline extracts documentation
---
## Overview

<div class="portrait">

![Overview](./_img/extract-documentation-overview.svg)
TODO portrait image of multiple repos/sites, extract into current data set, can audit (right), used for check on changes (left), used by AI (mcp lower left), used by org (lower right)

Hyaline has the ability to extract documentation into a current data set that can be used to build systems and products as well as verify that the documentation is accurate and complete.

In this example you can see documentation spread over multiple repositories and documentation sites. Hyaline can extract documentation from each of these, (optionally) [merge](./merge.md) them together into a unified documentation set, and then use Hyaline to [check](./check.md) and [audit](./audit.md) the extracted documentation or use it via an [MCP server](./mcp.md) or referencing the [current data set](../reference/data-set.md). 

In Hyaline repository or documentation site is a documentation source, or source for short. 

</div>

<div class="portrait">

![Extract Phases](./_img/extract-documentation-phases.svg)
TODO square image of crawl -> extract -> add metadata w/ nothing highlighted.

Extracting documentation is broken up into 3 phases: Crawling, Extracting, and Adding Metadata

TODO talk about tags

</div>

## Crawling Documentation

<div class="portrait">

![Crawling Documentation](./_img/extract-documentation-crawling.svg)
TODO square image of crawl -> extract -> add metadata w/ crawl highlighted.

Hyaline can be configured to crawl a documentation source and extract documentation. Hyaline supports a number of different crawlers, each with their own capabilities and configuration.

</div>

### Crawling Documentation - fs

<div class="portrait">

![Crawling Documentation fs](./_img/extract-documentation-fs.svg)
TODO square show config and directory structure, highlight files that were crawled

The `fs` crawler crawls a local file system starting at a path, and processes each document it encounters.

In this example you can see that Hyaline is configured to start in the TODO directory and process any documents that match TODO. It also does not process any documents in TODO as that directory and its children are excluded from the crawl.

</div>

### Crawling Documentation - git

<div class="portrait">

![Crawling Documentation git](./_img/extract-documentation-git.svg)
TODO square show config and repo structure, highlight files that were crawled

The `git` crawler crawls a git repository starting at its root, and processes each document it encounters.

In this example you can see that Hyaline is configured to clone the remote repo TODO into memory and process any documents that match TODO. It also does not process any documents in TODO as that directory and its children are excluded from the crawl.

</div>

### Crawling Documentation - http

<div class="portrait">

![Crawling Documentation http](./_img/extract-documentation-http.svg)
TODO square show config and site structure, highlight files that were crawled

The `http` crawler crawls a HTTP or HTTPS website starting at a configured starting url, and processes each document it encounters.

In this example you can see that Hyaline is configured to start crawling at TODO and process any documents that match TODO. It also does not process any documents in TODO as that directory and its children are excluded from the crawl.

Note that Hyaline will not crawl outside of the specified domain, so you don't need to worry about it getting lost in the internet.

</div>

## Extracting Documentation

<div class="portrait">

![Extracting Documentation](./_img/extract-documentation-extracting.svg)
TODO square image of crawl -> extract -> add metadata w/ extract highlighted.

Hyaline can be configured to extract documentation differently based on the type of documentation encountered. Hyaline supports a number of different extractors, each with their own capabilities and configuration.

</div>

### Extracting Documentation - md

<div class="portrait">

![Markdown](./_img/extract-documentation-markdown.svg)
TODO square image with markdown on the left, extracted document and child sections on the right (show sub-section example)

The `markdown` extractor extracts markdown documents.

In this example you can see a markdown document being extracted into a document and its sections.

</div>

### Extracting Documentation - html

<div class="portrait">

![HTML to Markdown](./_img/extract-documentation-html-to-markdown.svg)
TODO square image with html on the left, extracted document and child sections on the right (show sub-section example). Highlight selected area (main) with html tag

The `html` extractor extracts html documents by extracting the content of the documentation and transforming it into markdown.

In this example you can see an html document being extracted into a document and its sections. Hyaline is configured to select just the html in the `main` tag, which is then transformed into markdown and stored as a document and sections.

</div>

### A Note on Sections

<div class="portrait">

![Document Sections](./_img/extract-documentation-sections.svg)
TODO square image with markdown doc on left w/ 3 level section(s), extracted sections on the right with names and parents

TODO Hyaline scans the markdown document and extracts any sections it encounters. It identifies each section by name, and preserves any section level hierarchy it find.

Note that when storing the ID of the section it replaces any "/" characters with "_", as Hyaline uses "/" to separate sections in the ID if a sub-section (i.e. `Section1/Section1.1`).

</div>

## Adding Metadata

<div class="portrait">

![Adding Metadata](./_img/extract-documentation-metadata.svg)
TODO square image of crawl -> extract -> add metadata w/ add metadata highlighted.

Hyaline can be configured to add tags and purposes to each document and section that is extracted.

In this example you can see a set of documents that have been extracted. Based on the configuration TODO document and sections are tagged with TODO, and TODO section is tagged with TODO. The TODO document also has a purpose associated with it to help Hyaline ensure that it is updated when it needs to be.

</div>

## Next Steps
TODO
