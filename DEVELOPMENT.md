## Development Guide

### Prerequisites

- Go 1.21 or later
- Slack workspace for testing
- Slack app with appropriate bot tokens

### Initial Setup

1. **Install Go** (if not already installed):
   ```bash
   sudo ./install-go.sh
   ```

2. **Run setup script**:
   ```bash
   ./setup.sh
   ```

   This will:
   - Download dependencies
   - Run tests
   - Build local binaries

### Project Structure

```
slka/
├── cmd/
│   ├── slka-read/         # Read CLI entry point
│   └── slka-write/        # Write CLI entry point
├── internal/              # Internal packages (not importable)
│   ├── config/           # Configuration management
│   ├── auth/             # Authentication (future)
│   ├── approval/         # Human approval system
│   ├── links/            # Link parsing/formatting
│   ├── output/           # JSON output formatting
│   └── slack/            # Slack API wrapper
│       ├── client.go     # Client interface
│       ├── mock_client.go # Mock for testing
│       ├── channels.go   # Channel operations
│       └── users.go      # User operations
├── pkg/                  # Public packages
│   ├── read/             # slka-read commands
│   └── write/            # slka-write commands
├── go.mod                # Go module definition
├── Makefile             # Build automation
└── README.md            # User documentation
```

### Development Workflow

#### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/links/
go test -v ./internal/config/
```

#### Building

```bash
# Build for current platform
make build-local

# Build for all platforms
make build

# Install to GOPATH/bin
make install
```

#### Running Locally

```bash
# From source
go run ./cmd/slka-read channels list
go run ./cmd/slka-write message send general "Hello"

# From dist/
./dist/slka-read channels list
./dist/slka-write --help
```

### Test-Driven Development

This project follows TDD principles:

1. **Write tests first** - Define expected behavior
2. **Run tests** - Confirm they fail
3. **Implement** - Write minimum code to pass
4. **Refactor** - Clean up while keeping tests green
5. **Repeat** - Add more test cases

Example test structure:

```go
// internal/example/example_test.go
package example

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExampleFunction(t *testing.T) {
    result := ExampleFunction("input")
    assert.Equal(t, "expected", result)
}
```

### Mocking Slack API

Tests use a mock Slack client to avoid real API calls:

```go
// In your test
mockClient := new(slack.MockClient)
mockClient.On("GetConversations", mock.Anything).Return(
    []slack.Channel{{Name: "general"}},
    "",
    nil,
)

svc := slack.NewChannelService(mockClient)
result, err := svc.List(slack.ListChannelsOptions{})

assert.NoError(t, err)
mockClient.AssertExpectations(t)
```

### Adding New Commands

1. **Create test file** (e.g., `pkg/read/foo_test.go`)
2. **Write tests** for the new command
3. **Create implementation** (e.g., `pkg/read/foo.go`)
4. **Register command** in `init()` function
5. **Run tests** to verify

Example:

```go
// pkg/read/foo.go
var fooCmd = &cobra.Command{
    Use:   "foo",
    Short: "Do something",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    RootCmd.AddCommand(fooCmd)
}
```

### Configuration

Test configuration is loaded from `~/.config/slka/config.json`:

```json
{
  "read_token": "xoxb-test-read-token",
  "write_token": "xoxb-test-write-token",
  "user_token": "xoxp-test-user-token",
  "require_approval": true
}
```

For testing, you can override with environment variables:

```bash
export SLKA_READ_TOKEN="xoxb-test-token"
export SLKA_WRITE_TOKEN="xoxb-test-token"
go test ./...
```

### Debugging

Enable verbose output:

```bash
# Run with verbose logging
go run ./cmd/slka-read channels list --output-pretty

# Debug tests
go test -v -run TestSpecificTest ./internal/config/
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run linter: `make lint`
- Keep functions small and focused
- Write descriptive test names

### Common Issues

#### Go not found
```bash
sudo ./install-go.sh
source ~/.profile  # or restart terminal
```

#### Import errors
```bash
go mod tidy
go mod download
```

#### Test failures
- Check Slack token configuration
- Verify mock expectations match implementation
- Use `-v` flag for verbose test output

### Release Process

1. Update version in code
2. Tag release: `git tag v1.0.0`
3. Build all platforms: `make build`
4. Create GitHub release with binaries
5. Update README with release notes

### Contributing

1. Write tests first
2. Keep changes focused
3. Update documentation
4. Ensure all tests pass
5. Run linter

### Resources

- [Slack API Documentation](https://api.slack.com/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Testify Testing Library](https://github.com/stretchr/testify)
- [Go Testing](https://golang.org/pkg/testing/)
