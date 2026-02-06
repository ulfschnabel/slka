# slka - AI Agent Quick Reference

One-page cheat sheet for AI agents using slka.

## Core Concept

All commands return JSON. Parse it. Check `"ok": true/false`.

## Read Commands (Safe, No Approval)

### Channels - **Use --filter for token efficiency!**
```bash
# List channels (filtered - RECOMMENDED)
slka channels list --filter engineering
slka channels list --filter backend --type private

# List all (expensive, avoid if possible)
slka channels list

# Get channel info
slka channels info general

# Get recent messages
slka channels history general --limit 50
slka channels history general --since TIMESTAMP
```

### Direct Messages
```bash
# List DMs (filtered - RECOMMENDED)
slka dm list --filter alice

# List all DMs (1-on-1 and groups)
slka dm list

# Get DM history
slka dm history alice
slka dm history alice,bob,charlie  # Group DM
```

### Reactions
```bash
# Check if message was acknowledged
slka reaction check-acknowledged general 1234567890.123456

# List reactions on message
slka reaction list general 1234567890.123456
```

### Users
```bash
# List users
slka users list

# Find user
slka users lookup alice@example.com
slka users lookup alice
```

## Write Commands (May Require Approval)

### Messages
```bash
# Send message
slka message send general "text"

# Reply to thread
slka message reply general THREAD_TS "text"

# Edit message
slka message edit general TIMESTAMP "new text"
```

### Direct Messages
```bash
# Send DM (1-on-1)
slka dm send alice "Hello!"

# Send to group DM
slka dm send alice,bob,charlie "Team meeting at 3"

# Reply in DM thread
slka dm reply alice THREAD_TS "Got it"
```

### Reactions
```bash
# Add reaction
slka reaction add general TIMESTAMP thumbsup

# Remove reaction
slka reaction remove general TIMESTAMP eyes
```

### Channels
```bash
# Create channel
slka channels create newchannel [--private] [--description "..."]

# Archive channel
slka channels archive oldchannel

# Invite users
slka channels invite general alice,bob
```

### Dry Run (Preview Without Executing)
```bash
slka message send general "text" --dry-run
slka dm send alice "hi" --dry-run
```

## JSON Response Format

### Success
```json
{
  "ok": true,
  "data": { ... }
}
```

### Error
```json
{
  "ok": false,
  "error": "error_code",
  "error_description": "Human readable error",
  "suggestion": "What to do about it"
}
```

### Approval Required
```json
{
  "ok": false,
  "requires_approval": true,
  "action": "send_message",
  "description": "Send message to #general",
  "payload": { ... }
}
```

## Python Template

```python
import json, subprocess

def slka(cmd):
    """Run slka command and return JSON response"""
    result = subprocess.run(
        ["slka"] + cmd.split(),
        capture_output=True, text=True
    )
    return json.loads(result.stdout)

# Read examples (no approval)
channels = slka("channels list --filter eng")
dms = slka("dm list --filter alice")
messages = slka("channels history general --limit 10")

# Write examples (may need approval)
response = slka("message send general 'Hello!'")
if response["ok"]:
    print("Message sent!")
elif response.get("requires_approval"):
    print("Waiting for human approval...")
```

## Token Efficiency (Critical for AI Agents!)

```python
# ❌ BAD: Returns 100+ channels, wastes tokens
data = slka("channels list")  # ~10,000 tokens

# ✅ GOOD: Returns 2-3 channels, efficient
data = slka("channels list --filter backend")  # ~300 tokens

# ❌ BAD: Returns all DMs
data = slka("dm list")  # ~5,000 tokens

# ✅ GOOD: Returns only DMs with alice
data = slka("dm list --filter alice")  # ~200 tokens
```

**Always use --filter when you know what you're looking for!**

## Common Patterns

### Find Channel and Send Message
```python
# Find channel
channels = slka("channels list --filter project-alpha")
if channels["ok"] and channels["data"]["channels"]:
    channel_id = channels["data"]["channels"][0]["id"]

    # Send message
    response = slka(f"message send {channel_id} 'Update: Build passed'")
```

### Check if Message Was Acknowledged
```python
# Send message
sent = slka("message send general 'Please review PR #123'")
if sent["ok"]:
    timestamp = sent["data"]["timestamp"]

    # Check for acknowledgment later
    ack = slka(f"reaction check-acknowledged general {timestamp}")
    if ack["ok"] and ack["data"]["acknowledgment"]["is_acknowledged"]:
        print("Someone acknowledged!")
```

### Send DM to User
```python
# Find user
user = slka("users lookup alice@example.com")
if user["ok"]:
    # Send DM
    dm = slka("dm send alice 'Quick question about the deploy'")
```

### Get Yesterday's Messages (Filtered)
```python
from datetime import datetime, timedelta

yesterday = int((datetime.now() - timedelta(days=1)).timestamp())
data = slka(f"channels history engineering --since {yesterday} --limit 100")

if data["ok"]:
    for msg in data["data"]["messages"]:
        print(f"{msg['user_name']}: {msg['text']}")
```

### Dry Run First Pattern
```python
# Preview what would happen
preview = slka("message send general 'Important!' --dry-run")
print(f"Would: {preview['description']}")

# If looks good, execute
if user_confirms():
    actual = slka("message send general 'Important!'")
```

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Auth error (bad token)
- `3` - Permission error (missing scope)
- `4` - Not found (channel/user)
- `5` - Approval required but not given
- `6` - Rate limited

## Best Practices

✅ **DO:**
- **Use --filter flags** (saves massive amounts of tokens!)
- Parse JSON output
- Check `data["ok"]` before using data
- Use `--dry-run` to preview writes
- Cache channel/user IDs
- Handle approval gracefully
- Space out requests (rate limits)

❌ **DON'T:**
- List all channels/DMs without filtering
- Log tokens in plaintext
- Ignore error messages
- Spam the API
- Assume commands succeed

## Debug Tips

```bash
# Pretty print output
slka channels list --output-pretty

# Show config (tokens masked)
slka config show

# Test token
slka users list --limit 1

# Test write (safe)
slka message send general "Test" --dry-run
```

## Message Links

Both formats work (converted automatically):

```bash
# Markdown (converted to Slack format)
slka message send general "Check [the docs](https://example.com)"

# Slack native (passed through)
slka message send general "Check <https://example.com|the docs>"
```

## Acknowledgment Tracking Example

```python
# Send message and track acknowledgment
def send_and_wait_for_ack(channel, message, timeout_seconds=300):
    # Send
    sent = slka(f"message send {channel} '{message}'")
    if not sent["ok"]:
        return False

    timestamp = sent["data"]["timestamp"]

    # Poll for acknowledgment
    import time
    start = time.time()
    while time.time() - start < timeout_seconds:
        ack = slka(f"reaction check-acknowledged {channel} {timestamp}")
        if ack["ok"] and ack["data"]["acknowledgment"]["is_acknowledged"]:
            return True
        time.sleep(10)  # Check every 10 seconds

    return False  # Timeout

# Usage
if send_and_wait_for_ack("general", "Deploy ready, please review"):
    print("Someone acknowledged!")
else:
    print("No acknowledgment within timeout")
```

## Full Documentation

See **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** for:
- Complete command reference
- Workflow examples
- Integration patterns
- Error handling
- Security best practices
