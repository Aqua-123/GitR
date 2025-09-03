package output

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CommitTUI handles the interactive commit message editor
type CommitTUI struct {
	message string
}

// NewCommitTUI creates a new commit TUI
func NewCommitTUI(message string) *CommitTUI {
	return &CommitTUI{
		message: message,
	}
}

// ShowCommitEditor displays the commit message and allows user to accept, edit, or reject
func (t *CommitTUI) ShowCommitEditor() (string, bool, error) {
	fmt.Println("Generated Commit Message:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(t.message)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("")

	for {
		fmt.Println("Options:")
		fmt.Println(" [a] Accept and commit")
		fmt.Println(" [e] Edit message")
		fmt.Println(" [r] Reject (don't commit)")
		fmt.Println("")
		fmt.Print("Choose an option [a/e/r]: ")

		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			return "", false, err
		}

		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "a", "accept":
			return t.message, true, nil
		case "e", "edit":
			editedMessage, err := t.editMessage()
			if err != nil {
				fmt.Printf("Error editing message: %v\n", err)
				continue
			}
			if editedMessage != "" {
				t.message = editedMessage
				fmt.Println("")
				fmt.Println("Updated Commit Message:")
				fmt.Println(strings.Repeat("=", 50))
				fmt.Println(t.message)
				fmt.Println(strings.Repeat("=", 50))
				fmt.Println("")
			}
		case "r", "reject":
			return "", false, nil
		default:
			fmt.Println("Invalid option. Please choose 'a', 'e', or 'r'.")
			fmt.Println("")
		}
	}
}

// editMessage allows the user to edit the commit message
func (t *CommitTUI) editMessage() (string, error) {
	fmt.Println("")
	fmt.Println("️  Edit Commit Message:")
	fmt.Println("(Press Enter to keep current message, or type new message)")
	fmt.Printf("Current: %s\n", t.message)
	fmt.Print("New message: ")

	reader := bufio.NewReader(os.Stdin)
	newMessage, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	newMessage = strings.TrimSpace(newMessage)
	if newMessage == "" {
		return t.message, nil // Keep current message
	}

	return newMessage, nil
}
