package write

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	slackpkg "github.com/ulf/slka/internal/slack"
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Channel management operations",
	Long:  `Create, archive, and manage Slack channels`,
}

var channelsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		private, _ := cmd.Flags().GetBool("private")
		description, _ := cmd.Flags().GetString("description")
		topic, _ := cmd.Flags().GetString("topic")

		// Prepare payload
		payload := map[string]interface{}{
			"name":    name,
			"private": private,
		}
		if description != "" {
			payload["description"] = description
		}
		if topic != "" {
			payload["topic"] = topic
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("create_channel", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("create_channel", payload); err != nil {
			result := output.ApprovalRequired("create_channel", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Create channel
		params := slack.CreateConversationParams{
			ChannelName: name,
			IsPrivate:   private,
		}

		channel, err := slackClient.CreateConversation(params)
		if err != nil {
			result := output.Error("create_failed", err.Error(), "Check channel name and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Set description/topic if provided
		if description != "" {
			_, _ = slackClient.SetPurposeOfConversation(channel.ID, description)
		}
		if topic != "" {
			_, _ = slackClient.SetTopicOfConversation(channel.ID, topic)
		}

		result := output.Success(map[string]interface{}{
			"channel": map[string]interface{}{
				"id":         channel.ID,
				"name":       channel.Name,
				"is_private": channel.IsPrivate,
			},
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsArchiveCmd = &cobra.Command{
	Use:   "archive <channel>",
	Short: "Archive a channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		payload := map[string]interface{}{
			"channel": channelID,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("archive_channel", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("archive_channel", payload); err != nil {
			result := output.ApprovalRequired("archive_channel", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Archive channel
		err = slackClient.ArchiveConversation(channelID)
		if err != nil {
			result := output.Error("archive_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"archived":   true,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <channel>",
	Short: "Unarchive a channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		payload := map[string]interface{}{
			"channel": channelID,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("unarchive_channel", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("unarchive_channel", payload); err != nil {
			result := output.ApprovalRequired("unarchive_channel", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Unarchive channel
		err = slackClient.UnArchiveConversation(channelID)
		if err != nil {
			result := output.Error("unarchive_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"archived":   false,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsRenameCmd = &cobra.Command{
	Use:   "rename <channel> <new_name>",
	Short: "Rename a channel",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		newName := args[1]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		payload := map[string]interface{}{
			"channel":  channelID,
			"new_name": newName,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("rename_channel", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("rename_channel", payload); err != nil {
			result := output.ApprovalRequired("rename_channel", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Rename channel
		channel, err := slackClient.RenameConversation(channelID, newName)
		if err != nil {
			result := output.Error("rename_failed", err.Error(), "Check your permissions and name availability")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel": map[string]interface{}{
				"id":   channel.ID,
				"name": channel.Name,
			},
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsSetTopicCmd = &cobra.Command{
	Use:   "set-topic <channel> <topic>",
	Short: "Set channel topic",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		topic := args[1]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		payload := map[string]interface{}{
			"channel": channelID,
			"topic":   topic,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("set_topic", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("set_topic", payload); err != nil {
			result := output.ApprovalRequired("set_topic", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Set topic
		_, err = slackClient.SetTopicOfConversation(channelID, topic)
		if err != nil {
			result := output.Error("set_topic_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"topic":      topic,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsSetDescriptionCmd = &cobra.Command{
	Use:   "set-description <channel> <description>",
	Short: "Set channel description/purpose",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		description := args[1]

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		payload := map[string]interface{}{
			"channel":     channelID,
			"description": description,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("set_description", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("set_description", payload); err != nil {
			result := output.ApprovalRequired("set_description", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Set purpose
		_, err = slackClient.SetPurposeOfConversation(channelID, description)
		if err != nil {
			result := output.Error("set_description_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id":  channelID,
			"description": description,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsMarkReadCmd = &cobra.Command{
	Use:   "mark-read <channel> [timestamp]",
	Short: "Mark a channel as read up to a specific message",
	Long: `Mark a channel as read up to a specific timestamp. If no timestamp is provided, marks all messages as read.

This is useful for bots that process messages and want to clear the unread indicator.

Examples:
  # Mark all messages as read in a channel
  slka channels mark-read general

  # Mark as read up to a specific message timestamp
  slka channels mark-read general 1234567890.000000`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		var timestamp string
		if len(args) > 1 {
			timestamp = args[1]
		}

		// Resolve channel
		svc := slackpkg.NewChannelService(slackClient)
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// If no timestamp provided, get the latest message timestamp
		if timestamp == "" {
			history, err := svc.GetHistory(channelID, slackpkg.HistoryOptions{Limit: 1})
			if err != nil {
				result := output.Error("get_history_failed", err.Error(), "Could not fetch latest message")
				result.Print(outputPretty)
				return fmt.Errorf("exit code %d", result.ExitCode())
			}
			if len(history) == 0 {
				// No messages, nothing to mark
				result := output.Success(map[string]interface{}{
					"channel_id": channelID,
					"marked":     false,
					"reason":     "no messages in channel",
				})
				result.Print(outputPretty)
				return nil
			}
			timestamp = history[0].Timestamp
		}

		payload := map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("mark_read", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval (usually auto-approved for mark read)
		if err := approver.Require("mark_read", payload); err != nil {
			result := output.ApprovalRequired("mark_read", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Mark as read
		err = svc.MarkAsRead(channelID, timestamp)
		if err != nil {
			result := output.Error("mark_read_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"timestamp":  timestamp,
			"marked":     true,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add channels commands
	RootCmd.AddCommand(channelsCmd)
	channelsCmd.AddCommand(channelsCreateCmd)
	channelsCmd.AddCommand(channelsArchiveCmd)
	channelsCmd.AddCommand(channelsUnarchiveCmd)
	channelsCmd.AddCommand(channelsRenameCmd)
	channelsCmd.AddCommand(channelsSetTopicCmd)
	channelsCmd.AddCommand(channelsSetDescriptionCmd)
	channelsCmd.AddCommand(channelsMarkReadCmd)

	// Create flags
	channelsCreateCmd.Flags().Bool("private", false, "Create as private channel (default: false, creates public channel)")
	channelsCreateCmd.Flags().String("description", "", "Set channel purpose/description text (optional)")
	channelsCreateCmd.Flags().String("topic", "", "Set channel topic text (optional)")
}
