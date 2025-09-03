package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsRepository checks if the current directory is a git repository
func IsRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetStagedChanges returns the staged changes using git diff --cached
func GetStagedChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Commit creates a git commit with the given message
func Commit(message string) error {
	// Escape the message for shell execution
	escapedMessage := strings.ReplaceAll(message, `"`, `\"`)
	cmd := exec.Command("git", "commit", "-m", escapedMessage)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", string(output))
	}

	return nil
}

// HasStagedChanges checks if there are any staged changes
func HasStagedChanges() (bool, error) {
	changes, err := GetStagedChanges()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(changes) != "", nil
}
