package llm

import (
	"strings"

	llm "humancli-server/internal/domain/provider"
	"humancli-server/internal/infra/config"
)

// NewProvider cria a implementação adequada de LLM de acordo com a configuração.
func NewProvider(cfg *config.Config) llm.Provider {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" {
		provider = "ollama"
	}

	switch provider {
	case "groq":
		return NewGroqClient(cfg.LLMBaseURL, cfg.Model, cfg.LLMAPIKey)
	case "ollama":
		return NewOllamaClient(cfg.OllamaURL, cfg.Model)
	default:
		return NewOllamaClient(cfg.OllamaURL, cfg.Model)
	}
}
