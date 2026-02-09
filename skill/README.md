# slka Agent Skill

This directory contains the source files for the slka agent skill, which helps AI agents use slka safely and effectively.

## What's Included

- **SKILL.md** - Main skill instructions with workflows and safety guidelines
- **references/commands.md** - Complete command reference (loaded on-demand)
- **../slka.skill** - Packaged skill file (ready to install)

## Installing the Skill

To install the pre-packaged skill:

```bash
# For Claude Code CLI
claude code skills install slka.skill

# Or copy to your skills directory
cp slka.skill ~/.claude/skills/
```

## Modifying the Skill

If you need to update the skill:

1. Edit the source files in `skill/` directory
2. Rebuild the package:
   ```bash
   cd skill
   ./package.sh
   ```

## Key Features

The skill ensures agents:
- **Only use slka binary** (never official Slack CLI or keychain tools)
- **Test with --dry-run** before any write operations
- **Get user approval** for all write operations
- **Use filters** for token efficiency
- Follow established safety workflows

## Security

This skill includes critical security guardrails to prevent agents from:
- Installing or using the official Slack CLI (`slack` command)
- Using OAuth flows or accessing macOS keychain
- Installing Python libraries like `slack-sdk`
- Performing write operations without user approval

These safeguards prevent unauthorized access to system credentials and ensure all Slack operations go through the controlled slka interface.
