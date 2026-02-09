# slka Testing Framework

This directory contains the testing framework for slka, including a mock Slack API server and integration tests.

## Structure

```
test/
├── mockserver/        # Mock Slack API server
│   └── server.go      # HTTP server that simulates Slack API
├── fixtures/          # Test data
│   ├── channels.go    # Test channels with various states
│   ├── users.go       # Test users
│   └── messages.go    # Test messages
├── integration/       # Integration tests
│   ├── unread_test.go    # Tests for unread command
│   ├── channels_test.go  # Tests for channels commands
│   └── ...
├── helpers.go         # Test utilities and CLI helpers
└── README.md          # This file
```

## Running Tests

### Unit Tests (Fast, No Network)

```bash
# Run unit tests only
make test-unit

# Or directly with go
go test ./internal/... ./pkg/...
```

### Integration Tests (With Mock Server)

```bash
# Run integration tests
make test-integration

# Run specific integration test
go test ./test/integration -run TestUnreadList -v

# Run all tests
make test-all
```

## How It Works

### Mock Slack Server

The `mockserver` package provides an HTTP server that simulates Slack's API:

```go
// Start mock server
mockServer := mockserver.New()
defer mockServer.Close()

// Get mock server URL and token
url := mockServer.URL()
token := mockServer.Token
```

The mock server handles these endpoints:
- `conversations.list` - Returns test channels/DMs
- `conversations.info` - Returns detailed channel info with unread counts
- `conversations.history` - Returns test messages
- `users.list` - Returns test users
- `users.info` - Returns specific user info
- `chat.postMessage` - Simulates sending messages
- `auth.test` - Validates authentication

### Test Fixtures

Pre-defined test data includes:
- **Channels**: `general`, `engineering`, `random`, `secret-project` (with varying unread counts)
- **DMs**: 1-on-1 DMs with `alice` and `bob`
- **Group DMs**: Multi-party DMs
- **Users**: `alice`, `bob`, `charlie`, `slackbot`
- **Messages**: Sample messages in each channel

### Test Helpers

The `TestEnv` helper provides a complete test environment:

```go
func TestUnreadList(t *testing.T) {
    // Create test environment (mock server + temp config)
    env := test.NewTestEnv(t)
    defer env.Cleanup()

    // Run command and assert on results
    result := env.RunCommand("unread", "list", "--output-pretty")
    result.AssertSuccess().AssertJSONOK()

    // Assert on JSON fields
    result.AssertJSONField("data.channels[0].name", "engineering")

    // Get JSON arrays
    channels := result.GetJSONArray("data.channels")
    if len(channels) == 0 {
        t.Fatal("Expected channels")
    }
}
```

### Available Assertions

```go
// Command execution
result.AssertSuccess()              // Exit code 0
result.AssertFailure()              // Exit code != 0

// Output content
result.AssertContains("text")       // Output contains substring
result.AssertNotContains("text")    // Output doesn't contain substring

// JSON assertions
result.AssertJSONOK()               // {"ok": true}
result.AssertJSONField("path", val) // Specific field value
data := result.ParseJSON()          // Get full JSON object
arr := result.GetJSONArray("path")  // Get array field
val := result.GetJSONField("path")  // Get any field
```

## Writing New Tests

### 1. Create Test File

```bash
touch test/integration/myfeature_test.go
```

### 2. Write Test

```go
package integration

import (
    "testing"
    "github.com/ulf/slka/test"
)

func TestMyFeature(t *testing.T) {
    env := test.NewTestEnv(t)
    defer env.Cleanup()

    result := env.RunCommand("mycommand", "--flag", "value")
    result.AssertSuccess().AssertJSONOK()

    // Add assertions
}
```

### 3. Run Test

```bash
go test ./test/integration -run TestMyFeature -v
```

## Adding New Test Fixtures

### Add New Channel

Edit `test/fixtures/channels.go`:

```go
{
    ID:                 "C999",
    Name:               "my-test-channel",
    IsChannel:          true,
    UnreadCount:        42,
    UnreadCountDisplay: 42,
    NumMembers:         10,
}
```

### Add New User

Edit `test/fixtures/users.go`:

```go
{
    ID:       "U999",
    Name:     "dave",
    RealName: "Dave Wilson",
    Email:    "dave@example.com",
}
```

## Extending the Mock Server

To add support for new Slack API endpoints, edit `test/mockserver/server.go`:

```go
// In New()
mux.HandleFunc("/api/new.endpoint", m.handleNewEndpoint)

// Add handler
func (m *MockSlackServer) handleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    if !m.checkAuth(r) {
        m.writeError(w, "invalid_auth")
        return
    }

    // Handle endpoint logic
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "ok": true,
        // ... response data
    })
}
```

## Environment Variables

- `SLACK_API_URL` - Custom Slack API URL (automatically set by TestEnv to point to mock server)

## Tips

1. **Fast Iteration**: Use `-run` to run specific tests while developing
2. **Verbose Output**: Add `-v` flag to see detailed test output
3. **Test Coverage**: Run `make test-coverage` to see what code is tested
4. **Debugging**: Add `t.Logf("debug: %v", data)` to print debug info during tests

## Common Issues

### Tests Fail with "connection refused"

The mock server may not have started properly. Check that:
- The `NewTestEnv()` function is called at the start of each test
- The `defer env.Cleanup()` is present to clean up after tests

### Binary not found

Run `go build ./cmd/slka` first, or let the test framework build it automatically.

### Wrong API responses

Check that the mock server handlers in `mockserver/server.go` match the expected Slack API format.
