package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var dmCmd = &cobra.Command{
	Use:   "dm",
	Short: "Direct message operations",
	Long:  `Query, send, and manage direct messages (1-on-1 and group DMs)`,
}

var dmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DM conversations",
	Long: `List all direct message conversations (both 1-on-1 and group DMs).

Examples:
  slka dm list
  slka dm list --limit 50
  slka dm list --filter alice
  slka dm list --filter alice@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		filter, _ := cmd.Flags().GetString("filter")

		svc := slack.NewDMService(slackClient)
		dms, err := svc.List(limit)

		if err != nil {
			result := output.Error("dm_list_failed", err.Error(), "Check your token and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Filter by user if specified
		if filter != "" {
			// Resolve filter user to ID
			filterUserID, err := svc.ResolveUser(filter)
			if err != nil {
				result := output.Error("user_not_found", err.Error(), "Check the user ID, email, or username")
				result.Print(outputPretty)
				return fmt.Errorf("exit code %d", result.ExitCode())
			}

			// Filter conversations to only those containing the user
			filteredDMs := make([]slack.DMInfo, 0)
			for _, dm := range dms {
				for _, userID := range dm.UserIDs {
					if userID == filterUserID {
						filteredDMs = append(filteredDMs, dm)
						break
					}
				}
			}
			dms = filteredDMs
		}

		result := output.Success(map[string]interface{}{
			"conversations": dms,
		})
		result.Print(outputPretty)
		return nil
	},
}

var dmHistoryCmd = &cobra.Command{
	Use:   "history <users>",
	Short: "Get DM message history",
	Long: `Get the message history of a direct message conversation.

Users can be specified as:
- Single user: alice, user@example.com, U123456
- Multiple users (group DM): alice,bob,charlie

Examples:
  slka dm history alice
  slka dm history user@example.com
  slka dm history alice,bob,charlie
  slka dm history U123456,U789012 --limit 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		usersArg := args[0]
		sinceStr, _ := cmd.Flags().GetString("since")
		untilStr, _ := cmd.Flags().GetString("until")
		limit, _ := cmd.Flags().GetInt("limit")

		svc := slack.NewDMService(slackClient)

		opts := slack.HistoryOptions{
			Limit: limit,
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

		messages, err := svc.GetHistory(usersArg, opts)
		if err != nil {
			result := output.Error("dm_history_failed", err.Error(), "Check users exist and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"users":    usersArg,
			"messages": messages,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add dm commands
	RootCmd.AddCommand(dmCmd)
	dmCmd.AddCommand(dmListCmd)
	dmCmd.AddCommand(dmHistoryCmd)

	// List flags
	dmListCmd.Flags().Int("limit", 0, "Maximum number of DM conversations to return")
	dmListCmd.Flags().String("filter", "", "Filter conversations by user (name, email, or ID)")

	// History flags
	dmHistoryCmd.Flags().String("since", "", "Only messages after this timestamp")
	dmHistoryCmd.Flags().String("until", "", "Only messages before this timestamp")
	dmHistoryCmd.Flags().Int("limit", 100, "Maximum number of messages")
}
