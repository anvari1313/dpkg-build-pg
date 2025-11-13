# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple HTTP server written in Go that can be packaged as a Debian package (.deb).
The server reads its configuration from a YAML file and responds with a configurable message.
The project demonstrates how to package Go applications for Debian-based systems.

## Development Commands

### Building
```bash
go build -o dpkg-build-pg .
```

### Running
```bash
go run .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestName ./path/to/package

# Run tests with coverage
go test -cover ./...
```

### Dependency Management
```bash
# Add a new dependency
go get <package>

# Update dependencies
go get -u ./...

# Tidy up go.mod and go.sum
go mod tidy

# Verify dependencies
go mod verify
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Vet code for suspicious constructs
go vet ./...
```

## Architecture

### Application Structure

- **main.go**: HTTP server that reads configuration and serves a single endpoint
- **config.yaml**: Configuration file with server settings and response message
- **debian/**: Debian packaging files for creating .deb packages
- **.github/workflows/**: GitHub Actions for automated builds and releases

### Configuration

The application looks for configuration in these locations (in order):
1. `/etc/dpkg-build-pg/config.yaml` (production)
2. `./config.yaml` (development)

Configuration format:
```yaml
server:
  port: 8080
  host: "0.0.0.0"
message: "Your message here"
```

### Package Structure

When installed via .deb package:
- Binary: `/usr/bin/dpkg-build-pg`
- Config: `/etc/dpkg-build-pg/config.yaml`
- Service: `/usr/lib/systemd/system/dpkg-build-pg.service`
