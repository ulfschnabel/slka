package commands

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	slackpkg "github.com/ulf/slka/internal/slack"
)

// AddChannelsCommands adds all channel-related commands to the root
func AddChannelsCommands(rootCmd *cobra.Command, getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) {
	channelsCmd := &cobra.Command{
		Use:   "channels",
		Short: "Channel operations",
		Long:  `Query and manage Slack channels`,
	}

	// Read commands
	channelsCmd.AddCommand(createChannelsListCmd(getClient, getPretty))
	channelsCmd.AddCommand(createChannelsInfoCmd(getClient, getPretty))
	channelsCmd.AddCommand(createChannelsHistoryCmd(getClient, getPretty))
	channelsCmd.AddCommand(createChannelsMembersCmd(getClient, getPretty))

	// Write commands
	channelsCmd.AddCommand(createChannelsCreateCmd(getClient, getApprover, getPretty, getDryRun))
	channelsCmd.AddCommand(createChannelsArchiveCmd(getClient, getApprover, getPretty, getDryRun))
	channelsCmd.AddCommand(createChannelsUnarchiveCmd(getClient, getApprover, getPretty, getDryRun))
	channelsCmd.AddCommand(createChannelsRenameCmd(getClient, getApprover, getPretty, getDryRun))
	channelsCmd.AddCommand(createChannelsSetTopicCmd(getClient, getApprover, getPretty, getDryRun))
	channelsCmd.AddCommand(createChannelsSetDescriptionCmd(getClient, getApprover, getPretty, getDryRun))

	rootCmd.AddCommand(channelsCmd)
}

// Approver interface to avoid circular dependency
type Approver interface {
	Require(action string, payload map[string]interface{}) error
}

func createChannelsListCmd(getClient func() slackpkg.Client, getPretty func() bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all channels",
		RunE: func(cmd *cobra.Command, args []string) error {
			includeArchived, _ := cmd.Flags().GetBool("include-archived")
			channelType, _ := cmd.Flags().GetString("type")
			limit, _ := cmd.Flags().GetInt("limit")

			svc := slackpkg.NewChannelService(getClient())
			channels, err := svc.List(slackpkg.ListChannelsOptions{
				IncludeArchived: includeArchived,
				Type:            channelType,
				Limit:           limit,
			})

			if err != nil {
				result := output.Error("channels_list_failed", err.Error(), "Check your token and permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channels": channels,
			})
			result.Print(getPretty())
			return nil
		},
	}

	cmd.Flags().Bool("include-archived", false, "Include archived channels")
	cmd.Flags().String("type", "all", "Filter by type: public, private, all")
	cmd.Flags().Int("limit", 0, "Maximum number of channels to return")

	return cmd
}

func createChannelsInfoCmd(getClient func() slackpkg.Client, getPretty func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "info <channel>",
		Short: "Get channel information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			channel, err := svc.GetInfo(channelID)
			if err != nil {
				result := output.Error("channel_info_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel": channel,
			})
			result.Print(getPretty())
			return nil
		},
	}
}

func createChannelsHistoryCmd(getClient func() slackpkg.Client, getPretty func() bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history <channel>",
		Short: "Get channel message history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]
			sinceStr, _ := cmd.Flags().GetString("since")
			untilStr, _ := cmd.Flags().GetString("until")
			limit, _ := cmd.Flags().GetInt("limit")
			includeThreads, _ := cmd.Flags().GetBool("include-threads")

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			opts := slackpkg.HistoryOptions{
				Limit:          limit,
				IncludeThreads: includeThreads,
			}

			if sinceStr != "" {
				ts, err := parseTimestamp(sinceStr)
				if err != nil {
					result := output.Error("invalid_timestamp", err.Error(), "Use Unix timestamp or ISO8601 format")
					result.Print(getPretty())
					return fmt.Errorf("exit code %d", result.ExitCode())
				}
				opts.Since = ts
			}

			if untilStr != "" {
				ts, err := parseTimestamp(untilStr)
				if err != nil {
					result := output.Error("invalid_timestamp", err.Error(), "Use Unix timestamp or ISO8601 format")
					result.Print(getPretty())
					return fmt.Errorf("exit code %d", result.ExitCode())
				}
				opts.Until = ts
			}

			messages, err := svc.GetHistory(channelID, opts)
			if err != nil {
				result := output.Error("history_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id": channelID,
				"messages":   messages,
			})
			result.Print(getPretty())
			return nil
		},
	}

	cmd.Flags().String("since", "", "Only messages after this timestamp")
	cmd.Flags().String("until", "", "Only messages before this timestamp")
	cmd.Flags().Int("limit", 100, "Maximum number of messages")
	cmd.Flags().Bool("include-threads", false, "Include thread replies inline")

	return cmd
}

func createChannelsMembersCmd(getClient func() slackpkg.Client, getPretty func() bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members <channel>",
		Short: "List channel members",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]
			limit, _ := cmd.Flags().GetInt("limit")

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			members, err := svc.GetMembers(channelID, limit)
			if err != nil {
				result := output.Error("members_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id": channelID,
				"members":    members,
			})
			result.Print(getPretty())
			return nil
		},
	}

	cmd.Flags().Int("limit", 0, "Maximum number of members to return")

	return cmd
}

func createChannelsCreateCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			private, _ := cmd.Flags().GetBool("private")
			description, _ := cmd.Flags().GetString("description")
			topic, _ := cmd.Flags().GetString("topic")

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

			if getDryRun() {
				result := output.DryRun("create_channel", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("create_channel", payload); err != nil {
				result := output.ApprovalRequired("create_channel", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			params := slack.CreateConversationParams{
				ChannelName: name,
				IsPrivate:   private,
			}

			channel, err := getClient().CreateConversation(params)
			if err != nil {
				result := output.Error("create_failed", err.Error(), "Check channel name and permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			if description != "" {
				_, _ = getClient().SetPurposeOfConversation(channel.ID, description)
			}
			if topic != "" {
				_, _ = getClient().SetTopicOfConversation(channel.ID, topic)
			}

			result := output.Success(map[string]interface{}{
				"channel": map[string]interface{}{
					"id":         channel.ID,
					"name":       channel.Name,
					"is_private": channel.IsPrivate,
				},
			})
			result.Print(getPretty())
			return nil
		},
	}

	cmd.Flags().Bool("private", false, "Create as private channel")
	cmd.Flags().String("description", "", "Set channel purpose/description")
	cmd.Flags().String("topic", "", "Set channel topic")

	return cmd
}

func createChannelsArchiveCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <channel>",
		Short: "Archive a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			payload := map[string]interface{}{
				"channel": channelID,
			}

			if getDryRun() {
				result := output.DryRun("archive_channel", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("archive_channel", payload); err != nil {
				result := output.ApprovalRequired("archive_channel", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			err = getClient().ArchiveConversation(channelID)
			if err != nil {
				result := output.Error("archive_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id": channelID,
				"archived":   true,
			})
			result.Print(getPretty())
			return nil
		},
	}
}

func createChannelsUnarchiveCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "unarchive <channel>",
		Short: "Unarchive a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			payload := map[string]interface{}{
				"channel": channelID,
			}

			if getDryRun() {
				result := output.DryRun("unarchive_channel", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("unarchive_channel", payload); err != nil {
				result := output.ApprovalRequired("unarchive_channel", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			err = getClient().UnArchiveConversation(channelID)
			if err != nil {
				result := output.Error("unarchive_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id": channelID,
				"archived":   false,
			})
			result.Print(getPretty())
			return nil
		},
	}
}

func createChannelsRenameCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "rename <channel> <new_name>",
		Short: "Rename a channel",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]
			newName := args[1]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			payload := map[string]interface{}{
				"channel":  channelID,
				"new_name": newName,
			}

			if getDryRun() {
				result := output.DryRun("rename_channel", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("rename_channel", payload); err != nil {
				result := output.ApprovalRequired("rename_channel", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			channel, err := getClient().RenameConversation(channelID, newName)
			if err != nil {
				result := output.Error("rename_failed", err.Error(), "Check your permissions and name availability")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel": map[string]interface{}{
					"id":   channel.ID,
					"name": channel.Name,
				},
			})
			result.Print(getPretty())
			return nil
		},
	}
}

func createChannelsSetTopicCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "set-topic <channel> <topic>",
		Short: "Set channel topic",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]
			topic := args[1]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			payload := map[string]interface{}{
				"channel": channelID,
				"topic":   topic,
			}

			if getDryRun() {
				result := output.DryRun("set_topic", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("set_topic", payload); err != nil {
				result := output.ApprovalRequired("set_topic", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			_, err = getClient().SetTopicOfConversation(channelID, topic)
			if err != nil {
				result := output.Error("set_topic_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id": channelID,
				"topic":      topic,
			})
			result.Print(getPretty())
			return nil
		},
	}
}

func createChannelsSetDescriptionCmd(getClient func() slackpkg.Client, getApprover func() Approver, getPretty func() bool, getDryRun func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "set-description <channel> <description>",
		Short: "Set channel description/purpose",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			channelArg := args[0]
			description := args[1]

			svc := slackpkg.NewChannelService(getClient())
			channelID, err := svc.ResolveChannel(channelArg)
			if err != nil {
				result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			payload := map[string]interface{}{
				"channel":     channelID,
				"description": description,
			}

			if getDryRun() {
				result := output.DryRun("set_description", payload)
				result.Print(getPretty())
				return nil
			}

			if err := getApprover().Require("set_description", payload); err != nil {
				result := output.ApprovalRequired("set_description", payload)
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			_, err = getClient().SetPurposeOfConversation(channelID, description)
			if err != nil {
				result := output.Error("set_description_failed", err.Error(), "Check your permissions")
				result.Print(getPretty())
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			result := output.Success(map[string]interface{}{
				"channel_id":  channelID,
				"description": description,
			})
			result.Print(getPretty())
			return nil
		},
	}
}
