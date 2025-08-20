---
title: "Reference: Export"
description: The output formats supported by the export command in Hyaline
purpose: Detail the output formats available to be used with the export documentation command
---
## Overview
This documents the output formats available in the `export documentation` command.

## fs
The file system export format (`--format fs`). This format will output documentation to a file system in the folder structure shown below:

```txt
output-path/ # The path specified by --output
  source1/ # separate directories for each source
    /path/to/document1.md # 1 file for each document exported for a source
    /path/to/document2.md
  source2
    /path/to/document3.md
    ...
  README.md # Metadata about the export
```
