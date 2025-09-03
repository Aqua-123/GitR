package config

import (
	"encoding/json"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

// Config represents the configuration structure
type Config struct {
	OpenAI         OpenAIConfig         `json:"openai"`
	CommitTemplate CommitTemplateConfig `json:"commit_template"`
}

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	// Azure OpenAI Service configuration (optional)
	ByAzure    bool   `json:"by_azure"`
	BaseURL    string `json:"base_url"`
	APIVersion string `json:"api_version"`

	// Basic configuration
	APIKey  string `json:"api_key"`
	Timeout int    `json:"timeout"` // in seconds

	// Model parameters
	Model            string   `json:"model"`
	MaxTokens        *int     `json:"max_tokens"`
	Temperature      *float32 `json:"temperature"`
	TopP             *float32 `json:"top_p"`
	Stop             []string `json:"stop"`
	PresencePenalty  *float32 `json:"presence_penalty"`
	FrequencyPenalty *float32 `json:"frequency_penalty"`

	// Advanced parameters
	ResponseFormat *openai.ChatCompletionResponseFormat `json:"response_format"`
	Seed           *int                                 `json:"seed"`
	LogitBias      map[string]int                       `json:"logit_bias"`
	User           *string                              `json:"user"`
}

// CommitTemplateConfig represents commit template configuration
type CommitTemplateConfig struct {
	Style        string `json:"style"`
	MaxLength    int    `json:"max_length"`
	IncludeScope bool   `json:"include_scope"`
}

// Load loads configuration from JSON file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
