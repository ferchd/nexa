#!/bin/bash
set -e

echo "Nexa Quick Installer"
echo "============================"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="armv7" ;;
    *) ARCH="386" ;;
esac

LATEST_VERSION=$(curl -s https://api.github.com/repos/ferchd/nexa/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
DOWNLOAD_URL="https://github.com/ferchd/nexa/releases/download/${LATEST_VERSION}/nexa-${OS}-${ARCH}"

echo "Downloading Nexa ${LATEST_VERSION} for ${OS}/${ARCH}..."

# Download and install
curl -L -o nexa ${DOWNLOAD_URL}
chmod +x nexa
sudo mv nexa /usr/local/bin/

echo "Nexa installed successfully!"
echo ""
echo "Usage:"
echo "  nexa --external 8.8.8.8:53 --stdout-json"
echo "  nexa --help"