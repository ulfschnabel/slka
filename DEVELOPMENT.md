## Development Guide

### Prerequisites

- Go 1.21 or later
- Slack workspace for testing
- Slack app with appropriate tokens (user or bot)

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
   - Build local binary

### Project Structure

```
slka/
├── cmd/
│   └── slka/              # Main CLI entry point
├── internal/              # Internal packages (not importable)
│   ├── config/           # Configuration management
│   ├── approval/         # Human approval system
│   ├── links/            # Link parsing/formatting
│   ├── output/           # JSON output formatting
│   └── slack/            # Slack API wrapper
│       ├── client.go     # Client interface
│       ├── mock_client.go # Mock for testing
│       ├── channels.go   # Channel operations
│       ├── dms.go        # DM operations
│       ├── reactions.go  # Reaction operations
│       └── users.go      # User operations
├── pkg/                  # Public packages
│   ├── read/             # Read commands
│   │   ├── channels.go   # Channel queries
│   │   ├── dm.go         # DM queries
│   │   ├── reaction.go   # Reaction queries
│   │   └── users.go      # User queries
│   └── write/            # Write commands
│       ├── channels.go   # Channel management
│       ├── dm.go         # DM sending
│       ├── message.go    # Message operations
│       ├── reaction.go   # Reaction management
│       └── config.go     # Config management
├── go.mod                # Go module definition
├── Makefile             # Build automation
├── .goreleaser.yaml     # Release configuration
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
go test -v ./internal/slack/
go test -v ./pkg/read/
```

#### Building

```bash
# Build for current platform
make build-local

# Test snapshot build with goreleaser
goreleaser build --snapshot --clean

# The binary will be in dist/
./dist/slka channels list
```

#### Running Locally

```bash
# From source
go run ./cmd/slka channels list
go run ./cmd/slka message send general "Hello"

# From dist/
./dist/slka channels list
./dist/slka --help
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
go run ./cmd/slka channels list --output-pretty

# Debug tests
go test -v -run TestSpecificTest ./internal/config/
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run linter: `make lint` (if configured)
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

Releases are automated with GoReleaser and GitHub Actions:

1. **Tag a release:**
   ```bash
   git tag v0.3.1
   git push origin v0.3.1
   ```

2. **GitHub Actions automatically:**
   - Builds for all platforms
   - Creates release on GitHub
   - Uploads binaries and archives
   - Generates changelog

3. **Manual release (if needed):**
   ```bash
   # Test locally first
   goreleaser build --snapshot --clean

   # Release (requires GITHUB_TOKEN)
   goreleaser release --clean
   ```

See [RELEASING.md](RELEASING.md) for details.

### Testing New Features

When adding new features (like DMs, reactions):

1. **Service layer first** - Implement in `internal/slack/`
2. **Write comprehensive tests** - Cover all edge cases
3. **Add read commands** - In `pkg/read/`
4. **Add write commands** - In `pkg/write/` with approval
5. **Update documentation** - All relevant .md files

Example workflow:
```bash
# 1. Create service with tests
vim internal/slack/newfeature.go
vim internal/slack/newfeature_test.go
go test ./internal/slack -v

# 2. Add commands
vim pkg/read/newfeature.go
vim pkg/write/newfeature.go

# 3. Build and test
go build ./cmd/slka
./slka newfeature --help
```

### Contributing

1. Write tests first (TDD)
2. Keep changes focused
3. Update documentation
4. Ensure all tests pass: `make test`
5. Build successfully: `make build-local`

### Project Testing Status

All core modules have comprehensive test suites:

- ✅ `internal/links/` - Link parsing/formatting
- ✅ `internal/output/` - Output formatting
- ✅ `internal/config/` - Configuration
- ✅ `internal/approval/` - Approval system
- ✅ `internal/slack/channels.go` - Channel operations
- ✅ `internal/slack/dms.go` - DM operations
- ✅ `internal/slack/reactions.go` - Reaction operations

### Resources

- [Slack API Documentation](https://api.slack.com/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Testify Testing Library](https://github.com/stretchr/testify)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Go Testing](https://golang.org/pkg/testing/)
