package output

import (
	"regexp"
	"strings"
)

// ParseCommitMessage extracts the actual commit message from LLM response
// Handles various response formats including markdown code blocks
func ParseCommitMessage(response string) string {
	// Remove leading/trailing whitespace
	response = strings.TrimSpace(response)
	
	// If response is empty, return as is
	if response == "" {
		return response
	}
	
	// Try to extract content from markdown code blocks
	// Pattern: ```[language]\ncontent\n```
	codeBlockPattern := regexp.MustCompile("```(?:\\w+)?\\s*\\n([\\s\\S]*?)\\n```")
	matches := codeBlockPattern.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	
	// Try to extract content from single backticks
	// Pattern: `content`
	singleBacktickPattern := regexp.MustCompile("`([^`]+)`")
	matches = singleBacktickPattern.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	
	// If no code blocks found, return the response as is (trimmed)
	return response
}

// FormatCommitMessageOutput formats the commit message for display
func FormatCommitMessageOutput(commitMessage string) string {
	if commitMessage == "" {
		return "No commit message generated"
	}
	
	// Clean up the message
	commitMessage = strings.TrimSpace(commitMessage)
	
	// Create formatted output
	var output strings.Builder
	output.WriteString("Generated commit message:\n")
	output.WriteString(strings.Repeat("=", 50))
	output.WriteString("\n")
	output.WriteString(commitMessage)
	output.WriteString("\n")
	
	return output.String()
}
