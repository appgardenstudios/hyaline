# How to Add Articles to www

This guide explains how to add new articles to the Hyaline website at [`/www/content/articles/`](../../www/content/articles/). The live articles can be viewed at [https://hyaline.dev/articles](https://hyaline.dev/articles) once deployed.

## Adding Articles

### 1. File Location and Structure

Articles are stored in `/www/content/articles/` with a `.md` extension. Each article is a single Markdown file.

### 2. Frontmatter

**Required frontmatter properties:**
- `title`: The main title of the article
- `subtitle`: A brief subtitle describing the article content
- `author`: The author's name
- `date`: Publication date in YYYY-MM-DD format
- `url`: Custom URL path (e.g., `articles/purpose`)
- `thumbnail`: Path to thumbnail image (144x144 pixels)

Example:
```yaml
---
title: Document with Purpose
subtitle: How to purposefully write documentation.
author: John Clark
date: 2025-06-23
url: articles/purpose
thumbnail: ./_img/purpose-thumbnail.svg
---
```

### 3. Images and Visual Content

#### Thumbnail Images:
- Must be exactly **144x144 pixels**
- Use relative paths from the article file
- Prefer SVG format for scalability

#### Portrait Images:
- Wrap sections containing portrait images with `<div class="portrait">`
- Portrait images should be designed with **1920x1080 aspect ratio**
- Use relative paths from the document

Example:
```markdown
<div class="portrait">

![Purpose of Documentation](./_img/purpose-purpose.svg)

Your content here...

</div>
```