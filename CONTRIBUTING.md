# Contributing to TesselBox Game

Thank you for your interest in contributing to TesselBox Game! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project adheres to our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/TesselBox-game.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`

## How to Contribute

### Reporting Bugs

- Check if the issue already exists
- Include Go version and OS information
- Provide steps to reproduce
- Include relevant code snippets or error messages

### Suggesting Features

- Open an issue with the "feature request" label
- Describe the feature and its use case
- Discuss implementation approach

### Code Contributions

1. **Pick an issue** - Look for issues labeled `good first issue` or `help wanted`
2. **Discuss** - Comment on the issue to let others know you're working on it
3. **Implement** - Write your code following our standards
4. **Test** - Ensure your changes work correctly
5. **Submit** - Create a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Building

```bash
cd cmd
go build -o tesselbox .
```

### Running Tests

```bash
go test ./...
```

## Coding Standards

### Go Code Style

- Follow standard Go conventions (gofmt, golint)
- Use meaningful variable names
- Keep functions focused and small
- Add comments for exported functions

### Project Structure

```
cmd/          - Main application entry point
pkg/          - Library packages
  chest/      - Chest/inventory system
  gui/        - UI components
  save/       - Save/load system
  world/      - World management
  config/     - Configuration
  game/       - Game management
```

### Key Principles

1. **Single Responsibility** - Each function/package does one thing well
2. **Error Handling** - Always handle errors, never ignore them
3. **Resource Management** - Close files, release resources properly
4. **Thread Safety** - Use mutexes where concurrent access occurs

## Commit Guidelines

### Commit Message Format

```
type(scope): short description

Longer description if needed

Closes #issue-number
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code restructuring
- `test`: Adding or fixing tests
- `chore`: Maintenance tasks

### Examples

```
feat(gui): add scrollable inventory view

fix(save): resolve data corruption on exit

docs(readme): update installation instructions
```

## Pull Request Process

1. **Branch naming**: `feature/description`, `fix/description`, `hotfix/description`
2. **One concern per PR** - Don't mix unrelated changes
3. **Update documentation** - If behavior changes, update docs
4. **Keep commits clean** - Squash if needed before review
5. **Request review** - Tag relevant maintainers

### PR Checklist

- [ ] Code compiles without errors
- [ ] Follows coding standards
- [ ] Includes tests if applicable
- [ ] Documentation updated
- [ ] Commit messages are clear

## Reporting Issues

When reporting issues, please include:

- **Title**: Clear, concise summary
- **Description**: Detailed explanation
- **Environment**: Go version, OS, game version
- **Steps to Reproduce**: Numbered list
- **Expected vs Actual**: What you expected vs what happened
- **Logs/Errors**: Relevant output or screenshots

## Questions?

- Open a discussion in GitHub
- Tag maintainers in issues

## License

By contributing, you agree that your contributions will be licensed under the project's license.

---

Thank you for contributing to TesselBox Game!
