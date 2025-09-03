package cmd

import (
	"fmt"
	"os"

	"gitr/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitr",
	Short: "AI-powered Git commit message generator",
	Long: `GitR is a tool that uses LLMs to generate conventional commit messages based on your staged changes.

Usage:
  gitr                    Generate and display commit message
  gitr --commit, -c       Generate commit message and commit with confirmation
  gitr -c --bypass, -b    Generate message and commit without confirmation
  gitr --help, -h         Show help information
  gitr config             Open configuration editor

Examples:
  gitr                    # Just generate and show the commit message
  gitr -c                 # Generate message and commit with TUI confirmation
  gitr -c -b              # Generate message and commit without confirmation
  gitr --commit --bypass  # Same as -c -b
  gitr --help             # Show help information
  gitr config             # Edit configuration settings`,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure GitR settings",
	Long:  `Open an interactive configuration editor to set up your GitR preferences.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfigEditor()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runConfigEditor() {
	// Find or create config file
	configPath, err := config.FindConfigFile()
	if err != nil {
		fmt.Printf("Error finding config file: %v\n", err)
		os.Exit(1)
	}

	// Load existing config or create default
	var cfg *config.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config and run first-time setup
		cfg = config.CreateDefaultConfig()
		fmt.Println("No configuration found. Starting first-time setup...")
		fmt.Println("")

		err = config.RunFirstTimeSetup(cfg, configPath)
		if err != nil {
			fmt.Printf("Error during setup: %v\n", err)
			os.Exit(1)
		}
		return
	} else {
		// Load existing config
		cfg, err = config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	}

	// Run the configuration editor
	err = config.RunConfigEditor(cfg, configPath)
	if err != nil {
		fmt.Printf("Error in configuration editor: %v\n", err)
		os.Exit(1)
	}
}
