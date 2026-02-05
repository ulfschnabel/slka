#!/bin/bash
# Setup and build script for slka

set -e

echo "====================================="
echo "slka - Slack CLI Setup"
echo "====================================="
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    echo "Please run: sudo ./install-go.sh"
    exit 1
fi

echo "Go version:"
go version
echo

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies downloaded"
echo

# Run tests
echo "Running tests..."
go test -v ./...
TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✓ All tests passed"
else
    echo "✗ Some tests failed (exit code: $TEST_EXIT_CODE)"
    echo "This is expected if you haven't set up Slack tokens yet."
fi
echo

# Build binaries
echo "Building binaries..."
make build-local
echo "✓ Binaries built in dist/"
echo

# Show what was built
echo "Built binaries:"
ls -lh dist/
echo

echo "====================================="
echo "Setup complete!"
echo "====================================="
echo
echo "Next steps:"
echo "1. Install the binaries:"
echo "   make install"
echo
echo "2. Or run directly from dist/:"
echo "   ./dist/slka-read --help"
echo "   ./dist/slka-write --help"
echo
echo "3. Configure your Slack tokens:"
echo "   ./dist/slka-write config init"
echo
echo "4. To build for all platforms:"
echo "   make build"
echo
