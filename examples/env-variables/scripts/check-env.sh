#!/bin/bash

echo "Checking Environment Variables"
echo "====================================="
echo ""

# Load variables from previous step
if [ -f /tmp/plexr-env-demo.txt ]; then
    source /tmp/plexr-env-demo.txt
    echo "Loaded saved variables from previous step"
    echo ""
fi

# Function to check if variable exists
check_var() {
    local var_name=$1
    if [ -z "${!var_name}" ]; then
        echo "❌ $var_name is not set"
    else
        echo "✓ $var_name is set to: ${!var_name}"
    fi
}

echo "Checking standard variables:"
check_var "HOME"
check_var "USER"
check_var "PATH"
echo ""

echo "Checking custom variables:"
check_var "PLEXR_DEMO_VAR"
check_var "PLEXR_VERSION"
check_var "PLEXR_TIMESTAMP"
echo ""

# Clean up
rm -f /tmp/plexr-env-demo.txt
echo "Cleanup complete!"