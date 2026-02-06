# Using slka with User Tokens

Guide for using slka to control **your own Slack account** (not a bot).

## User Token vs Bot Token

| Feature | User Token (xoxp-) | Bot Token (xoxb-) |
|---------|-------------------|-------------------|
| Messages appear from | Your account | Bot user |
| Channel access | Channels you're in | Channels bot is invited to |
| Token prefix | `xoxp-` | `xoxb-` |
| Use case | Personal automation | Team bots |
| Scope names | User scopes | Bot scopes |

**Use user tokens when:**
- You want messages to appear from your account
- You want to automate your personal workflow
- You're already in the channels you need to access

**Use bot tokens when:**
- You want a separate bot identity
- Multiple people will use the same automation
- You need the bot to be clearly identified

## Getting a User Token

### Step 1: Create a Slack App

1. Go to: **https://api.slack.com/apps**
2. Click **"Create New App"** → **"From an app manifest"**
3. Select your workspace
4. Copy the contents of `slack-manifest-user-token.yaml` (included in release)
5. Click **"Create"**

**Or create from scratch:**
1. Go to: **https://api.slack.com/apps**
2. Click **"Create New App"** → **"From scratch"**
3. App Name: `slka-personal` (or any name)
4. Select your workspace
5. Click **"Create App"**
6. Continue to Step 2 to add scopes manually

### Step 2: Add User Token Scopes

**Important:** Use **"User Token Scopes"** (NOT "Bot Token Scopes")

1. Click **"OAuth & Permissions"** in left sidebar
2. Scroll down to **"User Token Scopes"** section
3. Click **"Add an OAuth Scope"** and add:

#### Required Scopes:
```
channels:read       - List public channels
channels:history    - Read channel messages
channels:manage     - Create/archive channels
groups:read         - List private channels
groups:history      - Read private messages
groups:write        - Manage private channels
im:read            - List DMs
im:history         - Read DM messages
im:write           - Send DMs
mpim:read          - List group DMs
mpim:history       - Read group DM messages
mpim:write         - Send group DMs
users:read         - List users
users:read.email   - Look up users by email
chat:write         - Send messages
reactions:read     - Read reactions on messages
reactions:write    - Add/remove reactions
```

#### Scope Descriptions:
```
channels:read       - List channels
channels:history    - Read channel messages
channels:manage     - Create/archive channels
groups:read         - List private channels
groups:history      - Read private messages
groups:write        - Manage private channels
im:read            - List DMs (1-on-1)
im:history         - Read DM messages
im:write           - Send DMs
mpim:read          - List group DMs
mpim:history       - Read group DM messages
mpim:write         - Send group DMs
users:read         - List users
users:read.email   - Look up users by email
chat:write         - Send messages
reactions:read     - Read reactions
reactions:write    - Add/remove reactions
```

### Step 3: Install and Authorize

1. Scroll to top of "OAuth & Permissions" page
2. Click **"Install to Workspace"**
3. **Review carefully:** This will give the app access to act as YOU
4. Click **"Allow"**

### Step 4: Copy Your User Token

Look for: **"User OAuth Token"**
- It starts with `xoxp-` (NOT `xoxb-`)
- Example: `xoxp-1234567890-1234567890-1234567890-abcdefghijk`

**Copy the entire token.**

### Step 5: Configure slka

```bash
# Interactive setup
slka config init
```

When prompted:
- **Read token:** Paste your `xoxp-...` token
- **Write token:** Paste the SAME `xoxp-...` token
- **Require approval:** Choose yes or no (recommended: yes)

**Option 2: Manual config:**

Create `~/.config/slka/config.json`:
```json
{
  "read_token": "xoxp-your-token-here",
  "write_token": "xoxp-your-token-here",
  "require_approval": true
}
```

Set permissions:
```bash
chmod 600 ~/.config/slka/config.json
```

**Option 3: Environment variables:**
```bash
export SLKA_READ_TOKEN="xoxp-your-token-here"
export SLKA_WRITE_TOKEN="xoxp-your-token-here"
```

### Step 6: Verify Token Type

Check that slka recognizes your user token:

```bash
slka config show
```

Output should show:
```json
{
  "ok": true,
  "data": {
    "read_token": "xoxp-***",
    "read_token_type": "User Token",
    "write_token": "xoxp-***",
    "write_token_type": "User Token",
    ...
  }
}
```

### Step 7: Test It

```bash
# List channels (you're already in them)
slka channels list --filter general

# Read messages
slka channels history general --limit 10

# List your DMs
slka dm list

# Send a test message (dry run)
slka message send general "Test from my account" --dry-run

# Actually send (approve when prompted if approval is enabled)
slka message send general "Hello from slka!"
```

## Key Differences with User Tokens

### 1. No Need to Invite Bot

With user tokens, you can access any channel you're already in. No `/invite` needed.

### 2. Messages Appear from You

```bash
slka message send general "Hello!"
```

In Slack, this appears as:
```
You: Hello!
```

Not from a bot.

### 3. Rate Limits

User tokens have different (sometimes stricter) rate limits than bot tokens. Be respectful:

```python
import time

for channel in channels:
    result = slka(f"channels history {channel}")
    time.sleep(2)  # Wait 2 seconds between requests
```

### 4. Permissions

You can only do what your Slack account can do:
- ✅ Can send to channels you're in
- ✅ Can read channels you have access to
- ✅ Can send DMs to anyone
- ❌ Cannot access channels you're not in
- ❌ Cannot do things your role doesn't allow

## Security Considerations

### User Tokens are Powerful

A user token acts as YOU. Anyone with this token can:
- Read your messages
- Send messages as you
- Access channels you're in
- Send DMs on your behalf
- Do anything you can do in Slack

**Protect it like your password!**

### Best Practices

1. **Store securely**
   ```bash
   chmod 600 ~/.config/slka/config.json
   ```

2. **Enable approval mode**
   ```json
   {
     "require_approval": true
   }
   ```
   This requires you to approve each write action.

3. **Use environment variables for automation**
   ```bash
   # In a secure CI/CD environment
   export SLKA_WRITE_TOKEN="xoxp-..."
   ```

4. **Rotate tokens regularly**
   - Go to https://api.slack.com/apps
   - Select your app
   - "OAuth & Permissions" → "Reinstall to Workspace"
   - Update your config with the new token

5. **Revoke if compromised**
   - Go to https://api.slack.com/apps
   - Select your app
   - "OAuth & Permissions" → "Revoke Token"

6. **Don't commit to git**
   ```bash
   # Already in .gitignore:
   config.json
   .config/
   ```

## Common Use Cases

### Personal Daily Summary

```bash
#!/bin/bash
# Get messages from channels you follow
YESTERDAY=$(date -d 'yesterday 9am' +%s)

slka channels history general --since $YESTERDAY > daily.json
slka channels history team --since $YESTERDAY >> daily.json

# Process with AI and post summary
# (appears as you posting it)
```

### Auto-respond to Mentions

```python
# Monitor for @your-name mentions
result = slka("channels history support --limit 50")
messages = result["data"]["messages"]

for msg in messages:
    if "@yourname" in msg["text"]:
        # Respond as yourself
        slka(f"message reply support {msg['ts']} 'Thanks, looking into it!'")
```

### Personal Channel Management

```bash
# Archive old channels you created
slka channels archive old-project

# Set topics on your channels
slka channels set-topic my-project "Updated project status"
```

### Send Direct Messages

```bash
# Send 1-on-1 DM
slka dm send alice "Quick question about the project"

# Send group DM
slka dm send alice,bob,charlie "Team sync at 3pm"
```

## Troubleshooting

### "Missing scope" Error

**Problem:**
```json
{
  "error": "missing_scope",
  "error_description": "The token is missing required scope: chat:write"
}
```

**Solution:**
1. Go to https://api.slack.com/apps → your app
2. "OAuth & Permissions" → "User Token Scopes" (NOT Bot!)
3. Add the missing scope
4. **"Reinstall to Workspace"** (important!)
5. Copy the NEW user token
6. Update your slka config

### "Not in channel" Error

**Problem:**
```json
{
  "error": "not_in_channel"
}
```

**Solution:**
Join the channel first:
1. Open Slack
2. Join the channel
3. Try again

(Unlike bot tokens, user tokens can't post to channels you're not in)

### Token Shows as "Bot Token"

**Problem:**
```bash
slka config show
# Shows: "read_token_type": "Bot Token"
```

**Solution:**
You used a bot token (`xoxb-`) instead of a user token (`xoxp-`).

Go back to Step 2 and make sure you're adding **User Token Scopes**, not Bot Token Scopes.

### Messages Don't Appear from Me

**Problem:**
Messages appear from a bot user, not your account.

**Solution:**
You're using a bot token. Check:
```bash
slka config show
```

Should show `"User Token"`, not `"Bot Token"`.

## Comparison: Same Task, Different Token Types

### With User Token (xoxp-)

```bash
# Setup
export SLKA_WRITE_TOKEN="xoxp-123-456-789-abc"

# Send message
slka message send general "Meeting at 3pm"

# In Slack:
# You: Meeting at 3pm
```

### With Bot Token (xoxb-)

```bash
# Setup
export SLKA_WRITE_TOKEN="xoxb-123-456-abc"

# Send message
slka message send general "Meeting at 3pm"

# In Slack:
# slka [BOT]: Meeting at 3pm
```

## When to Use Each

| Scenario | Use User Token | Use Bot Token |
|----------|---------------|---------------|
| Personal automation | ✅ | ❌ |
| Appear as yourself | ✅ | ❌ |
| Team automation | ❌ | ✅ |
| Clear bot identity | ❌ | ✅ |
| Share with team | ❌ | ✅ |
| Access your channels | ✅ | Use either |
| Send personal DMs | ✅ | Use either |

## Next Steps

- **[MANIFEST_SETUP.md](MANIFEST_SETUP.md)** - Easy 2-minute setup with manifest
- **[QUICKSTART.md](QUICKSTART.md)** - General getting started guide
- **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** - Build AI automations

## Support

Questions about user tokens:
- Slack API docs: https://api.slack.com/authentication/token-types
- File an issue: https://github.com/ulfschnabel/slka/issues
