package write

import (
	"fmt"

	"github.com/spf13/cobra"
	slackpkg "github.com/ulf/slka/internal/slack"
	"github.com/ulf/slka/internal/output"
)

var dmCmd = &cobra.Command{
	Use:   "dm",
	Short: "Send direct messages",
	Long:  `Send direct messages to Slack users`,
}

var dmSendCmd = &cobra.Command{
	Use:   "send <users> <text>",
	Short: "Send a direct message",
	Long: `Send a direct message to one or more users.

Users can be specified as:
- Single user: alice, user@example.com, U123456
- Multiple users (group DM): alice,bob,charlie

Examples:
  slka dm send alice "Hello there!"
  slka dm send user@example.com "Quick question..."
  slka dm send alice,bob,charlie "Hello team!" --dry-run`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		usersArg := args[0]
		text := args[1]
		unfurlLinks, _ := cmd.Flags().GetBool("unfurl-links")
		unfurlMedia, _ := cmd.Flags().GetBool("unfurl-media")

		// Resolve users to IDs (handles single or comma-separated)
		svc := slackpkg.NewDMService(slackClient)
		userIDs, err := svc.ResolveUsers(usersArg)
		if err != nil {
			result := output.Error("user_not_found", err.Error(), "Check the user IDs, emails, or usernames")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"users":    usersArg,
			"user_ids": userIDs,
			"text":     text,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("send_dm", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("send_dm", payload); err != nil {
			result := output.ApprovalRequired("send_dm", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Send DM (works for single user or group)
		channelID, timestamp, err := svc.SendDM(userIDs, text, unfurlLinks, unfurlMedia)
		if err != nil {
			result := output.Error("send_dm_failed", err.Error(), "Check your permissions and user IDs")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"users":     usersArg,
			"user_ids":  userIDs,
			"channel":   channelID,
			"timestamp": timestamp,
			"text":      text,
		})
		result.Print(outputPretty)
		return nil
	},
}

var dmReplyCmd = &cobra.Command{
	Use:   "reply <users> <timestamp> <text>",
	Short: "Reply to a message in a DM thread",
	Long: `Reply to a specific message in a direct message thread.

Users can be specified as:
- Single user: alice, user@example.com, U123456
- Multiple users (group DM): alice,bob,charlie

Examples:
  slka dm reply alice 1706123456.789000 "Good point!"
  slka dm reply user@example.com 1706123456.789000 "Thanks!"
  slka dm reply alice,bob,charlie 1706123456.789000 "Got it" --dry-run`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		usersArg := args[0]
		threadTS := args[1]
		text := args[2]
		unfurlLinks, _ := cmd.Flags().GetBool("unfurl-links")
		unfurlMedia, _ := cmd.Flags().GetBool("unfurl-media")

		// Resolve users to IDs
		svc := slackpkg.NewDMService(slackClient)
		userIDs, err := svc.ResolveUsers(usersArg)
		if err != nil {
			result := output.Error("user_not_found", err.Error(), "Check the user IDs, emails, or usernames")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"users":     usersArg,
			"user_ids":  userIDs,
			"thread_ts": threadTS,
			"text":      text,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("reply_dm", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("reply_dm", payload); err != nil {
			result := output.ApprovalRequired("reply_dm", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Reply in DM
		channelID, timestamp, err := svc.ReplyInDM(userIDs, threadTS, text, unfurlLinks, unfurlMedia)
		if err != nil {
			result := output.Error("reply_dm_failed", err.Error(), "Check your permissions and message timestamp")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"users":     usersArg,
			"user_ids":  userIDs,
			"channel":   channelID,
			"thread_ts": threadTS,
			"timestamp": timestamp,
			"text":      text,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add dm commands
	RootCmd.AddCommand(dmCmd)
	dmCmd.AddCommand(dmSendCmd)
	dmCmd.AddCommand(dmReplyCmd)

	// Send flags
	dmSendCmd.Flags().Bool("unfurl-links", false, "Unfurl links in the message")
	dmSendCmd.Flags().Bool("unfurl-media", false, "Unfurl media in the message")

	// Reply flags
	dmReplyCmd.Flags().Bool("unfurl-links", false, "Unfurl links in the message")
	dmReplyCmd.Flags().Bool("unfurl-media", false, "Unfurl media in the message")
}
