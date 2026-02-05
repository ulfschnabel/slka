# slka — Slack CLI for Agentic Workflows

Two separate CLI tools for interacting with Slack, with a clear separation between read and write operations for security purposes.

- **`slka-read`** — Read-only operations (low risk)
- **`slka-write`** — Write operations (high risk, includes human approval mode)

## Features

- **Read/write separation** — Two distinct tools with separate Slack app tokens/scopes
- **JSON output** — All commands output JSON for easy parsing by LLMs and scripts
- **Human approval mode** — `slka-write` can require explicit human confirmation before executing
- **Test-driven development** — Built with comprehensive test coverage
- **Link handling** — Properly handles Slack's `<url|text>` link format in both directions

## Installation

### Prerequisites

- Go 1.21 or later
- Slack workspace with appropriate bot tokens

### Build from source

```bash
# Clone the repository
git clone https://github.com/ulf/slka
cd slka

# Install dependencies
make deps

# Build for your platform
make build-local

# Install to GOPATH/bin
make install
```

### Build for all platforms

```bash
make build
```

Binaries will be in the `dist/` directory.

## Configuration

### Get Slack Tokens

**Choose your setup:**

- **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** - Control **your own account** (messages appear as you) ⭐ Personal use
- **[SLACK_SETUP.md](SLACK_SETUP.md)** - Set up a **bot account** (messages appear from bot) ⭐ Team automation

Both token types work with slka. User tokens (`xoxp-`) are for personal automation, bot tokens (`xoxb-`) are for team bots.

### Configure slka

All configuration is stored in `~/.config/slka/config.json`:

```json
{
  "read_token": "xoxb-...",
  "write_token": "xoxb-...",
  "user_token": "xoxp-...",
  "require_approval": true
}
```

**Initialize configuration interactively:**

```bash
slka-write config init
```

This will walk you through setting up your tokens and preferences.

### Environment variables

You can override config file values with environment variables:

- `SLKA_READ_TOKEN` — Override read token
- `SLKA_WRITE_TOKEN` — Override write token
- `SLKA_USER_TOKEN` — Override user token

Note: `require_approval` can only be set in the config file.

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes
- **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** - Complete guide for AI agents and LLMs
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Developer guide for contributing
- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Architecture and implementation details

## Usage

### Read operations

```bash
# List channels
slka-read channels list

# Get channel history
slka-read channels history general --since 2024-01-01

# List users
slka-read users list

# Look up a user
slka-read users lookup john@example.com
```

### Write operations

```bash
# Send a message (with approval if configured)
slka-write message send general "Hello world"

# Create a channel
slka-write channels create newchannel --description "New channel for project"

# Archive a channel
slka-write channels archive old-project
```

## For AI Agents

This tool is designed specifically for AI agents and LLMs. All commands output JSON for easy parsing:

```python
import json
import subprocess

result = subprocess.run(
    ["slka-read", "channels", "history", "general", "--limit", "50"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)
if data["ok"]:
    messages = data["data"]["messages"]
    # Process messages...
```

See **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** for complete integration examples, workflows, and best practices.

## Human Approval Mode

When `require_approval` is set to `true` in the config file, `slka-write` will:

1. Print a description of the action to be taken
2. Print the full JSON payload that would be sent to Slack
3. Prompt for confirmation: `Execute this action? [y/N]`
4. Only proceed if the user explicitly types `y` or `yes`

This is a safety measure to prevent accidental or malicious actions.

## Development

### Run tests

```bash
make test
```

### Generate coverage report

```bash
make test-coverage
```

### Run linter

```bash
make lint
```

## License

MIT
