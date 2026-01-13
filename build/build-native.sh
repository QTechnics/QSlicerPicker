#!/bin/bash

# Build script for native platform only (simplest, no cross-compilation)

set -e

CURRENT_OS=$(go env GOOS)
CURRENT_ARCH=$(go env GOARCH)

echo "Building for native platform: $CURRENT_OS/$CURRENT_ARCH..."

# Set output name
if [ "$CURRENT_OS" = "windows" ]; then
    OUTPUT="qslicerpicker.exe"
else
    OUTPUT="qslicerpicker"
fi

# Build (CGO is automatically enabled)
go build -o "$OUTPUT" ../

echo "Build complete: $OUTPUT"
