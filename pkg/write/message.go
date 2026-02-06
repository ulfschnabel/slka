package write

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/links"
	"github.com/ulf/slka/internal/output"
	slackpkg "github.com/ulf/slka/internal/slack"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Message operations",
	Long:  `Send and manage Slack messages`,
}

var messageSendCmd = &cobra.Command{
	Use:   "send <channel> <text>",
	Short: "Send a message to a channel",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		text := args[1]

		unfurlLinks, _ := cmd.Flags().GetBool("unfurl-links")
		unfurlMedia, _ := cmd.Flags().GetBool("unfurl-media")

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Format links for Slack
		formattedText := links.FormatLinksForSlack(text)

		// Prepare payload
		payload := map[string]interface{}{
			"channel":      channelID,
			"text":         formattedText,
			"unfurl_links": unfurlLinks,
			"unfurl_media": unfurlMedia,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("send_message", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("send_message", payload); err != nil {
			result := output.ApprovalRequired("send_message", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Send message
		msgOptions := []slack.MsgOption{
			slack.MsgOptionText(formattedText, false),
		}
		if !unfurlLinks {
			msgOptions = append(msgOptions, slack.MsgOptionDisableLinkUnfurl())
		}
		if !unfurlMedia {
			msgOptions = append(msgOptions, slack.MsgOptionDisableMediaUnfurl())
		}

		channel, ts, err := slackClient.PostMessage(channelID, msgOptions...)

		if err != nil {
			result := output.Error("send_failed", err.Error(), "Check your permissions and token")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel": channel,
			"ts":      ts,
			"message": map[string]interface{}{
				"text": formattedText,
			},
		})
		result.Print(outputPretty)
		return nil
	},
}

var messageReplyCmd = &cobra.Command{
	Use:   "reply <channel> <thread_ts> <text>",
	Short: "Reply to a thread",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		threadTS := args[1]
		text := args[2]

		broadcast, _ := cmd.Flags().GetBool("broadcast")

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Format links
		formattedText := links.FormatLinksForSlack(text)

		// Prepare payload
		payload := map[string]interface{}{
			"channel":   channelID,
			"thread_ts": threadTS,
			"text":      formattedText,
			"broadcast": broadcast,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("reply_message", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("reply_message", payload); err != nil {
			result := output.ApprovalRequired("reply_message", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Send reply
		options := []slack.MsgOption{
			slack.MsgOptionText(formattedText, false),
			slack.MsgOptionTS(threadTS),
		}

		if broadcast {
			options = append(options, slack.MsgOptionBroadcast())
		}

		channel, ts, err := slackClient.PostMessage(channelID, options...)

		if err != nil {
			result := output.Error("reply_failed", err.Error(), "Check your permissions and thread_ts")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel":   channel,
			"ts":        ts,
			"thread_ts": threadTS,
		})
		result.Print(outputPretty)
		return nil
	},
}

var messageEditCmd = &cobra.Command{
	Use:   "edit <channel> <timestamp> <text>",
	Short: "Edit an existing message",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		timestamp := args[1]
		text := args[2]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Format links
		formattedText := links.FormatLinksForSlack(text)

		// Prepare payload
		payload := map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
			"text":      formattedText,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("edit_message", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("edit_message", payload); err != nil {
			result := output.ApprovalRequired("edit_message", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Update message
		channel, ts, _, err := slackClient.UpdateMessage(
			channelID,
			timestamp,
			slack.MsgOptionText(formattedText, false),
		)

		if err != nil {
			result := output.Error("edit_failed", err.Error(), "Check your permissions and timestamp")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel": channel,
			"ts":      ts,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add message commands
	RootCmd.AddCommand(messageCmd)
	messageCmd.AddCommand(messageSendCmd)
	messageCmd.AddCommand(messageReplyCmd)
	messageCmd.AddCommand(messageEditCmd)

	// Send flags
	messageSendCmd.Flags().Bool("unfurl-links", true, "Enable link previews (shows website preview cards, default: true)")
	messageSendCmd.Flags().Bool("unfurl-media", true, "Enable media previews (shows embedded images/videos, default: true)")

	// Reply flags
	messageReplyCmd.Flags().Bool("broadcast", false, "Also post reply to main channel (not just in thread, default: false)")
}
