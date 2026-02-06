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
	Use:   "send <user> <text>",
	Short: "Send a direct message to a user",
	Long: `Send a direct message to a user.

User can be specified as:
- User ID (U123456)
- Email address (user@example.com)
- Username (alice)

Examples:
  slka dm send alice "Hello there!"
  slka dm send user@example.com "Quick question..."
  slka dm send U123456 "Hi!" --dry-run`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userArg := args[0]
		text := args[1]
		unfurlLinks, _ := cmd.Flags().GetBool("unfurl-links")
		unfurlMedia, _ := cmd.Flags().GetBool("unfurl-media")

		// Resolve user to ID
		svc := slackpkg.NewDMService(slackClient)
		userID, err := svc.ResolveUser(userArg)
		if err != nil {
			result := output.Error("user_not_found", err.Error(), "Check the user ID, email, or username")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"user":    userArg,
			"user_id": userID,
			"text":    text,
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

		// Send DM
		channelID, timestamp, err := svc.SendDM(userID, text, unfurlLinks, unfurlMedia)
		if err != nil {
			result := output.Error("send_dm_failed", err.Error(), "Check your permissions and user ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"user":      userArg,
			"user_id":   userID,
			"channel":   channelID,
			"timestamp": timestamp,
			"text":      text,
		})
		result.Print(outputPretty)
		return nil
	},
}

var dmReplyCmd = &cobra.Command{
	Use:   "reply <user> <timestamp> <text>",
	Short: "Reply to a message in a DM thread",
	Long: `Reply to a specific message in a direct message thread.

User can be specified as:
- User ID (U123456)
- Email address (user@example.com)
- Username (alice)

Examples:
  slka dm reply alice 1706123456.789000 "Good point!"
  slka dm reply user@example.com 1706123456.789000 "Thanks!"
  slka dm reply U123456 1706123456.789000 "Got it" --dry-run`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		userArg := args[0]
		threadTS := args[1]
		text := args[2]
		unfurlLinks, _ := cmd.Flags().GetBool("unfurl-links")
		unfurlMedia, _ := cmd.Flags().GetBool("unfurl-media")

		// Resolve user to ID
		svc := slackpkg.NewDMService(slackClient)
		userID, err := svc.ResolveUser(userArg)
		if err != nil {
			result := output.Error("user_not_found", err.Error(), "Check the user ID, email, or username")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"user":      userArg,
			"user_id":   userID,
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
		channelID, timestamp, err := svc.ReplyInDM(userID, threadTS, text, unfurlLinks, unfurlMedia)
		if err != nil {
			result := output.Error("reply_dm_failed", err.Error(), "Check your permissions and message timestamp")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"user":      userArg,
			"user_id":   userID,
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
