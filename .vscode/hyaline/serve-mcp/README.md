# Input Setup
```bash
cd ${workspaceFolder}/cli/

# Extract documentation using filesystem crawler
rm -f ../.vscode/hyaline/serve-mcp/documentation.sqlite
./hyaline extract documentation --config ../.vscode/hyaline/serve-mcp/extract-config.yml --output ../.vscode/hyaline/serve-mcp/documentation.sqlite
```