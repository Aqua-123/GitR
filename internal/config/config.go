package config

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

// Config represents the configuration structure
type Config struct {
	XMLName        xml.Name             `xml:"config"`
	OpenAI         OpenAIConfig         `xml:"openai"`
	CommitTemplate CommitTemplateConfig `xml:"commit_template"`
}

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	XMLName xml.Name `xml:"openai"`

	// Azure OpenAI Service configuration (optional)
	ByAzure    bool   `xml:"by_azure"`
	BaseURL    string `xml:"base_url"`
	APIVersion string `xml:"api_version"`

	// Basic configuration
	APIKey  string `xml:"api_key"`
	Timeout int    `xml:"timeout"` // in seconds

	// Model parameters
	Model            string   `xml:"model"`
	MaxTokens        *int     `xml:"max_tokens"`
	Temperature      *float32 `xml:"temperature"`
	TopP             *float32 `xml:"top_p"`
	Stop             []string `xml:"stop"`
	PresencePenalty  *float32 `xml:"presence_penalty"`
	FrequencyPenalty *float32 `xml:"frequency_penalty"`

	// Advanced parameters
	ResponseFormat *openai.ChatCompletionResponseFormat `xml:"response_format"`
	Seed           *int                                 `xml:"seed"`
	LogitBias      map[string]int                       `xml:"-"`
	User           *string                              `xml:"user"`
}

// CommitTemplateConfig represents commit template configuration
type CommitTemplateConfig struct {
	XMLName                   xml.Name `xml:"commit_template"`
	Style                     string   `xml:"style"`
	MaxLength                 int      `xml:"max_length"`
	IncludeScope              bool     `xml:"include_scope"`
	CommitWithoutConfirmation bool     `xml:"commit_without_confirmation"`
}

// Load loads configuration from XML file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = xml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves configuration to XML file
func (c *Config) Save(filename string) error {
	data, err := xml.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	// Add XML declaration
	xmlData := []byte(xml.Header + string(data))

	return os.WriteFile(filename, xmlData, 0644)
}

// FindConfigFile searches for config file in standard locations
// Priority: 1. .gitr_config in current directory, 2. ~/.gitr_config, 3. ~/.config/gitr/config
func FindConfigFile() (string, error) {
	// Check current directory first
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	localConfig := filepath.Join(pwd, ".gitr_config")
	if _, err := os.Stat(localConfig); err == nil {
		return localConfig, nil
	}

	// Check home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check ~/.gitr_config
	homeConfig := filepath.Join(homeDir, ".gitr_config")
	if _, err := os.Stat(homeConfig); err == nil {
		return homeConfig, nil
	}

	// Check ~/.config/gitr/config
	configDir := filepath.Join(homeDir, ".config", "gitr")
	configFile := filepath.Join(configDir, "config")
	if _, err := os.Stat(configFile); err == nil {
		return configFile, nil
	}

	// Return default location (current directory)
	return localConfig, nil
}

// CreateDefaultConfig creates a default configuration
func CreateDefaultConfig() *Config {
	return &Config{
		OpenAI: OpenAIConfig{
			ByAzure:          false,
			BaseURL:          "https://api.openai.com/v1",
			APIVersion:       "2023-05-15",
			APIKey:           "",
			Timeout:          30,
			Model:            "gpt-3.5-turbo",
			MaxTokens:        intPtr(150),
			Temperature:      float32Ptr(0.7),
			TopP:             float32Ptr(1.0),
			Stop:             []string{},
			PresencePenalty:  float32Ptr(0.0),
			FrequencyPenalty: float32Ptr(0.0),
			ResponseFormat:   nil,
			Seed:             nil,
			LogitBias:        map[string]int{},
			User:             nil,
		},
		CommitTemplate: CommitTemplateConfig{
			Style:                     "conventional",
			MaxLength:                 72,
			IncludeScope:              true,
			CommitWithoutConfirmation: false,
		},
	}
}

// Validate checks if the configuration has all required fields
func (c *Config) Validate() error {
	if c.OpenAI.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}
	if c.OpenAI.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}
	if c.OpenAI.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}

// IsFirstTimeSetup checks if this is a first-time setup (missing required fields)
func (c *Config) IsFirstTimeSetup() bool {
	return c.OpenAI.BaseURL == "" || c.OpenAI.APIKey == "" || c.OpenAI.Model == ""
}

// Helper functions for pointer creation
func intPtr(i int) *int {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}
