package llm

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadPrompt(t *testing.T) {
	// Test parsing YAML prompt structure
	yamlContent := `name: Test Prompt
description: A test prompt
model: openai/gpt-4o
modelParameters:
  temperature: 0.7
  topP: 0.9
messages:
  - role: system
    content: You are a test assistant
  - role: user
    content: Hello {{name}}
`

	var config PromptConfig
	err := yaml.Unmarshal([]byte(yamlContent), &config)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	// Validate parsed structure
	if config.Name != "Test Prompt" {
		t.Errorf("Name = %v, want 'Test Prompt'", config.Name)
	}
	if config.Description != "A test prompt" {
		t.Errorf("Description = %v, want 'A test prompt'", config.Description)
	}
	if config.Model != "openai/gpt-4o" {
		t.Errorf("Model = %v, want 'openai/gpt-4o'", config.Model)
	}
	if config.ModelParameters.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", config.ModelParameters.Temperature)
	}
	if config.ModelParameters.TopP != 0.9 {
		t.Errorf("TopP = %v, want 0.9", config.ModelParameters.TopP)
	}
	if len(config.Messages) != 2 {
		t.Fatalf("Messages length = %d, want 2", len(config.Messages))
	}
	if config.Messages[0].Role != "system" {
		t.Errorf("Message[0] Role = %v, want 'system'", config.Messages[0].Role)
	}
	if config.Messages[1].Content != "Hello {{name}}" {
		t.Errorf("Message[1] Content = %v, want 'Hello {{name}}'", config.Messages[1].Content)
	}
}

func TestLoadPrompt_FromFile(t *testing.T) {
	// Test loading from testdata file
	path := "../testdata/sample_prompt.yml"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("Skipping file test, testdata not found: %v", err)
		return
	}

	var config PromptConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		t.Fatalf("Failed to parse YAML from file: %v", err)
	}

	if config.Name == "" {
		t.Error("Name is empty")
	}
	if config.Model == "" {
		t.Error("Model is empty")
	}
	if len(config.Messages) == 0 {
		t.Error("Messages is empty")
	}
}

func TestLoadPrompt_NotFound(t *testing.T) {
	_, err := LoadPrompt("nonexistent_prompt_file")
	if err == nil {
		t.Error("LoadPrompt() expected error for non-existent file, got nil")
	}
}

func TestRenderPrompt(t *testing.T) {
	tests := []struct {
		name        string
		config      *PromptConfig
		vars        map[string]string
		wantContent string
		wantErr     bool
	}{
		{
			name: "Simple variable replacement",
			config: &PromptConfig{
				Name:  "Test",
				Model: "gpt-4",
				Messages: []PromptMessage{
					{Role: "user", Content: "Hello {{name}}!"},
				},
			},
			vars:        map[string]string{"name": "Alice"},
			wantContent: "Hello Alice!",
			wantErr:     false,
		},
		{
			name: "Multiple variables",
			config: &PromptConfig{
				Name:  "Test",
				Model: "gpt-4",
				Messages: []PromptMessage{
					{Role: "user", Content: "{{greeting}} {{name}}, welcome to {{place}}!"},
				},
			},
			vars: map[string]string{
				"greeting": "Hello",
				"name":     "Bob",
				"place":    "GitHub",
			},
			wantContent: "Hello Bob, welcome to GitHub!",
			wantErr:     false,
		},
		{
			name: "Missing variable",
			config: &PromptConfig{
				Name:  "Test",
				Model: "gpt-4",
				Messages: []PromptMessage{
					{Role: "user", Content: "Hello {{name}}!"},
				},
			},
			vars:    map[string]string{},
			wantErr: true,
		},
		{
			name: "No variables to replace",
			config: &PromptConfig{
				Name:  "Test",
				Model: "gpt-4",
				Messages: []PromptMessage{
					{Role: "user", Content: "Hello World!"},
				},
			},
			vars:        map[string]string{},
			wantContent: "Hello World!",
			wantErr:     false,
		},
		{
			name: "Multiple messages",
			config: &PromptConfig{
				Name:  "Test",
				Model: "gpt-4",
				Messages: []PromptMessage{
					{Role: "system", Content: "You are {{role}}"},
					{Role: "user", Content: "Analyze {{repo}}"},
				},
			},
			vars: map[string]string{
				"role": "an assistant",
				"repo": "myrepo",
			},
			wantContent: "Analyze myrepo",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderPrompt(tt.config, tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Check that the last message contains the expected content
			if len(got.Messages) > 0 {
				lastMsg := got.Messages[len(got.Messages)-1]
				if !strings.Contains(lastMsg.Content, tt.wantContent) {
					t.Errorf("RenderPrompt() content = %v, want to contain %v", lastMsg.Content, tt.wantContent)
				}
			}

			// Verify no unreplaced variables remain
			for _, msg := range got.Messages {
				if strings.Contains(msg.Content, "{{") && strings.Contains(msg.Content, "}}") {
					t.Errorf("RenderPrompt() contains unreplaced variable in: %v", msg.Content)
				}
			}

			// Verify original config is not modified
			if tt.config.Messages[0].Content != got.Messages[0].Content {
				// This is expected - they should be different after rendering
				for _, msg := range tt.config.Messages {
					if !strings.Contains(msg.Content, "{{") {
						continue // Skip messages without variables
					}
					// Original should still have variables
					if !strings.Contains(msg.Content, "{{") || !strings.Contains(msg.Content, "}}") {
						t.Error("RenderPrompt() modified original config")
					}
				}
			}
		})
	}
}

func TestRenderPrompt_PreservesStructure(t *testing.T) {
	original := &PromptConfig{
		Name:        "Original",
		Description: "Test prompt",
		Model:       "gpt-4",
		ModelParameters: ModelParameters{
			Temperature: 0.5,
			TopP:        0.8,
		},
		Messages: []PromptMessage{
			{Role: "system", Content: "You are {{role}}"},
			{Role: "user", Content: "Hello {{name}}"},
		},
	}

	vars := map[string]string{
		"role": "helper",
		"name": "user",
	}

	rendered, err := RenderPrompt(original, vars)
	if err != nil {
		t.Fatalf("RenderPrompt() error = %v", err)
	}

	// Check that structure is preserved
	if rendered.Name != original.Name {
		t.Errorf("RenderPrompt() Name = %v, want %v", rendered.Name, original.Name)
	}
	if rendered.Description != original.Description {
		t.Errorf("RenderPrompt() Description = %v, want %v", rendered.Description, original.Description)
	}
	if rendered.Model != original.Model {
		t.Errorf("RenderPrompt() Model = %v, want %v", rendered.Model, original.Model)
	}
	if rendered.ModelParameters.Temperature != original.ModelParameters.Temperature {
		t.Errorf("RenderPrompt() Temperature = %v, want %v", rendered.ModelParameters.Temperature, original.ModelParameters.Temperature)
	}
	if len(rendered.Messages) != len(original.Messages) {
		t.Errorf("RenderPrompt() Messages length = %d, want %d", len(rendered.Messages), len(original.Messages))
	}
	for i := range rendered.Messages {
		if rendered.Messages[i].Role != original.Messages[i].Role {
			t.Errorf("RenderPrompt() Message[%d] Role = %v, want %v", i, rendered.Messages[i].Role, original.Messages[i].Role)
		}
	}
}
