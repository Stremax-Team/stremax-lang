# Contributing to Stremax-Lang

Thank you for your interest in contributing to Stremax-Lang! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and considerate of others.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue on GitHub with the following information:

1. A clear, descriptive title
2. A detailed description of the bug
3. Steps to reproduce the bug
4. Expected behavior
5. Actual behavior
6. Screenshots (if applicable)
7. Environment information (OS, Go version, etc.)

### Suggesting Enhancements

If you have an idea for an enhancement, please create an issue on GitHub with the following information:

1. A clear, descriptive title
2. A detailed description of the enhancement
3. The motivation behind the enhancement
4. Any potential implementation details

### Pull Requests

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes
4. Run tests to ensure your changes don't break existing functionality
5. Submit a pull request

Please include a clear description of the changes and reference any related issues.

## Development Setup

### Prerequisites

- Go 1.16 or higher
- Git

### Setting Up the Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/Stremax-Team/stremax-lang.git
   cd stremax-lang
   ```

2. Build the project:
   ```bash
   go build -o stremax ./cmd/stremax
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

## Coding Standards

- Follow Go's official [style guide](https://golang.org/doc/effective_go)
- Write clear, descriptive commit messages
- Add comments to your code where necessary
- Write tests for new features

## Documentation

- Update documentation when changing functionality
- Use clear, concise language
- Include examples where appropriate

## License

By contributing to Stremax-Lang, you agree that your contributions will be licensed under the project's MIT License. 