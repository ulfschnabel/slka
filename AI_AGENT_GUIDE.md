# slka - AI Agent Guide

This guide explains how AI agents (like Claude, GPT, or custom LLMs) should use the slka CLI tools to interact with Slack.

## Overview

`slka` is designed specifically for AI agents with:
- **JSON output** for all commands (easy parsing)
- **Read/write separation** for safety
- **Human approval mode** for sensitive operations
- **Structured error responses** with suggestions
- **Both user and bot token support** for flexibility

## Token Types

slka works with both token types:

| Token Type | Prefix | Messages Appear From | Best For |
|------------|--------|---------------------|----------|
| **User Token** | `xoxp-` | The user's account | Personal AI assistants |
| **Bot Token** | `xoxb-` | A bot user | Team automation |

**For AI agents:**
- Use **user tokens** if the AI acts on behalf of a specific person
- Use **bot tokens** if the AI acts as a separate entity

Setup guides:
- **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** - User tokens
- **[SLACK_SETUP.md](SLACK_SETUP.md)** - Bot tokens

## Core Principles for AI Agents

### 1. Always Parse JSON Output

Every command returns JSON. Parse it to extract data:

```json
{
  "ok": true,
  "data": {
    "channels": [...]
  }
}
```

### 2. Check the `ok` Field

Before using data, verify the operation succeeded:

```json
{
  "ok": false,
  "error": "channel_not_found",
  "error_description": "The specified channel does not exist",
  "suggestion": "Check the channel ID or ensure the bot is invited"
}
```

### 3. Use Read Operations Freely

`slka-read` commands are safe to use without approval:
- List channels, users, messages
- Get channel information
- Fetch message history

### 4. Request Approval for Write Operations

`slka-write` commands may require human approval if configured:
- Sending messages
- Creating/modifying channels
- Adding reactions

## Command Reference for AI Agents

### Reading Slack Data

#### Get Recent Messages from a Channel

```bash
slka-read channels history general --limit 50
```

Returns:
```json
{
  "ok": true,
  "data": {
    "channel_id": "C123456",
    "messages": [
      {
        "ts": "1706123456.789000",
        "user": "U123456",
        "user_name": "johndoe",
        "text": "Hello world",
        "reactions": [
          {"name": "thumbsup", "count": 2, "users": ["U111", "U222"]}
        ]
      }
    ]
  }
}
```

**AI Agent Pattern:**
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
else:
    print(f"Error: {data['error_description']}")
```

#### Get Messages Since a Specific Time

```bash
# Unix timestamp
slka-read channels history general --since 1706123456

# ISO8601 format
slka-read channels history general --since 2024-01-24T09:00:00
```

#### List All Channels

```bash
slka-read channels list
```

Returns:
```json
{
  "ok": true,
  "data": {
    "channels": [
      {
        "id": "C123456",
        "name": "general",
        "is_private": false,
        "is_archived": false,
        "topic": "Company-wide announcements",
        "member_count": 42
      }
    ]
  }
}
```

#### Find a User

```bash
# By email (auto-detected)
slka-read users lookup john@example.com

# By name
slka-read users lookup johndoe --by name
```

Returns:
```json
{
  "ok": true,
  "data": {
    "user": {
      "id": "U123456",
      "name": "johndoe",
      "real_name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### Writing to Slack

#### Send a Message

```bash
slka-write message send general "Deployment to production completed successfully ✓"
```

**With approval enabled**, the human will see:
```
Send message to general: "Deployment to production completed successfully ✓"

Payload:
{
  "channel": "C123456",
  "text": "Deployment to production completed successfully ✓"
}

Execute this action? [y/N]:
```

**If approval is denied**, returns:
```json
{
  "ok": false,
  "requires_approval": true,
  "action": "send_message",
  "description": "Send message to general",
  "payload": {...}
}
```

**AI Agent Pattern:**
```python
result = subprocess.run(
    ["slka-write", "message", "send", "general", "Hello from AI"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)
if data["ok"]:
    print(f"Message sent: {data['data']['ts']}")
elif data.get("requires_approval"):
    print("Waiting for human approval...")
    # Human needs to approve in their terminal
else:
    print(f"Error: {data['error_description']}")
```

#### Use Dry Run to Preview Actions

```bash
slka-write message send general "Test" --dry-run
```

Returns:
```json
{
  "ok": false,
  "dry_run": true,
  "action": "send_message",
  "description": "Send message to general: \"Test\"",
  "payload": {
    "channel": "C123456",
    "text": "Test"
  }
}
```

**AI Agent Use Case:** Show the human what you plan to do before requesting approval.

#### Format Links Properly

The tool automatically converts Markdown links to Slack format:

```bash
# Markdown format (converted automatically)
slka-write message send general "Check out [our docs](https://example.com)"

# Slack native format (passed through)
slka-write message send general "Check out <https://example.com|our docs>"
```

Both work correctly.

### Channel Management

#### Create a Channel

```bash
slka-write channels create project-alpha --description "Alpha project workspace"
```

#### Archive Inactive Channels

```bash
slka-write channels archive old-project
```

## AI Agent Workflows

### Workflow 1: Daily Standup Summary

**Goal:** Summarize yesterday's messages from key channels

```bash
#!/bin/bash
# Get timestamp for yesterday 9am
YESTERDAY=$(date -d 'yesterday 9am' +%s)

# Fetch messages from multiple channels
slka-read channels history general --since $YESTERDAY > general.json
slka-read channels history engineering --since $YESTERDAY > engineering.json
slka-read channels history product --since $YESTERDAY > product.json

# AI agent processes the JSON files and generates summary
# Then posts summary back to Slack
slka-write message send daily-standup "$(cat summary.txt)"
```

**AI Agent Implementation:**
```python
import json
import subprocess
from datetime import datetime, timedelta

# Get yesterday's timestamp
yesterday = int((datetime.now() - timedelta(days=1)).replace(hour=9).timestamp())

channels = ["general", "engineering", "product"]
all_messages = []

for channel in channels:
    result = subprocess.run(
        ["slka-read", "channels", "history", channel, "--since", str(yesterday)],
        capture_output=True,
        text=True
    )
    data = json.loads(result.stdout)
    if data["ok"]:
        all_messages.extend(data["data"]["messages"])

# Process messages with AI (your logic here)
summary = generate_summary(all_messages)

# Post summary
subprocess.run(
    ["slka-write", "message", "send", "daily-standup", summary],
    capture_output=True,
    text=True
)
```

### Workflow 2: Automated Response to Mentions

**Goal:** Monitor channel for mentions and respond appropriately

```bash
# Get recent messages
slka-read channels history support --limit 100 > messages.json

# AI agent checks for mentions of the bot
# Generates appropriate responses
# Posts replies to threads
```

**AI Agent Implementation:**
```python
def monitor_and_respond(channel, bot_user_id):
    result = subprocess.run(
        ["slka-read", "channels", "history", channel, "--limit", "100"],
        capture_output=True,
        text=True
    )

    data = json.loads(result.stdout)
    if not data["ok"]:
        return

    for msg in data["data"]["messages"]:
        # Check if bot was mentioned
        if bot_user_id in msg["text"]:
            # Generate response with AI
            response = generate_response(msg["text"])

            # Reply in thread
            subprocess.run([
                "slka-write", "message", "reply",
                channel,
                msg["ts"],
                response
            ])
```

### Workflow 3: Channel Cleanup

**Goal:** Find and archive inactive channels

```bash
# List all channels
slka-read channels list --include-archived > all_channels.json

# AI agent analyzes activity, identifies stale channels
# Presents list to human for approval
# Archives approved channels
```

**AI Agent Implementation:**
```python
import json
from datetime import datetime, timedelta

# Get all channels
result = subprocess.run(
    ["slka-read", "channels", "list"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)
channels = data["data"]["channels"]

stale_channels = []
threshold = datetime.now() - timedelta(days=90)

for channel in channels:
    # Get recent history
    history_result = subprocess.run(
        ["slka-read", "channels", "history", channel["id"], "--limit", "1"],
        capture_output=True,
        text=True
    )

    history_data = json.loads(history_result.stdout)
    messages = history_data["data"]["messages"]

    if not messages or is_older_than(messages[0]["ts"], threshold):
        stale_channels.append(channel)

# Present to human
print(f"Found {len(stale_channels)} stale channels:")
for ch in stale_channels:
    print(f"  - #{ch['name']}")

# Archive with approval
for ch in stale_channels:
    subprocess.run([
        "slka-write", "channels", "archive", ch["id"]
    ])
```

### Workflow 4: User Onboarding

**Goal:** Automatically invite new user to relevant channels

```bash
# Find user
slka-read users lookup newuser@example.com > user.json

# AI agent determines relevant channels based on role
# Sends welcome message
slka-write message send general "Welcome @newuser to the team!"
```

## Error Handling for AI Agents

### Parse Exit Codes

```python
result = subprocess.run(["slka-read", "channels", "info", "nonexistent"])

if result.returncode == 0:
    # Success
    pass
elif result.returncode == 2:
    # Authentication error - check tokens
    pass
elif result.returncode == 3:
    # Permission error - missing Slack scope
    pass
elif result.returncode == 4:
    # Not found - channel/user doesn't exist
    pass
elif result.returncode == 5:
    # Approval required but not given
    pass
elif result.returncode == 6:
    # Rate limited - back off and retry
    pass
```

### Handle Common Errors

#### Channel Not Found

```json
{
  "ok": false,
  "error": "channel_not_found",
  "error_description": "The specified channel does not exist or the bot does not have access",
  "suggestion": "Check the channel ID or ensure the bot is invited to the channel"
}
```

**AI Agent Action:**
- Verify channel name/ID
- Suggest inviting the bot: `/invite @slka` in the channel

#### Missing Permissions

```json
{
  "ok": false,
  "error": "missing_scope",
  "error_description": "The token is missing required scope: chat:write",
  "suggestion": "Add the chat:write scope to your Slack app and reinstall"
}
```

**AI Agent Action:**
- Inform human that Slack app needs additional permissions
- Provide link to Slack app settings

#### Rate Limited

```json
{
  "ok": false,
  "error": "rate_limited",
  "error_description": "Rate limited by Slack API. Retry after 60 seconds.",
  "retry_after": 60
}
```

**AI Agent Action:**
- Wait for `retry_after` seconds
- Retry the request
- Consider reducing request frequency

## Best Practices for AI Agents

### 1. Always Use `--output-pretty` for Debugging

When debugging, use pretty-printed JSON:

```bash
slka-read channels list --output-pretty
```

### 2. Cache Channel/User IDs

Don't look up the same channel ID repeatedly:

```python
# BAD: Looks up "general" every time
for i in range(10):
    subprocess.run(["slka-read", "channels", "history", "general"])

# GOOD: Look up once, cache the ID
result = subprocess.run(
    ["slka-read", "channels", "info", "general"],
    capture_output=True,
    text=True
)
general_id = json.loads(result.stdout)["data"]["channel"]["id"]

# Use the ID directly
for i in range(10):
    subprocess.run(["slka-read", "channels", "history", general_id])
```

### 3. Use Dry Run Before Executing

Show users what you plan to do:

```python
# First: dry run
result = subprocess.run(
    ["slka-write", "message", "send", "general", "Hello", "--dry-run"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)
print(f"Planning to: {data['description']}")
print(f"Payload: {json.dumps(data['payload'], indent=2)}")

# Then: execute if approved
subprocess.run(["slka-write", "message", "send", "general", "Hello"])
```

### 4. Handle Approval Gracefully

When approval is required:

```python
result = subprocess.run(
    ["slka-write", "message", "send", "general", "Hello"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)

if data.get("requires_approval"):
    print("⏳ Waiting for human approval in terminal...")
    print(f"Action: {data['description']}")
    # Human must approve in the terminal where slka-write is running
elif data["ok"]:
    print("✓ Message sent successfully")
else:
    print(f"✗ Error: {data['error_description']}")
```

### 5. Respect Rate Limits

Slack has rate limits. Space out requests:

```python
import time

channels = ["general", "random", "engineering", "product"]

for channel in channels:
    subprocess.run(["slka-read", "channels", "history", channel])
    time.sleep(1)  # Wait 1 second between requests
```

### 6. Provide Context in Messages

When posting on behalf of a user, make it clear:

```bash
slka-write message send general "🤖 AI Assistant: Based on recent activity, I recommend..."
```

### 7. Use Thread Replies for Context

Reply in threads instead of posting to main channel:

```bash
# Reply to a specific message
slka-write message reply general 1706123456.789000 "Here's the analysis you requested..."
```

## Integration Examples

### Python

```python
import json
import subprocess

class SlkaClient:
    def read_messages(self, channel, limit=100):
        result = subprocess.run(
            ["slka-read", "channels", "history", channel, "--limit", str(limit)],
            capture_output=True,
            text=True,
            check=False
        )

        data = json.loads(result.stdout)
        if not data["ok"]:
            raise Exception(f"Failed to read messages: {data.get('error_description')}")

        return data["data"]["messages"]

    def send_message(self, channel, text, dry_run=False):
        cmd = ["slka-write", "message", "send", channel, text]
        if dry_run:
            cmd.append("--dry-run")

        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=False
        )

        data = json.loads(result.stdout)
        return data

# Usage
client = SlkaClient()
messages = client.read_messages("general", limit=50)
response = client.send_message("general", "Hello from Python!")
```

### Node.js

```javascript
const { execSync } = require('child_process');

class SlkaClient {
    readMessages(channel, limit = 100) {
        const result = execSync(
            `slka-read channels history ${channel} --limit ${limit}`,
            { encoding: 'utf-8' }
        );

        const data = JSON.parse(result);
        if (!data.ok) {
            throw new Error(`Failed to read messages: ${data.error_description}`);
        }

        return data.data.messages;
    }

    sendMessage(channel, text, dryRun = false) {
        const cmd = `slka-write message send ${channel} "${text}"${dryRun ? ' --dry-run' : ''}`;
        const result = execSync(cmd, { encoding: 'utf-8' });
        return JSON.parse(result);
    }
}

// Usage
const client = new SlkaClient();
const messages = client.readMessages('general', 50);
const response = client.sendMessage('general', 'Hello from Node.js!');
```

### Shell Script

```bash
#!/bin/bash

# Read messages and extract text
slka-read channels history general --limit 10 | \
    jq -r '.data.messages[] | "\(.user_name): \(.text)"'

# Send message with error handling
send_message() {
    local channel=$1
    local message=$2

    result=$(slka-write message send "$channel" "$message")

    if echo "$result" | jq -e '.ok' > /dev/null; then
        echo "✓ Message sent"
        return 0
    else
        error=$(echo "$result" | jq -r '.error_description')
        echo "✗ Error: $error"
        return 1
    fi
}

send_message "general" "Hello from bash!"
```

## Security Considerations for AI Agents

### 1. Never Log Tokens

Don't include tokens in logs or error messages:

```python
# BAD
print(f"Using token: {token}")

# GOOD
print(f"Using token: {mask_token(token)}")
```

### 2. Validate User Input

If users provide channel names or message content, validate it:

```python
def sanitize_channel_name(name):
    # Remove special characters
    return ''.join(c for c in name if c.isalnum() or c in ['-', '_'])
```

### 3. Use Read-Only Operations When Possible

Prefer `slka-read` for queries. Only use `slka-write` when necessary.

### 4. Enable Approval Mode in Production

For production AI agents, ensure `require_approval: true` in config.

### 5. Audit Actions

Log all write operations for audit purposes:

```python
import logging

logging.info(f"AI Agent posting message to {channel}: {message}")
subprocess.run(["slka-write", "message", "send", channel, message])
```

## Troubleshooting for AI Agents

### Check Configuration

```bash
slka-write config show
```

### Test Connectivity

```bash
# Test read token
slka-read users list --limit 1

# Test write token (dry run)
slka-write message send general "Test" --dry-run
```

### Debug with Pretty Output

```bash
slka-read channels list --output-pretty | less
```

### Verify Bot Permissions

If you get "not in channel" errors:
1. Go to the Slack channel
2. Type: `/invite @slka`
3. Retry the operation

## Summary

**For AI Agents:**
- ✅ Use `slka-read` freely for queries
- ✅ Parse JSON output from all commands
- ✅ Check the `ok` field before using data
- ✅ Use `--dry-run` to preview write operations
- ✅ Handle approval requests gracefully
- ✅ Respect rate limits
- ✅ Provide clear context in messages
- ✅ Cache channel/user IDs
- ✅ Handle errors with suggestions from JSON
- ✅ Use thread replies for related messages

The tool is designed to make Slack interaction simple and safe for AI agents while maintaining human oversight for sensitive operations.
