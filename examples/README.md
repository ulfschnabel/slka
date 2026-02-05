# slka Examples - AI Agent Scripts

Practical examples showing how AI agents can use slka to automate Slack workflows.

## Examples

### 1. Daily Summary Bot (`daily_summary.py`)

Fetches messages from yesterday across multiple channels and generates a summary.

**Use case:** Morning briefings, team updates, activity reports

**What it demonstrates:**
- Reading messages with time filters
- Processing multiple channels in parallel
- Aggregating data from multiple sources
- Posting summary with approval

**Usage:**
```bash
cd examples
python3 daily_summary.py
```

**Customize:**
- Edit `channels` list to monitor different channels
- Edit `summary_channel` to change where summary is posted
- Replace `generate_summary()` with your LLM integration

---

### 2. Mention Responder (`mention_responder.py`)

Monitors channels for mentions of the bot and responds automatically.

**Use case:** Support bots, Q&A assistants, interactive agents

**What it demonstrates:**
- Monitoring channels for specific patterns
- Detecting mentions
- Generating contextual responses
- Replying in threads
- Tracking processed messages

**Usage:**
```bash
cd examples
python3 mention_responder.py
```

**Customize:**
- Edit `bot_keywords` to change what triggers responses
- Edit `channels` to monitor different channels
- Replace `generate_response()` with your LLM integration

---

### 3. Channel Cleanup Bot (`channel_cleanup.py`)

Finds inactive channels and suggests archiving them.

**Use case:** Workspace hygiene, reducing channel clutter

**What it demonstrates:**
- Listing all channels
- Analyzing channel activity over time
- Making batch recommendations
- Safe execution with dry run mode
- Requesting approval for destructive actions

**Usage:**
```bash
cd examples
python3 channel_cleanup.py
```

**Customize:**
- Change `THRESHOLD_DAYS` to adjust staleness definition
- Edit the skip list to protect important channels
- Set `DRY_RUN = False` to actually archive channels

---

## Setup

### 1. Install slka

```bash
cd ..
./setup.sh
make install
```

### 2. Configure Slack tokens

```bash
slka-write config init
```

### 3. Make examples executable

```bash
chmod +x examples/*.py
```

### 4. Run an example

```bash
./examples/daily_summary.py
```

## Creating Your Own AI Agent

### Basic Template

```python
#!/usr/bin/env python3
import json
import subprocess

def slka_read(cmd):
    result = subprocess.run(
        ["slka-read"] + cmd,
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

def slka_write(cmd):
    result = subprocess.run(
        ["slka-write"] + cmd,
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

def main():
    # Your agent logic here
    data = slka_read(["channels", "list"])

    if data["ok"]:
        channels = data["data"]["channels"]
        # Process channels...

if __name__ == "__main__":
    main()
```

### Integration with LLMs

```python
import openai  # or anthropic, etc.

def generate_ai_response(messages):
    """Use an LLM to generate intelligent responses"""
    # Format messages for LLM context
    context = "\n".join([
        f"{msg['user_name']}: {msg['text']}"
        for msg in messages
    ])

    # Call your LLM
    response = openai.ChatCompletion.create(
        model="gpt-4",
        messages=[
            {"role": "system", "content": "You are a helpful Slack bot."},
            {"role": "user", "content": f"Summarize:\n{context}"}
        ]
    )

    return response.choices[0].message.content

def main():
    # Get messages
    data = slka_read(["channels", "history", "general", "--limit", "50"])
    messages = data["data"]["messages"]

    # Generate AI summary
    summary = generate_ai_response(messages)

    # Post to Slack
    slka_write(["message", "send", "summary", summary])
```

## Best Practices

### 1. Handle Errors Gracefully

```python
data = slka_read(["channels", "info", "nonexistent"])
if not data["ok"]:
    print(f"Error: {data['error_description']}")
    print(f"Suggestion: {data.get('suggestion', '')}")
    return
```

### 2. Use Dry Run Mode

```python
# Preview first
preview = slka_write(["message", "send", "general", "Hello", "--dry-run"])
print(f"Would do: {preview['description']}")

# Then execute
result = slka_write(["message", "send", "general", "Hello"])
```

### 3. Cache Channel IDs

```python
# BAD: Looks up channel every time
for i in range(100):
    slka_read(["channels", "history", "general"])

# GOOD: Look up once, cache ID
channel_info = slka_read(["channels", "info", "general"])
channel_id = channel_info["data"]["channel"]["id"]

for i in range(100):
    slka_read(["channels", "history", channel_id])
```

### 4. Respect Rate Limits

```python
import time

for channel in channels:
    data = slka_read(["channels", "history", channel])
    # Process data...
    time.sleep(1)  # Wait between requests
```

### 5. Handle Approval Gracefully

```python
result = slka_write(["message", "send", "general", "Hello"])

if result["ok"]:
    print("✓ Sent")
elif result.get("requires_approval"):
    print("⏳ Approval needed (check terminal)")
    # Human must approve in the terminal where slka-write runs
else:
    print(f"✗ {result['error_description']}")
```

## Scheduling

### Using Cron

```bash
# Edit crontab
crontab -e

# Run daily at 9am
0 9 * * * /path/to/examples/daily_summary.py >> /var/log/slka-bot.log 2>&1

# Run every hour
0 * * * * /path/to/examples/mention_responder.py
```

### Using systemd Timer

Create `/etc/systemd/system/slka-daily-summary.service`:
```ini
[Unit]
Description=Slack Daily Summary Bot

[Service]
Type=oneshot
ExecStart=/path/to/examples/daily_summary.py
User=your-user
```

Create `/etc/systemd/system/slka-daily-summary.timer`:
```ini
[Unit]
Description=Run Slack Daily Summary Bot daily

[Timer]
OnCalendar=daily
OnCalendar=09:00

[Install]
WantedBy=timers.target
```

Enable:
```bash
sudo systemctl enable slka-daily-summary.timer
sudo systemctl start slka-daily-summary.timer
```

## Monitoring

### Logging

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/var/log/slka-bot.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)

def main():
    logger.info("Starting daily summary bot")
    # ... your code ...
    logger.info("Completed successfully")
```

### Health Checks

```python
def health_check():
    """Verify slka is working"""
    result = slka_read(["users", "list", "--limit", "1"])
    return result["ok"]

if not health_check():
    send_alert("slka is not responding!")
```

## Further Reading

- **[AI_AGENT_GUIDE.md](../AI_AGENT_GUIDE.md)** - Complete guide for AI agents
- **[AI_QUICK_REFERENCE.md](../AI_QUICK_REFERENCE.md)** - Quick reference cheat sheet
- **[DEVELOPMENT.md](../DEVELOPMENT.md)** - Contributing and extending

## Support

If you build something cool with these examples, share it!

Issues and questions: https://github.com/ulf/slka/issues
