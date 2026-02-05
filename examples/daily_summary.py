#!/usr/bin/env python3
"""
Daily Summary Bot - AI Agent Example

Fetches messages from yesterday across multiple channels and generates
a summary. Demonstrates:
- Reading messages with time filters
- Processing multiple channels
- Posting summary with approval
"""

import json
import subprocess
from datetime import datetime, timedelta


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


def get_messages_since(channel, timestamp):
    """Get messages from a channel since a specific timestamp"""
    data = slka_read([
        "channels", "history", channel,
        "--since", str(timestamp),
        "--limit", "1000"
    ])

    if not data["ok"]:
        print(f"Warning: Failed to get messages from {channel}: {data.get('error_description')}")
        return []

    return data["data"]["messages"]


def generate_summary(messages_by_channel):
    """
    Generate a summary from messages.

    In a real AI agent, this would use an LLM to generate
    a natural language summary. For this example, we just
    count messages.
    """
    summary = "📊 Daily Summary\n\n"

    total_messages = 0
    for channel, messages in messages_by_channel.items():
        count = len(messages)
        total_messages += count
        summary += f"• #{channel}: {count} messages\n"

        # Get unique users
        users = set(msg.get("user_name", "unknown") for msg in messages)
        summary += f"  Active users: {', '.join(sorted(users))}\n\n"

    summary += f"\n**Total: {total_messages} messages across {len(messages_by_channel)} channels**"

    # In a real implementation, you would:
    # 1. Extract key topics using NLP
    # 2. Identify important decisions
    # 3. Summarize action items
    # 4. Use an LLM to generate natural language summary

    return summary


def main():
    # Configuration
    channels = ["general", "engineering", "product"]
    summary_channel = "daily-summary"

    # Get timestamp for yesterday 9am
    yesterday = datetime.now() - timedelta(days=1)
    yesterday_9am = yesterday.replace(hour=9, minute=0, second=0, microsecond=0)
    timestamp = int(yesterday_9am.timestamp())

    print(f"Fetching messages since {yesterday_9am.strftime('%Y-%m-%d %H:%M:%S')}...")

    # Collect messages from all channels
    messages_by_channel = {}
    for channel in channels:
        print(f"  Fetching from #{channel}...")
        messages = get_messages_since(channel, timestamp)
        if messages:
            messages_by_channel[channel] = messages
            print(f"    Found {len(messages)} messages")

    if not messages_by_channel:
        print("No messages found in any channel.")
        return

    # Generate summary
    print("\nGenerating summary...")
    summary = generate_summary(messages_by_channel)
    print(summary)

    # Post summary (with approval)
    print(f"\nPosting to #{summary_channel}...")
    response = slka_write([
        "message", "send", summary_channel, summary
    ])

    if response["ok"]:
        print("✓ Summary posted successfully!")
    elif response.get("requires_approval"):
        print("⏳ Waiting for human approval to post summary...")
    else:
        print(f"✗ Failed to post: {response.get('error_description')}")


if __name__ == "__main__":
    main()
