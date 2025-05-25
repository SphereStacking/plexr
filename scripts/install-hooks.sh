#!/bin/bash
# Install git hooks for the project

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "Installing git hooks..."

# Create pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash
# Pre-commit hook for Plexr

echo "Running pre-commit checks..."

# Format check
echo "Checking code formatting..."
if ! make fmt-check > /dev/null 2>&1; then
    echo "❌ Code is not properly formatted. Run 'make fmt' to fix."
    exit 1
fi

# Lint check
echo "Running linters..."
if ! make lint; then
    echo "❌ Linting failed. Please fix the issues above."
    exit 1
fi

# Test check (optional - can be slow)
# echo "Running tests..."
# if ! make test; then
#     echo "❌ Tests failed. Please fix the failing tests."
#     exit 1
# fi

echo "✅ All pre-commit checks passed!"
EOF

chmod +x "$HOOKS_DIR/pre-commit"

echo "✅ Git hooks installed successfully!"
echo ""
echo "The following hooks were installed:"
echo "  - pre-commit: Runs formatting and linting checks"
echo ""
echo "To skip hooks temporarily, use: git commit --no-verify"