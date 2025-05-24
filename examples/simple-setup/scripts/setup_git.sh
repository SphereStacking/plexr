#!/bin/bash

echo "Setting up Git configuration..."

# Create .gitignore if it doesn't exist
if [ ! -f .gitignore ]; then
    cat > .gitignore << EOF
# Dependencies
node_modules/
vendor/

# Build output
build/
dist/
*.exe

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
EOF
    echo "✓ Created .gitignore"
else
    echo "✓ .gitignore already exists"
fi

# Initialize git if not already initialized
if [ ! -d .git ]; then
    git init
    echo "✓ Initialized git repository"
else
    echo "✓ Git repository already initialized"
fi

echo "Git setup completed!"
