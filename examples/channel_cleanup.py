#!/usr/bin/env python3
"""
Channel Cleanup Bot - AI Agent Example

Finds inactive channels and suggests archiving them.
Demonstrates:
- Listing channels
- Analyzing channel activity
- Batch operations with approval
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


def get_all_channels():
    """Get all channels (excluding archived)"""
    data = slka_read(["channels", "list"])

    if not data["ok"]:
        print(f"Error: {data.get('error_description')}")
        return []

    return data["data"]["channels"]


def get_last_message_time(channel_id):
    """Get the timestamp of the last message in a channel"""
    data = slka_read([
        "channels", "history", channel_id,
        "--limit", "1"
    ])

    if not data["ok"] or not data["data"]["messages"]:
        return None

    # Parse Slack timestamp (e.g., "1706123456.789000")
    ts = float(data["data"]["messages"][0]["ts"])
    return datetime.fromtimestamp(ts)


def is_stale(last_message_time, threshold_days):
    """Check if a channel is stale based on last message time"""
    if not last_message_time:
        return True  # No messages = stale

    threshold = datetime.now() - timedelta(days=threshold_days)
    return last_message_time < threshold


def analyze_channels(threshold_days=90):
    """Analyze all channels and find stale ones"""
    print(f"Analyzing channels (threshold: {threshold_days} days)...\n")

    channels = get_all_channels()
    print(f"Found {len(channels)} active channels")

    stale_channels = []

    for channel in channels:
        # Skip important channels
        if channel["name"] in ["general", "random", "announcements"]:
            continue

        print(f"  Checking #{channel['name']}...", end=" ")

        last_message = get_last_message_time(channel["id"])

        if is_stale(last_message, threshold_days):
            if last_message:
                days_ago = (datetime.now() - last_message).days
                print(f"STALE (last message {days_ago} days ago)")
            else:
                print("STALE (no messages)")

            stale_channels.append({
                "id": channel["id"],
                "name": channel["name"],
                "member_count": channel.get("member_count", 0),
                "last_message": last_message
            })
        else:
            print("active")

    return stale_channels


def present_findings(stale_channels):
    """Present findings to user"""
    if not stale_channels:
        print("\n✓ No stale channels found!")
        return

    print(f"\n📋 Found {len(stale_channels)} stale channels:\n")

    for ch in stale_channels:
        last_msg = ch["last_message"].strftime("%Y-%m-%d") if ch["last_message"] else "never"
        print(f"  • #{ch['name']}")
        print(f"    Members: {ch['member_count']}")
        print(f"    Last message: {last_msg}")
        print()


def archive_channels(channels):
    """Archive channels with approval"""
    print(f"\nArchiving {len(channels)} channels...")

    archived = 0
    failed = 0
    pending_approval = 0

    for ch in channels:
        result = slka_write([
            "channels", "archive", ch["id"]
        ])

        if result["ok"]:
            print(f"  ✓ Archived #{ch['name']}")
            archived += 1
        elif result.get("requires_approval"):
            print(f"  ⏳ #{ch['name']} waiting for approval")
            pending_approval += 1
        else:
            print(f"  ✗ Failed to archive #{ch['name']}: {result.get('error_description')}")
            failed += 1

    print(f"\nResults:")
    print(f"  Archived: {archived}")
    print(f"  Pending approval: {pending_approval}")
    print(f"  Failed: {failed}")


def main():
    # Configuration
    THRESHOLD_DAYS = 90
    DRY_RUN = True  # Set to False to actually archive

    print("=" * 50)
    print("Channel Cleanup Bot")
    print("=" * 50)
    print()

    # Analyze channels
    stale_channels = analyze_channels(THRESHOLD_DAYS)

    # Present findings
    present_findings(stale_channels)

    if not stale_channels:
        return

    # Ask for confirmation
    if DRY_RUN:
        print("DRY RUN MODE: Set DRY_RUN=False to actually archive channels")
        return

    print("\n⚠️  This will archive the channels listed above.")
    response = input("Continue? [y/N]: ").strip().lower()

    if response not in ["y", "yes"]:
        print("Cancelled.")
        return

    # Archive channels
    archive_channels(stale_channels)

    print("\n✓ Channel cleanup complete!")


if __name__ == "__main__":
    main()
