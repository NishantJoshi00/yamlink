# Yamlink

A fast, lightweight URL mapping and redirection system that dynamically manages URL shortcuts through YAML configuration files. Perfect for creating and managing custom URL shorteners for internal tools, documentation, or any web resources.

[![asciicast](https://asciinema.org/a/vmPkgMpJbqfnVwiliutXswU92.svg)](https://asciinema.org/a/vmPkgMpJbqfnVwiliutXswU92)

## Description

Yamlink solves the common challenge of managing and accessing numerous URLs within an organization or personal workflow by:

- Creating human-readable shortcuts for complex URLs
- Supporting hierarchical URL organization through YAML
- Providing real-time configuration updates without service restart
- Offering both CLI and server components for flexible usage

## Installation

1. **Prerequisites**

   - Go 1.22.4 or higher
   - Make (optional, for using Makefile commands)

2. **Build from Source**

   ```bash
   # Clone the repository
   git clone https://github.com/NishantJoshi00/yamlink
   cd yamlink

   # Build the binaries
   go build ./cmd/yamlink    # Server component
   go build ./cmd/shelinks   # CLI component
   ```

3. **Install System-wide (Optional)**
   ```bash
   # Move binaries to system path
   sudo mv yamlink /usr/local/bin/
   sudo mv shelinks /usr/local/bin/
   ```

## Usage

### Server Mode (yamlink)

1. **Create Configuration Files**

   ```yaml
   # config.yaml
   host: localhost
   port: 8080
   map_file: map.yaml
   refresh_interval: 5 # seconds
   ```

   ```yaml
   # map.yaml
   github:
     profile: https://github.com/NishantJoshi00
   docs:
     - https://docs.example.com
   ```

2. **Start the Server**

   ```bash
   CONFIG_FILE=config.yaml ./yamlink
   ```

3. **Access URLs**
   - Visit `http://localhost:8080/github/profile`
   - Visit `http://localhost:8080/docs/0`

### CLI Mode (shelinks)


1. **Set up Configuration**

   ```yaml
   # ~/.shelinks.yaml
   gs: git status
   gp: git push
   ```

2. **Use in Shell**

   ```bash
   # For ZSH
   source scripts/init.zsh /path/to/shelinks ~/.shelinks.yaml

   # For Fish
   source scripts/init.fish /path/to/shelinks ~/.shelinks.yaml
   ```

3. **Use Shortcuts**
   ```bash
   s/gs    # Expands to 'git status'
   s/gp    # Expands to 'git push'
   ```

## Features

- **Dynamic URL Mapping**: Support for nested URL structures up to 3 levels deep
- **Real-time Updates**: Configuration changes are automatically detected and applied
- **Multiple Access Methods**:
  - Server mode for web-based access
  - CLI mode for shell integration
- **Flexible Configuration**:
  - Support for both single URLs and arrays of URLs
  - Custom refresh intervals for configuration updates
- **Shell Integration**: Native support for zsh and fish shells
- **Logging**: Structured JSON logging with configurable log levels

## Contributing Guidelines

1. **Issue First**: Create or find an issue before starting work
2. **Issue Tags**: Use descriptive tags:
   - [BUG] for bug reports
   - [FEATURE] for feature requests
   - [DOCS] for documentation improvements
3. **Work Assignment**: Don't work on issues already assigned to others
4. **Testing**: Ensure all tests pass by running `make test`
5. **Code Style**: Follow Go standard formatting guidelines

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
