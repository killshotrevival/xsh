# Contributing to XSH

Thank you for your interest in contributing to XSH! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

When filing an issue, include:

- **Clear title** describing the problem
- **Steps to reproduce** the behavior
- **Expected behavior** vs. actual behavior
- **Environment details** (OS, Go version, xsh version)
- **Relevant logs** with `--debug` flag enabled

### Suggesting Features

Feature requests are welcome! Please include:

- **Use case** — What problem does this solve?
- **Proposed solution** — How should it work?
- **Alternatives considered** — Other approaches you've thought about

### Pull Requests

1. **Fork** the repository
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes following our coding standards
4. Write tests for new functionality
5. Run checks before submitting:

```bash
# For installing golint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Linting the code
make lint
```

6. Commit with clear messages following Conventional Commits:

```bash
feat: add custom port support for jump hosts
fix: handle empty region gracefully
docs: update README with new examples
```

7. Push and create a Pull Request


### Development Setup

#### Prerequisites
- Go 1.21 or higher
- SQLite3
- Git

#### Building from Source
```bash
git clone https://github.com/killshotrevival/xsh.git
cd xsh
make build
```

#### Running Tests
```bash
make test
```

#### Project Structure

```
xsh/
├── cmd/           # CLI commands (Cobra)
├── internal/      # Internal packages
│   ├── config/    # Configuration management
│   ├── db/        # Database operations
│   ├── host/      # Host management
│   ├── identity/  # SSH identity management
│   ├── region/    # Region management
│   ├── table/     # Table output formatting
│   └── tag/       # Tag operations
└── docs/ 
```

#### Coding Standards
- Follow standard Go conventions (Effective Go)
- Use gofmt for formatting
- Add comments for exported functions
- Handle errors explicitly — don't ignore them
- Keep functions focused and small


#### Commit Guidelines
- Use present tense ("add feature" not "added feature")
- Use imperative mood ("move cursor to..." not "moves cursor to...")
- Reference issues and PRs where appropriate


## Questions?
Feel free to open an issue with the question label if you need help.

--- 
Thank you for contributing! 🎉