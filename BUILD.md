# Building and Installing the Debian Package

This document explains how to build and install the dpkg-build-pg Debian package.

## Prerequisites

On a Debian-based system (Debian, Ubuntu, etc.), install the required build tools:

```bash
sudo apt-get update
sudo apt-get install dpkg-dev debhelper golang-go build-essential devscripts
```

## Building the Package

### Option 1: Using the build script (Recommended)

Simply run the provided build script:

```bash
./build-deb.sh
```

This will:
1. Check for required tools (Go and dpkg-buildpackage)
2. Clean previous builds
3. Compile the Go binary with optimizations
4. Build the Debian package using the pre-built binary
5. Display installation instructions

### Option 2: Manual build

First, compile the Go binary:

```bash
go build -o dpkg-build-pg -ldflags="-s -w" .
chmod +x dpkg-build-pg
```

Then build the package:

```bash
dpkg-buildpackage -us -uc -b
```

The flags mean:
- `-us`: Do not sign the source package
- `-uc`: Do not sign the changes file
- `-b`: Build binary package only

## Installing the Package

After building, the `.deb` file will be created in the parent directory:

```bash
sudo dpkg -i ../dpkg-build-pg_1.0.0_amd64.deb
```

Note: The architecture suffix (`amd64`) may vary depending on your system.

## Post-Installation

### Check Service Status

```bash
sudo systemctl status dpkg-build-pg
```

### View Logs

```bash
sudo journalctl -u dpkg-build-pg -f
```

### Test the Server

```bash
curl http://localhost:8080
```

You should see: `Hello from dpkg-build-pg server!`

## Managing the Service

### Start/Stop/Restart

```bash
sudo systemctl start dpkg-build-pg
sudo systemctl stop dpkg-build-pg
sudo systemctl restart dpkg-build-pg
```

### Enable/Disable Auto-start

```bash
sudo systemctl enable dpkg-build-pg   # Start on boot
sudo systemctl disable dpkg-build-pg  # Don't start on boot
```

## Uninstalling

To remove the package:

```bash
sudo apt-get remove dpkg-build-pg
```

To remove including configuration files:

```bash
sudo apt-get purge dpkg-build-pg
```

## Package Contents

The package installs:
- Binary: `/usr/bin/dpkg-build-pg`
- Systemd service: `/lib/systemd/system/dpkg-build-pg.service`

## Troubleshooting

### Service won't start

Check the logs:
```bash
sudo journalctl -u dpkg-build-pg -n 50
```

### Port already in use

If port 8080 is already in use, you'll need to:
1. Stop the conflicting service, or
2. Modify the port in `main.go` and rebuild the package

### Permission issues

The service runs as the `www-data` user. Ensure this user exists:
```bash
id www-data
```

## Build Process

The build follows a two-stage process:

1. **Stage 1 - Compile Go Binary**: The Go source code is compiled into a static binary with optimizations (`-ldflags="-s -w"` strips debug info and symbol table to reduce size).

2. **Stage 2 - Package Binary**: The pre-compiled binary is packaged into a `.deb` file along with the systemd service configuration.

This approach has several advantages:
- Faster package builds (no recompilation needed)
- Consistent binaries across different package builds
- Smaller package build dependencies (only needs `debhelper`, not `golang-go`)
- Better control over Go compilation flags

## Development

To make changes:
1. Edit the source code (e.g., `main.go`)
2. Update `debian/changelog` with your changes
3. Rebuild the package using `./build-deb.sh`
4. Reinstall with `sudo dpkg -i ../dpkg-build-pg_*.deb`
