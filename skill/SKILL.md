---
name: slka
description: "Slack CLI for agentic workflows with token-based authentication. Use this skill when interacting with Slack to: (1) Read channels, DMs, or messages, (2) Send messages to channels or users, (3) Monitor unread messages, (4) Look up users, (5) Manage reactions, (6) Send semi-personalized messages to multiple people. CRITICAL - This skill uses the slka binary only, never use official Slack CLI, slack-sdk, or any tool that accesses macOS keychain or OAuth flows."
---

# slka - Slack CLI for Agentic Workflows

## Overview

slka provides command-line access to Slack using simple token-based authentication. All commands output JSON for easy parsing.

**CRITICAL SECURITY REQUIREMENT**: Always use the `slka` binary. Never install or use:
- Official Slack CLI (`slack` command)
- Python `slack-sdk` or similar libraries
- Any tool that attempts OAuth flows or keychain access

## Workflow Decision Tree

```
┌─ User asks to interact with Slack
│
├─ Check if slka is configured
│  ├─ No config → Guide setup (see Setup section)
│  └─ Has config → Continue
│
├─ Determine operation type
│  ├─ READ (channels, DMs, users, reactions, unread)
│  │  └─ Execute command directly
│  │
│  └─ WRITE (send, create, edit, archive, reactions)
│     ├─ 1. Show what will be done
│     ├─ 2. Run with --dry-run flag first
│     ├─ 3. Ask user for approval
│     └─ 4. Execute if approved
│
└─ Parse JSON output and respond to user
```

## Setup

Check if slka is configured:

```bash
slka config show
```

If not configured, guide the user through setup:

1. **User needs Slack tokens** - Point them to the project documentation:
   - Quick setup: `/path/to/slka/MANIFEST_SETUP.md` (2 minutes)
   - User tokens: `/path/to/slka/USER_TOKEN_SETUP.md`
   - Bot tokens: `/path/to/slka/SLACK_SETUP.md`

2. **Run config init**:
   ```bash
   slka config init
   ```

   This prompts for:
   - Read token (xoxb-... or xoxp-...)
   - Write token (can be same as read token)
   - User token (optional)
   - Require approval (recommend: yes)

## Read Operations

Read operations execute immediately without approval. **Always use filters** to reduce token usage.

**Common patterns**:

```bash
# Monitor unread messages
slka unread list

# Find specific channels (token efficient!)
slka channels list --filter engineering

# Get recent messages
slka channels history general --limit 20

# Find DMs with a user
slka dm list --filter alice
slka dm history alice --limit 10

# Look up users
slka users lookup alice@company.com
```

For complete command reference, see [commands.md](references/commands.md).

## Write Operations

Write operations require explicit user approval. Always follow this workflow:

### 1. Test with --dry-run

```bash
slka message send general "Hello team!" --dry-run
```

Shows exactly what would happen without executing.

### 2. Show user what will happen

Present the dry-run output and explain the action.

### 3. Ask for approval

"I will send this message to #general. Proceed?"

### 4. Execute if approved

```bash
slka message send general "Hello team!"
```

If approval mode is enabled in config, slka will prompt for confirmation.

## Common Workflows

### Monitor and Respond to Unread Messages

```bash
# 1. Find unread channels/DMs
slka unread list --order-by oldest

# 2. Get messages from specific channel
slka channels history general --limit 10

# 3. Draft response and get approval
slka message send general "Response text" --dry-run

# 4. Send if approved
slka message send general "Response text"
```

### Send Semi-Personalized Messages

When sending similar messages to multiple people:

```bash
# 1. Look up users to get correct usernames
slka users list --filter relevant-filter

# 2. For each person, customize message
# 3. Test one with --dry-run
slka dm send alice "Hi Alice, [personalized content]" --dry-run

# 4. Ask user to approve approach
# 5. Send to each person after approval
slka dm send alice "Hi Alice, [personalized content]"
slka dm send bob "Hi Bob, [personalized content]"
slka dm send charlie "Hi Charlie, [personalized content]"
```

### Broadcast to Channel

```bash
# 1. Verify channel name
slka channels list --filter team-updates

# 2. Draft message and show to user
slka message send team-updates "Weekly update: ..." --dry-run

# 3. Get approval
# 4. Send
slka message send team-updates "Weekly update: ..."
```

## Token Efficiency

slka is designed for AI agents and token efficiency:

**❌ Inefficient**:
```bash
slka channels list  # Returns 100+ channels (~10k tokens)
slka dm list        # Returns all DMs
```

**✅ Efficient**:
```bash
slka channels list --filter engineering  # Returns 2-3 channels (~300 tokens)
slka dm list --filter alice             # Returns specific DMs
slka unread list                        # Only channels with unread
```

Always use filters when possible.

## JSON Output Parsing

All commands return JSON. Example:

```bash
slka channels info general
```

Returns:
```json
{
  "ok": true,
  "data": {
    "id": "C123456",
    "name": "general",
    "is_private": false,
    "num_members": 50
  }
}
```

Parse the JSON to extract needed information.

## Error Handling

Common errors:

**"Missing scope" error**:
- Token lacks required permissions
- User needs to recreate Slack app with correct scopes
- Point to MANIFEST_SETUP.md

**"Channel not found"**:
- Use channel name without #: `general` not `#general`
- Or use channel ID: `C123456`

**"No token configured"**:
- Run `slka config init` to set up tokens

## Safety Checklist

Before any operation:

- [ ] Using `slka` binary (not other Slack tools)
- [ ] For writes: Used `--dry-run` first
- [ ] For writes: Showed user what will happen
- [ ] For writes: Got explicit approval
- [ ] Used filters for efficiency
- [ ] Parsed JSON output correctly

## Reference

For complete command syntax and options, see [commands.md](references/commands.md).
