#!/bin/sh
set -e

REPO="killshotrevival/xsh"
INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="xsh"

# Detect OS
detect_os() {
    os="$(uname -s)"
    case "$os" in
        Darwin) echo "mac" ;;
        Linux)  echo "linux" ;;
        *)
            echo "Unsupported OS: $os" >&2
            exit 1
            ;;
    esac
}

# Fetch the latest release tag from GitHub API
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4
    else
        echo "Error: curl or wget is required" >&2
        exit 1
    fi
}

# Download the binary for the detected OS
download_binary() {
    version="$1"
    os_suffix="$2"
    url="https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}-${os_suffix}"

    echo "Downloading ${BINARY_NAME}-${os_suffix} (${version})..."

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "${INSTALL_DIR}/${BINARY_NAME}" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${INSTALL_DIR}/${BINARY_NAME}" "$url"
    fi
}

main() {
    os_suffix="$(detect_os)"

    echo "Fetching latest release... for ${os_suffix} OS"
    version="$(get_latest_version)"
    if [ -z "$version" ]; then
        echo "Error: could not determine latest release version" >&2
        exit 1
    fi
    echo "Latest version: ${version}"

    mkdir -p "$INSTALL_DIR"

    download_binary "$version" "$os_suffix"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    echo "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"

    # Ensure INSTALL_DIR is on PATH for this session
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*) ;;
        *) export PATH="${INSTALL_DIR}:${PATH}" ;;
    esac

    echo "Running '${BINARY_NAME} init'..."
    "${INSTALL_DIR}/${BINARY_NAME}" init

    echo ""
    echo "Installation complete!"
    echo "Make sure ${INSTALL_DIR} is in your PATH. Add this to your shell profile if needed:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
}

main
