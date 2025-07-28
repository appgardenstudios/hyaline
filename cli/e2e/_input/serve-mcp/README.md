# Input Setup
```bash
cd ./cli/

# Extract documentation using filesystem crawler
rm -f ./e2e/_input/serve-mcp/documentation.sqlite
./hyaline extract documentation --config ./e2e/_input/serve-mcp/extract-config.yml --output ./e2e/_input/serve-mcp/documentation.sqlite
```