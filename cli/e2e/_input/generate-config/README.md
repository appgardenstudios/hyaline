# Input Setup
```bash
cd ./cli/
./hyaline --debug extract current --config ./_example/config.yml --system generate-config --output ./current.db
cp ./current.db ./e2e/_input/generate-config-with-purpose/current.sqlite
```