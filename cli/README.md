# Hyaline CLI

# Developing

## Dependencies

* `make`
* `go` (v1.24+)

## Running Locally
```sh
$ make
$ ./hyaline
```

## Debugging
There is a `.vscode/launch.json` file checked in that has various debugger launch configurations. They use the config stored in `cli/_example/config.yml` and rely on a `cli/.env` file being present to work (which is .gitignored). The .env file should look like:

```
HYALINE_ANTHROPIC_KEY= #The Anthropic API key
HYALINE_GITHUB_PAT= #A GitHub Personal Access Token that has read access to github.com/appgardenstudios/hyaline-example
HYALINE_SSH_PEM= #A SSH key that has pull access to github.com/appgardenstudios/hyaline-example. Note that this will need to be ""'d and newlines replaced with \n
HYALINE_SSH_PASSWORD= #A password for the PEM above (blank if PEM is not password protected)
```

## Testing
Unit tests are run with `make test`, and there are e2e tests that invoke the actual hyaline binary that you can run with `make e2e`.

Note that the following env vars must be set for the `e2e` tests to work and pass:
* **HYALINE_SSH_PEM** (with access to `github.com/appgardenstudios/hyaline-example`)
* **HYALINE_GITHUB_PAT** (with access to `github.com/appgardenstudios/hyaline-example`)