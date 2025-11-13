#!/bin/bash
set -e

echo "Building Debian package for dpkg-build-pg..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install it with:"
    echo "  sudo apt-get install golang-go"
    exit 1
fi

# Check if required tools are installed
if ! command -v dpkg-buildpackage &> /dev/null; then
    echo "Error: dpkg-buildpackage not found. Please install it with:"
    echo "  sudo apt-get install dpkg-dev"
    exit 1
fi

if ! command -v debuild &> /dev/null; then
    echo "Warning: debuild not found. Installing build-essential and devscripts is recommended:"
    echo "  sudo apt-get install build-essential devscripts"
fi

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf debian/dpkg-build-pg
rm -f dpkg-build-pg
rm -f ../dpkg-build-pg_*.deb ../dpkg-build-pg_*.changes ../dpkg-build-pg_*.buildinfo ../dpkg-build-pg_*.tar.xz

# Build the Go binary
echo "Compiling Go binary..."
go build -o dpkg-build-pg -ldflags="-s -w" .
chmod +x dpkg-build-pg
echo "Binary compiled successfully"

# Build the package
echo "Building Debian package..."
dpkg-buildpackage -us -uc -b

echo ""
echo "Build complete!"
echo "Package created: ../dpkg-build-pg_1.0.0_$(dpkg --print-architecture).deb"
echo ""
echo "To install the package, run:"
echo "  sudo dpkg -i ../dpkg-build-pg_1.0.0_$(dpkg --print-architecture).deb"
echo ""
echo "To check the service status after installation:"
echo "  sudo systemctl status dpkg-build-pg"
echo ""
echo "To view logs:"
echo "  sudo journalctl -u dpkg-build-pg -f"
