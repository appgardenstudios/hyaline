# Input Setup
```bash
cd ./cli/
rm -f ./e2e/_input/generate-config-with-purpose/current.sqlite
./hyaline --debug extract current --config ./_example/config.yml --system generate-config --output ./e2e/_input/generate-config-with-purpose/current.sqlite
```