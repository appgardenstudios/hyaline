# How to Add Documentation to www

This guide explains how to add new documentation to the Hyaline website at [`/www/content/documentation/`](../../www/content/documentation/). The live documentation can be viewed at [https://hyaline.dev/documentation](https://hyaline.dev/documentation) once it is deployed.

## Adding Documentation

### 1. File Naming and Ordering

Files use numbered prefixes to control their order in navigation. Example:

```
01-overview.md
03-how-to/
  01-install-cli.md
  02-run-cli.md
  03-check-pr.md
04-explanation/
  01-hyaline.md
  02-extract-current.md
```

**To insert between existing files:**
1. Add your file with the appropriate number
2. Renumber all following files to maintain sequence

Example: To add a new how-to guide between "Install CLI" and "Run CLI":
- Add `02-my-new-guide.md`
- Rename `02-run-cli.md` to `03-run-cli.md`
- Rename `03-check-pr.md` to `04-check-pr.md`

### 2. Directory Structure

- **Top-level files** (like `01-overview.md`) appear at the root documentation level
- **Subdirectories** automatically create categories in navigation
- **`_img/` directories** within categories store images for that category

### 3. Frontmatter

**Frontmatter properties:**
- `title`: Include category prefix for categorized docs (e.g., "How To: ", "Explanation: ")
- `linkTitle`: For categorized docs, omit the category prefix (used in navigation)
- `purpose`: Brief description of the document's purpose
- `url`: Custom URL path (removes numbered prefixes from URLs)


Example For Top-Level Documents:
```yaml
---
title: Overview
linkTitle: Overview
purpose: Give an overview of Hyaline and the documentation
url: documentation/overview
---
```

Example For Categorized Documents:
```yaml
---
title: "How To: Install the Hyaline CLI"
linkTitle: Install the Hyaline CLI
purpose: Document how to install the Hyaline CLI
url: documentation/how-to/install-cli
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
[How To Install the CLI](./03-how-to/01-install-cli.md)
[CLI Reference](../05-reference/02-cli.md)
```