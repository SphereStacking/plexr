#!/bin/bash

echo "Checking required tools..."

# Check for git
if ! command -v git &> /dev/null; then
    echo "Error: git is not installed"
    exit 1
fi
echo "✓ git $(git --version | cut -d' ' -f3)"

# Check for node (optional)
if command -v node &> /dev/null; then
    echo "✓ node $(node --version)"
else
    echo "⚠ node is not installed (optional)"
fi

echo "Tool check completed!"
