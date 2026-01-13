#!/bin/bash
# Setup vendor directory for offline builds

set -e

echo "Setting up vendor directory for offline builds..."

# Generate go.sum if missing
if [ ! -f "go.sum" ]; then
    echo "Generating go.sum..."
    go mod tidy
fi

# Create vendor directory
echo "Creating vendor directory..."
go mod vendor

echo "Vendor directory created successfully!"
echo "You can now build offline using: go build -mod=vendor"

