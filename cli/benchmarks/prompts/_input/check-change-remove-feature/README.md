# Input Data Generation

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-change-remove-feature/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-change-remove-feature/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-change-remove-feature/current.sqlite

# Generate change state (comparing url-shortener branch to PR #5 branch)
rm -f ./benchmarks/prompts/_input/check-change-remove-feature/change.sqlite
./hyaline extract change \
  --config ./benchmarks/prompts/_input/check-change-remove-feature/config.yml \
  --system url-shortener \
  --base url-shortener \
  --head url-shortener-remove-click-stats \
  --pull-request appgardenstudios/hyaline-example/5 \
  --output ./benchmarks/prompts/_input/check-change-remove-feature/change.sqlite
```