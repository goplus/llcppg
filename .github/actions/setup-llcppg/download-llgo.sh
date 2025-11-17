#!/bin/bash
set -e

# Download and install LLGo from pre-built releases
# Usage: download-llgo.sh <version>
# Example: download-llgo.sh v0.11.6

VERSION="${1:-v0.11.6}"

# Determine platform and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
  PLATFORM="darwin"
elif [ "$OS" = "linux" ]; then
  PLATFORM="linux"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Download and extract llgo
# Remove 'v' prefix from version for filename
VERSION_NUM="${VERSION#v}"
FILENAME="llgo${VERSION_NUM}.${PLATFORM}-${ARCH}.tar.gz"
URL="https://github.com/goplus/llgo/releases/download/${VERSION}/${FILENAME}"

echo "Downloading llgo from: $URL"
wget -q "$URL" -O llgo.tar.gz

# Extract to llgo directory
mkdir -p llgo
tar -xzf llgo.tar.gz -C llgo --strip-components=1
rm llgo.tar.gz

# Add to PATH
echo "$GITHUB_WORKSPACE/llgo/bin" >> "$GITHUB_PATH"

echo "LLGo ${VERSION} installed successfully to $GITHUB_WORKSPACE/llgo"
