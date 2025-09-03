package main

import (
	"flag"
	"fmt"
	"os"

	"gitr/cmd"
	"gitr/internal/config"
	"gitr/internal/git"
	"gitr/internal/llm"
	"gitr/internal/output"
)

func main() {
	// Parse flags
	var commitFlag = flag.Bool("commit", false, "Generate commit message and commit with confirmation")
	var commitShortFlag = flag.Bool("c", false, "Generate commit message and commit with confirmation")
	var bypassFlag = flag.Bool("bypass", false, "Bypass confirmation dialog (overrides config)")
	var bypassShortFlag = flag.Bool("b", false, "Bypass confirmation dialog (overrides config)")
	var helpFlag = flag.Bool("help", false, "Show help information")
	var helpShortFlag = flag.Bool("h", false, "Show help information")
	flag.Parse()

	// Check for help flags first
	if *helpFlag || *helpShortFlag {
		showHelp()
		return
	}

	// Check if we have command line arguments (excluding flags)
	args := flag.Args()
	if len(args) > 0 {
		// Use cobra command structure for subcommands
		cmd.Execute()
		return
	}

	// Check if commit flag is set
	if *commitFlag || *commitShortFlag {
		bypassConfirmation := *bypassFlag || *bypassShortFlag
		generateAndCommit(bypassConfirmation)
		return
	}

	// Default behavior: generate commit message
	generateCommitMessage()
}

func generateCommitMessage() {
	// Check if we're in a git repository
	if !git.IsRepository() {
		fmt.Println("Error: Not in a git repository")
		os.Exit(1)
	}

	// Find and load configuration
	configPath, err := config.FindConfigFile()
	if err != nil {
		fmt.Printf("Error finding config file: %v\n", err)
		os.Exit(1)
	}

	var cfg *config.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// First time setup - run guided configuration
		cfg = config.CreateDefaultConfig()
		fmt.Println("First time setup detected!")
		fmt.Println("")

		err = config.RunFirstTimeSetup(cfg, configPath)
		if err != nil {
			fmt.Printf("Error during setup: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load existing config
		cfg, err = config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Check if this is an incomplete configuration (first time setup)
		if cfg.IsFirstTimeSetup() {
			fmt.Println("Configuration incomplete. Running setup...")
			fmt.Println("")

			err = config.RunFirstTimeSetup(cfg, configPath)
			if err != nil {
				fmt.Printf("Error during setup: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Validate configuration before proceeding
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		fmt.Println("Run 'gitr config' to fix your configuration.")
		os.Exit(1)
	}

	// Get staged changes
	stagedChanges, err := git.GetStagedChanges()
	if err != nil {
		fmt.Printf("Error getting staged changes: %v\n", err)
		os.Exit(1)
	}

	if stagedChanges == "" {
		fmt.Println("No staged changes found")
		return
	}

	// Generate commit message using LLM
	rawResponse, err := llm.GenerateCommitMessage(cfg, stagedChanges)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// Parse and format the response
	commitMessage := output.ParseCommitMessage(rawResponse)
	formattedOutput := output.FormatCommitMessageOutput(commitMessage)

	fmt.Print(formattedOutput)
}

func generateAndCommit(bypassConfirmation bool) {
	// Check if we're in a git repository
	if !git.IsRepository() {
		fmt.Println("Error: Not in a git repository")
		os.Exit(1)
	}

	// Check if there are staged changes
	hasStaged, err := git.HasStagedChanges()
	if err != nil {
		fmt.Printf("Error checking staged changes: %v\n", err)
		os.Exit(1)
	}

	if !hasStaged {
		fmt.Println("No staged changes found")
		fmt.Println("Please stage some changes first with 'git add'")
		os.Exit(1)
	}

	// Find and load configuration
	configPath, err := config.FindConfigFile()
	if err != nil {
		fmt.Printf("Error finding config file: %v\n", err)
		os.Exit(1)
	}

	var cfg *config.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// First time setup - run guided configuration
		cfg = config.CreateDefaultConfig()
		fmt.Println("First time setup detected!")
		fmt.Println("")

		err = config.RunFirstTimeSetup(cfg, configPath)
		if err != nil {
			fmt.Printf("Error during setup: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load existing config
		cfg, err = config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if cfg.IsFirstTimeSetup() {
			fmt.Println("️  Configuration incomplete. Running setup...")
			fmt.Println("")

			err = config.RunFirstTimeSetup(cfg, configPath)
			if err != nil {
				fmt.Printf("Error during setup: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Validate configuration before proceeding
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		fmt.Println("Run 'gitr config' to fix your configuration.")
		os.Exit(1)
	}

	// Get staged changes
	stagedChanges, err := git.GetStagedChanges()
	if err != nil {
		fmt.Printf("Error getting staged changes: %v\n", err)
		os.Exit(1)
	}

	// Generate commit message using LLM
	rawResponse, err := llm.GenerateCommitMessage(cfg, stagedChanges)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	commitMessage := output.ParseCommitMessage(rawResponse)
	if commitMessage == "" {
		fmt.Println("Failed to generate commit message")
		os.Exit(1)
	}

	// Determine if we should bypass confirmation
	shouldBypass := bypassConfirmation || cfg.CommitTemplate.CommitWithoutConfirmation

	var finalMessage string
	var shouldCommit bool

	if shouldBypass {
		// Bypass confirmation - use the generated message directly
		finalMessage = commitMessage
		shouldCommit = true
		fmt.Println("Bypassing confirmation (using generated message)...")
	} else {
		// Show commit editor TUI
		commitTUI := output.NewCommitTUI(commitMessage)
		finalMessage, shouldCommit, err = commitTUI.ShowCommitEditor()
		if err != nil {
			fmt.Printf("Error in commit editor: %v\n", err)
			os.Exit(1)
		}

		if !shouldCommit {
			fmt.Println("Commit cancelled")
			return
		}
	}

	// Perform the commit
	fmt.Println("")
	fmt.Println("Committing changes...")
	err = git.Commit(finalMessage)
	if err != nil {
		fmt.Printf("Commit failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Commit successful!")
	fmt.Printf("Committed with message: %s\n", finalMessage)
}

func showHelp() {
	fmt.Println("GitR - AI-powered Git commit message generator")
	fmt.Println("==============================================")
	fmt.Println("")
	fmt.Println("DESCRIPTION:")
	fmt.Println(" GitR uses AI to generate conventional commit messages based on your staged changes.")
	fmt.Println(" It can generate messages for display or commit directly with confirmation.")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println(" gitr [OPTIONS]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println(" -c, --commit")
	fmt.Println("       Generate commit message and commit with confirmation dialog")
	fmt.Println("")
	fmt.Println(" -b, --bypass")
	fmt.Println("       Bypass confirmation dialog (overrides config setting)")
	fmt.Println("       Can be combined with --commit flag")
	fmt.Println("")
	fmt.Println(" -h, --help")
	fmt.Println("       Show this help information")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println(" gitr config")
	fmt.Println("       Open interactive configuration editor")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println(" gitr                    # Generate and display commit message")
	fmt.Println(" gitr -c                 # Generate message and commit with confirmation")
	fmt.Println(" gitr -c -b              # Generate message and commit without confirmation")
	fmt.Println(" gitr --commit --bypass  # Same as -c -b")
	fmt.Println(" gitr config             # Edit configuration settings")
	fmt.Println(" gitr --help             # Show this help")
	fmt.Println("")
	fmt.Println("CONFIGURATION:")
	fmt.Println(" GitR looks for configuration in the following order:")
	fmt.Println(" 1. .gitr_config in current directory")
	fmt.Println(" 2. ~/.gitr_config in home directory")
	fmt.Println(" 3. ~/.config/gitr/config in standard config location")
	fmt.Println("")
	fmt.Println(" Run 'gitr config' to set up or modify your configuration.")
	fmt.Println("")
	fmt.Println("REQUIREMENTS:")
	fmt.Println(" - Must be run in a git repository")
	fmt.Println(" - Must have staged changes (for commit operations)")
	fmt.Println(" - Requires OpenAI API key and model configuration")
}
