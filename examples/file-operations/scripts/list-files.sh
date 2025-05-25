#!/bin/bash

# Read the temp directory from state
TEMP_DIR=$(cat /tmp/plexr-demo-tempdir.txt)

echo "Listing files in temporary directory:"
echo "====================================="
echo ""

# List files with details
ls -la "$TEMP_DIR"

echo ""
echo "File count: $(ls -1 "$TEMP_DIR" | wc -l)"
echo ""