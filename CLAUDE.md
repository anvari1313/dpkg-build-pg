# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project for building Debian packages (dpkg) for PostgreSQL. The project uses Go 1.25.

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

This project is in its initial phase. As the codebase grows, this section will be updated with architectural details about how the Debian package building system works.
