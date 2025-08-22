---
title: "Reference: Export"
description: The output formats supported by the export command in Hyaline
purpose: Detail the output formats available to be used with the export documentation command
---
## Overview
This documents the output formats available in the `export documentation` command.

## File System
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

## llms-full.txt
The llms-full.txt export format (`--format llmsfulltxt`). This format will output documentation to the output path as a text file using the structure shown below:

```txt
# <Title> (Name of first section, or document ID if none found)
Source: <Document URI> (document://<source>/<document>)

<Document Contents>


# <Title>
...
```

**Note**: The documentation is sorted by source ID ascending, document ID ascending