package config

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TUI state
type configTUI struct {
	config     *Config
	configPath string
	cursor     int
	fields     []field
	width      int
	height     int
}

type field struct {
	name        string
	value       string
	description string
	category    string
	editFunc    func(string) error
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	fieldStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#FAFAFA")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// NewConfigTUI creates a new configuration TUI
func NewConfigTUI(config *Config, configPath string) *configTUI {
	tui := &configTUI{
		config:     config,
		configPath: configPath,
		fields:     make([]field, 0),
	}

	tui.setupFields()
	return tui
}

func (t *configTUI) setupFields() {
	// OpenAI Configuration
	t.fields = append(t.fields, []field{
		{
			name:        "API Key",
			value:       t.config.OpenAI.APIKey,
			description: "Your OpenAI API key",
			category:    "OpenAI",
			editFunc:    func(v string) error { t.config.OpenAI.APIKey = v; return nil },
		},
		{
			name:        "Base URL",
			value:       t.config.OpenAI.BaseURL,
			description: "OpenAI API base URL",
			category:    "OpenAI",
			editFunc:    func(v string) error { t.config.OpenAI.BaseURL = v; return nil },
		},
		{
			name:        "Model",
			value:       t.config.OpenAI.Model,
			description: "Model to use for generating commit messages",
			category:    "OpenAI",
			editFunc:    func(v string) error { t.config.OpenAI.Model = v; return nil },
		},
		{
			name:        "Max Tokens",
			value:       fmt.Sprintf("%d", *t.config.OpenAI.MaxTokens),
			description: "Maximum tokens for response",
			category:    "OpenAI",
			editFunc: func(v string) error {
				val, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				t.config.OpenAI.MaxTokens = &val
				return nil
			},
		},
		{
			name:        "Temperature",
			value:       fmt.Sprintf("%.2f", *t.config.OpenAI.Temperature),
			description: "Temperature for response generation (0.0-2.0)",
			category:    "OpenAI",
			editFunc: func(v string) error {
				val, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return err
				}
				fval := float32(val)
				t.config.OpenAI.Temperature = &fval
				return nil
			},
		},
		{
			name:        "Timeout",
			value:       fmt.Sprintf("%d", t.config.OpenAI.Timeout),
			description: "Request timeout in seconds",
			category:    "OpenAI",
			editFunc: func(v string) error {
				val, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				t.config.OpenAI.Timeout = val
				return nil
			},
		},
	}...)

	// Commit Template Configuration
	t.fields = append(t.fields, []field{
		{
			name:        "Style",
			value:       t.config.CommitTemplate.Style,
			description: "Commit message style (conventional, simple, etc.)",
			category:    "Commit Template",
			editFunc:    func(v string) error { t.config.CommitTemplate.Style = v; return nil },
		},
		{
			name:        "Max Length",
			value:       fmt.Sprintf("%d", t.config.CommitTemplate.MaxLength),
			description: "Maximum commit message length",
			category:    "Commit Template",
			editFunc: func(v string) error {
				val, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				t.config.CommitTemplate.MaxLength = val
				return nil
			},
		},
		{
			name:        "Include Scope",
			value:       fmt.Sprintf("%t", t.config.CommitTemplate.IncludeScope),
			description: "Include scope in commit messages",
			category:    "Commit Template",
			editFunc: func(v string) error {
				val, err := strconv.ParseBool(v)
				if err != nil {
					return err
				}
				t.config.CommitTemplate.IncludeScope = val
				return nil
			},
		},
	}...)
}

func (t configTUI) Init() tea.Cmd {
	return nil
}

func (t configTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
		return t, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit

		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}
			return t, nil

		case "down", "j":
			if t.cursor < len(t.fields)-1 {
				t.cursor++
			}
			return t, nil

		case "enter", "":
			// Edit the selected field
			return t.editField(t.cursor), nil

		case "s":
			// Save configuration
			err := t.config.Save(t.configPath)
			if err != nil {
				return t, tea.Quit
			}
			return t, tea.Quit
		}
	}

	return t, nil
}

func (t configTUI) View() string {
	if t.width == 0 {
		return "Loading..."
	}

	var content strings.Builder

	// Title
	title := titleStyle.Render("GitR Configuration Editor")
	content.WriteString(title + "\n\n")

	// Instructions
	help := helpStyle.Render("↑/↓: Navigate • Enter: Edit • s: Save & Exit • q: Quit without saving")
	content.WriteString(help + "\n\n")

	// Fields
	for i, field := range t.fields {
		var style lipgloss.Style
		if i == t.cursor {
			style = selectedStyle
		} else {
			style = fieldStyle
		}

		// Category header
		if i == 0 || t.fields[i-1].category != field.category {
			content.WriteString(fmt.Sprintf("\n%s:\n", field.category))
		}

		// Field
		fieldText := fmt.Sprintf(" %s: %s", field.name, field.value)
		content.WriteString(style.Render(fieldText) + "\n")
		content.WriteString(fmt.Sprintf("   %s\n", helpStyle.Render(field.description)))
	}

	// Footer
	content.WriteString("\n" + helpStyle.Render("Press 's' to save changes or 'q' to quit without saving"))

	return content.String()
}

func (t configTUI) editField(index int) tea.Model {
	if index >= len(t.fields) {
		return t
	}

	field := t.fields[index]

	// Simple input prompt
	fmt.Printf("\nEnter new value for %s (current: %s): ", field.name, field.value)

	var input string
	fmt.Scanln(&input)

	if input != "" {
		err := field.editFunc(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			// Update the field value
			t.fields[index].value = input
		}
	}

	return t
}

// RunConfigEditor runs the configuration editor
func RunConfigEditor(config *Config, configPath string) error {
	// For now, use a simple interactive editor
	fmt.Println("=== GitR Configuration Editor ===")
	fmt.Printf("Configuration file: %s\n\n", configPath)

	// Display current configuration
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")
	fmt.Printf("API Key: %s\n", maskAPIKey(config.OpenAI.APIKey))
	fmt.Printf("Base URL: %s\n", config.OpenAI.BaseURL)
	fmt.Printf("Model: %s\n", config.OpenAI.Model)
	fmt.Printf("Max Tokens: %d\n", *config.OpenAI.MaxTokens)
	fmt.Printf("Temperature: %.2f\n", *config.OpenAI.Temperature)
	fmt.Printf("Timeout: %d seconds\n", config.OpenAI.Timeout)
	fmt.Printf("Style: %s\n", config.CommitTemplate.Style)
	fmt.Printf("Max Length: %d\n", config.CommitTemplate.MaxLength)
	fmt.Printf("Include Scope: %t\n", config.CommitTemplate.IncludeScope)
	fmt.Printf("Commit Without Confirmation: %t\n", config.CommitTemplate.CommitWithoutConfirmation)

	fmt.Println("\nOptions:")
	fmt.Println("1. Edit API Key")
	fmt.Println("2. Edit Base URL")
	fmt.Println("3. Edit Model")
	fmt.Println("4. Edit Max Tokens")
	fmt.Println("5. Edit Temperature")
	fmt.Println("6. Edit Timeout")
	fmt.Println("7. Edit Style")
	fmt.Println("8. Edit Max Length")
	fmt.Println("9. Edit Include Scope")
	fmt.Println("10. Edit Commit Without Confirmation")
	fmt.Println("s. Save and exit")
	fmt.Println("q. Quit without saving")

	for {
		fmt.Print("\nEnter choice: ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			editStringField("API Key", &config.OpenAI.APIKey)
		case "2":
			editStringField("Base URL", &config.OpenAI.BaseURL)
		case "3":
			editStringField("Model", &config.OpenAI.Model)
		case "4":
			editIntField("Max Tokens", config.OpenAI.MaxTokens)
		case "5":
			editFloatField("Temperature", config.OpenAI.Temperature)
		case "6":
			editIntField("Timeout", &config.OpenAI.Timeout)
		case "7":
			editStringField("Style", &config.CommitTemplate.Style)
		case "8":
			editIntField("Max Length", &config.CommitTemplate.MaxLength)
		case "9":
			editBoolField("Include Scope", &config.CommitTemplate.IncludeScope)
		case "10":
			editBoolField("Commit Without Confirmation", &config.CommitTemplate.CommitWithoutConfirmation)
		case "s":
			err := config.Save(configPath)
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				continue
			}
			fmt.Println("Configuration saved successfully!")
			return nil
		case "q":
			fmt.Println("Exiting without saving...")
			return nil
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func editStringField(name string, value *string) {
	fmt.Printf("Enter new value for %s (current: %s): ", name, *value)
	var input string
	fmt.Scanln(&input)
	if input != "" {
		*value = input
		fmt.Printf("%s updated to: %s\n", name, *value)
	}
}

func editIntField(name string, value *int) {
	fmt.Printf("Enter new value for %s (current: %d): ", name, *value)
	var input string
	fmt.Scanln(&input)
	if input != "" {
		if val, err := strconv.Atoi(input); err == nil {
			*value = val
			fmt.Printf("%s updated to: %d\n", name, *value)
		} else {
			fmt.Printf("Invalid number: %s\n", input)
		}
	}
}

func editFloatField(name string, value *float32) {
	fmt.Printf("Enter new value for %s (current: %.2f): ", name, *value)
	var input string
	fmt.Scanln(&input)
	if input != "" {
		if val, err := strconv.ParseFloat(input, 32); err == nil {
			*value = float32(val)
			fmt.Printf("%s updated to: %.2f\n", name, *value)
		} else {
			fmt.Printf("Invalid number: %s\n", input)
		}
	}
}

func editBoolField(name string, value *bool) {
	fmt.Printf("Enter new value for %s (current: %t) [true/false]: ", name, *value)
	var input string
	fmt.Scanln(&input)
	if input != "" {
		if val, err := strconv.ParseBool(input); err == nil {
			*value = val
			fmt.Printf("%s updated to: %t\n", name, *value)
		} else {
			fmt.Printf("Invalid boolean: %s\n", input)
		}
	}
}

// RunFirstTimeSetup runs the first-time configuration setup
func RunFirstTimeSetup(config *Config, configPath string) error {
	fmt.Println("Welcome to GitR!")
	fmt.Println("==================")
	fmt.Println("Let's set up your configuration for AI-powered commit messages.")
	fmt.Println("")

	// Required fields setup
	fmt.Println("Required Configuration:")
	fmt.Println("-------------------------")

	// API Key
	for {
		fmt.Print("Enter your OpenAI API Key: ")
		var apiKey string
		fmt.Scanln(&apiKey)
		if apiKey != "" {
			config.OpenAI.APIKey = apiKey
			break
		}
		fmt.Println("API Key is required!")
	}

	// Base URL
	fmt.Printf("Enter API Base URL (default: %s): ", config.OpenAI.BaseURL)
	var baseURL string
	fmt.Scanln(&baseURL)
	if baseURL != "" {
		config.OpenAI.BaseURL = baseURL
	}

	// Model
	fmt.Printf("Enter Model name (default: %s): ", config.OpenAI.Model)
	var model string
	fmt.Scanln(&model)
	if model != "" {
		config.OpenAI.Model = model
	}

	fmt.Println("")
	fmt.Println("️  Optional Configuration:")
	fmt.Println("-------------------------")

	// Max Tokens
	fmt.Printf("Enter Max Tokens (default: %d): ", *config.OpenAI.MaxTokens)
	var maxTokensStr string
	fmt.Scanln(&maxTokensStr)
	if maxTokensStr != "" {
		if val, err := strconv.Atoi(maxTokensStr); err == nil {
			config.OpenAI.MaxTokens = &val
		}
	}

	// Temperature
	fmt.Printf("Enter Temperature (default: %.2f): ", *config.OpenAI.Temperature)
	var tempStr string
	fmt.Scanln(&tempStr)
	if tempStr != "" {
		if val, err := strconv.ParseFloat(tempStr, 32); err == nil {
			fval := float32(val)
			config.OpenAI.Temperature = &fval
		}
	}

	// Commit Style
	fmt.Printf("Enter Commit Style (default: %s): ", config.CommitTemplate.Style)
	var style string
	fmt.Scanln(&style)
	if style != "" {
		config.CommitTemplate.Style = style
	}

	// Max Length
	fmt.Printf("Enter Max Commit Length (default: %d): ", config.CommitTemplate.MaxLength)
	var maxLengthStr string
	fmt.Scanln(&maxLengthStr)
	if maxLengthStr != "" {
		if val, err := strconv.Atoi(maxLengthStr); err == nil {
			config.CommitTemplate.MaxLength = val
		}
	}

	// Include Scope
	fmt.Printf("Include scope in commits? (default: %t) [true/false]: ", config.CommitTemplate.IncludeScope)
	var scopeStr string
	fmt.Scanln(&scopeStr)
	if scopeStr != "" {
		if val, err := strconv.ParseBool(scopeStr); err == nil {
			config.CommitTemplate.IncludeScope = val
		}
	}

	// Commit Without Confirmation
	fmt.Printf("Commit without confirmation? (default: %t) [true/false]: ", config.CommitTemplate.CommitWithoutConfirmation)
	var bypassStr string
	fmt.Scanln(&bypassStr)
	if bypassStr != "" {
		if val, err := strconv.ParseBool(bypassStr); err == nil {
			config.CommitTemplate.CommitWithoutConfirmation = val
		}
	}

	fmt.Println("")
	fmt.Println("Saving configuration...")

	// Save the configuration
	err := config.Save(configPath)
	if err != nil {
		return fmt.Errorf("failed to save configuration: %v", err)
	}

	fmt.Printf("Configuration saved to: %s\n", configPath)
	fmt.Println("")
	fmt.Println("Setup complete! You can now use GitR to generate commit messages.")
	fmt.Println("  Run 'gitr config' anytime to modify your settings.")
	fmt.Println("")

	return nil
}
