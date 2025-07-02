---
title: "How To: Build the Hyaline CLI"
linkTitle: Build the CLI
purpose: Document how to build the Hyaline CLI
url: documentation/how-to/build-cli
---
## Purpose
Build the hyaline cli locally or on a remote machine.

## Prerequisite(s)
* [Go Toolchain](https://go.dev/) (version 1.24+)

## Steps

### 1. Clone/Checkout Hyaline Repo
Ensure that the [Hyaline Repository](https://github.com/appgardenstudios/hyaline) is cloned and checked out to the version you wish to build.

### 2. Determine the Version
Hyaline uses the commit date and short hash to specify the version being build. This can be determined by running the following commands in the root of the repository with the desired commit checked out.

```bash
# Format: YYYY-MM-DD-HASH
DATE=`git log -n1 --pretty='%cd' --date=format:'%Y-%m-%d'`
HASH=`git rev-parse --short HEAD`
VERSION="$DATE-$HASH"
```

### 3. Run Go Build
Once the version is specified you can build hyaline using the following command (from the root of the repo):

```bash
$ go build -o ./dist/hyaline -ldflags="-X 'main.Version=$VERSION'" ./cmd/hyaline.go
```

You can specify the OS and architecture to use by setting the appropriate [GOOS/GOARCH environment variables](https://go.dev/doc/install/source#environment). For example, to build for the 64bit ARM version of MacOS:

```bash
$ GOOS=darwin GOARCH=arm64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
```

### 4. Test Executable
Make sure you test that the resulting executable was built sucessfully and at the right version by running:

```bash
$ ./hyaline version
```

## Next Steps
Visit [CLI Reference](../04-reference/02-cli.md).