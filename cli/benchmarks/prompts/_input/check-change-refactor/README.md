# Input Data Generation

To generate the current.sqlite and change.sqlite files for this scenario:

```bash
cd ./cli/

# Generate current state (url-shortener branch)
rm -f ./benchmarks/prompts/_input/check-change-refactor/current.sqlite
./hyaline extract current \
  --config ./benchmarks/prompts/_input/check-change-refactor/config.yml \
  --system url-shortener \
  --output ./benchmarks/prompts/_input/check-change-refactor/current.sqlite

# Generate change state (comparing url-shortener branch to PR #6 branch)
rm -f ./benchmarks/prompts/_input/check-change-refactor/change.sqlite
./hyaline extract change \
  --config ./benchmarks/prompts/_input/check-change-refactor/config.yml \
  --system url-shortener \
  --base url-shortener \
  --head url-shortener-refactor \
  --pull-request appgardenstudios/hyaline-example/6 \
  --output ./benchmarks/prompts/_input/check-change-refactor/change.sqlite
```