package write

import (
	"fmt"

	"github.com/spf13/cobra"
	slackpkg "github.com/ulf/slka/internal/slack"
	"github.com/ulf/slka/internal/output"
)

var reactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Reaction operations",
	Long:  `Add and remove reactions to/from Slack messages`,
}

var reactionAddCmd = &cobra.Command{
	Use:   "add <channel> <timestamp> <emoji>",
	Short: "Add a reaction to a message",
	Long: `Add a reaction emoji to a message.

Examples:
  slka reaction add general 1706123456.789000 thumbsup
  slka reaction add general 1706123456.789000 :thumbsup:
  slka reaction add C01234567 1706123456.789000 eyes`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		timestamp := args[1]
		emoji := args[2]

		// Resolve channel name to ID
		channelSvc := slackpkg.NewChannelService(slackClient)
		channelID, err := channelSvc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
			"emoji":     emoji,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("add_reaction", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("add_reaction", payload); err != nil {
			result := output.ApprovalRequired("add_reaction", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Execute via service
		svc := slackpkg.NewReactionService(slackClient)
		err = svc.AddReaction(channelID, timestamp, emoji)
		if err != nil {
			result := output.Error("add_reaction_failed", err.Error(), "Check your permissions and message timestamp")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
			"emoji":     emoji,
			"added":     true,
		})
		result.Print(outputPretty)
		return nil
	},
}

var reactionRemoveCmd = &cobra.Command{
	Use:   "remove <channel> <timestamp> <emoji>",
	Short: "Remove a reaction from a message",
	Long: `Remove a reaction emoji from a message.

Examples:
  slka reaction remove general 1706123456.789000 thumbsup
  slka reaction remove general 1706123456.789000 :eyes:
  slka reaction remove C01234567 1706123456.789000 tada`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelArg := args[0]
		timestamp := args[1]
		emoji := args[2]

		// Resolve channel name to ID
		channelSvc := slackpkg.NewChannelService(slackClient)
		channelID, err := channelSvc.ResolveChannel(channelArg)
		if err != nil {
			result := output.Error("channel_not_found", err.Error(), "Check the channel name or ID")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Prepare payload for approval
		payload := map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
			"emoji":     emoji,
		}

		// Check for dry run
		if dryRun {
			result := output.DryRun("remove_reaction", payload)
			result.Print(outputPretty)
			return nil
		}

		// Request approval
		if err := approver.Require("remove_reaction", payload); err != nil {
			result := output.ApprovalRequired("remove_reaction", payload)
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Execute via service
		svc := slackpkg.NewReactionService(slackClient)
		err = svc.RemoveReaction(channelID, timestamp, emoji)
		if err != nil {
			result := output.Error("remove_reaction_failed", err.Error(), "Check your permissions and message timestamp")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Return success
		result := output.Success(map[string]interface{}{
			"channel":   channelID,
			"timestamp": timestamp,
			"emoji":     emoji,
			"removed":   true,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add reaction commands
	RootCmd.AddCommand(reactionCmd)
	reactionCmd.AddCommand(reactionAddCmd)
	reactionCmd.AddCommand(reactionRemoveCmd)
}
