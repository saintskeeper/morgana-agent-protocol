package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Loader handles loading and caching of agent prompts
type Loader struct {
	agentDir string
	cache    map[string]string
	mu       sync.RWMutex
}

// NewPromptLoader creates a new Loader instance
func NewPromptLoader(agentDir string) *Loader {
	return &Loader{
		agentDir: agentDir,
		cache:    make(map[string]string),
	}
}

// Load retrieves an agent prompt, using cache if available
func (l *Loader) Load(agentType string) (string, error) {
	// Check cache first
	l.mu.RLock()
	if prompt, ok := l.cache[agentType]; ok {
		l.mu.RUnlock()
		return prompt, nil
	}
	l.mu.RUnlock()

	// Load from file
	path := filepath.Join(l.agentDir, agentType+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Provide a fallback prompt
			fallback := fmt.Sprintf("You are a %s specialist. Follow best practices and project conventions.", agentType)
			l.cachePrompt(agentType, fallback)
			return fallback, nil
		}
		return "", fmt.Errorf("reading agent file: %w", err)
	}

	// Parse YAML frontmatter if present
	prompt := l.extractPrompt(string(content))

	// Cache the result
	l.cachePrompt(agentType, prompt)

	return prompt, nil
}

// extractPrompt removes YAML frontmatter and returns the content
func (l *Loader) extractPrompt(content string) string {
	// Check if content starts with YAML frontmatter
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			// Return content after frontmatter
			return strings.TrimSpace(parts[2])
		}
	}
	// No frontmatter, return entire content
	return strings.TrimSpace(content)
}

// cachePrompt safely adds a prompt to the cache
func (l *Loader) cachePrompt(agentType, prompt string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache[agentType] = prompt
}

// ClearCache removes all cached prompts
func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[string]string)
}
