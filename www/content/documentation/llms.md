---
title: LLMs
description: "Learn how Hyaline uses LLMs and what LLM providers are supported"
purpose: Document how Hyaline uses LLMs and list all of the LLM providers that Hyaline supports
---
## Overview

Hyaline performs many checks to determine when documentation is outdated or out of compliance. Some of these checks are objective, such as whether a document has been updated or contains a certain phrase. Other types of checks are more subjective, such as whether code changes will make documentation outdated or whether documentation matches its intended purpose.

Hyaline uses an LLM to perform subjective checks by providing it a targeted prompt and relevant context. Due to the non-deterministic nature of LLMs and subjectivity of the checks, recommendations and audit results are presented to humans for review and resolution.

## Bring Your Own LLM

Hyaline requires that users provide their own access to an LLM. This has several advantages:

- **Cost Control** - You have complete visibility and control over your LLM usage. You can monitor costs directly, set spending limits, and optimize usage based on your needs.
- **Privacy** - This approach minimizes the data that flows through Hyaline's systems. Your code and documentation are processed directly by your chosen LLM provider.
- **Model Control** - You have full control over which model Hyaline uses. You can use the right model for your needs and have confidence that you're not being downgraded to a lower performing model.
- **Simplicity** - Many organizations already have enterprise agreements with LLM providers. You can leverage these existing relationships, compliance approvals, and negotiated rates rather than managing another vendor relationship.

## Supported LLM Providers

The following LLM providers are currently supported:
- [Anthropic](https://www.anthropic.com)
- [GitHub Models](https://github.com/features/models) (Offers a free tier)
- [OpenAI](https://openai.com)

We are working to add support for more LLM providers. Please [send us feedback](https://github.com/appgardenstudios/hyaline/discussions/categories/feedback) if there is an LLM provider you would like us to add support for.

## Next Steps

Read more about how Hyaline [generates recommendations](./explanation/check/) and [audits documentation](./explanation/audit/).