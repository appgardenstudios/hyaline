#!/bin/bash
#
# Setup script for hyaline repository
# This script installs git hooks and sets up the development environment

set -e

echo "Setting up hyaline development environment..."

# Install pre-commit hook
echo "Installing git hooks..."

cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
#
# Pre-commit hook that runs go fmt on .go files in the cli directory

echo "Running go fmt on cli directory..."

# Change to cli directory and run go fmt on all .go files
cd cli
if ! go fmt ./...; then
    echo "Error: go fmt failed"
    exit 1
fi

# Add any formatted files to staging area
git add .

echo "go fmt completed successfully"
exit 0
EOF

# Make the hook executable
chmod +x .git/hooks/pre-commit

echo "✓ Pre-commit hook installed"
echo "✓ Setup complete!"