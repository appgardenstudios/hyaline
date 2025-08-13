#!/bin/bash
#
# Setup script for hyaline repository
# This script installs git hooks and sets up the development environment

set -e

echo "Setting up hyaline development environment..."

# Install pre-commit hook
echo "Installing git hooks..."

# Copy pre-commit hook from scripts/hooks to .git/hooks
cp scripts/hooks/pre-commit .git/hooks/pre-commit

echo "✓ Pre-commit hook installed"
echo "✓ Setup complete!"