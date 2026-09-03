#!/usr/bin/env bash
set -e

# Dock-Diet 1-Line Installer Script
# Usage: curl -fsSL https://raw.githubusercontent.com/AsmatZahra-code/dock-diet/main/install.sh | bash

REPO="AsmatZahra-code/dock-diet"
BINARY_NAME="dock-diet"

echo "🥗 Installing Dock-Diet..."

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux*)  OS="linux" ;;
  darwin*) OS="darwin" ;;
  *)
    echo "Error: Unsupported OS: $OS. Dock-Diet supports Linux and macOS via this installer."
    exit 1
    ;;
esac

# 2. Detect CPU Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported Architecture: $ARCH"
    exit 1
    ;;
esac

# 3. Determine latest version release tag
TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
  TAG="v1.1.4"
fi

ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET_NAME}"

TMP_DIR="$(mktemp -d)"
TMP_BINARY="${TMP_DIR}/${BINARY_NAME}"

echo "Downloading Dock-Diet ${TAG} for ${OS}/${ARCH}..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_BINARY"
chmod +x "$TMP_BINARY"

# 4. Install binary to target directory
INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
else
  echo "Elevated permissions required to install to ${INSTALL_DIR}..."
  sudo mv "$TMP_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
fi

rm -rf "$TMP_DIR"

echo "Success! Dock-Diet ${TAG} has been installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Built with ❤️ by Asmat Zahra."
echo "Try running: dock-diet --help"
