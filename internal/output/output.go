package output

import (
	"regexp"
	"strings"
)

func ParseCommitMessage(response string) string {
	response = strings.TrimSpace(response)
	if response == "" {
		return response
	}

	codeBlockPattern := regexp.MustCompile("```(?:\\w+)?\\s*\\n([\\s\\S]*?)\\n```")
	matches := codeBlockPattern.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	singleBacktickPattern := regexp.MustCompile("`([^`]+)`")
	matches = singleBacktickPattern.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return response
}

// FormatCommitMessageOutput formats the commit message for display
func FormatCommitMessageOutput(commitMessage string) string {
	if commitMessage == "" {
		return "No commit message generated"
	}

	commitMessage = strings.TrimSpace(commitMessage)

	var output strings.Builder
	output.WriteString("Generated commit message:\n")
	output.WriteString(strings.Repeat("=", 50))
	output.WriteString("\n")
	output.WriteString(commitMessage)
	output.WriteString("\n")

	return output.String()
}
