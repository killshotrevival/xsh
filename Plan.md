# XSH — Roadmap

This document outlines the planned features and improvements for XSH. Each initiative is designed to extend the tool's capabilities while staying true to its core mission: simplifying SSH management at scale.

---

## Remote Backend Support

**Status:** Planned

Currently, XSH stores all connection data in a local SQLite database. While this works well for individual users, teams managing shared infrastructure often need a single source of truth for host configurations.

Remote backend support will introduce a centralized server that teams can use to store and distribute SSH connection metadata. With this feature:

- A team administrator can define hosts, regions, and connection parameters once on the remote backend.
- Team members pull the shared configuration automatically and only need to configure their own identity file mappings locally.
- Changes to the shared configuration propagate to all team members, eliminating configuration drift across machines.
- The local database will continue to work as a cache and fallback, ensuring XSH remains functional even when the remote backend is unreachable.

This feature is aimed at organizations operating large fleets of servers where consistency and collaboration around SSH access are critical.

---

## SCP Integration

**Status:** Planned

XSH already manages the details needed to establish SSH connections — host addresses, users, ports, identity files, and jump hosts. Extending this to support SCP (Secure Copy Protocol) is a natural next step.

With SCP integration, users will be able to:

- Transfer files to and from remote hosts using the same simple host identifiers they already use for SSH connections.
- Leverage jump host configurations transparently during file transfers, without manually constructing complex SCP commands.
- Perform bulk file operations across multiple hosts in a single command.

The goal is to make file transfers as effortless as connecting — no need to look up addresses, ports, or identity files.

---

## Direct SSH Config Management

**Status:** Planned

Many users maintain an `~/.ssh/config` file to define connection shortcuts, proxy rules, and identity mappings. Editing this file manually is error-prone, especially as the number of hosts grows, and requires familiarity with the SSH config syntax.

XSH will offer the ability to generate and update `~/.ssh/config` directly from its database:

- Export all or a filtered subset of hosts into properly formatted SSH config entries.
- Keep the config file in sync with the XSH database through a single command.
- Preserve any manually added entries in the config file that are not managed by XSH.

This bridges the gap between XSH's managed workflow and tools or scripts that rely on the standard SSH config file.

---

## Resource Tagging

**Status:** Planned

As the number of managed hosts grows, efficient filtering and organization become essential. Resource tagging will allow users to assign arbitrary key-value tags to hosts and identities.

Planned capabilities include:

- Assign one or more tags to any host or identity (e.g., `env:production`, `team:platform`, `os:ubuntu`).
- Filter and list resources by tag during lookups and connections.
- Combine tags with existing region-based organization for fine-grained grouping.
- Use tags in scripting workflows to target specific subsets of infrastructure.

Tagging provides a flexible, user-defined taxonomy that adapts to any team's organizational model without imposing a rigid structure.

---

*This roadmap is subject to change based on community feedback and project priorities. Contributions and feature requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) for details.*