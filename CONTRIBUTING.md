<!--
Copyright (c) 1898 & Co.
SPDX-License-Identifier: Apache-2.0
-->

# Contributing

Thank you for your interest in contributing to the Armis Go SDK!

## Project Structure

```
.
├── armis/                      # SDK source code
│   ├── client.go              # Client initialization and HTTP helpers
│   ├── auth.go                # Authentication logic
│   ├── errors.go              # Sentinel errors
│   ├── model_*.go             # API request/response structs
│   ├── *_test.go              # Unit tests
│   └── ...                    # Resource-specific files (users.go, policies.go, etc.)
├── tools/                      # Development tooling
├── Makefile                    # Build and test commands
└── go.mod                      # Module definition
```

## Development Setup

### Prerequisites

- Go 1.22 or later
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [pre-commit](https://pre-commit.com/)

### Getting Started

Clone the repository:

```sh
git clone https://github.com/1898andCo/armis-sdk-go.git
cd armis-sdk-go
```

Install dependencies and pre-commit hooks:

```sh
make deps
make hooks
```

## Testing

> [!NOTE]
> We recommend using a tool such as Postman or Paw to quickly develop and test the Armis API.
> This enables you to quickly debug requests to and responses from the API before implementing them in Go.
>
> For more information on the Armis Centrix platform, refer to the Armis user guide.

### Running Tests

```sh
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific package tests
go test ./armis -v

# Run specific test function
go test ./armis -v -run TestAuthenticate
```

### Writing Tests

All tests follow Go best practices:

- **Table-driven tests** for comprehensive coverage
- **Parallel execution** with `t.Parallel()` for performance
- **Clear test names** describing what is being tested
- **Test helpers** in `testhelpers_test.go` for common setup

Example:

```go
func TestMyFunction(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        input    MyInput
        expected MyOutput
    }{
        {
            name:     "valid input",
            input:    MyInput{Field: "value"},
            expected: MyOutput{Result: "expected"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Adding New API Resources

When adding support for a new Armis API resource:

1. **Create model structs** in a new `model_<resource>.go` file defining the API request/response structures
2. **Implement CRUD operations** in a new `<resource>.go` file with methods on `*Client`
3. **Add sentinel errors** to `errors.go` for resource-specific validation errors
4. **Write comprehensive tests** in `<resource>_test.go` using the test helpers

## Code Quality

### Linting

```sh
# Run linter
make lint

# Format code
make fmt
```

### Pre-commit Hooks

Pre-commit hooks run automatically on commit. To run manually:

```sh
make pre-commit
```

## Pull Request Process

1. Fork the repository and create a feature branch
2. Make your changes with appropriate tests
3. Ensure all tests pass: `make test`
4. Ensure linting passes: `make lint`
5. Submit a pull request with a clear description of the changes

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
