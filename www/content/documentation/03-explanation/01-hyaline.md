---
title: "Explanation: Hyaline"
linkTitle: Hyaline
purpose: Explain the overall data flow and architecture of Hyaline
url: documentation/explanation/hyaline
---
## Introduction
Hyaline is intended to help software development teams use their documentation to build and maintain their products. To that end Hyaline has two primary objectives: 1) help teams create, update, and maintain their documentation so that they can 2) use their documentation to create, maintain, and ship their products.

## Use Cases
To help you understand a bit more about what Hyaline does and does not do consider the following use cases that Hyaline is intended to support:
* Identify documentation that needs to be updated when a software change is implemented - Hyaline can examine your code changes and help identify what documentation needs to be updated and why
* Identify documentation that does not match its intended purpose or is incomplete - Hyaline can scan your existing documentation and ensure that each document and section matches your intended purpose and is complete.
* Ensure that certain documentation, including documentation that is required for regulatory or compliance purposes, is consistently created and maintained - Hyaline can be configured to use centralized rules to ensure certain documentation is present and updated across all of your products and systems.
* Allow an LLM to scan, search, and use your documentation to help you build your product(s) - Hyaline extracts and indexes all of your documentation, and makes that information available through an [MCP server](https://modelcontextprotocol.io).

And now consider the following use cases that Hyaline is not intended to support:
* Creating and/or updating documentation without human involvement - Hyaline is intended to _augment_ team members, not replace them.
* Extracting or storing 3rd party documentation such as API/Library documentation.

## Concepts
The following concepts are important to understand before talking about Hyaline.

### Product
A product is any discrete set of features delivered to a set of users. It is defined recursively so that a product may be made up of one or more products, and a product may belong to a larger product.

### System
A system is any discrete set of components, either built in-house, purchased, or otherwise used by an organization to deliver one or more products. It is defined recursively such that a system can be composed of one or more systems, and a system may belong to a larger system.

### Person/People
A person is any individual within the organization. People is a discrete set of persons within an organization.

### Process
A process is any discrete set of tasks that are completed to accomplish a goal. Processes are used by people to build and maintain products and systems.

### Team
A team is a set of people that use processes to build, maintain, and ship one or more products and/or systems.

![Team](/documentation/03-explanation/_img/hyaline-team.svg)

## Workflow
Hyaline is built to support the conceptual workflow as follows:

![Team](/documentation/03-explanation/_img/hyaline-workflow.svg)

In this workflow people, assisted by AI, build products and systems. While doing so they create and update Documentation. That documentation is then read and used by people to build products and systems, and the cycle continues. Hyaline sits in between People/AI and Documentation, and is intended to assist in both creating/updating documentation and reading/using documentation to build products and systems.

## Use Documentation
Zooming in a bit on the People/AI using Documentation part of the diagram above, we see the following:

![Team](/documentation/03-explanation/_img/hyaline-use.svg)

In the image above you can see the Hyaline extracts code, documentation, and other metadata from a variety of sources. Each system can have multiple code and documentation sources, and Hyaline supports an unlimited number of systems. In the future Hyaline will also support pulling product, process, and team documentation and metadata.

The result of this extraction is a current data set containing the organization's code, documentation, and other metadata. This current data set is stored as a single sqlite database, making it useful for Hyaline as well as other organizational uses.

Hyaline then has the ability to use this current data set to provide an MCP server. This MCP server provides a set of tools that can be used by AI assistants to read, search, and enumerate documentation and other organizational metadata. That means that you can setup Hyaline to extract all of your organization's documentation and provide an MCP server to your organization that always has the latest internal documentation and metadata. This MCP server can then help your people build and maintain your products and systems.

Hyaline also has the ability to scan the current set of an organizations documentation and verify that it is complete and meets its stated purpose. It can then make those results available to people and AI to act on those results.

## Update Documentation
Zooming in a bit on the People/AI updating Documentation part of the diagram above, we see the following:

![Team](/documentation/03-explanation/_img/hyaline-update.svg)

In the image above you can see that Hyaline looks at system changes to determine what documentation should be updated. It does this by extracting the change and associated metadata (tickets, issues, PRs, etc.) and then determining what documentation in your current data set needs to be updated. It then makes those results available in your existing processes (usually via a comment on your PR) and helps both people and AI make the appropriate changes to the appropriate documentation both inside and outside the repository.

## Components
The following items are components of Hyaline:

![Team](/documentation/03-explanation/_img/hyaline-components.svg)

## Documentation
Your organization's documentation. It can be plain text, markdown, or html, and can be extracted from a file system, git repo, or HTTP(s) server.

## Code
Your organization's code. It can be extracted from a file system or git repo.

## Metadata
Your organization's metadata, such as issues and pull requests. They can be extracted from GitHub at the moment, but there are plans to support a wider range of sources in the future.

## System
A conceptual boundary or unit that can contain any number of code and documentation sources. A system is the primary unit of operation within Hyaline, and there are no limits to the number of systems you can have. It is up to you and your team to split up your code and documentation in a way that makes the most sense to you.

## Extract
The process of extracting code, documentation, and other metadata from their respective sources. This extraction process results in a Data Set stored in an SQLite database that can be used by Hyaline or other processes that need access to the code, documentation, and other metadata. For more information on the extraction process please see [hyaline extract current]({{< relref "/documentation/04-reference/02-cli.md#extract-current" >}}) and/or [hyaline extract change]({{< relref "/documentation/04-reference/02-cli.md#extract-change" >}})

## Data Set (Current & Change)
The result of an extraction, the Data Set holds all of the extracted code, documentation, and other metadata. It consists of a single SQLite database. Please see the [Data Set reference]({{< relref "/documentation/04-reference/03-data-set.md" >}}) for more information.

## Config
The configuration that Hyaline uses. It is currently supplied via a yaml file. Please see the [Config reference]({{< relref "/documentation/04-reference/01-config.md" >}}) for more information.

## Check
The process of checking either a specific change or an entire set of documentation for issues, recommendations, or suggestions. This produces one or more results with information specific to the actual check that was performed. Please see the documentation for each process linked below:

* [hyaline check current]({{< relref "/documentation/04-reference/02-cli.md#check-current" >}})
* [hyaline check change]({{< relref "/documentation/04-reference/02-cli.md#check-change" >}})

## Result(s)
The outcome of a check, the result holds actionable information on what was found and how to address it. Please see the links to the reference documentation for more information:

* [hyaline check current results]({{< relref "/documentation/04-reference/04-results.md#check-current" >}})
* [hyaline check change results]({{< relref "/documentation/04-reference/04-results.md#check-change" >}})

## Next Steps
Continue reading about various Hyaline concepts such as [extract current]({{< relref "/documentation/03-explanation/02-extract-current.md" >}}), or get started by visiting [how to install the cli]({{< relref "/documentation/02-how-to/01-install-cli.md" >}}).