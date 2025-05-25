#!/bin/bash

echo "Creating test file in current directory..."
echo "This file was created in the global work directory" > global-test.txt
echo "✓ Created file: $(pwd)/global-test.txt"
echo "Content: $(cat global-test.txt)"
echo ""