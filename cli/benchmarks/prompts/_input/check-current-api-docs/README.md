# Input Data Generation

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-current-api-docs/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-current-api-docs/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-current-api-docs/current.sqlite
```

## Test Execution

The benchmark will run:
```bash
./hyaline check current \
  --config ./benchmarks/prompts/_input/check-current-api-docs/config.yml \
  --system url-shortener \
  --current ./benchmarks/prompts/_input/check-current-api-docs/current.sqlite \
  --output results.json \
  --check-purpose \
  --check-completeness
```