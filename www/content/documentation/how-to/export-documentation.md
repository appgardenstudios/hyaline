---
title: "How To: Export Documentation"
description: "How to export documentation extracted by Hyaline."
purpose: Document how to export documentation using the Hyaline CLI
---
## Purpose
Export documentation using the Hyaline CLI.

## Prerequisite(s)
* [Installed the CLI](./install-cli.md)
* Extracted and (optionally) merged documentation in a documentation data set (e.g. `./documentation.db`)

## Steps

### 1. Select Output Format
Hyaline supports exporting extracted documentation in a variety of formats. See [CLI](../reference/cli.md) or [Export](../reference/export.md) reference documentation to see the available options and choose a format.

### 2. Select Included Documentation
By default Hyaline will export all available documentation. If you wish, you can include or exclude specific documentation by passing in `--include` or `--exclude` document URIs in the form of `document://<source>/<path/to/document>(?tagValue=tagKey)`.

For example:
- `--include 'document://*/**/*'` will include every document (`**/*`) in every source (`*`) in the data set (this is the default if no includes are specified)
- `--include 'document://my-app/**/*'` will include every document (`**/*`) in the `my-app` source
- `--exclude 'document://my-app/old/README.md'` will exclude the document `old/README.md` in the `my-app` source
- `--exclude 'document://*/**?type=customer` will exclude any document that has the tag `type=customer`

Note that you can include any number of `--include` and `--exclude` in the export command. Hyaline will export any document that matches at least one include and does not match any exclude.

### 3. Run Export
Run `hyaline export documentation` with your desired arguments to export your documentation.

```
$ hyaline export documentation --documentation ./documentation.db /
  --format json --output ./export.json /
  --include 'document://frontend/**/*' /
  --exclude 'document://*/**/*?type=customer'
```
For example, the command above exports the documentation in `./documentation.db` to the file `./export.json` in JSON format. It only include documentation from the `frontend` source, and excludes any documentation with the tag `type=customer`.

## Next Steps
Visit [How To Run the Hyaline MCP Server](./run-mcp.md) or the [CLI Reference](../reference/cli.md).