# XSH

> A powerful CLI tool for managing SSH connections across massive clusters of machines.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

XSH extends SSH functionality by providing a unified interface for storing, managing, and executing SSH connections. It eliminates the need to remember complex SSH commands, IP addresses, and configuration details by storing everything in a local SQLite database.

## Features

- **Centralized Configuration** — Store all SSH connection details in a structured SQLite database
- **Simple Identifiers** — Connect to hosts using easy-to-remember names instead of IP addresses
- **Jump Host Support** — Built-in support for SSH jump hosts (bastion servers)
- **Identity Management** — Manage multiple SSH identity files with ease
- **Region Tagging** — Organize hosts by region with custom slugs
- **Flexible Output** — Retrieve connection details in various formats for scripting

## Installation

```bash
git clone https://github.com/yourusername/xsh.git
cd xsh
go build -o xsh .
```

## Usage

### Connect to a Host

```bash
# Connect using host name
xsh ssh webserver-01

# Connect using IP address
xsh ssh 192.168.1.100

# Print the SSH command without executing
xsh ssh webserver-01 --dry-run
```

### Store Host Information

```bash
# Add a new host
xsh put host --name webserver-01 --address 192.168.1.100 --user root

# Add a host with jump host
xsh put host --name private-db --address 10.0.0.50 --user admin --jumphost bastion-01

# Add an identity file
xsh put identity --name production-key --path ~/.ssh/prod_rsa
```

### Retrieve Information

```bash
# Get host details
xsh get host webserver-01

# Use in scripts with other SSH commands
ssh -AJ $(xsh get jumphost bastion-01) root@10.0.0.50
```

## Database Schema

XSH uses SQLite to store configuration with the following structure:

| Table | Description |
|-------|-------------|
| `hosts` | Host connection details (address, user, identity, jumphost) |
| `identities` | SSH identity files (name, path) |
| `regions` | Geographic regions for organization (name, slug) |

## Configuration

The database is stored at `~/.xsh/config.db` by default. Override with:

```bash
export XSH_DB_PATH=/custom/path/config.db
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
