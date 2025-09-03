package llm

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"

	"gitr/internal/config"
)

// GenerateCommitMessage generates a commit message using OpenAI
func GenerateCommitMessage(cfg *config.Config, stagedChanges string) (string, error) {
	ctx := context.Background()

	// Create OpenAI chat model
	chatModel, err := createOpenAIModel(ctx, &cfg.OpenAI)
	if err != nil {
		return "", err
	}

	// Create prompt template for commit message generation
	template := createCommitMessageTemplate(cfg)

	// Prepare scope instruction
	scopeInstruction := ""
	if cfg.CommitTemplate.IncludeScope {
		scopeInstruction = "Include a scope in the commit message when appropriate."
	} else {
		scopeInstruction = "Do not include a scope in the commit message."
	}

	// Format the prompt with staged changes
	messages, err := template.Format(ctx, map[string]any{
		"staged_changes":            stagedChanges,
		"style":                     cfg.CommitTemplate.Style,
		"max_length":                cfg.CommitTemplate.MaxLength,
		"include_scope_instruction": scopeInstruction,
	})

	if err != nil {
		return "", err
	}

	// Generate response
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	return result.Content, nil
}

// createOpenAIModel creates an OpenAI chat model
func createOpenAIModel(ctx context.Context, cfg *config.OpenAIConfig) (model.ChatModel, error) {
	// Use API key from config, fallback to environment variable
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// Set default timeout if not specified
	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout == 0 {
		timeout = 30 * time.Second
	}

	// Initialize logit bias if nil
	logitBias := cfg.LogitBias
	if logitBias == nil {
		logitBias = map[string]int{}
	}

	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		// Azure OpenAI Service configuration (optional)
		ByAzure:    cfg.ByAzure,
		BaseURL:    cfg.BaseURL,
		APIVersion: cfg.APIVersion,

		// Basic configuration
		APIKey:  apiKey,
		Timeout: timeout,

		// Model parameters
		Model:            cfg.Model,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             cfg.TopP,
		Stop:             cfg.Stop,
		PresencePenalty:  cfg.PresencePenalty,
		FrequencyPenalty: cfg.FrequencyPenalty,

		// Advanced parameters
		ResponseFormat: cfg.ResponseFormat,
		Seed:           cfg.Seed,
		LogitBias:      logitBias,
		User:           cfg.User,
	})
}

// createCommitMessageTemplate creates a prompt template for commit message generation
func createCommitMessageTemplate(cfg *config.Config) *prompt.DefaultChatTemplate {
	return prompt.FromMessages(schema.FString,
		// System message template
		schema.SystemMessage(`You are an expert Git commit message generator. Generate clear, concise, and conventional commit messages based on the staged changes. Follow the {style} style and keep the message under {max_length} characters. {include_scope_instruction}
		Be less specific about the changes and only include the most important changes or the general change or broader concept that the user is trying to convey.`),

		// User message template
		schema.UserMessage("Please generate a commit message for the following staged changes:\n\n{staged_changes}"),
	)
}
