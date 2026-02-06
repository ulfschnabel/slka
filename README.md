# slka — Slack CLI for Agentic Workflows

A unified CLI tool for interacting with Slack, designed for AI agents and automation workflows.

## Features

- **Unified binary** — Single `slka` command with both read and write operations
- **JSON output** — All commands output JSON for easy parsing by LLMs and scripts
- **Human approval mode** — Write operations can require explicit human confirmation
- **Token-efficient filtering** — Filter channels and DMs to reduce API calls and token usage
- **Unread tracking** — Find channels/DMs needing attention with flexible ordering
- **Reaction tracking** — Check if messages have been acknowledged by others
- **Direct messages** — Support for 1-on-1 and group DMs
- **Link handling** — Properly handles Slack's `<url|text>` link format

## Quick Start

### Installation

Download the latest release from [GitHub Releases](https://github.com/ulfschnabel/slka/releases):

```bash
# Linux AMD64
tar -xzf slka-v0.3.0-linux-amd64.tar.gz
chmod +x slka
sudo mv slka /usr/local/bin/

# macOS ARM64 (Apple Silicon)
tar -xzf slka-v0.3.0-darwin-arm64.tar.gz
chmod +x slka
sudo mv slka /usr/local/bin/
```

### Setup

1. **Get Slack token** — Use our manifest files for easy setup:
   - **[MANIFEST_SETUP.md](MANIFEST_SETUP.md)** - 2-minute setup with pre-configured scopes ⭐ Recommended
   - **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** - Detailed guide for user tokens (personal automation)
   - **[SLACK_SETUP.md](SLACK_SETUP.md)** - Bot token setup (team automation)

2. **Configure slka**:
   ```bash
   slka config init
   ```

## Available Commands

### Unread Tracking
```bash
# Find what needs attention (NEW!)
slka unread list
slka unread list --channels-only
slka unread list --dms-only
slka unread list --min-unread 5
slka unread list --order-by oldest  # Process oldest first
```

### Channels
```bash
# List channels (with optional filtering)
slka channels list
slka channels list --filter engineering
slka channels list --type private

# Get channel info and history
slka channels info general
slka channels history general --limit 50

# Manage channels (requires approval)
slka channels create new-project
slka channels archive old-project
```

### Direct Messages
```bash
# List DM conversations (1-on-1 and group)
slka dm list
slka dm list --filter alice

# Get DM history
slka dm history alice
slka dm history alice,bob,charlie  # Group DM

# Send DMs (requires approval)
slka dm send alice "Hello!"
slka dm send alice,bob,charlie "Team meeting at 3pm"
slka dm reply alice 1234567890.123456 "Got it!"
```

### Messages
```bash
# Send messages (requires approval)
slka message send general "Hello team!"
slka message reply general 1234567890.123456 "Reply text"
slka message edit general 1234567890.123456 "Updated text"
```

### Reactions
```bash
# List reactions on a message
slka reaction list general 1234567890.123456

# Check if message was acknowledged
slka reaction check-acknowledged general 1234567890.123456

# Add/remove reactions (requires approval)
slka reaction add general 1234567890.123456 thumbsup
slka reaction remove general 1234567890.123456 eyes
```

### Users
```bash
# List all users
slka users list

# Look up a user
slka users lookup alice@example.com
slka users lookup alice
```

## Configuration

Config is stored in `~/.config/slka/config.json`:

```json
{
  "read_token": "xoxp-...",
  "write_token": "xoxp-...",
  "require_approval": true
}
```

**Token types:**
- **User tokens** (`xoxp-`) - Messages appear as you (personal automation)
- **Bot tokens** (`xoxb-`) - Messages appear from bot (team automation)

Use the same token for both read and write, or separate them for added security.

## For AI Agents

All commands output JSON for easy parsing:

```python
import json
import subprocess

# List channels matching "eng"
result = subprocess.run(
    ["slka", "channels", "list", "--filter", "eng"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)
if data["ok"]:
    channels = data["data"]["channels"]
    # Process channels...
```

**Token-efficient filtering:**
```bash
# Bad: Returns all 100+ channels (~10k tokens)
slka channels list

# Good: Returns 2-3 channels (~300 tokens)
slka channels list --filter backend

# Find all DMs with a specific user
slka dm list --filter alice
```

See **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** for complete integration guide.

## Human Approval Mode

When `require_approval: true` in config, write operations require confirmation:

1. Shows action description and payload
2. Prompts: `Execute this action? [y/N]`
3. Only proceeds if user types `y` or `yes`

Test safely with `--dry-run`:
```bash
slka message send general "test" --dry-run
slka dm send alice "hello" --dry-run
```

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes
- **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** - Complete guide for AI agents
- **[AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md)** - Quick command reference
- **[MANIFEST_SETUP.md](MANIFEST_SETUP.md)** - Easy token setup with manifests
- **[RELEASING.md](RELEASING.md)** - How to create releases

## Development

```bash
# Build
go build ./cmd/slka

# Run tests
go test ./...

# Test with goreleaser
goreleaser build --snapshot --clean
```

See **[DEVELOPMENT.md](DEVELOPMENT.md)** for detailed development guide.

## License

MIT
