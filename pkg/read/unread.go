package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var unreadCmd = &cobra.Command{
	Use:   "unread",
	Short: "Query unread conversations",
	Long:  `Find channels and DMs that have unread messages requiring attention`,
}

var unreadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all unread conversations",
	Long: `List all channels and DMs with unread messages, ordered by urgency.

Use this to answer "what needs my attention?" - perfect for AI agents
monitoring Slack for items requiring action.

Examples:
  # List all unread (ordered by most unread first)
  slka unread list

  # List only unread channels
  slka unread list --channels-only

  # List only unread DMs
  slka unread list --dms-only

  # Show only items with 5+ unread messages
  slka unread list --min-unread 5

  # Order by oldest unread first (process old items first)
  slka unread list --order-by oldest

  # Order by highest unread count first (most urgent first, default)
  slka unread list --order-by count`,
	RunE: func(cmd *cobra.Command, args []string) error {
		channelsOnly, _ := cmd.Flags().GetBool("channels-only")
		dmsOnly, _ := cmd.Flags().GetBool("dms-only")
		minUnread, _ := cmd.Flags().GetInt("min-unread")
		orderBy, _ := cmd.Flags().GetString("order-by")

		// Validate order-by
		if orderBy != "" && orderBy != "count" && orderBy != "oldest" {
			result := output.Error("invalid_order_by", "order-by must be 'count' or 'oldest'", "Use --order-by count or --order-by oldest")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		svc := slack.NewUnreadService(slackClient)
		unreads, err := svc.ListUnread(slack.UnreadOptions{
			ChannelsOnly:   channelsOnly,
			DMsOnly:        dmsOnly,
			MinUnreadCount: minUnread,
			OrderBy:        orderBy,
		})

		if err != nil {
			result := output.Error("unread_list_failed", err.Error(), "Check your token and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"unread_conversations": unreads,
			"total_count":          len(unreads),
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add unread commands
	RootCmd.AddCommand(unreadCmd)
	unreadCmd.AddCommand(unreadListCmd)

	// List flags
	unreadListCmd.Flags().Bool("channels-only", false, "Only show unread channels (not DMs)")
	unreadListCmd.Flags().Bool("dms-only", false, "Only show unread DMs (1-on-1 and groups)")
	unreadListCmd.Flags().Int("min-unread", 0, "Minimum number of unread messages to show (0 = show all unread)")
	unreadListCmd.Flags().String("order-by", "count", "Sort order: 'count' (most unread first, default) or 'oldest' (oldest unread first for processing old items)")
}
