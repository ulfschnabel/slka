package write

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/approval"
	"github.com/ulf/slka/internal/config"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
	"golang.org/x/term"
)

var (
	cfgFile     string
	token       string
	outputPretty bool
	dryRun      bool
	cfg         *config.Config
	slackClient slack.Client
	approver    *approval.Approver
)

// RootCmd is the root command for slka-write
var RootCmd = &cobra.Command{
	Use:   "slka-write",
	Short: "Slack CLI for write operations",
	Long:  `slka-write provides write access to Slack for sending messages, managing channels, and more.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		var err error
		if cfgFile == "" {
			cfgFile = config.DefaultConfigPath()
		}

		cfg, err = config.Load(cfgFile)
		if err != nil {
			// Config file not found is ok, we'll use environment variables or flags
			cfg = &config.Config{}
		}

		// Override with flag if provided
		if token != "" {
			cfg.WriteToken = token
		}

		// Validate we have a token
		if cfg.WriteToken == "" {
			return fmt.Errorf("no write token configured. Set SLKA_WRITE_TOKEN environment variable, use --token flag, or run 'slka-write config init'")
		}

		// Create Slack client
		slackClient = slack.NewClient(cfg.WriteToken)

		// Create approver
		isatty := term.IsTerminal(int(os.Stdin.Fd()))
		approver = approval.NewApprover(isatty, nil)
		approver.SetRequired(cfg.RequireApproval)

		// If dry run, disable approval
		if dryRun {
			approver.SetRequired(false)
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/slka/config.json)")
	RootCmd.PersistentFlags().StringVar(&token, "token", "", "Slack write token (overrides config)")
	RootCmd.PersistentFlags().BoolVar(&outputPretty, "output-pretty", false, "Pretty print JSON output")
	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without executing")
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(output.ExitGeneralError)
	}
}

// Initialize sets up the write package with the given config file path
// This is used by the unified slka binary to initialize write commands
func Initialize(configFile string) error {
	cfgFile = configFile

	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		// Config file not found is ok
		cfg = &config.Config{}
	}

	// Validate we have a write token
	if cfg.WriteToken == "" {
		return fmt.Errorf("no write token configured")
	}

	// Create Slack client
	slackClient = slack.NewClient(cfg.WriteToken)

	// Create approver
	isatty := term.IsTerminal(int(os.Stdin.Fd()))
	approver = approval.NewApprover(isatty, nil)
	approver.SetRequired(cfg.RequireApproval)

	return nil
}
