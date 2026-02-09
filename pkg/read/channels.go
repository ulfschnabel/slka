package read

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Channel operations",
	Long:  `Query and inspect Slack channels`,
}

var channelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all channels",
	Long: `List all channels, optionally filtered by name.

Examples:
  slka channels list
  slka channels list --filter general
  slka channels list --filter eng --type public`,
	RunE: func(cmd *cobra.Command, args []string) error {
		includeArchived, _ := cmd.Flags().GetBool("include-archived")
		channelType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")
		filter, _ := cmd.Flags().GetString("filter")

		svc := slack.NewChannelService(slackClient)
		channels, err := svc.List(slack.ListChannelsOptions{
			IncludeArchived: includeArchived,
			Type:            channelType,
			Limit:           limit,
		})

		if err != nil {
			result := output.Error("channels_list_failed", err.Error(), "Check your token and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Filter by name if specified (case-insensitive substring match)
		if filter != "" {
			filteredChannels := make([]slack.ChannelInfo, 0)
			filterLower := strings.ToLower(filter)
			for _, ch := range channels {
				if strings.Contains(strings.ToLower(ch.Name), filterLower) {
					filteredChannels = append(filteredChannels, ch)
				}
			}
			channels = filteredChannels
		}

		result := output.Success(map[string]interface{}{
			"channels": channels,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsInfoCmd = &cobra.Command{
	Use:   "info <channel>",
	Short: "Get channel information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]

		svc := slack.NewChannelService(slackClient)

		// Resolve channel name to ID if needed
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		channel, err := svc.GetInfo(channelID)
		if err != nil {
			result := output.Error("channel_info_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return channel fields directly (not wrapped in "channel" object)
		result := output.Success(channel)
		result.Print(outputPretty)
		return nil
	},
}

var channelsHistoryCmd = &cobra.Command{
	Use:   "history <channel>",
	Short: "Get channel message history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		sinceStr, _ := cmd.Flags().GetString("since")
		untilStr, _ := cmd.Flags().GetString("until")
		limit, _ := cmd.Flags().GetInt("limit")
		includeThreads, _ := cmd.Flags().GetBool("include-threads")

		svc := slack.NewChannelService(slackClient)

		// Resolve channel
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		opts := slack.HistoryOptions{
			Limit:          limit,
			IncludeThreads: includeThreads,
		}

		// Parse timestamps
		if sinceStr != "" {
			ts, err := parseTimestamp(sinceStr)
			if err != nil {
				result := output.Error("invalid_timestamp", err.Error(), "Use Unix timestamp or ISO8601 format")
				result.Print(outputPretty)
				return fmt.Errorf("exit code %d", result.ExitCode())
			}
			opts.Since = ts
		}

		if untilStr != "" {
			ts, err := parseTimestamp(untilStr)
			if err != nil {
				result := output.Error("invalid_timestamp", err.Error(), "Use Unix timestamp or ISO8601 format")
				result.Print(outputPretty)
				return fmt.Errorf("exit code %d", result.ExitCode())
			}
			opts.Until = ts
		}

		messages, err := svc.GetHistory(channelID, opts)
		if err != nil {
			result := output.Error("history_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"messages":   messages,
		})
		result.Print(outputPretty)
		return nil
	},
}

var channelsMembersCmd = &cobra.Command{
	Use:   "members <channel>",
	Short: "List channel members",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		svc := slack.NewChannelService(slackClient)

		// Resolve channel
		channelID, err := svc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		members, err := svc.GetMembers(channelID, limit)
		if err != nil {
			result := output.Error("members_failed", err.Error(), "Check your permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"channel_id": channelID,
			"members":    members,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add channels commands
	RootCmd.AddCommand(channelsCmd)
	channelsCmd.AddCommand(channelsListCmd)
	channelsCmd.AddCommand(channelsInfoCmd)
	channelsCmd.AddCommand(channelsHistoryCmd)
	channelsCmd.AddCommand(channelsMembersCmd)

	// List flags
	channelsListCmd.Flags().Bool("include-archived", false, "Include archived channels in results")
	channelsListCmd.Flags().String("type", "all", "Filter by type: public, private, or all (default: all)")
	channelsListCmd.Flags().Int("limit", 0, "Maximum number of channels to return (0 = unlimited)")
	channelsListCmd.Flags().String("filter", "", "Filter channels by name substring (case-insensitive, e.g., 'eng' matches 'engineering')")

	// History flags
	channelsHistoryCmd.Flags().String("since", "", "Only messages after this timestamp (Unix timestamp or ISO8601: 1706123456 or 2024-01-25)")
	channelsHistoryCmd.Flags().String("until", "", "Only messages before this timestamp (Unix timestamp or ISO8601: 1706123456 or 2024-01-25)")
	channelsHistoryCmd.Flags().Int("limit", 100, "Maximum number of messages")
	channelsHistoryCmd.Flags().Bool("include-threads", false, "Include thread replies inline")

	// Members flags
	channelsMembersCmd.Flags().Int("limit", 0, "Maximum number of members to return")
}

// parseTimestamp parses a timestamp string (Unix timestamp or ISO8601)
func parseTimestamp(s string) (int64, error) {
	// Try parsing as Unix timestamp
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return ts, nil
	}

	// Try parsing as ISO8601
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("invalid timestamp format: %s", s)
}
