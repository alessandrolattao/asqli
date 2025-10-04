#!/bin/bash
set -e

# ASQLI Installer
# This script automatically detects your OS and architecture, then downloads
# and installs the latest version of ASQLI.

REPO="alessandrolattao/asqli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="asqli"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored message
print_message() {
    echo -e "${GREEN}==>${NC} $1"
}

print_error() {
    echo -e "${RED}Error:${NC} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        MINGW*|MSYS*|CYGWIN*)    echo "windows";;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64";;
        aarch64|arm64)  echo "arm64";;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    if command -v curl &> /dev/null; then
        curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
    elif command -v wget &> /dev/null; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download file
download_file() {
    local url=$1
    local output=$2

    if command -v curl &> /dev/null; then
        curl -L -o "$output" "$url"
    elif command -v wget &> /dev/null; then
        wget -O "$output" "$url"
    fi
}

# Main installation function
install_asqli() {
    print_message "Starting ASQLI installation..."

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    print_message "Detected OS: ${OS}, Architecture: ${ARCH}"

    # Get latest version
    print_message "Fetching latest version..."
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi

    print_message "Latest version: ${VERSION}"

    # Construct download URL
    if [ "$OS" = "windows" ]; then
        ARCHIVE_NAME="${BINARY_NAME}-${OS}-${ARCH}-${VERSION}.zip"
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        ARCHIVE_NAME="${BINARY_NAME}-${OS}-${ARCH}-${VERSION}.tar.gz"
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE_NAME}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    print_message "Downloading ${ARCHIVE_NAME}..."
    download_file "$DOWNLOAD_URL" "${TMP_DIR}/${ARCHIVE_NAME}"

    # Extract archive
    print_message "Extracting archive..."
    if [ "$OS" = "windows" ]; then
        unzip -q "${TMP_DIR}/${ARCHIVE_NAME}" -d "$TMP_DIR"
    else
        tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "$TMP_DIR"
    fi

    # Check if we need sudo for installation
    if [ -w "$INSTALL_DIR" ]; then
        SUDO=""
    else
        if command -v sudo &> /dev/null; then
            print_warning "Installation directory requires elevated privileges"
            SUDO="sudo"
        else
            print_error "Cannot write to ${INSTALL_DIR} and sudo is not available"
            print_message "Try running with: INSTALL_DIR=~/.local/bin $0"
            exit 1
        fi
    fi

    # Install binary
    print_message "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    $SUDO mkdir -p "$INSTALL_DIR"
    $SUDO mv "${TMP_DIR}/${BINARY_FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    # Verify installation
    if command -v ${BINARY_NAME} &> /dev/null; then
        print_message "Installation successful! 🎉"
        print_message "Run '${BINARY_NAME} --version' to verify"
        print_message "Run '${BINARY_NAME} --help' to get started"
    else
        print_warning "Installation completed, but ${BINARY_NAME} is not in PATH"
        print_message "Add ${INSTALL_DIR} to your PATH or run: export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Run installation
install_asqli
