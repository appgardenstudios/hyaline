#!/usr/bin/env bash

# Ensure dist directory exists and is empty
mkdir -p ./dist/
rm -f ./dist/*

# Build and test hyaline
make build
make test
make e2e

# Make sure we are on the main branch
CURRENT_BRANCH=`git rev-parse --abbrev-ref HEAD`
if [ "$CURRENT_BRANCH" != "main" ]; then
  echo "Must be on main branch to release."
  exit 1
fi

# Calculate tag YYYY-MM-DD-HASH
DATE=`git log -n1 --pretty='%cd' --date=format:'%Y-%m-%d'`
HASH=`git rev-parse --short HEAD`
TAG="$DATE-$HASH"

# Build/zip binaries
GOOS=darwin GOARCH=amd64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
zip -9 ./dist/hyaline-darwin-amd64.zip ./dist/hyaline
rm -f ./dist/hyaline

GOOS=darwin GOARCH=arm64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
zip -9 ./dist/hyaline-darwin-arm64.zip ./dist/hyaline
rm -f ./dist/hyaline

GOOS=linux GOARCH=amd64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
zip -9 ./dist/hyaline-linux-amd64.zip ./dist/hyaline
rm -f ./dist/hyaline

GOOS=linux GOARCH=arm64 go build -o ./dist/hyaline -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
zip -9 ./dist/hyaline-linux-arm64.zip ./dist/hyaline
rm -f ./dist/hyaline

GOOS=windows GOARCH=amd64 go build -o ./dist/hyaline.exe -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
zip -9 ./dist/hyaline-windows-amd64.zip ./dist/hyaline.exe
rm -f ./dist/hyaline.exe

# Create Tag
git tag $TAG

# Push tag to GitHub
git push origin $TAG

# Create Draft Release (will print link to release when done)
gh release create $TAG --draft --verify-tag --fail-on-no-commits --generate-notes --latest ./dist/*.zip
