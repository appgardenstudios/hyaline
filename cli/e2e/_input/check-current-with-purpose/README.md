# Input Setup
```bash
cd ./cli/
rm -f ./e2e/_input/check-current-with-purpose/current.sqlite
./hyaline --debug extract current --config ./_example/config.yml --system check-current --output ./e2e/_input/check-current-with-purpose/current.sqlite
```