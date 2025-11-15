# Contributing to Configly

Thank you for your interest in contributing to Configly! We welcome contributions from the community.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples (code snippets)
- Describe the behavior you observed and what you expected
- Include your environment details (Go version, OS, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Include code examples showing how the feature would be used
- Explain why this enhancement would be useful

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write tests** for your changes
3. **Ensure tests pass**: `go test ./...`
4. **Update documentation** if needed
5. **Follow the Go coding style**
6. **Write clear commit messages**

#### Pull Request Process

1. Update the README.md with details of changes if applicable
2. Add tests for new functionality
3. Ensure all tests pass
4. Update documentation/examples if needed
5. The PR will be merged once reviewed and approved

## Development Setup

### Prerequisites

- Go 1.24.1 or later
- Git

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/configly.git
cd configly

# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Project Structure

```
configly/
├── sources/          # Configuration source implementations
│   ├── source.go    # Source interface
│   ├── env.go       # Environment variable source
│   ├── file.go      # File source (JSON, YAML, .env)
│   └── mock.go      # Mock source for testing
├── load.go          # Main loader implementation
├── load_test.go     # Loader tests
└── cmd/             # Example applications
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...

# Run specific test
go test -run TestLoaderLoad

# Run tests in a specific package
go test ./sources
```

### Writing Tests

- Write table-driven tests when possible
- Use subtests for multiple test cases
- Test both success and error cases
- Aim for high test coverage
- Include edge cases

Example test structure:
```go
func TestFeature(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // test code
    })

    t.Run("error case", func(t *testing.T) {
        // test code
    })
}
```

## Code Style

### Go Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format your code
- Run `go vet` to catch common issues
- Follow standard Go naming conventions
- Write descriptive variable and function names
- Add comments for exported functions and types

### Documentation

- Document all exported types, functions, and methods
- Use complete sentences in comments
- Start comments with the name of the thing being described
- Include examples in documentation when helpful

Example:
```go
// LoaderConfig contains configuration options for creating a new Loader.
// It specifies which sources to use and in what priority order.
type LoaderConfig struct {
    TagKey  string           // The struct tag key to use (defaults to "configly")
    Sources []sources.Source // Configuration sources in priority order
}
```

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests when relevant

Examples:
```
Add support for custom validation functions

Fix panic when parsing malformed YAML files (#123)

Update documentation for file source
```

## Adding New Features

When adding new features:

1. **Open an issue first** to discuss the feature
2. **Wait for approval** before starting work
3. **Write tests** that cover the new functionality
4. **Update documentation** including README.md
5. **Add examples** if applicable
6. **Ensure backward compatibility** when possible

## Adding New Configuration Sources

To add a new configuration source:

1. Implement the `sources.Source` interface in `sources/` directory
2. Add comprehensive tests in a `*_test.go` file
3. Update the main README.md with usage examples
4. Add the source to the "Supported Configuration Sources" section

Example:
```go
// sources/consul.go
package sources

type ConsulSource struct {
    client *consul.Client
}

func FromConsul(address string) (*ConsulSource, error) {
    // implementation
}

func (s *ConsulSource) Name() string {
    return "consul:" + s.address
}

func (s *ConsulSource) GetValue(key string) (string, bool, error) {
    // implementation
}
```

## Questions?

- Open a [Discussion](https://github.com/zanedma/configly/discussions) for general questions
- Open an [Issue](https://github.com/zanedma/configly/issues) for bugs or feature requests

## Recognition

Contributors will be recognized in:
- The repository's contributor graph
- Release notes when applicable
- Our gratitude!

Thank you for contributing to Configly!
