# Input Data Generation

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-change-api-contract/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-change-api-contract/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-change-api-contract/current.sqlite

# Generate change state (comparing url-shortener branch to PR #7 branch)
rm -f ./benchmarks/prompts/_input/check-change-api-contract/change.sqlite
./hyaline extract change \
  --config ./benchmarks/prompts/_input/check-change-api-contract/config.yml \
  --system url-shortener \
  --base url-shortener \
  --head url-shortener-break-error-api-contract \
  --pull-request appgardenstudios/hyaline-example/7 \
  --output ./benchmarks/prompts/_input/check-change-api-contract/change.sqlite
```