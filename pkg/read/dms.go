package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var dmsCmd = &cobra.Command{
	Use:   "dms",
	Short: "Direct message operations",
	Long:  `Query and inspect direct message conversations`,
}

var dmsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DM conversations",
	Long: `List all direct message conversations.

Examples:
  slka dms list
  slka dms list --limit 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")

		svc := slack.NewDMService(slackClient)
		dms, err := svc.List(limit)

		if err != nil {
			result := output.Error("dms_list_failed", err.Error(), "Check your token and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"dms": dms,
		})
		result.Print(outputPretty)
		return nil
	},
}

var dmsHistoryCmd = &cobra.Command{
	Use:   "history <user>",
	Short: "Get DM message history with a user",
	Long: `Get the message history of a direct message conversation with a user.

User can be specified as:
- User ID (U123456)
- Email address (user@example.com)
- Username (alice)

Examples:
  slka dms history alice
  slka dms history user@example.com
  slka dms history U123456 --limit 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userArg := args[0]
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

		messages, err := svc.GetHistory(userArg, opts)
		if err != nil {
			result := output.Error("dm_history_failed", err.Error(), "Check user exists and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"user":     userArg,
			"messages": messages,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add dms commands
	RootCmd.AddCommand(dmsCmd)
	dmsCmd.AddCommand(dmsListCmd)
	dmsCmd.AddCommand(dmsHistoryCmd)

	// List flags
	dmsListCmd.Flags().Int("limit", 0, "Maximum number of DM conversations to return")

	// History flags
	dmsHistoryCmd.Flags().String("since", "", "Only messages after this timestamp")
	dmsHistoryCmd.Flags().String("until", "", "Only messages before this timestamp")
	dmsHistoryCmd.Flags().Int("limit", 100, "Maximum number of messages")
}
