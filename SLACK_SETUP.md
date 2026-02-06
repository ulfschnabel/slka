# Getting Bot Tokens for slka

Step-by-step guide to create a Slack **bot** and get bot tokens.

> **Want to control your own account instead?** See **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** for user tokens.

## Overview

This guide shows you how to set up a **bot user** for slka. Messages will appear from the bot (not from your personal account).

You need to create a Slack app in your workspace and configure it with the right permissions (scopes). You'll get a bot token (`xoxb-...`) that can be used for both read and write operations.

**Bot tokens vs User tokens:**
- **Bot tokens** (`xoxb-`): Messages appear from a bot user ← This guide
- **User tokens** (`xoxp-`): Messages appear from your account → [USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)

## Quick Setup with Manifest (Recommended)

**Easiest method:** Use our pre-configured manifest file.

See **[MANIFEST_SETUP.md](MANIFEST_SETUP.md)** for 2-minute setup with all scopes included.

## Manual Setup

### Step 1: Create a Slack App

1. **Go to Slack API Apps page**
   - Visit: https://api.slack.com/apps
   - Sign in with your Slack account

2. **Click "Create New App"**
   - Choose **"From scratch"**
   - App Name: `slka` (or any name you prefer)
   - Pick your workspace
   - Click **"Create App"**

### Step 2: Add Bot Token Scopes

You need to add permissions (scopes) to your app so it can read and write to Slack.

1. **Go to "OAuth & Permissions"** (in the left sidebar)

2. **Scroll down to "Scopes" → "Bot Token Scopes"**

3. **Click "Add an OAuth Scope"** and add these scopes:

#### Required Scopes

```
channels:read         - List public channels
channels:history      - Read messages in public channels
channels:manage       - Create, archive, rename channels
groups:read          - List private channels
groups:history       - Read messages in private channels
groups:write         - Manage private channels
im:read              - List DMs
im:history           - Read DM messages
im:write             - Send DMs
mpim:read            - List group DMs
mpim:history         - Read group DM messages
mpim:write           - Send group DMs
users:read           - List users, get user info
users:read.email     - Look up users by email
chat:write           - Send messages
chat:write.public    - Send to channels without joining
reactions:read       - Read reactions on messages
reactions:write      - Add/remove reactions
```

**Tip:** Add all scopes to ensure slka has full functionality.

### Step 3: Install App to Workspace

1. **Scroll to the top of "OAuth & Permissions" page**

2. **Click "Install to Workspace"**

3. **Review the permissions** and click **"Allow"**

4. **Copy the "Bot User OAuth Token"**
   - It starts with `xoxb-`
   - This is your token!
   - Example: `xoxb-1234567890123-1234567890123-abcdefghijklmnopqrstuvwx`

5. **Save this token securely**
   - You'll need it for slka configuration
   - Don't share it publicly (it's like a password)

### Step 4: Configure slka

#### Option 1: Interactive Setup (Recommended)

```bash
slka config init
```

This will prompt you for:
- Read token (paste your `xoxb-...` token)
- Write token (paste the same token)
- Enable approval mode? (recommended: yes)

#### Option 2: Manual Configuration

Create `~/.config/slka/config.json`:

```bash
mkdir -p ~/.config/slka
nano ~/.config/slka/config.json
```

Add:
```json
{
  "read_token": "xoxb-your-token-here",
  "write_token": "xoxb-your-token-here",
  "require_approval": true
}
```

Save and set permissions:
```bash
chmod 600 ~/.config/slka/config.json
```

#### Option 3: Environment Variables

For quick testing:
```bash
export SLKA_READ_TOKEN="xoxb-your-token-here"
export SLKA_WRITE_TOKEN="xoxb-your-token-here"
```

### Step 5: Test Your Setup

#### Test Token

```bash
slka users list --limit 1
```

Expected output:
```json
{
  "ok": true,
  "data": {
    "users": [...]
  }
}
```

#### Test Write (Dry Run)

```bash
slka message send general "Test" --dry-run
```

Expected output:
```json
{
  "ok": false,
  "dry_run": true,
  "action": "send_message",
  "description": "Send message to general: \"Test\"",
  "payload": {...}
}
```

### Step 6: Invite Bot to Channels

Before the bot can read/write to a channel, it needs to be added:

1. **Go to the Slack channel** you want the bot to access

2. **Type and send:**
   ```
   /invite @slka
   ```
   (Replace `slka` with your app name)

3. **The bot will join the channel** and can now read/write there

## Troubleshooting

### "Missing scope" Error

**Problem:**
```json
{
  "ok": false,
  "error": "missing_scope",
  "error_description": "The token is missing required scope: chat:write"
}
```

**Solution:**
1. Go to https://api.slack.com/apps
2. Select your app
3. Go to "OAuth & Permissions"
4. Add the missing scope
5. **Reinstall the app** (scroll up and click "Reinstall to Workspace")
6. Copy the new token
7. Update your slka config

### "Channel not found" Error

**Problem:**
```json
{
  "ok": false,
  "error": "channel_not_found"
}
```

**Solution:**
- Invite the bot to the channel: `/invite @slka`
- Or use the channel ID instead of name

### "Invalid auth" Error

**Problem:**
```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

**Solution:**
- Check that your token is correct
- Make sure you copied the full token (starts with `xoxb-`)
- Try reinstalling the app to get a fresh token

### Token Expired

Tokens can be revoked or expire. To get a new one:
1. Go to https://api.slack.com/apps
2. Select your app
3. Go to "OAuth & Permissions"
4. Click "Reinstall to Workspace"
5. Copy the new token

## Security Best Practices

### 1. Keep Tokens Secret

- ✅ Store in `~/.config/slka/config.json` with permissions `0600`
- ✅ Use environment variables for automation
- ❌ Don't commit tokens to git
- ❌ Don't share tokens in Slack or public places

### 2. Enable Approval Mode

Set `"require_approval": true` in config to require human confirmation for write operations.

### 3. Rotate Tokens Regularly

Reinstall your app periodically to get fresh tokens.

### 4. Limit Scope Access

Only add the scopes your bot actually needs. For read-only bots, omit write scopes.

## Quick Reference

| What | Where | Starts with |
|------|-------|-------------|
| Create app | https://api.slack.com/apps | - |
| Bot token | OAuth & Permissions → Bot User OAuth Token | `xoxb-` |
| Add scopes | OAuth & Permissions → Scopes | - |
| Reinstall app | OAuth & Permissions → Install to Workspace | - |

## All Available Commands

Once configured, slka can:

**Channels:**
- List, filter, get info, history, members
- Create, archive, rename, set topic/description

**Direct Messages:**
- List DMs (1-on-1 and groups)
- Send DMs, send group DMs
- Get DM history

**Messages:**
- Send, reply, edit
- Add/remove reactions
- Check if acknowledged

**Users:**
- List users
- Lookup by name, email, ID

**Configuration:**
- Show, set, init config

See **[AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md)** for complete command reference.

## Next Steps

Once you have your tokens configured:

1. **Read the quickstart:** [QUICKSTART.md](QUICKSTART.md)
2. **Try some commands:** [README.md](README.md#available-commands)
3. **Build an AI agent:** [AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)

## Support

If you get stuck:
- Check the troubleshooting section above
- Review Slack API docs: https://api.slack.com/docs
- File an issue: https://github.com/ulfschnabel/slka/issues
