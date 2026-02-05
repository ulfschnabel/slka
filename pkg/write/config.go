package write

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/config"
	"github.com/ulf/slka/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage slka configuration`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration (tokens masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := config.DefaultConfigPath()
		cfg, err := config.Load(cfgPath)
		if err != nil {
			result := output.Error("config_not_found", err.Error(), "Run 'slka-write config init' to create a config file")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Detect token types
		readTokenType := config.DetectTokenType(cfg.ReadToken)
		writeTokenType := config.DetectTokenType(cfg.WriteToken)

		result := output.Success(map[string]interface{}{
			"config_file":      cfgPath,
			"read_token":       config.MaskToken(cfg.ReadToken),
			"read_token_type":  config.GetTokenTypeName(readTokenType),
			"write_token":      config.MaskToken(cfg.WriteToken),
			"write_token_type": config.GetTokenTypeName(writeTokenType),
			"user_token":       config.MaskToken(cfg.UserToken),
			"require_approval": cfg.RequireApproval,
		})
		result.Print(outputPretty)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Valid keys: read_token, write_token, user_token, require_approval`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfgPath := config.DefaultConfigPath()
		cfg, err := config.Load(cfgPath)
		if err != nil {
			// Create new config if it doesn't exist
			cfg = &config.Config{}
		}

		// Set the value
		switch key {
		case "read_token":
			cfg.ReadToken = value
		case "write_token":
			cfg.WriteToken = value
		case "user_token":
			cfg.UserToken = value
		case "require_approval":
			if value == "true" {
				cfg.RequireApproval = true
			} else if value == "false" {
				cfg.RequireApproval = false
			} else {
				result := output.Error("invalid_value", "require_approval must be 'true' or 'false'", "")
				result.Print(outputPretty)
				return fmt.Errorf("exit code %d", result.ExitCode())
			}
		default:
			result := output.Error("invalid_key", fmt.Sprintf("unknown config key: %s", key), "Valid keys: read_token, write_token, user_token, require_approval")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Save config
		if err := cfg.Save(cfgPath); err != nil {
			result := output.Error("save_failed", err.Error(), "Check file permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"config_file": cfgPath,
			"key":         key,
			"updated":     true,
		})
		result.Print(outputPretty)
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("slka Configuration Setup")
		fmt.Println("========================")
		fmt.Println()

		cfg := &config.Config{}

		// Read token
		fmt.Print("Read token (xoxb-...): ")
		if scanner.Scan() {
			cfg.ReadToken = strings.TrimSpace(scanner.Text())
		}

		// Write token
		fmt.Print("Write token (xoxb-...): ")
		if scanner.Scan() {
			cfg.WriteToken = strings.TrimSpace(scanner.Text())
		}

		// User token (optional)
		fmt.Print("User token (xoxp-..., optional): ")
		if scanner.Scan() {
			cfg.UserToken = strings.TrimSpace(scanner.Text())
		}

		// Require approval
		fmt.Print("Require approval for write operations? [Y/n]: ")
		if scanner.Scan() {
			response := strings.ToLower(strings.TrimSpace(scanner.Text()))
			cfg.RequireApproval = response != "n" && response != "no"
		} else {
			cfg.RequireApproval = true // default to true
		}

		// Validate
		if err := cfg.Validate(); err != nil {
			result := output.Error("invalid_config", err.Error(), "Read and write tokens are required")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		// Save
		cfgPath := config.DefaultConfigPath()
		if err := cfg.Save(cfgPath); err != nil {
			result := output.Error("save_failed", err.Error(), "Check file permissions")
			result.Print(outputPretty)
			return fmt.Errorf("exit code %d", result.ExitCode())
		}

		result := output.Success(map[string]interface{}{
			"config_file":      cfgPath,
			"require_approval": cfg.RequireApproval,
			"initialized":      true,
		})
		result.Print(outputPretty)
		return nil
	},
}

func init() {
	// Add config commands
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)
}
