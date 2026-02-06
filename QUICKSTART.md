# Quick Start Guide

Get up and running with slka in 5 minutes.

## 1. Installation

Download from [GitHub Releases](https://github.com/ulfschnabel/slka/releases):

```bash
# Linux (AMD64)
wget https://github.com/ulfschnabel/slka/releases/latest/download/slka-linux-amd64.tar.gz
tar -xzf slka-linux-amd64.tar.gz
chmod +x slka
sudo mv slka /usr/local/bin/

# macOS (Apple Silicon)
wget https://github.com/ulfschnabel/slka/releases/latest/download/slka-darwin-arm64.tar.gz
tar -xzf slka-darwin-arm64.tar.gz
chmod +x slka
sudo mv slka /usr/local/bin/
```

## 2. Get Slack Token

**Fastest method** - Use our manifest files (2 minutes):

1. Go to https://api.slack.com/apps
2. Click "Create New App" → "From an app manifest"
3. Select your workspace
4. Copy content from `slack-manifest-user-token.yaml` (included in release)
5. Click "Create" → "Install to Workspace"
6. Copy the **User OAuth Token** (starts with `xoxp-`)

See [MANIFEST_SETUP.md](MANIFEST_SETUP.md) for detailed instructions.

## 3. Configure

```bash
slka config init
```

When prompted, paste your token for both read and write (same token).

## 4. Test It!

```bash
# List channels
slka channels list

# List with filter (token efficient!)
slka channels list --filter general

# Get channel history
slka channels history general --limit 10

# List your DMs
slka dm list

# Test sending (dry-run, no approval needed)
slka message send general "Hello from slka!" --dry-run
```

## Common Commands

### Unread Tracking

```bash
# Find what needs attention
slka unread list
slka unread list --channels-only
slka unread list --order-by oldest
```

### Channels

```bash
# Find specific channels
slka channels list --filter engineering
slka channels list --type private

# Get info and history
slka channels info general
slka channels history general --limit 50 --since 2024-01-01

# Manage (requires approval)
slka channels create new-project
slka channels archive old-project
slka channels invite new-project alice,bob
```

### Direct Messages

```bash
# List DMs (1-on-1 and groups)
slka dm list
slka dm list --filter alice  # Find all DMs with alice

# View history
slka dm history alice
slka dm history alice,bob,charlie  # Group DM

# Send (requires approval)
slka dm send alice "Hey!"
slka dm send alice,bob "Team sync at 3"
```

### Messages

```bash
# Send (requires approval)
slka message send general "Hello team!"
slka message reply general 1234567890.123456 "Reply"
slka message edit general 1234567890.123456 "Updated"
```

### Reactions

```bash
# Check if acknowledged
slka reaction check-acknowledged general 1234567890.123456

# List reactions
slka reaction list general 1234567890.123456

# Add/remove (requires approval)
slka reaction add general 1234567890.123456 thumbsup
slka reaction remove general 1234567890.123456 eyes
```

### Users

```bash
# List all users
slka users list

# Lookup
slka users lookup alice@example.com
slka users lookup alice
```

## Configuration Options

Edit `~/.config/slka/config.json`:

```json
{
  "read_token": "xoxp-...",
  "write_token": "xoxp-...",
  "require_approval": true
}
```

**Options:**
- `read_token` - Token for read operations
- `write_token` - Token for write operations (can be same as read_token)
- `require_approval` - Require confirmation before write operations (default: true)

**Use separate tokens?**
- Same token: Simpler, one app to manage
- Separate tokens: More security (read-only app can't write)

## Approval Mode

When `require_approval: true`:

```bash
$ slka message send general "test"

Send message to general
Payload:
{
  "channel": "C123456",
  "text": "test"
}
Execute this action? [y/N]:
```

Type `y` or `yes` to proceed.

## Dry Run Mode

Test commands without executing:

```bash
slka message send general "test" --dry-run
slka dm send alice "hi" --dry-run
slka channels create test-channel --dry-run
```

Shows exactly what would happen, no approval needed.

## For AI Agents

All commands output JSON:

```bash
$ slka channels list --filter eng --output-pretty
{
  "ok": true,
  "data": {
    "channels": [
      {
        "id": "C123",
        "name": "engineering",
        "is_private": false,
        "num_members": 15
      }
    ]
  }
}
```

**Token efficiency tips:**
```bash
# ❌ Bad: Returns all 100 channels
slka channels list

# ✅ Good: Returns 2-3 matching channels
slka channels list --filter backend

# ✅ Good: Find specific user's DMs
slka dm list --filter alice
```

See [AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md) for integration examples.

## Next Steps

- **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** - Build AI automations
- **[AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md)** - Quick command reference
- **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** - Detailed token setup

## Troubleshooting

**"Missing scope" error?**
- Your token needs additional scopes
- Recreate app using manifest files (easiest)
- Or manually add scopes in app settings

**Commands hanging?**
- Check token is valid: `slka users list`
- Verify network access to Slack API

**"Channel not found"?**
- Use channel name without `#`: `slka channels info general`
- Or use channel ID: `slka channels info C123456`

**Approval not working?**
- Set `require_approval: true` in config
- Check you're using write operations (not read)
