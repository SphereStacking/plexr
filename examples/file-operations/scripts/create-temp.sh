#!/bin/bash

echo "Creating temporary files for demonstration..."
echo ""

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
echo "Created temporary directory: $TEMP_DIR"

# Store the temp directory in state for later steps
echo "$TEMP_DIR" > /tmp/plexr-demo-tempdir.txt

# Create some sample files
echo "Hello from Plexr!" > "$TEMP_DIR/hello.txt"
echo "This is a demo file" > "$TEMP_DIR/demo.txt"
echo "Line 1
Line 2
Line 3" > "$TEMP_DIR/multiline.txt"

echo ""
echo "✓ Created hello.txt"
echo "✓ Created demo.txt"
echo "✓ Created multiline.txt"
echo ""