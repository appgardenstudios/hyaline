# Hyaline CLI

# Developing

## Dependencies

* `make`
* `go` (v1.24+)
* `diff` (gnu version) for testing (`brew install diffutils`)
* `gh` (github cli) for testing and releasing (`brew install gh`)
* `sqlc` for compiling sql queries (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0`)

## Generate DB Queries
```sh
$ make db
```

## Running Locally
```sh
$ make
$ ./hyaline
```

Note that you can also use `make install` to build and install Hyaline locally for testing.

## Debugging
There is a `.vscode/launch.json` file checked in that has various debugger launch configurations. They use the config stored in `cli/_example/config.yml` and rely on a `cli/.env` file being present to work (which is .gitignored). The .env file should look like:

```
HYALINE_ANTHROPIC_KEY= #The Anthropic API key
HYALINE_GITHUB_PAT= #A GitHub Personal Access Token that has read access to github.com/appgardenstudios/hyaline-example
HYALINE_SSH_PEM= #A SSH key that has pull access to github.com/appgardenstudios/hyaline-example. Note that this will need to be ""'d and newlines replaced with \n
HYALINE_SSH_PASSWORD= #A password for the PEM above (blank if PEM is not password protected)
```

Note that you must have the `github.com/appgardenstudios/hyaline-example` repository cloned as a sibling directory to hyaline for some of the launch configurations to work properly.

## Testing
Unit tests are run with `make test`, and there are e2e tests that invoke the actual hyaline binary that you can run with `make e2e`.

Note that the following env vars must be set for the `e2e` tests to work and pass:
* **HYALINE_SSH_PEM** (with access to `github.com/appgardenstudios/hyaline-example`)
* **HYALINE_GITHUB_PAT** (with access to `github.com/appgardenstudios/hyaline-example`)

Note that the e2e test [updatePR_test.go](./e2e/updatePR_test.go) creates a comment on a PR and cleans it up using the GitHub CLI (`gh`).

## Releasing
Release by checking out the appropriate commit on the main branch and then running `make release`.

Note that test will be run, so the env vars required for testing must be set (see above for details)