#!/bin/bash

# Read the temp directory from state
TEMP_DIR=$(cat /tmp/plexr-demo-tempdir.txt)

echo "Cleaning up temporary files..."
echo ""

# Remove the temporary directory and its contents
if [ -d "$TEMP_DIR" ]; then
    rm -rf "$TEMP_DIR"
    echo "✓ Removed temporary directory: $TEMP_DIR"
fi

# Remove the state file
rm -f /tmp/plexr-demo-tempdir.txt
echo "✓ Removed state file"

echo ""
echo "Cleanup complete! No permanent changes were made to your system."
echo ""