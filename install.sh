#!/bin/bash

set -e

# Check root
if [ "$EUID" -ne 0 ]; then
    echo "Run with: sudo ./install.sh"
    exit 1
fi

# Project directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Docker
if command -v docker >/dev/null 2>&1; then
    echo "Docker: found"
else
    echo "Docker: installing"

    apt-get update
    apt-get install -y docker.io

    systemctl enable docker
    systemctl start docker
fi

# Caddy
if command -v caddy >/dev/null 2>&1; then
    echo "Caddy: found"
else
    echo "Caddy: installing"

    apt-get update
    apt-get install -y \
        debian-keyring \
        debian-archive-keyring \
        apt-transport-https \
        curl \
        gnupg

    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
        | gpg --dearmor \
        | tee /usr/share/keyrings/caddy-stable-archive-keyring.gpg >/dev/null

    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
        | tee /etc/apt/sources.list.d/caddy-stable.list

    apt-get update
    apt-get install -y caddy

    systemctl enable caddy
    systemctl start caddy
fi

# Go
if ! command -v go >/dev/null 2>&1; then
    echo "Go is required"
    exit 1
fi

# Build
cd "$SCRIPT_DIR"

echo "Building Dockyard..."
go build -o dockyard .

# Install
install -m 755 dockyard /usr/local/bin/dockyard

echo "Dockyard installed"

# Verify
dockyard --help >/dev/null 2>&1 || true

echo "Done"