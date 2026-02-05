#!/usr/bin/env python3
"""
Mention Responder - AI Agent Example

Monitors channels for mentions of the bot and responds automatically.
Demonstrates:
- Monitoring channels
- Detecting mentions
- Replying in threads
- Handling approval
"""

import json
import subprocess
import sys


def slka_read(cmd):
    """Execute slka-read command and parse JSON output"""
    result = subprocess.run(
        ["slka-read"] + cmd,
        capture_output=True,
        text=True,
        check=False
    )
    return json.loads(result.stdout)


def slka_write(cmd):
    """Execute slka-write command and parse JSON output"""
    result = subprocess.run(
        ["slka-write"] + cmd,
        capture_output=True,
        text=True,
        check=False
    )
    return json.loads(result.stdout)


def get_recent_messages(channel, limit=100):
    """Get recent messages from a channel"""
    data = slka_read([
        "channels", "history", channel,
        "--limit", str(limit)
    ])

    if not data["ok"]:
        print(f"Error: {data.get('error_description')}", file=sys.stderr)
        return []

    return data["data"]["messages"]


def find_bot_mentions(messages, bot_keywords):
    """Find messages that mention the bot"""
    mentions = []
    for msg in messages:
        text = msg.get("text", "").lower()
        if any(keyword in text for keyword in bot_keywords):
            mentions.append(msg)
    return mentions


def generate_response(message_text):
    """
    Generate a response to the message.

    In a real AI agent, this would use an LLM to generate
    a contextual response. For this example, we use simple rules.
    """
    text_lower = message_text.lower()

    # Simple keyword-based responses
    if "help" in text_lower or "?" in text_lower:
        return "I'm here to help! I can assist with:\n• Daily summaries\n• Channel analysis\n• Q&A about the project\n\nWhat would you like to know?"

    elif "status" in text_lower or "update" in text_lower:
        return "I'm running normally and monitoring channels. Last check: just now."

    elif "thanks" in text_lower or "thank you" in text_lower:
        return "You're welcome! Let me know if you need anything else."

    else:
        return "I received your message. How can I help you?"

    # In a real implementation, you would:
    # 1. Use an LLM to understand context
    # 2. Search relevant documentation
    # 3. Generate a helpful, contextual response
    # 4. Include relevant links or resources


def reply_to_mention(channel, message, response_text):
    """Reply to a mention in a thread"""
    thread_ts = message.get("thread_ts", message["ts"])

    result = slka_write([
        "message", "reply",
        channel,
        thread_ts,
        f"🤖 {response_text}"
    ])

    return result


def main():
    # Configuration
    channels = ["general", "support", "engineering"]
    bot_keywords = ["@bot", "hey bot", "bot help"]

    # You could also find your bot's user ID and look for that
    # bot_user_id = "U123456"

    print("Monitoring channels for mentions...")
    print(f"Channels: {', '.join(channels)}")
    print(f"Keywords: {', '.join(bot_keywords)}")
    print()

    processed_messages = set()

    for channel in channels:
        print(f"Checking #{channel}...")

        # Get recent messages
        messages = get_recent_messages(channel, limit=50)

        # Find mentions
        mentions = find_bot_mentions(messages, bot_keywords)

        if not mentions:
            print(f"  No new mentions in #{channel}")
            continue

        print(f"  Found {len(mentions)} mention(s)")

        for msg in mentions:
            msg_id = f"{channel}:{msg['ts']}"

            # Skip if already processed (in a real bot, use persistent storage)
            if msg_id in processed_messages:
                continue

            user = msg.get("user_name", "someone")
            text = msg.get("text", "")[:50]  # First 50 chars
            print(f"    @{user}: {text}...")

            # Generate response
            response = generate_response(msg["text"])

            # Reply (with approval)
            result = reply_to_mention(channel, msg, response)

            if result["ok"]:
                print(f"      ✓ Replied")
                processed_messages.add(msg_id)
            elif result.get("requires_approval"):
                print(f"      ⏳ Waiting for approval to reply")
            else:
                print(f"      ✗ Failed: {result.get('error_description')}")

    print("\nDone!")


if __name__ == "__main__":
    main()
