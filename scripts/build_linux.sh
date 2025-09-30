#!/bin/bash

set -e

VERSION=${1:-$(git describe --tags --abbrev=0)}
OUTPUT_DIR="dist/linux"
mkdir -p $OUTPUT_DIR

echo "Building NetCheck for Linux - Version: $VERSION"

ARCHITECTURES=("amd64" "386" "arm64" "arm")

for arch in "${ARCHITECTURES[@]}"; do
    echo "Building for $arch..."
    
    export GOOS=linux
    export GOARCH=$arch
    
    go build -ldflags="-s -w -X main.version=$VERSION" \
             -o "$OUTPUT_DIR/netcheck-linux-$arch" \
             ./cmd/netcheck
    
    # Create tarball
    tar -czf "$OUTPUT_DIR/netcheck-linux-$arch.tar.gz" \
        -C "$OUTPUT_DIR" \
        "netcheck-linux-$arch"
    
    echo "✅ Built $OUTPUT_DIR/netcheck-linux-$arch.tar.gz"
done

cd $OUTPUT_DIR
sha256sum *.tar.gz > checksums.txt
echo "✅ Created checksums"

echo "Linux builds completed successfully!"