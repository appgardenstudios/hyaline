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

# Build binaries
GOOS=darwin GOARCH=amd64 go build -o ./dist/hyaline-darwin-amd64 -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
GOOS=darwin GOARCH=arm64 go build -o ./dist/hyaline-darwin-arm64 -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
GOOS=linux GOARCH=amd64 go build -o ./dist/hyaline-linux-amd64 -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
GOOS=linux GOARCH=arm64 go build -o ./dist/hyaline-linux-arm64 -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go
GOOS=windows GOARCH=amd64 go build -o ./dist/hyaline-windows-amd64.exe -ldflags="-X 'main.Version=$TAG'" ./cmd/hyaline.go

# Create Tag
git tag $TAG

# Push tag to GitHub
git push origin $TAG

# Create Release (will print link to release when done)
gh release create $TAG --verify-tag --fail-on-no-commits --generate-notes --latest ./dist/*
