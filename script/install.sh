#!/usr/bin/env bash
set -e

# Installation script for gh-repomon GitHub CLI extension
# This script automatically detects OS and architecture and installs the appropriate binary

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# GitHub repository information
OWNER="hazadus"
REPO="gh-repomon"

# Function to print colored messages
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     OS="linux";;
        Darwin*)    OS="darwin";;
        MINGW*|MSYS*|CYGWIN*) OS="windows";;
        *)
            error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   ARCH="amd64";;
        aarch64|arm64)  ARCH="arm64";;
        *)
            error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Get the latest release version
get_latest_version() {
    info "Fetching latest release version..."
    VERSION=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        error "Could not determine latest version"
        exit 1
    fi

    info "Latest version: ${VERSION}"
}

# Download and install binary
install_binary() {
    local binary_name="gh-repomon-${OS}-${ARCH}"

    if [ "$OS" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi

    local download_url="https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${binary_name}"

    info "Downloading ${binary_name}..."
    info "URL: ${download_url}"

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf ${TEMP_DIR}" EXIT

    # Download binary
    if ! curl -L -o "${TEMP_DIR}/gh-repomon" "${download_url}"; then
        error "Failed to download binary"
        exit 1
    fi

    # Make binary executable
    chmod +x "${TEMP_DIR}/gh-repomon"

    # Determine installation directory
    # For GitHub CLI extensions, the binary should be in the current directory
    # or in a location specified by GH_EXTENSION_PATH
    if [ -n "$GH_EXTENSION_PATH" ]; then
        INSTALL_DIR="$GH_EXTENSION_PATH"
    else
        INSTALL_DIR="$(pwd)"
    fi

    info "Installing to ${INSTALL_DIR}..."

    # Copy binary to installation directory
    cp "${TEMP_DIR}/gh-repomon" "${INSTALL_DIR}/gh-repomon"

    info "Installation complete!"
    info "You can now use: gh repomon --help"
}

# Verify installation
verify_installation() {
    info "Verifying installation..."

    if command -v gh-repomon &> /dev/null; then
        info "✓ gh-repomon is installed and available in PATH"
        gh-repomon --version || true
    else
        warn "gh-repomon is not in PATH, but installed successfully"
        info "You can use it via: gh repomon"
    fi
}

# Main installation flow
main() {
    info "Starting gh-repomon installation..."

    detect_os
    info "Detected OS: ${OS}"

    detect_arch
    info "Detected architecture: ${ARCH}"

    get_latest_version

    install_binary

    verify_installation

    info "Installation successful! 🎉"
}

# Run main function
main
