#!/bin/bash

# Script to download and extract LLGo release
# Usage: ./download-llgo.sh <version> <install_dir>
# Example: ./download-llgo.sh v0.11.6 ./llgo

set -e

VERSION=$1
INSTALL_DIR=$2

if [ -z "$VERSION" ] || [ -z "$INSTALL_DIR" ]; then
    echo "Usage: $0 <version> <install_dir>"
    echo "Example: $0 v0.11.6 ./llgo"
    exit 1
fi

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Construct download URL
# Format: llgo{version}.{os}-{arch}.tar.gz
# Example: llgo0.11.6.darwin-arm64.tar.gz or llgo0.11.6.linux-amd64.tar.gz
# Remove 'v' prefix from version if present
VERSION_NUMBER="${VERSION#v}"
FILENAME="llgo${VERSION_NUMBER}.${OS}-${ARCH}.tar.gz"
URL="https://github.com/goplus/llgo/releases/download/${VERSION}/${FILENAME}"

echo "Downloading LLGo ${VERSION} for ${OS}-${ARCH}..."
echo "URL: $URL"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download and extract
curl -L -o "/tmp/${FILENAME}" "$URL"
tar -xzf "/tmp/${FILENAME}" -C "$INSTALL_DIR"
rm "/tmp/${FILENAME}"

echo "LLGo ${VERSION} has been installed to ${INSTALL_DIR}"
echo "Binary location: ${INSTALL_DIR}/bin/llgo"

# Verify installation
if [ -f "${INSTALL_DIR}/bin/llgo" ]; then
    echo "Installation verified successfully"
    ls -lh "${INSTALL_DIR}/bin/llgo"
else
    echo "Error: llgo binary not found at ${INSTALL_DIR}/bin/llgo"
    exit 1
fi
