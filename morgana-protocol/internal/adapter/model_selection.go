package adapter

import (
	"fmt"
	"log"
	"strings"
)

// ModelSelector handles model selection logic based on agent type, retries, and complexity
type ModelSelector struct {
	defaultModels map[string]string
	escalationMap map[string]map[string]string
}

// NewModelSelector creates a new model selector with default configurations
func NewModelSelector() *ModelSelector {
	return &ModelSelector{
		defaultModels: map[string]string{
			"code-implementer":  "claude-3-7-sonnet-20250219",
			"test-specialist":   "claude-3-7-sonnet-20250219",
			"validation-expert": "claude-3-7-sonnet-20250219",
			"sprint-planner":    "claude-3-7-sonnet-20250219",
		},
		escalationMap: map[string]map[string]string{
			"code-implementer": {
				"retry_1":            "claude-4-sonnet",
				"retry_2":            "claude-4-opus",
				"validation_failure": "claude-4-sonnet",
				"complexity_high":    "claude-4-opus",
			},
			"test-specialist": {
				"retry_1":           "claude-4-sonnet",
				"retry_2":           "claude-4-opus",
				"complex_testing":   "claude-4-sonnet",
				"performance_tests": "claude-4-opus",
			},
			"validation-expert": {
				"retry_1":          "claude-4-sonnet",
				"retry_2":          "claude-4-opus",
				"security_audit":   "claude-4-opus",
				"complex_analysis": "claude-4-opus",
			},
			"sprint-planner": {
				"retry_1":             "claude-4-sonnet",
				"retry_2":             "gemini-2.5-pro",
				"complex_planning":    "gemini-2.5-pro",
				"architecture_design": "o3",
			},
		},
	}
}

// SelectModel determines the appropriate model based on task context
func (ms *ModelSelector) SelectModel(task Task) string {
	// Check if model hint is provided
	if task.ModelHint != "" {
		log.Printf("Using model hint: %s for agent: %s", task.ModelHint, task.AgentType)
		return task.ModelHint
	}

	// Get escalation rules for this agent type
	escalationRules, exists := ms.escalationMap[task.AgentType]
	if !exists {
		log.Printf("No escalation rules found for agent: %s, using default", task.AgentType)
		return ms.getDefaultModel(task.AgentType)
	}

	// Check retry-based escalation first
	if task.RetryCount > 0 {
		retryKey := fmt.Sprintf("retry_%d", task.RetryCount)
		if model, exists := escalationRules[retryKey]; exists {
			log.Printf("Escalating to %s for %s (retry %d)", model, task.AgentType, task.RetryCount)
			return model
		}

		// Fall back to highest retry level if specific retry not found
		if task.RetryCount >= 2 {
			if model, exists := escalationRules["retry_2"]; exists {
				log.Printf("Escalating to %s for %s (retry %d, using retry_2 fallback)", model, task.AgentType, task.RetryCount)
				return model
			}
		}
	}

	// Check complexity-based escalation
	if task.Complexity != "" {
		complexityKey := fmt.Sprintf("complexity_%s", strings.ToLower(task.Complexity))
		if model, exists := escalationRules[complexityKey]; exists {
			log.Printf("Escalating to %s for %s (complexity: %s)", model, task.AgentType, task.Complexity)
			return model
		}

		// Check for specific complexity mappings
		switch strings.ToLower(task.Complexity) {
		case "high", "complex":
			if task.AgentType == "code-implementer" {
				if model, exists := escalationRules["complexity_high"]; exists {
					log.Printf("Escalating to %s for %s (high complexity)", model, task.AgentType)
					return model
				}
			}
		case "security":
			if task.AgentType == "validation-expert" {
				if model, exists := escalationRules["security_audit"]; exists {
					log.Printf("Escalating to %s for %s (security complexity)", model, task.AgentType)
					return model
				}
			}
		case "performance":
			if task.AgentType == "test-specialist" {
				if model, exists := escalationRules["performance_tests"]; exists {
					log.Printf("Escalating to %s for %s (performance complexity)", model, task.AgentType)
					return model
				}
			}
		case "planning":
			if task.AgentType == "sprint-planner" {
				if model, exists := escalationRules["complex_planning"]; exists {
					log.Printf("Escalating to %s for %s (planning complexity)", model, task.AgentType)
					return model
				}
			}
		case "architecture":
			if task.AgentType == "sprint-planner" {
				if model, exists := escalationRules["architecture_design"]; exists {
					log.Printf("Escalating to %s for %s (architecture complexity)", model, task.AgentType)
					return model
				}
			}
		}
	}

	// Check for validation failure escalation (could be passed via Options)
	if options := task.Options; options != nil {
		if validationFailed, exists := options["validation_failed"]; exists && validationFailed == true {
			if model, exists := escalationRules["validation_failure"]; exists {
				log.Printf("Escalating to %s for %s (validation failure)", model, task.AgentType)
				return model
			}
		}
	}

	// Default model
	defaultModel := ms.getDefaultModel(task.AgentType)
	log.Printf("Using default model %s for agent %s", defaultModel, task.AgentType)
	return defaultModel
}

// getDefaultModel returns the default model for an agent type
func (ms *ModelSelector) getDefaultModel(agentType string) string {
	if model, exists := ms.defaultModels[agentType]; exists {
		return model
	}
	return "claude-3-7-sonnet-20250219" // Global fallback
}

// GetModelCapabilities returns information about model capabilities
func (ms *ModelSelector) GetModelCapabilities(model string) map[string]interface{} {
	capabilities := map[string]interface{}{
		"token_efficient": false,
		"reasoning_level": "standard",
		"cost_tier":       "standard",
	}

	switch model {
	case "claude-3-7-sonnet-20250219":
		capabilities["token_efficient"] = true
		capabilities["reasoning_level"] = "fast"
		capabilities["cost_tier"] = "low"
	case "claude-4-sonnet":
		capabilities["reasoning_level"] = "high"
		capabilities["cost_tier"] = "medium"
	case "claude-4-opus":
		capabilities["reasoning_level"] = "maximum"
		capabilities["cost_tier"] = "high"
	case "gemini-2.5-pro":
		capabilities["reasoning_level"] = "high"
		capabilities["cost_tier"] = "medium"
		capabilities["specialization"] = "planning"
	case "o3":
		capabilities["reasoning_level"] = "systematic"
		capabilities["cost_tier"] = "high"
		capabilities["specialization"] = "architecture"
	}

	return capabilities
}
