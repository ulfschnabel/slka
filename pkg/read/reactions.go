package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var reactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Query reaction information",
	Long:  `List reactions and check acknowledgment status of messages`,
}

var reactionListCmd = &cobra.Command{
	Use:   "list <channel> <timestamp>",
	Short: "List all reactions on a message",
	Long: `List all reactions on a specific message, including who reacted.

Examples:
  slka reaction list general 1706123456.789000
  slka reaction list C01234567 1706123456.789000`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		timestamp := args[1]

		// Resolve channel
		channelSvc := slack.NewChannelService(slackClient)
		channelID, err := channelSvc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// List reactions
		svc := slack.NewReactionService(slackClient)
		reactions, err := svc.ListReactions(channelID, timestamp)
		if err != nil {
			result := output.Error("list_reactions_failed", err.Error(), "Check message exists and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"reactions": reactions,
		})
		result.Print(outputPretty)
		return nil
	},
}

var reactionCheckAckCmd = &cobra.Command{
	Use:   "check-acknowledged <channel> <timestamp>",
	Short: "Check if a message has been acknowledged",
	Long: `Check if anyone (other than the message author) has reacted to or replied to the message.

A message is considered acknowledged if:
- Someone other than the author added a reaction, OR
- Someone other than the author replied in the thread

Examples:
  slka reaction check-acknowledged general 1706123456.789000
  slka reaction check-acknowledged C01234567 1706123456.789000`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		timestamp := args[1]

		// Resolve channel
		channelSvc := slack.NewChannelService(slackClient)
		channelID, err := channelSvc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Check acknowledgment
		svc := slack.NewReactionService(slackClient)
		ackInfo, err := svc.CheckAcknowledgment(channelID, timestamp)
		if err != nil {
			result := output.Error("check_ack_failed", err.Error(), "Check message exists and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"acknowledgment": ackInfo,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add reaction commands
	RootCmd.AddCommand(reactionCmd)
	reactionCmd.AddCommand(reactionListCmd)
	reactionCmd.AddCommand(reactionCheckAckCmd)
}
