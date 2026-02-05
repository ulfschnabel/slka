package read

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ulf/slka/internal/config"
	"github.com/ulf/slka/internal/output"
	"github.com/ulf/slka/internal/slack"
)

var (
	cfgFile     string
	token       string
	outputPretty bool
	cfg         *config.Config
	slackClient slack.Client
)

// RootCmd is the root command for slka-read
var RootCmd = &cobra.Command{
	Use:   "slka-read",
	Short: "Slack CLI for read-only operations",
	Long:  `slka-read provides read-only access to Slack for querying channels, messages, users, and more.`,
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
			cfg.ReadToken = token
		}

		// Validate we have a token
		if cfg.ReadToken == "" {
			return fmt.Errorf("no read token configured. Set SLKA_READ_TOKEN environment variable, use --token flag, or run 'slka-write config init'")
		}

		// Create Slack client
		slackClient = slack.NewClient(cfg.ReadToken)

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/slka/config.json)")
	RootCmd.PersistentFlags().StringVar(&token, "token", "", "Slack read token (overrides config)")
	RootCmd.PersistentFlags().BoolVar(&outputPretty, "output-pretty", false, "Pretty print JSON output")
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(output.ExitGeneralError)
	}
}
