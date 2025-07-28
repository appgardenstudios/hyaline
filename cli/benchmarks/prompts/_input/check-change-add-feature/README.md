# Input Data Generation

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-change-add-feature/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-change-add-feature/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-change-add-feature/current.sqlite

# Generate change state (comparing url-shortener branch to url-shortener-expiration branch)
rm -f ./benchmarks/prompts/_input/check-change-add-feature/change.sqlite
./hyaline extract change \
  --config ./benchmarks/prompts/_input/check-change-add-feature/config.yml \
  --system url-shortener \
  --base url-shortener \
  --head url-shortener-expiration \
  --pull-request appgardenstudios/hyaline-example/4 \
  --output ./benchmarks/prompts/_input/check-change-add-feature/change.sqlite
```