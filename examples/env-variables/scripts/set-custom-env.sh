#!/bin/bash

echo "Setting Custom Environment Variables"
echo "====================================="
echo ""

# Set some custom variables for this session
export PLEXR_DEMO_VAR="Hello from Plexr!"
export PLEXR_VERSION="1.0.0"
export PLEXR_TIMESTAMP=$(date +%s)

echo "Set the following variables:"
echo "  PLEXR_DEMO_VAR: $PLEXR_DEMO_VAR"
echo "  PLEXR_VERSION: $PLEXR_VERSION"
echo "  PLEXR_TIMESTAMP: $PLEXR_TIMESTAMP"
echo ""

# Demonstrate variable expansion
MESSAGE="Demo running with Plexr version $PLEXR_VERSION"
echo "Message: $MESSAGE"
echo ""

# Save variables for next step (simulating state)
echo "PLEXR_DEMO_VAR=$PLEXR_DEMO_VAR" > /tmp/plexr-env-demo.txt
echo "PLEXR_VERSION=$PLEXR_VERSION" >> /tmp/plexr-env-demo.txt
echo "PLEXR_TIMESTAMP=$PLEXR_TIMESTAMP" >> /tmp/plexr-env-demo.txt