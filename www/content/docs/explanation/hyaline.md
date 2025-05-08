---
title: Hyaline
purpose: Explain the overall data flow and architecture of Hyaline
---
# Introduction
Hyaline is intended to help software development teams establish? and use documentation to build and maintain their products. To that end Hyaline has two primary objectives: 1) help teams create, update, and maintain their documentation so that they can 2) use their documentation to create, maintain, and ship their products.

## Use Cases
To help you understand a bit more about what Hyaline does and does not do consider the following use cases that Hyaline is intended to support:
* Identify documentation that needs to be updated when a software change is implemented - Hyaline can examine your code changes and help identify what documentation needs to be updated and why
* Identify documentation that does not match its intended purpose or is incomplete - Hyaline can scan your existing documentation and ensure that each document and section matches your intended purpose and is complete.
* Ensure that certain documentation, including documentation that is required for regulatory or compliance purposes, is consistently created and maintained - Hyaline can be configured to use centralized rules to ensure certain documentation is present and updated across all of your products and systems.
* Allow an LLM to scan, search, and use your documentation to help you build your product(s) - Hyaline extracts and indexes all of your documentation, and makes that information available through an [MCP server](https://modelcontextprotocol.io).
* TODO other use cases?

And now consider the following use cases that Hyaline is not intended to support:
* Creating and/or updating documentation without human involvement - Hyaline is intended to _augment_ team members, not replace them.
* TODO other use cases?








* Introduction - The purpose of hyaline, what Hyaline is and is not
* Use Cases - What use cases Hyaline is and is not intended to support
* Concepts - The Concepts and terminology used by Hyaline
  * System
  * Product
  * Process
  * Person
  * Team
* Data Flow - Conceptual data flow(s) supported by Hyaline
  * User/AI -> Updates -> Documentation -> Reads -> User/AI
  * Update Documentation (extract current, extract/check change)
  * Use Documentation (extract current, merge, MCP server)
* Components - Document the core components of Hyaline
  * Image on top and then (1) style references to documentation below for each component. More detailed than the 3 data flows above.
  * Components
    * Documentation
    * Code
    * Metadata
    * Extract
    * Data Set (Current & Change)
    * Config
    * Results of check