package llm

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed prompts/*.prompt.yml
var promptsFS embed.FS

// ModelParameters represents LLM model parameters
type ModelParameters struct {
	Temperature float64 `yaml:"temperature"`
	TopP        float64 `yaml:"topP"`
}

// PromptMessage represents a single message in a prompt
type PromptMessage struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

// PromptConfig represents a YAML prompt configuration
type PromptConfig struct {
	Name            string          `yaml:"name"`
	Description     string          `yaml:"description"`
	Model           string          `yaml:"model"`
	ModelParameters ModelParameters `yaml:"modelParameters"`
	Messages        []PromptMessage `yaml:"messages"`
}

// LoadPrompt loads a YAML prompt configuration from file.
// It first tries to load from an external file (for development/customization),
// and if that fails, loads from the embedded filesystem (production).
func LoadPrompt(name string) (*PromptConfig, error) {
	var data []byte
	var err error

	// Try to load from external file first (for development and customization)
	externalPath := filepath.Join("internal", "llm", "prompts", name+".prompt.yml")
	data, err = os.ReadFile(externalPath)

	// If external file not found, load from embedded filesystem
	if err != nil {
		embeddedPath := filepath.Join("prompts", name+".prompt.yml")
		data, err = promptsFS.ReadFile(embeddedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read prompt file %s (tried external and embedded): %w", name, err)
		}
	}

	// Parse YAML
	var config PromptConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from %s: %w", name, err)
	}

	return &config, nil
}

// RenderPrompt renders a prompt config by replacing variables in messages
func RenderPrompt(config *PromptConfig, vars map[string]string) (*PromptConfig, error) {
	// Create a copy of the config
	rendered := &PromptConfig{
		Name:            config.Name,
		Description:     config.Description,
		Model:           config.Model,
		ModelParameters: config.ModelParameters,
		Messages:        make([]PromptMessage, len(config.Messages)),
	}

	// Render each message
	for i, msg := range config.Messages {
		renderedContent := msg.Content

		// Replace all variables in the content
		for key, value := range vars {
			placeholder := "{{" + key + "}}"
			if strings.Contains(renderedContent, placeholder) {
				renderedContent = strings.ReplaceAll(renderedContent, placeholder, value)
			}
		}

		// Check if there are any unreplaced variables
		if strings.Contains(renderedContent, "{{") && strings.Contains(renderedContent, "}}") {
			// Find the first unreplaced variable
			start := strings.Index(renderedContent, "{{")
			end := strings.Index(renderedContent[start:], "}}") + start + 2
			missingVar := renderedContent[start:end]
			return nil, fmt.Errorf("missing variable in vars map: %s", missingVar)
		}

		rendered.Messages[i] = PromptMessage{
			Role:    msg.Role,
			Content: renderedContent,
		}
	}

	return rendered, nil
}
