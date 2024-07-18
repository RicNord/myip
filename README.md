# MyIp CLI Tool

## Overview

The `MyIp` CLI tool helps you manage and monitor your public IP address. It can
fetch your current public IP, display known IP aliases from a configuration
file, and notify you when your IP address changes. The tool supports
notifications using `notify-send` for linux and simple console output for other
platforms.

## Features

- Fetch and display your public IP address.
- Display the alias for your public IP if it exists in the configuration file.
- List all known IP aliases.
- Monitor for IP changes and notify when a new IP or alias is detected.

## Installation

1. Clone the repository

2. Build the project:

    ```sh
    make build
    ```

## Usage

### Configuration

Ensure the configuration file `.myip.json` exists in your home directory with
the following format:

```json
{
    "aliases": [
        {"alias": "Home", "ip": "203.0.113.1"},
        {"alias": "Office", "ip": "203.0.113.2"}
    ]
}
```

### Commands

```sh
myip --help
```
