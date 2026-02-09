package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	readpkg "github.com/ulf/slka/pkg/read"
	writepkg "github.com/ulf/slka/pkg/write"
)

var Version = "dev"

func main() {
	// Use read package as base (it has the full setup)
	rootCmd := readpkg.RootCmd
	rootCmd.Use = "slka"
	rootCmd.Short = "Slack CLI for Agentic Workflows"
	rootCmd.Long = `slka provides command-line access to Slack for automation and AI agents. Supports both user and bot tokens.`

	// Store original read PreRunE
	originalReadPreRunE := rootCmd.PersistentPreRunE

	// Wrap PersistentPreRunE to handle both read and write initialization
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Always run read initialization
		if err := originalReadPreRunE(cmd, args); err != nil {
			return err
		}

		// If this is a write command, also run write initialization
		if isWriteCommand(cmd) {
			// Get config file path from flags
			configFile, _ := cmd.Root().PersistentFlags().GetString("config")
			if configFile == "" {
				configFile = os.Getenv("HOME") + "/.config/slka/config.json"
			}
			if err := writepkg.Initialize(configFile); err != nil {
				return err
			}
		}

		return nil
	}

	// Merge write commands into root
	for _, writeCmd := range writepkg.RootCmd.Commands() {
		// Check if command already exists (like "channels")
		existingCmd := findCommand(rootCmd, writeCmd.Use)

		if existingCmd != nil {
			// Merge subcommands from write into existing command
			for _, writeSubCmd := range writeCmd.Commands() {
				existingCmd.AddCommand(writeSubCmd)
			}
		} else {
			// Add new command (like "message", "config")
			rootCmd.AddCommand(writeCmd)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func findCommand(rootCmd *cobra.Command, use string) *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == use {
			return cmd
		}
	}
	return nil
}

// isWriteCommand checks if a command is a write command
func isWriteCommand(cmd *cobra.Command) bool {
	// Check if command path contains write-only commands
	writeCommands := []string{"message", "config"}
	writeSubcommands := []string{"send", "reply", "edit", "create", "archive", "invite", "kick", "init"}

	cmdPath := cmd.CommandPath()

	// Check if it's a top-level write command
	for _, writeCmd := range writeCommands {
		if strings.Contains(cmdPath, " "+writeCmd+" ") || strings.HasSuffix(cmdPath, " "+writeCmd) {
			return true
		}
	}

	// Check if it's a write subcommand
	for _, writeSubCmd := range writeSubcommands {
		if strings.HasSuffix(cmdPath, " "+writeSubCmd) {
			return true
		}
	}

	return false
}
