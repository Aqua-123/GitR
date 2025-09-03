package git

import (
	"os/exec"
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
