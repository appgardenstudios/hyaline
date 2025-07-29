#!/bin/bash

set -e  # Exit on any error

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$(cd "$SCRIPT_DIR/../../../cli" && pwd)"

echo "Generating input databases for check diff e2e tests..."

cd "$CLI_DIR"

echo "Generating documentation.sqlite..."
./hyaline --debug extract documentation --config "$SCRIPT_DIR/hyaline.yml" --output "$SCRIPT_DIR/documentation.sqlite"

echo "Finished"