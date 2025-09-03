package main

import (
	"fmt"
	"os"

	"gitr/internal/config"
	"gitr/internal/git"
	"gitr/internal/llm"
	"gitr/internal/output"
)

func main() {
	// Check if we're in a git repository
	if !git.IsRepository() {
		fmt.Println("Error: Not in a git repository")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
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
