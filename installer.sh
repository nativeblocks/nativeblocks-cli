#!/bin/bash
# Installer script for nativeblocks-cli

# Check for curl installation
if ! command -v curl > /dev/null; then
    echo "curl not found. Please install curl and try again."
    exit 1
fi

# Check OS and architecture
OS=$(uname)
ARCH=$(uname -m)
case "$OS" in
    Darwin)
        if [[ "$ARCH" == "arm64" ]]; then
            FILE="nativeblocks-cli_Darwin_arm64.tar.gz"
        elif [[ "$ARCH" == "x86_64" ]]; then
            FILE="nativeblocks-cli_Darwin_x86_64.tar.gz"
        else
            echo "Unsupported architecture for Darwin: $ARCH"
            exit 1
        fi
        ;;
    Linux)
        if [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
            FILE="nativeblocks-cli_Linux_arm64.tar.gz"
        elif [[ "$ARCH" == "x86_64" ]]; then
            FILE="nativeblocks-cli_Linux_x86_64.tar.gz"
        elif [[ "$ARCH" == "i386" || "$ARCH" == "i686" ]]; then
            FILE="nativeblocks-cli_Linux_i386.tar.gz"
        else
            echo "Unsupported architecture for Linux: $ARCH"
            exit 1
        fi
        ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT)
        if [[ "$ARCH" == "arm64" ]]; then
            FILE="nativeblocks-cli_Windows_arm64.zip"
        elif [[ "$ARCH" == "x86_64" ]]; then
            FILE="nativeblocks-cli_Windows_x86_64.zip"
        elif [[ "$ARCH" == "i386" || "$ARCH" == "i686" ]]; then
            FILE="nativeblocks-cli_Windows_i386.zip"
        else
            echo "Unsupported architecture for Windows: $ARCH"
            exit 1
        fi
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Set directories
NATIVEBLOCKS_DIR="${HOME}/.nativeblocks"
NATIVEBLOCKS_BIN_DIR="${NATIVEBLOCKS_DIR}/bin"
mkdir -p "$NATIVEBLOCKS_BIN_DIR"
NATIVEBLOCKS_TMP="${NATIVEBLOCKS_DIR}/tmp"
mkdir -p "$NATIVEBLOCKS_TMP"

# Download binary
DOWNLOAD_URL="https://github.com/nativeblocks/nativeblocks-cli/releases/latest/download/$FILE"
echo "* Downloading $DOWNLOAD_URL"
curl -L "$DOWNLOAD_URL" -o "${NATIVEBLOCKS_TMP}/$FILE"
if [[ $? -ne 0 ]]; then
    echo "Failed to download $DOWNLOAD_URL"
    exit 1
fi

# Extract binary
echo "* Extracting..."
if [[ "$FILE" == *.tar.gz ]]; then
    tar -xzf "${NATIVEBLOCKS_TMP}/$FILE" -C "$NATIVEBLOCKS_BIN_DIR"
elif [[ "$FILE" == *.zip ]]; then
    unzip -o "${NATIVEBLOCKS_TMP}/$FILE" -d "$NATIVEBLOCKS_BIN_DIR"
else
    echo "Unsupported file format: $FILE"
    exit 1
fi

# Clean up
rm -rf "$NATIVEBLOCKS_TMP"

# Add to PATH if not already present
if ! grep -q 'nativeblocks/bin' <<< "$PATH"; then
    SHELL_CONFIG=""
    case "$SHELL" in
        */zsh)
            SHELL_CONFIG="${HOME}/.zshrc"
            ;;
        */bash)
            SHELL_CONFIG="${HOME}/.bashrc"
            ;;
        *)
            SHELL_CONFIG="${HOME}/.profile"
            ;;
    esac
    echo "export PATH=\"\$PATH:$NATIVEBLOCKS_BIN_DIR\"" >> "$SHELL_CONFIG"
    echo "* Added nativeblocks-cli to PATH in $SHELL_CONFIG"
fi

echo ""
echo "Installation was successful!"
echo "Please open a new terminal OR run the following to update PATH in the current terminal:"
echo ""
echo "    export PATH=\"\$PATH:$NATIVEBLOCKS_BIN_DIR\""
echo ""
echo "Then run 'nativeblocks' to start using it!"
