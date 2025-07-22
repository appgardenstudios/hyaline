# Input Data Generation

To regenerate the current.sqlite file for this scenario:

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-current-customer-docs/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-current-customer-docs/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-current-customer-docs/current.sqlite
```