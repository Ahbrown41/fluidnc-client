# FluidNC CLI Tool

A command-line tool for interacting with the FluidNC Web API, supporting file uploads, WebSocket commands, and interactive sessions.

## Features

- **File Upload**: Upload files to FluidNC flash storage or SD card
- **WebSocket Commands**: Send individual commands or start interactive sessions
- **Flexible Configuration**: YAML config files, environment variables, and CLI flags
- **Multiple Output Formats**: Text and JSON output options
- **Cross-Platform**: Builds for Linux, Windows, and macOS

## Installation

### From Source

```bash
git clone https://github.com/Ahbrown41/fluidnc-client
cd fluidnc-client
make build (go build -o build/fluidnc-cli .)
```
