# slka - AI Agent Guide

Complete guide for AI agents (Claude, GPT, custom LLMs) using slka to interact with Slack.

## Overview

`slka` is designed specifically for AI agents with:
- **JSON output** for all commands (easy parsing)
- **Unified binary** with both read and write operations
- **Activity-sorted results** — all lists sorted by most recent activity first
- **Conservative defaults** — lists return 50 items, history returns 20 messages
- **Human approval mode** for sensitive operations
- **Token-efficient filtering** to reduce API calls and context usage
- **Structured error responses** with actionable suggestions
- **User and bot token support** for flexibility

## Quick Start for AI Agents

```python
import json, subprocess

def slka(cmd):
    result = subprocess.run(["slka"] + cmd.split(), capture_output=True, text=True)
    return json.loads(result.stdout)

# List channels (filtered for efficiency)
channels = slka("channels list --filter engineering")

# Send message (may require approval)
response = slka("message send general 'Build complete!'")

# Check if message was acknowledged
ack = slka("reaction check-acknowledged general 1234567890.123")
```

## Token Types

| Token Type | Prefix | Messages From | Best For |
|------------|--------|---------------|----------|
| **User Token** | `xoxp-` | User's account | Personal AI assistants |
| **Bot Token** | `xoxb-` | Bot user | Team automation |

**For AI agents:**
- User tokens: AI acts on behalf of a specific person
- Bot tokens: AI acts as separate entity

Setup: See [MANIFEST_SETUP.md](MANIFEST_SETUP.md) for easy 2-minute setup.

## Core Principles

### 1. Use Filtering for Token Efficiency

**Critical for AI agents:** Always filter when searching.

```python
# ❌ BAD: Returns 100+ channels (~10,000 tokens)
channels = slka("channels list")

# ✅ GOOD: Returns 2-3 channels (~300 tokens)  
channels = slka("channels list --filter backend")
```

**Token savings:** 30x reduction by using `--filter`.

### 2. Always Parse JSON

Every command returns structured JSON:

```json
{
  "ok": true,
  "data": {
    "channels": [...]
  }
}
```

### 3. Check the `ok` Field

Verify success before using data:

```python
response = slka("channels info general")
if response["ok"]:
    channel = response["data"]["channel"]
else:
    print(f"Error: {response['error_description']}")
```

### 4. Handle Approval Workflow

Write operations may require human approval:

```python
response = slka("message send general 'Deploy complete'")
if response["ok"]:
    print("Sent!")
elif response.get("requires_approval"):
    print("Waiting for human approval...")
else:
    print(f"Failed: {response['error_description']}")
```

## Command Reference

### Unread Tracking

#### List Unread Conversations
```bash
slka unread list
slka unread list --channels-only
slka unread list --dms-only
slka unread list --min-unread 5
slka unread list --order-by oldest
```

Returns:
```json
{
  "ok": true,
  "data": {
    "unread_conversations": [
      {
        "id": "C123",
        "name": "engineering",
        "type": "channel",
        "is_channel": true,
        "unread_count": 15,
        "last_read": "1706123450.000000"
      },
      {
        "id": "D456",
        "type": "im",
        "is_im": true,
        "unread_count": 3,
        "user_id": "U789",
        "user_name": "alice"
      }
    ],
    "total_count": 2
  }
}
```

**Ordering Options:**
- `--order-by count` (default): Most unread first, for urgent items
- `--order-by oldest`: Oldest unread first, for FIFO processing

### Channels

#### List Channels (Filtered - Recommended)
```bash
slka channels list --filter engineering
slka channels list --filter backend --type private
```

Returns:
```json
{
  "ok": true,
  "data": {
    "channels": [
      {
        "id": "C123",
        "name": "engineering",
        "is_private": false,
        "num_members": 25
      }
    ]
  }
}
```

#### Get Channel History
```bash
slka channels history general --limit 50
slka channels history general --since 1706123456
```

### Direct Messages

#### List DMs (Filtered - Recommended)
```bash
slka dm list --filter alice
```

Returns both 1-on-1 and group DMs with that user:
```json
{
  "ok": true,
  "data": {
    "conversations": [
      {
        "id": "D123",
        "type": "im",
        "user_ids": ["U456"],
        "user_names": ["alice"]
      },
      {
        "id": "G789",
        "type": "mpim",
        "user_ids": ["U456", "U789", "U012"],
        "user_names": ["alice", "bob", "charlie"]
      }
    ]
  }
}
```

#### Send DM (1-on-1 or Group)
```bash
# Single user
slka dm send alice "Quick question..."

# Group DM
slka dm send alice,bob,charlie "Team meeting at 3"
```

#### Get DM History
```bash
slka dm history alice
slka dm history alice,bob,charlie  # Group DM
```

### Messages

#### Send Message
```bash
slka message send general "Build passed!"
slka message send general "Deploy ready" --dry-run
```

#### Reply to Thread
```bash
slka message reply general 1706123456.789 "Looks good!"
```

#### Edit Message
```bash
slka message edit general 1706123456.789 "Updated text"
```

### Reactions

#### Check if Message Was Acknowledged
```bash
slka reaction check-acknowledged general 1706123456.789
```

Returns:
```json
{
  "ok": true,
  "data": {
    "acknowledgment": {
      "is_acknowledged": true,
      "reacted_users": ["U456", "U789"],
      "reaction_count": 2,
      "reply_count": 1,
      "has_replies": true,
      "has_reactions": true,
      "message_author": "U123"
    }
  }
}
```

**Acknowledgment = any reaction OR reply from someone other than message author.**

#### List Reactions
```bash
slka reaction list general 1706123456.789
```

#### Add/Remove Reactions
```bash
slka reaction add general 1706123456.789 thumbsup
slka reaction remove general 1706123456.789 eyes
```

### Users

#### Lookup User
```bash
slka users lookup alice@example.com
slka users lookup alice
```

#### List Users
```bash
slka users list
slka users list --limit 50
```

## AI Agent Workflows

### Workflow 0: Check What Needs Attention

```python
def check_unread_and_prioritize():
    # Get all unread conversations
    result = slka("unread list")
    if not result["ok"]:
        return

    unreads = result["data"]["unread_conversations"]

    # Prioritize by type and count
    urgent_channels = [c for c in unreads if c["is_channel"] and c["unread_count"] >= 10]
    urgent_dms = [c for c in unreads if c["is_im"] and c["unread_count"] >= 3]

    # Handle urgent items first
    for channel in urgent_channels:
        handle_channel(channel["id"], channel["name"])

    for dm in urgent_dms:
        handle_dm(dm["id"], dm["user_name"])
```

### Workflow 1: Monitor Channel and Respond

```python
def monitor_channel(channel, keywords):
    # Get recent messages (filtered by time)
    from datetime import datetime, timedelta
    since = int((datetime.now() - timedelta(hours=1)).timestamp())
    
    result = slka(f"channels history {channel} --since {since}")
    if not result["ok"]:
        return
    
    for msg in result["data"]["messages"]:
        if any(kw in msg["text"].lower() for kw in keywords):
            # Respond to keyword match
            slka(f"message reply {channel} {msg['ts']} 'I can help with that!'")
```

### Workflow 2: Send and Wait for Acknowledgment

```python
def send_and_wait_for_ack(channel, message, timeout=300):
    import time
    
    # Send message
    sent = slka(f"message send {channel} '{message}'")
    if not sent["ok"]:
        return False, "Failed to send"
    
    ts = sent["data"]["timestamp"]
    
    # Poll for acknowledgment
    start = time.time()
    while time.time() - start < timeout:
        ack = slka(f"reaction check-acknowledged {channel} {ts}")
        if ack["ok"] and ack["data"]["acknowledgment"]["is_acknowledged"]:
            return True, "Acknowledged"
        time.sleep(10)
    
    return False, "Timeout"

# Usage
success, msg = send_and_wait_for_ack("general", "Deploy ready for review")
print(msg)
```

### Workflow 3: Find Channel and Send Update

```python
def send_project_update(project_name, message):
    # Find channel (token efficient!)
    result = slka(f"channels list --filter {project_name}")
    
    if not result["ok"] or not result["data"]["channels"]:
        return False, "Channel not found"
    
    channel_id = result["data"]["channels"][0]["id"]
    
    # Send message
    response = slka(f"message send {channel_id} '{message}'")
    return response["ok"], response.get("error_description", "Success")
```

### Workflow 4: DM Multiple Users

```python
def notify_team(users, message):
    # users can be list: ["alice", "bob", "charlie"]
    user_list = ",".join(users)
    
    result = slka(f"dm send {user_list} '{message}'")
    return result["ok"]

# Usage
notify_team(["alice", "bob"], "PR ready for review")
```

### Workflow 5: Check All DMs for Mentions

```python
def check_dms_for_mentions(keyword):
    # Get all DMs
    result = slka("dm list")
    if not result["ok"]:
        return []
    
    mentions = []
    for conv in result["data"]["conversations"]:
        # Get history
        user_ids = ",".join(conv["user_ids"]) if len(conv["user_ids"]) > 1 else conv["user_ids"][0]
        history = slka(f"dm history {user_ids} --limit 50")
        
        if history["ok"]:
            for msg in history["data"]["messages"]:
                if keyword.lower() in msg["text"].lower():
                    mentions.append({
                        "conversation": conv["user_names"],
                        "message": msg["text"],
                        "timestamp": msg["ts"]
                    })
    
    return mentions
```

## Error Handling

### Error Response Format

```json
{
  "ok": false,
  "error": "channel_not_found",
  "error_description": "The specified channel does not exist",
  "suggestion": "Check the channel ID or ensure the bot is invited"
}
```

### Common Errors

| Error Code | Description | Solution |
|------------|-------------|----------|
| `auth_error` | Invalid token | Check token, regenerate if needed |
| `permission_error` | Missing scope | Add scope to app, reinstall |
| `channel_not_found` | Channel doesn't exist | Use `channels list` to find it |
| `user_not_found` | User doesn't exist | Use `users lookup` to verify |
| `approval_required` | Human approval needed | Wait for user to approve |
| `rate_limited` | Too many requests | Back off, retry later |

### Robust Error Handling

```python
def robust_slka(cmd, max_retries=3):
    import time
    
    for attempt in range(max_retries):
        result = slka(cmd)
        
        if result["ok"]:
            return result
        
        # Handle specific errors
        if result["error"] == "rate_limited":
            time.sleep(5 * (attempt + 1))  # Exponential backoff
            continue
        elif result["error"] == "approval_required":
            print("Waiting for approval...")
            time.sleep(2)
            continue
        else:
            # Unrecoverable error
            print(f"Error: {result['error_description']}")
            return result
    
    return {"ok": False, "error": "max_retries_exceeded"}
```

## Token Efficiency Best Practices

### Always Filter Lists

```python
# ❌ BAD: Wastes 10,000 tokens
all_channels = slka("channels list")
engineering = [c for c in all_channels["data"]["channels"] if "eng" in c["name"]]

# ✅ GOOD: Uses 300 tokens
engineering = slka("channels list --filter eng")["data"]["channels"]
```

### Cache Channel/User IDs

```python
class SlackCache:
    def __init__(self):
        self.channels = {}
        self.users = {}
    
    def get_channel_id(self, name):
        if name not in self.channels:
            result = slka(f"channels list --filter {name}")
            if result["ok"] and result["data"]["channels"]:
                self.channels[name] = result["data"]["channels"][0]["id"]
        return self.channels.get(name)
    
    def get_user_id(self, identifier):
        if identifier not in self.users:
            result = slka(f"users lookup {identifier}")
            if result["ok"]:
                self.users[identifier] = result["data"]["user"]["id"]
        return self.users.get(identifier)

cache = SlackCache()
channel_id = cache.get_channel_id("engineering")
```

### Use Limits Appropriately

```python
# Recent activity only
recent = slka("channels history general --limit 20")

# Specific time range
from datetime import datetime, timedelta
since = int((datetime.now() - timedelta(hours=6)).timestamp())
result = slka(f"channels history general --since {since} --limit 50")
```

## Security Best Practices

### Never Log Tokens

```python
# ❌ BAD
print(f"Using token: {token}")
logging.info(f"Token: {token}")

# ✅ GOOD
print("Token configured")
logging.info("Authentication successful")
```

### Use Dry Run for Testing

```python
# Test command without executing
preview = slka("message send general 'Test' --dry-run")
print(f"Would: {preview['description']}")

# Execute if preview looks good
if user_approves(preview):
    actual = slka("message send general 'Test'")
```

### Validate Input

```python
def safe_send_message(channel, text):
    # Validate channel exists
    info = slka(f"channels info {channel}")
    if not info["ok"]:
        return False, "Invalid channel"
    
    # Sanitize message (remove any token-like strings)
    if "xox" in text.lower():
        return False, "Message contains sensitive data"
    
    # Send
    result = slka(f"message send {channel} '{text}'")
    return result["ok"], result.get("error_description", "Success")
```

## Response Formats

### Success Response

```json
{
  "ok": true,
  "data": {
    "timestamp": "1706123456.789",
    "channel": "C123456",
    "text": "Message sent"
  }
}
```

### Error Response

```json
{
  "ok": false,
  "error": "channel_not_found",
  "error_description": "The specified channel does not exist",
  "suggestion": "Check the channel ID"
}
```

### Approval Required Response

```json
{
  "ok": false,
  "requires_approval": true,
  "action": "send_message",
  "description": "Send message to #general",
  "payload": {
    "channel": "C123",
    "text": "Hello"
  }
}
```

## Exit Codes

Use exit codes for scripting:

- `0` - Success
- `1` - General error
- `2` - Authentication error
- `3` - Permission error
- `4` - Not found
- `5` - Approval required
- `6` - Rate limited

```python
import subprocess

result = subprocess.run(["slka", "channels", "info", "general"])
if result.returncode == 0:
    print("Success")
elif result.returncode == 4:
    print("Channel not found")
```

## Advanced Integration

### Async/Concurrent Operations

```python
import concurrent.futures

def check_multiple_channels(channel_names):
    def get_history(channel):
        return slka(f"channels history {channel} --limit 10")
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        futures = {executor.submit(get_history, ch): ch for ch in channel_names}
        results = {}
        
        for future in concurrent.futures.as_completed(futures):
            channel = futures[future]
            results[channel] = future.result()
        
        return results

# Usage
results = check_multiple_channels(["general", "engineering", "random"])
```

### Webhook Integration

```python
from flask import Flask, request

app = Flask(__name__)

@app.route("/slack-notification", methods=["POST"])
def handle_notification():
    data = request.json
    
    # Send to Slack
    result = slka(f"message send general '{data['message']}'")
    
    return {"success": result["ok"]}

app.run(port=5000)
```

## Debugging

### Enable Pretty Output

```bash
slka channels list --filter eng --output-pretty
```

### Check Configuration

```bash
slka config show
```

### Test Tokens

```bash
# Test read operations
slka users list --limit 1

# Test write operations (safe)
slka message send general "Test" --dry-run
```

## Next Steps

- **[AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md)** - One-page cheat sheet
- **[QUICKSTART.md](QUICKSTART.md)** - 5-minute setup guide
- **[README.md](README.md)** - Full feature overview

## Support

Questions or issues:
- GitHub: https://github.com/ulfschnabel/slka/issues
- Review error messages and suggestions in JSON responses
