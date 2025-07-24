#!/bin/bash

set -e  # Exit on any error

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$(cd "$SCRIPT_DIR/../../../" && pwd)"

echo "Generating input databases for merge documentation e2e tests..."

cd "$CLI_DIR"

echo "Generating input-1.sqlite..."
./hyaline extract documentation --config "$SCRIPT_DIR/extract-input-1.yml" --output "$SCRIPT_DIR/input-1.sqlite"

echo "Generating input-2.sqlite..."
./hyaline extract documentation --config "$SCRIPT_DIR/extract-input-2.yml" --output "$SCRIPT_DIR/input-2.sqlite"

echo "Generating input-3.sqlite..."
./hyaline extract documentation --config "$SCRIPT_DIR/extract-input-3.yml" --output "$SCRIPT_DIR/input-3.sqlite"

echo "Finished"