# slka Quickstart Guide

Get up and running with slka in 5 minutes.

## Step 1: Install Go

```bash
sudo ./install-go.sh
```

## Step 2: Build slka

```bash
./setup.sh
```

This will:
- Download dependencies
- Run tests
- Build both `slka-read` and `slka-write` binaries

## Step 3: Get Slack Tokens

**Choose your setup:**

### Option A: User Token (Personal - Recommended for Individual Use)
Messages appear as **you**. Perfect for personal automation.

**📖 Follow: [USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)**

Quick summary:
1. Create Slack app
2. Add **User Token Scopes** (NOT Bot Token Scopes)
3. Install and authorize
4. Copy "User OAuth Token" (`xoxp-...`)

### Option B: Bot Token (Team Automation)
Messages appear from a **bot**. Perfect for team tools.

**📖 Follow: [SLACK_SETUP.md](SLACK_SETUP.md)**

Quick summary:
1. Create Slack app
2. Add **Bot Token Scopes**
3. Install to workspace
4. Copy "Bot User OAuth Token" (`xoxb-...`)

---

**Both work with slka!** User tokens are simpler for personal use.

## Step 4: Configure slka

```bash
./dist/slka-write config init
```

Follow the prompts to enter your tokens.

Or manually create `~/.config/slka/config.json`:

```json
{
  "read_token": "xoxb-your-read-token",
  "write_token": "xoxb-your-write-token",
  "require_approval": true
}
```

## Step 5: Try It Out

### Read Operations

```bash
# List channels
./dist/slka-read channels list

# Get channel history
./dist/slka-read channels history general --limit 10

# List users
./dist/slka-read users list

# Look up a user
./dist/slka-read users lookup john@example.com
```

### Write Operations (with Approval)

```bash
# Send a message (will prompt for approval)
./dist/slka-write message send general "Hello from slka!"

# Create a channel
./dist/slka-write channels create test-channel --description "Test"

# Dry run (see what would happen without executing)
./dist/slka-write message send general "Test" --dry-run
```

## Common Tasks

### Morning Brief

```bash
# Get updates from yesterday
YESTERDAY=$(date -d 'yesterday 9am' +%s)
./dist/slka-read channels history general --since $YESTERDAY
```

### Send Notification

```bash
# With approval
./dist/slka-write message send announcements "Deploy complete ✓"
```

### Find All Channels

```bash
# Including archived
./dist/slka-read channels list --include-archived --output-pretty
```

### Channel Management

```bash
# Archive old channel
./dist/slka-write channels archive old-project

# Set channel topic
./dist/slka-write channels set-topic general "Welcome to our workspace"
```

## Tips

### Disable Approval (for Automation)

Edit `~/.config/slka/config.json`:

```json
{
  "read_token": "...",
  "write_token": "...",
  "require_approval": false
}
```

⚠️ **Warning**: Only disable approval in trusted automation contexts!

### Environment Variables

Override tokens without changing config:

```bash
export SLKA_WRITE_TOKEN="xoxb-automation-token"
./dist/slka-write message send general "Automated message"
```

### Pretty Print JSON

Add `--output-pretty` to any command:

```bash
./dist/slka-read channels list --output-pretty
```

### Install System-Wide

```bash
make install
# Now you can use 'slka-read' and 'slka-write' from anywhere
```

## Troubleshooting

### "No read/write token configured"

Run `./dist/slka-write config init` or set environment variables.

### "Channel not found"

Make sure the bot is invited to the channel:
```
/invite @slka
```

### "Missing scope" Error

Add the required scope in your Slack app settings, then reinstall the app.

### Tests Failing

Tests may fail if Slack tokens aren't configured. This is expected for a fresh install.

## Next Steps

- Read the full [README.md](README.md) for all commands
- Check [DEVELOPMENT.md](DEVELOPMENT.md) for contributing
- Review the [plan specification](text.txt) for architecture details

## Getting Help

- Check command help: `./dist/slka-read --help`
- View subcommand help: `./dist/slka-read channels --help`
- Read Slack API docs: https://api.slack.com/
