# slka Command Reference

All list commands return results **sorted by last activity** (most recent first). Default limits: 50 for lists, 20 for history. Use `--limit N` to override.

## Read Operations (No Approval Required)

### Unread Tracking
```bash
# Find what needs attention
slka unread list
slka unread list --channels-only
slka unread list --dms-only
slka unread list --min-unread 5
slka unread list --order-by oldest  # Process oldest first
```

### Channels
```bash
# List channels (with filtering - use filters to save tokens!)
slka channels list
slka channels list --filter engineering
slka channels list --type private

# Get channel info and history
slka channels info general
slka channels history general              # default: 20 messages
slka channels history general --limit 50
slka channels history general --since 2024-01-01
```

### Direct Messages
```bash
# List DM conversations
slka dm list
slka dm list --filter alice  # Find all DMs with alice

# Get DM history
slka dm history alice
slka dm history alice,bob,charlie  # Group DM
```

### Users
```bash
# List all users
slka users list

# Look up a user
slka users lookup alice@example.com
slka users lookup alice
```

### Reactions
```bash
# List reactions on a message
slka reaction list general 1234567890.123456

# Check if message was acknowledged
slka reaction check-acknowledged general 1234567890.123456
```

## Write Operations (Require Approval)

### Messages
```bash
# Send messages
slka message send general "Hello team!" --dry-run  # Test first
slka message send general "Hello team!"  # Actually send

slka message reply general 1234567890.123456 "Reply text"
slka message edit general 1234567890.123456 "Updated text"
```

### Direct Messages
```bash
# Send DMs
slka dm send alice "Hello!" --dry-run
slka dm send alice "Hello!"
slka dm send alice,bob,charlie "Team meeting at 3pm"

slka dm reply alice 1234567890.123456 "Got it!"
```

### Channels
```bash
# Manage channels
slka channels create new-project --dry-run
slka channels create new-project
slka channels archive old-project
slka channels invite engineering alice,bob
```

### Reactions
```bash
# Add/remove reactions
slka reaction add general 1234567890.123456 thumbsup --dry-run
slka reaction add general 1234567890.123456 thumbsup
slka reaction remove general 1234567890.123456 eyes
```

## JSON Output

All commands output JSON. Use `--output-pretty` for formatted output:

```bash
slka channels list --filter eng --output-pretty
```

Example response:
```json
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

## Token Efficiency Tips

```bash
# OK: Returns 50 most recently active channels (sorted by activity)
slka channels list

# ✅ Better: Returns only matching channels
slka channels list --filter backend

# ✅ Good: Find specific user's DMs
slka dm list --filter alice

# Override default limit if needed
slka channels list --limit 200
```
