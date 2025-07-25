# How to Add Documentation to www

This guide explains how to add new documentation to the Hyaline website at [`/www/content/documentation/`](../../www/content/documentation/). The live documentation can be viewed at [https://hyaline.dev/documentation](https://hyaline.dev/documentation) once it is deployed.

## Adding Documentation

### 1. Table of Contents Configuration

The documentation table of contents is controlled by the `www/data/documentation_toc.yml` file. This file defines the structure and order of documentation pages.

Example structure:
```yaml
items:
  - title: Overview
    url: /documentation-v1/overview # an item can have a link
  - title: How To # an item does not have to have a link
    items: # an item can have sub-items
      - title: Install the CLI
        url: /documentation-v1/how-to/install-cli
      - title: Run the CLI
        url: /documentation-v1/how-to/run-cli
  - title: Reference
    items:
      - title: CLI
        url: /documentation-v1/reference/cli
  - title: Roadmap
    url: /documentation-v1/roadmap
```

To add a new documentation page:
1. Create the markdown file in the appropriate directory
2. Add an entry in `documentation_toc.yml` with the title and URL

### 2. Directory Structure

- **Top-level files** (like `overview.md`) appear at the root documentation level
- **Subdirectories** organize related documentation (e.g., `how-to/`, `reference/`)
- **`_img/` directories** within categories store images for that category

### 3. Frontmatter

**Frontmatter properties:**
- `title`: The full title of the document
- `description`: A brief description of the document
- `purpose`: Brief description of the document's purpose

Example:
```yaml
---
title: Overview
description: "Learn how Hyaline helps teams maintain documentation"
purpose: Give an overview of Hyaline and the documentation
---
```

### 4. Images and Links

#### Images:
- Store images in `_img/` subdirectory within the document's category
- Use relative paths from the document
- Prefer SVG format for diagrams and illustrations

Example:
```markdown
![Team](_img/hyaline-team.svg)
```

#### Links:
- Use relative paths for internal documentation links
- Links can be relative to the current document

Example:
```markdown
[How To Install the CLI](./how-to/install-cli.md)
[CLI Reference](../reference/cli.md)
```