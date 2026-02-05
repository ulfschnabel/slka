# slka - AI Agent Quick Reference

One-page cheat sheet for AI agents using slka.

## Core Concept

All commands return JSON. Parse it. Check `"ok": true/false`.

## Read Commands (Safe, No Approval)

```bash
# Get recent messages
slka-read channels history CHANNEL --limit N

# Get messages since timestamp
slka-read channels history CHANNEL --since TIMESTAMP

# List all channels
slka-read channels list

# Get channel info
slka-read channels info CHANNEL

# List users
slka-read users list

# Find user
slka-read users lookup EMAIL_OR_NAME
```

## Write Commands (May Require Approval)

```bash
# Send message
slka-write message send CHANNEL "text"

# Reply to thread
slka-write message reply CHANNEL THREAD_TS "text"

# Edit message
slka-write message edit CHANNEL TIMESTAMP "new text"

# Create channel
slka-write channels create NAME [--private] [--description "..."]

# Archive channel
slka-write channels archive CHANNEL

# Dry run (preview without executing)
slka-write message send CHANNEL "text" --dry-run
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

def slka_read(cmd):
    result = subprocess.run(
        ["slka-read"] + cmd.split(),
        capture_output=True, text=True
    )
    return json.loads(result.stdout)

def slka_write(cmd):
    result = subprocess.run(
        ["slka-write"] + cmd.split(),
        capture_output=True, text=True
    )
    return json.loads(result.stdout)

# Usage
data = slka_read("channels history general --limit 10")
if data["ok"]:
    for msg in data["data"]["messages"]:
        print(f"{msg['user_name']}: {msg['text']}")

response = slka_write("message send general 'Hello!'")
if response["ok"]:
    print("Message sent!")
elif response.get("requires_approval"):
    print("Waiting for human approval...")
```

## Common Patterns

### Get Yesterday's Messages
```python
from datetime import datetime, timedelta
yesterday = int((datetime.now() - timedelta(days=1)).timestamp())
data = slka_read(f"channels history general --since {yesterday}")
```

### Send with Error Handling
```python
response = slka_write("message send general 'Hello'")
if response["ok"]:
    print("✓ Sent")
elif response.get("requires_approval"):
    print("⏳ Awaiting approval")
else:
    print(f"✗ {response['error_description']}")
```

### Dry Run First
```python
# Preview
preview = slka_write("message send general 'Hello' --dry-run")
print(f"Would do: {preview['description']}")

# Execute
actual = slka_write("message send general 'Hello'")
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
- Parse JSON output
- Check `data["ok"]` before using data
- Use `slka-read` freely
- Use `--dry-run` to preview writes
- Cache channel/user IDs
- Space out requests (rate limits)
- Handle approval gracefully

❌ **DON'T:**
- Log tokens in plaintext
- Ignore error messages
- Spam the API
- Assume commands succeed
- Use channel names repeatedly (cache IDs)

## Debug Tips

```bash
# Pretty print output
slka-read channels list --output-pretty

# Show config (tokens masked)
slka-write config show

# Test read token
slka-read users list --limit 1

# Test write token (dry run)
slka-write message send general "Test" --dry-run
```

## Links in Messages

Both formats work (converted automatically):

```bash
# Markdown (converted to Slack format)
slka-write message send general "Check [the docs](https://example.com)"

# Slack native (passed through)
slka-write message send general "Check <https://example.com|the docs>"
```

## Full Documentation

See **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** for:
- Complete command reference
- Workflow examples
- Integration patterns
- Error handling
- Security best practices
