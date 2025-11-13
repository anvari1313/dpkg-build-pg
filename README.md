# dpkg-build-pg

A simple HTTP server written in Go that can be packaged and distributed as a Debian package.

## Features

- Simple HTTP server with configurable response message
- YAML-based configuration
- Systemd service integration
- Debian package (.deb) for easy installation
- Automated builds via GitHub Actions

## Quick Start

### Development

1. Clone the repository
2. Run the server:
   ```bash
   go run .
   ```
3. Test:
   ```bash
   curl http://localhost:8080
   ```

### Configuration

Edit `config.yaml` to customize:

```yaml
server:
  port: 8080
  host: "0.0.0.0"
message: "Hello from dpkg-build-pg server!"
```

## Building and Installing

See [BUILD.md](BUILD.md) for detailed instructions on building and installing the Debian package.

### Quick Build

```bash
./build-deb.sh
sudo dpkg -i ../dpkg-build-pg_*.deb
```

## Usage

After installation, the service runs automatically.

### Check Status
```bash
sudo systemctl status dpkg-build-pg
```

### Customize Message
```bash
sudo nano /etc/dpkg-build-pg/config.yaml
sudo systemctl restart dpkg-build-pg
```

### Test
```bash
curl http://localhost:8080
```

## GitHub Actions

This project includes automated workflows:

- **Build**: Creates .deb package on every push to main
- **Release**: Creates GitHub releases when you push version tags

See [.github/workflows/README.md](.github/workflows/README.md) for details.

## Project Structure

```
dpkg-build-pg/
├── main.go              # HTTP server code
├── config.yaml          # Configuration file
├── go.mod               # Go module definition
├── debian/              # Debian packaging files
│   ├── control          # Package metadata
│   ├── changelog        # Version history
│   ├── rules            # Build instructions
│   └── ...
├── .github/workflows/   # GitHub Actions
└── build-deb.sh         # Build script
```

## License

This project is a demonstration of Debian package creation for Go applications.
