# XSH Implementation Plan

## Table of Contents

- [Overview](#overview)
- [Storage Layer](#storage-layer)
- [Database Schema](#database-schema)
- [CLI Operations](#cli-operations)

---

## Overview

XSH is a standalone CLI tool designed to simplify SSH connection management. The tool provides:

- Persistent storage of host configurations in a local database
- Dynamic retrieval of connection details in multiple formats
- Streamlined SSH connections with minimal user input

---

## Storage Layer

### SQLite Database

All configuration data is persisted using **SQLite**, chosen for its:

- Zero-configuration setup
- Single-file portability
- Reliable ACID compliance
- Excellent Go driver support (`modernc.org/sqlite`)

**Default location:** `~/.xsh/config.db`

---

## Database Schema

### Hosts Table

Stores SSH host connection details.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `name` | TEXT | User-friendly identifier for the host |
| `address` | TEXT | IP address or domain name |
| `user` | TEXT | SSH username |
| `region_id` | UUID | Foreign key → `regions.id` |
| `identity_id` | UUID | Foreign key → `identities.id` |
| `jumphost_id` | UUID | Foreign key → `hosts.id` (self-referencing) |

> **Note:** Jump hosts are regular hosts referenced via `jumphost_id`, enabling recursive jump host chains.

### Identities Table

Manages SSH identity files (private keys).

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `name` | TEXT | Descriptive name for the identity |
| `path` | TEXT | Absolute path to the identity file |

### Regions Table

Organizes hosts by geographic or logical regions.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `name` | TEXT | Full region name (e.g., "US East") |

---

## CLI Operations

### 1. `xsh ssh <identifier>`

Establishes an SSH connection using stored configuration.

**Behavior:**
- Accepts host `name` or `address` as the identifier
- Automatically resolves identity file, user, and jump host
- Constructs and executes the complete SSH command

**Flags:**
| Flag | Description |
|------|-------------|
| `--dry-run` | Print the SSH command without executing |

**Example:**
```bash
xsh ssh webserver-01
xsh ssh 192.168.1.100 --dry-run
```

### 2. `xsh put <resource>`

Stores configuration data in the database.

**Supported Resources:**
- `host` — Add or update host configuration
- `identity` — Register an SSH identity file
- `region` — Define a region

**Example:**
```bash
xsh put host --name bastion --address 10.0.0.1 --user admin
xsh put identity --name prod-key --path ~/.ssh/prod_rsa
```

### 3. `xsh get <resource> <identifier>`

Retrieves configuration data from the database.

**Design Goals:**
- Output format suitable for shell interpolation
- Enable seamless integration with native SSH commands

**Example:**
```bash
# Direct retrieval
xsh get host webserver-01

# Shell interpolation for complex SSH commands
ssh -AJ $(xsh get jumphost bastion-01) root@10.0.0.50
```

---

## Next Steps

- [ ] Initialize Go module structure
- [ ] Implement SQLite database layer
- [ ] Build CLI command parser
- [ ] Add unit and integration tests
- [ ] Create installation documentation
