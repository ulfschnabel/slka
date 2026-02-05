# Quick Setup with Slack Manifests

The easiest way to set up slka - use our pre-configured manifest files that automatically set all the right scopes!

## Choose Your Setup

### Option A: User Token (Personal Automation) ⭐ Recommended

Messages appear as **you**. Perfect for personal automation and AI assistants.

**Use:** `slack-manifest-user-token.yaml`

### Option B: Bot Token (Team Automation)

Messages appear from a **bot**. Perfect for team tools and shared automation.

**Use:** `slack-manifest-bot-token.yaml`

---

## Setup Steps

### 1. Create App from Manifest

Go to: **https://api.slack.com/apps**

1. Click **"Create New App"**
2. Choose **"From an app manifest"**
3. Select your workspace
4. Choose **YAML** tab
5. Copy and paste the content from one of our manifest files:
   - For personal use: Copy content from `slack-manifest-user-token.yaml`
   - For bot use: Copy content from `slack-manifest-bot-token.yaml`
6. Click **"Next"**
7. Review the settings (all scopes are pre-configured!)
8. Click **"Create"**

### 2. Install the App

1. Click **"Install to Workspace"**
2. Review the permissions
3. Click **"Allow"**

### 3. Get Your Token

**For User Token (Personal):**
- Go to **"OAuth & Permissions"**
- Copy the **"User OAuth Token"** (starts with `xoxp-`)

**For Bot Token:**
- Go to **"OAuth & Permissions"**
- Copy the **"Bot User OAuth Token"** (starts with `xoxb-`)

### 4. Configure slka

```bash
# Download slka (see releases)
chmod +x slka-*

# Configure
./slka-linux-amd64 config init
```

When prompted, paste your token for both read and write (same token).

### 5. Test It!

```bash
# List channels
./slka-linux-amd64 channels list

# Get channel history
./slka-linux-amd64 channels history general --limit 10

# Send a test message (with dry-run)
./slka-linux-amd64 message send general "Hello from slka!" --dry-run
```

---

## What's Included in the Manifests

### User Token Manifest Includes:

**Read Scopes:**
- `channels:read` - List channels
- `channels:history` - Read messages in channels
- `groups:read` - List private channels
- `groups:history` - Read private messages
- `im:read` - List DMs
- `im:history` - Read DMs
- `mpim:read` - List group DMs
- `mpim:history` - Read group DMs
- `users:read` - List users
- `reactions:read` - Read reactions on messages

**Write Scopes:**
- `chat:write` - Send messages as you
- `channels:manage` - Create/archive channels
- `groups:write` - Manage private channels
- `reactions:write` - Add/remove reactions

### Bot Token Manifest Includes:

All the same scopes, plus:
- `users:read.email` - Look up users by email
- `chat:write.public` - Post to channels without joining
- `im:write` - Send DMs
- `mpim:write` - Send group DMs

---

## Customizing the Manifest

Want to add or remove scopes? Edit the manifest before creating the app:

```yaml
oauth_config:
  scopes:
    user:  # or bot:
      - channels:read
      - channels:history
      # Add or remove scopes here
```

### Common Additional Scopes

```yaml
# For file operations
- files:read
- files:write

# For emoji
- emoji:read

# For reminders
- reminders:write

# For bookmarks
- bookmarks:read
- bookmarks:write
```

---

## Updating Scopes Later

If you need to add scopes after creating the app:

1. Go to https://api.slack.com/apps
2. Select your app
3. Go to **"OAuth & Permissions"**
4. Scroll to **"Scopes"**
5. Add new scopes
6. **Reinstall the app** (important!)
7. Copy the new token
8. Update your slka config

---

## Troubleshooting

### "Missing scope" Error

If you get a missing scope error, the manifest might not have included it. Add it manually:

1. Go to your app's "OAuth & Permissions"
2. Add the missing scope
3. Reinstall the app
4. Get the new token

### "Not in channel" (Bot Token)

Bot tokens need to be invited to channels:

```
/invite @slka-bot
```

User tokens don't need this - you can access any channel you're already in.

### Token Not Working

Make sure you:
- Copied the full token (they're long!)
- Used the right token type (User OAuth Token vs Bot User OAuth Token)
- Reinstalled the app after changing scopes

---

## Comparison: Manual vs Manifest Setup

| Step | Manual Setup | Manifest Setup |
|------|-------------|----------------|
| Create app | ✓ | ✓ |
| Add scopes | Add 10-15 scopes one by one 😫 | Pre-configured! ✅ |
| Configure settings | Manual | Automated ✅ |
| Install app | ✓ | ✓ |
| Get token | ✓ | ✓ |
| **Total time** | ~10 minutes | ~2 minutes ⚡ |

---

## Multiple Tokens

Want both read and write tokens separated?

1. Create two apps using the same manifest
2. Name them "slka-read" and "slka-write"
3. Install both
4. Use different tokens in your config:

```json
{
  "read_token": "xoxp-read-token-here",
  "write_token": "xoxp-write-token-here",
  "require_approval": true
}
```

---

## Next Steps

- **[QUICKSTART.md](QUICKSTART.md)** - Get started guide
- **[README.md](README.md)** - Full documentation
- **[AI_AGENT_GUIDE.md](AI_AGENT_GUIDE.md)** - Build AI automations
- **[USER_TOKEN_SETUP.md](USER_TOKEN_SETUP.md)** - Detailed user token guide

---

## FAQ

**Q: Can I use the same manifest for multiple workspaces?**
A: Yes! Create a new app from the manifest in each workspace.

**Q: Will this work with Slack Enterprise Grid?**
A: Yes, but you may need org admin approval.

**Q: Can I share this manifest with my team?**
A: Yes! They can use it to create their own apps with the same scopes.

**Q: Do I need both user and bot manifests?**
A: No, choose one based on your use case. User token for personal, bot token for team.

**Q: Can I automate this even more?**
A: Not really - Slack requires manual app creation and installation for security. But manifests make it as fast as possible!

---

## Support

Questions about manifests:
- Slack manifest docs: https://api.slack.com/reference/manifests
- File an issue: https://github.com/ulfschnabel/slka/issues
