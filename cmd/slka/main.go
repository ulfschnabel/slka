package main

import (
	"os"

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
