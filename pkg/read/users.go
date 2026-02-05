package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User operations",
	Long:  `Query and inspect Slack users`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		includeBots, _ := cmd.Flags().GetBool("include-bots")
		includeDeleted, _ := cmd.Flags().GetBool("include-deleted")
		limit, _ := cmd.Flags().GetInt("limit")

		svc := slack.NewUserService(slackClient)
		users, err := svc.List(slack.ListUsersOptions{
			IncludeBots:    includeBots,
			IncludeDeleted: includeDeleted,
			Limit:          limit,
		})

		if err != nil {
			result := output.Error("users_list_failed", err.Error(), "Check your token and permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"users": users,
		})
		result.Print(outputPretty)
		return nil
	},
}

var usersLookupCmd = &cobra.Command{
	Use:   "lookup <query>",
	Short: "Find a user by name or email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		byField, _ := cmd.Flags().GetString("by")

		svc := slack.NewUserService(slackClient)
		user, err := svc.Lookup(query, byField)

		if err != nil {
			result := output.Error("user_not_found", err.Error(), "Check the username or email")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"user": user,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add users commands
	RootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersLookupCmd)

	// List flags
	usersListCmd.Flags().Bool("include-bots", false, "Include bot users")
	usersListCmd.Flags().Bool("include-deleted", false, "Include deactivated users")
	usersListCmd.Flags().Int("limit", 0, "Maximum number of users")

	// Lookup flags
	usersLookupCmd.Flags().String("by", "auto", "Search by: name, email, auto")
}
