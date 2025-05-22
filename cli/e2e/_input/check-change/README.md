# Input Setup
```bash
cd ./cli/

rm -f ./e2e/_input/check-change/current.sqlite
./hyaline --debug extract current --config ./_example/config.yml --system check-change --output ./e2e/_input/check-change/current.sqlite

rm -f ./e2e/_input/check-change/change.sqlite
./hyaline --debug extract change --config ./_example/config.yml --system check-change --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./e2e/_input/check-change/change.sqlite
```