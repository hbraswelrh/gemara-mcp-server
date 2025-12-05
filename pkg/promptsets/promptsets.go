package promptsets

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Prompt represents a system prompt that can be used as a basis for dynamic prompts
type Prompt struct {
	// Name is a unique identifier for the prompt
	Name string

	// Description provides context about what this prompt is used for
	Description string

	// Content is the actual prompt text/template
	Content string

	// Variables is a map of variable names to their descriptions for dynamic substitution
	Variables map[string]string

	// Category groups related prompts together
	Category string

	// Tags help with searching and filtering prompts
	Tags []string
}

// PromptHandler is a function that can dynamically generate prompt content
type PromptHandler func(ctx context.Context, req PromptRequest) (PromptResponse, error)

// PromptRequest contains the context for generating a dynamic prompt
type PromptRequest struct {
	// PromptName identifies which prompt to use
	PromptName string

	// Variables contains values to substitute into the prompt template
	Variables map[string]interface{}

	// Context provides additional context for prompt generation
	Context map[string]interface{}
}

// PromptResponse contains the generated prompt content
type PromptResponse struct {
	// Content is the final prompt text after variable substitution
	Content string

	// Metadata contains additional information about the generated prompt
	Metadata map[string]interface{}
}

// PromptSet represents a collection of related prompts, similar to github-mcp-server's Toolset
type PromptSet struct {
	name        string
	description string
	prompts     map[string]*Prompt
	handlers    map[string]PromptHandler
	mu          sync.RWMutex
}

// NewPromptSet creates a new prompt set with the given name and description
func NewPromptSet(name, description string) *PromptSet {
	return &PromptSet{
		name:        name,
		description: description,
		prompts:     make(map[string]*Prompt),
		handlers:    make(map[string]PromptHandler),
	}
}

// Name returns the name of the prompt set
func (ps *PromptSet) Name() string {
	return ps.name
}

// Description returns the description of the prompt set
func (ps *PromptSet) Description() string {
	return ps.description
}

// AddPrompt adds a prompt to the prompt set
func (ps *PromptSet) AddPrompt(prompt *Prompt) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if prompt == nil {
		return fmt.Errorf("prompt cannot be nil")
	}

	if prompt.Name == "" {
		return fmt.Errorf("prompt name cannot be empty")
	}

	if _, exists := ps.prompts[prompt.Name]; exists {
		return fmt.Errorf("prompt %q already exists in prompt set %q", prompt.Name, ps.name)
	}

	ps.prompts[prompt.Name] = prompt
	return nil
}

// AddHandler adds a dynamic prompt handler to the prompt set
func (ps *PromptSet) AddHandler(name string, handler PromptHandler) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	if name == "" {
		return fmt.Errorf("handler name cannot be empty")
	}

	ps.handlers[name] = handler
	return nil
}

// GetPrompt retrieves a prompt by name
func (ps *PromptSet) GetPrompt(name string) (*Prompt, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	prompt, exists := ps.prompts[name]
	if !exists {
		return nil, fmt.Errorf("prompt %q not found in prompt set %q", name, ps.name)
	}

	return prompt, nil
}

// GetHandler retrieves a handler by name
func (ps *PromptSet) GetHandler(name string) (PromptHandler, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	handler, exists := ps.handlers[name]
	if !exists {
		return nil, fmt.Errorf("handler %q not found in prompt set %q", name, ps.name)
	}

	return handler, nil
}

// ListPrompts returns all prompts in the prompt set
func (ps *PromptSet) ListPrompts() []*Prompt {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	prompts := make([]*Prompt, 0, len(ps.prompts))
	for _, prompt := range ps.prompts {
		prompts = append(prompts, prompt)
	}

	return prompts
}

// GeneratePrompt generates a prompt dynamically using a handler if available,
// otherwise performs simple variable substitution
func (ps *PromptSet) GeneratePrompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	// First try to use a handler if one exists
	if handler, err := ps.GetHandler(req.PromptName); err == nil {
		return handler(ctx, req)
	}

	// Fall back to static prompt with variable substitution
	prompt, err := ps.GetPrompt(req.PromptName)
	if err != nil {
		return PromptResponse{}, err
	}

	content := prompt.Content
	// Simple variable substitution: {{variable_name}}
	for key, value := range req.Variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
	}

	return PromptResponse{
		Content: content,
		Metadata: map[string]interface{}{
			"prompt_name": prompt.Name,
			"category":    prompt.Category,
			"tags":        prompt.Tags,
		},
	}, nil
}

// PromptSetGroup manages multiple prompt sets
type PromptSetGroup struct {
	sets map[string]*PromptSet
	mu   sync.RWMutex
}

// NewPromptSetGroup creates a new prompt set group
func NewPromptSetGroup() *PromptSetGroup {
	return &PromptSetGroup{
		sets: make(map[string]*PromptSet),
	}
}

// AddPromptSet adds a prompt set to the group
func (psg *PromptSetGroup) AddPromptSet(ps *PromptSet) error {
	psg.mu.Lock()
	defer psg.mu.Unlock()

	if ps == nil {
		return fmt.Errorf("prompt set cannot be nil")
	}

	if _, exists := psg.sets[ps.name]; exists {
		return fmt.Errorf("prompt set %q already exists", ps.name)
	}

	psg.sets[ps.name] = ps
	return nil
}

// GetPromptSet retrieves a prompt set by name
func (psg *PromptSetGroup) GetPromptSet(name string) (*PromptSet, error) {
	psg.mu.RLock()
	defer psg.mu.RUnlock()

	ps, exists := psg.sets[name]
	if !exists {
		return nil, fmt.Errorf("prompt set %q not found", name)
	}

	return ps, nil
}

// ListPromptSets returns all prompt sets in the group
func (psg *PromptSetGroup) ListPromptSets() []*PromptSet {
	psg.mu.RLock()
	defer psg.mu.RUnlock()

	sets := make([]*PromptSet, 0, len(psg.sets))
	for _, ps := range psg.sets {
		sets = append(sets, ps)
	}

	return sets
}

// GeneratePromptFromGroup generates a prompt from any prompt set in the group
func (psg *PromptSetGroup) GeneratePromptFromGroup(ctx context.Context, promptSetName, promptName string, variables map[string]interface{}) (PromptResponse, error) {
	ps, err := psg.GetPromptSet(promptSetName)
	if err != nil {
		return PromptResponse{}, err
	}

	return ps.GeneratePrompt(ctx, PromptRequest{
		PromptName: promptName,
		Variables:  variables,
	})
}
