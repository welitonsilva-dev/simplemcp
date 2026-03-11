package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"simplemcp/internal/infra/logger"
)

// Client é o adaptador HTTP para o servidor Ollama.
type Client struct {
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

// NewClient retorna um novo Client configurado.
func NewClient(baseURL, model string) *Client {
	return &Client{baseURL: baseURL, model: model}
}

// generate envia um prompt ao Ollama e retorna a resposta bruta.
func (c *Client) generate(prompt string) (string, error) {
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
