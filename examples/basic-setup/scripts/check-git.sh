#!/bin/bash

echo "====================================="
echo "Git Configuration Check"
echo "====================================="
echo ""

# Check current Git configuration (read-only)
echo "Current Git configuration:"
echo ""

if command -v git &> /dev/null; then
    username=$(git config --global user.name 2>/dev/null || echo "not set")
    email=$(git config --global user.email 2>/dev/null || echo "not set")
    
    echo "User name: $username"
    echo "User email: $email"
    echo ""
    
    echo "Recommended configuration commands (not executed):"
    echo "  git config --global user.name \"Your Name\""
    echo "  git config --global user.email \"your.email@example.com\""
    echo "  git config --global init.defaultBranch main"
    echo "  git config --global core.editor vim"
else
    echo "Git is not installed."
    echo "Would install Git and configure with:"
    echo "  - User name and email"
    echo "  - Default branch: main"
    echo "  - Default editor: vim"
fi

echo ""
echo "Note: No changes are made to Git configuration."
echo ""