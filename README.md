# XSH

![image](./assets/spirit_animal.png)

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
git clone https://github.com/killshotrevival/xsh.git
cd xsh

make build
```

## Usage

### System Init
```bash
xsh init
```
This command will initilise the xsh environment as well as read the following files to populate the database for configruations:

- Identities in .ssh: Will look for all the identities files present in the .ssh directory and populate them in the database
- .ssh/config (TODO): Read the config file for populating the already present host configruation
- .zshrc / .bashrc (TODO): Read the config file for populating the already present host configruation

### Add New Resources
```bash
# This command will add a new region in the database that can be mapped to hosts
xsh put region us-east-1

# This command will add a new identity file using a unique name and its complete path
xsh put identity peeyush-development /Users/ptyagi/.ssh/development

# This command will create a new host in interactive mode
xsh put host -i

# This command can be used to create a host directly without interactive mode
# example format for host.json can be created using `xsh example host` command
xsh put host -f host.json


# This command can be used for viewing a list of all the hosts present
xsh get host
```

### Connecting To A Host
```bash
xsh connect host-1

# Add the debug flag to connect in verbose mode
xsh connect host-1 --debug
```

### For More Details
Please follow [this](./docs/xsh.md) for more information

## Database Schema

XSH uses SQLite to store configuration with the following structure:

| Table | Description |
|-------|-------------|
| `hosts` | Host connection details (address, user, identity, jumphost) |
| `identities` | SSH identity files (name, path) |
| `regions` | Geographic regions for organization (name, slug) |

## Configuration

The database is stored at `~/.xsh/xsh.db` by default. Override with:

```bash
export XSH_DB_PATH=/custom/path/config.db
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
