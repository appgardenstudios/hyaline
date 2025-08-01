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

## Prompt Benchmarks

Prompt benchmarks test hyaline's ability to use an LLM to correctly identify which documentation needs updates based on the desired documents and when code changes occur.

### Running Benchmarks

```sh
# Run all benchmark scenarios
make benchmark-prompts

# Run specific scenarios
make benchmark-prompts-add-feature
make benchmark-prompts-api-contract
make benchmark-prompts-refactor
# (see Makefile for all available benchmark-prompts-* targets)
```

### Environment Setup

The prompt benchmarks require the same environment variables as the e2e tests:
* **HYALINE_ANTHROPIC_KEY** - Anthropic API key for LLM calls
* **HYALINE_SSH_PEM** - SSH key with access to test repositories
* **HYALINE_GITHUB_PATS** - GitHub token for repository access

### Benchmark Architecture

Each benchmark scenario:
- Runs multiple iterations to account for LLM variability and provides statistical analysis
- Evaluate the results against golden expectations for objective scoring
- Generates markdown reports with collapsible sections showing results
- Scores using the formula: `(expected - missingWeight*missing - unexpectedWeight*unexpected) / expected`
  - `missingWeight` and `unexpectedWeight` can be adjusted in order to calibrate how heavily we weigh false negatives (missing) against false positives (unexpected).

### Generated Reports

Benchmarks generate several output files in `benchmarks/prompts/_output/`:
- Raw JSON results from Hyaline binary execution
- Raw JSON results for all runs
- Multi-run markdown reports with statistical summaries and individual run details

### Golden Expectations Format

Golden expectation files in `benchmarks/prompts/_golden/` define the expected behavior for each scenario:

```json
{
  "description": "Human-readable description of the scenario",
  "expectedRecommendations": [
    {
      "document": "path/to/document.md",
      "section": ["section titles", "that are expected"]
    }
  ]
}
```

- **document**: Full document path that should receive a recommendation
- **section**: A list of sections in hierarchically descending order. Or, an empty array to match the document itself 

## Releasing
Release by checking out the appropriate commit on the main branch and then running `make release`.

Note that test will be run, so the env vars required for testing must be set (see above for details)