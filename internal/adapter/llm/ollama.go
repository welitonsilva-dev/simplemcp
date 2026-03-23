package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"humancli-server/internal/domain/plan"
	"humancli-server/internal/infra/logger"
)

// Ollama é o adaptador HTTP para o servidor Ollama.
type OllamaClient struct {
	baseURL string
	model   string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

// NewOllamaClient retorna um novo OllamaClient configurado.
func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{baseURL: baseURL, model: model}
}

// Generate envia um prompt ao Ollama e retorna a resposta bruta.
func (c *OllamaClient) Generate(prompt string) (string, error) {
	body, err := json.Marshal(ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("falha ao serializar request: %w", err)
	}

	resp, err := http.Post(
		c.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		logger.Error("falha ao conectar ao ollama: %v", err)
		return "", fmt.Errorf("falha ao conectar ao ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama retornou status inesperado: %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("falha ao decodificar resposta do ollama: %v", err)
		return "", fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	if result.Response == "" {
		return "", fmt.Errorf("resposta vazia do ollama")
	}

	return result.Response, nil
}

// Plan cria o plano a partir do prompt de planejamento.
func (c *OllamaClient) Plan(history, tools string) (*plan.ExecutionPlan, error) {
	prompt := plannerPrompt(history, tools)
	raw, err := c.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar plano: %w", err)
	}
	return parsePlan(raw)
}

// Finalize gera resposta final em linguagem natural.
func (c *OllamaClient) Finalize(history string) (string, error) {
	prompt := finalizerPrompt(history)
	raw, err := c.Generate(prompt)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar resposta final: %w", err)
	}
	return strings.TrimSpace(raw), nil
}
