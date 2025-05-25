#!/bin/bash

echo "Checking if file exists in /tmp..."
if [ -f "/tmp/plexr-demo.txt" ]; then
    echo "✓ File found: /tmp/plexr-demo.txt"
    echo "Content: $(cat /tmp/plexr-demo.txt)"
    rm -f /tmp/plexr-demo.txt
    echo "✓ Cleaned up test file"
else
    echo "✗ File not found in /tmp"
fi
echo ""