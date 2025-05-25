#!/bin/bash

# Read the temp directory from state
TEMP_DIR=$(cat /tmp/plexr-demo-tempdir.txt)

echo "Reading file contents:"
echo "====================================="
echo ""

# Read each file
for file in "$TEMP_DIR"/*.txt; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        echo "--- Content of $filename ---"
        cat "$file"
        echo ""
    fi
done

echo "====================================="
echo ""