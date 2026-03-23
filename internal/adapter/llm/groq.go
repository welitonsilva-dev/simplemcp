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

type GroqClient struct {
	baseURL string
	model   string
	apiKey  string
}

// NewGroqClient retorna um novo GroqClient configurado.
func NewGroqClient(baseURL, model, apiKey string) *GroqClient {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}
	return &GroqClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		apiKey:  apiKey,
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Generate envia um prompt ao Groq e retorna a resposta bruta.
func (c *GroqClient) Generate(prompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("LLM_API_KEY não definido para Groq")
	}

	reqBody := groqRequest{
		Model: c.model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 1024,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("falha ao serializar request Groq: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("falha ao conectar Groq: %v", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Groq retornou status inesperado: %d", res.StatusCode)
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("falha ao decodificar Groq: %w", err)
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("resposta vazia do Groq")
	}

	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

// Plan cria o plano a partir do prompt de planejamento.
func (c *GroqClient) Plan(history, tools string) (*plan.ExecutionPlan, error) {
	prompt := plannerPrompt(history, tools)
	raw, err := c.Generate(prompt)
	if err != nil {
		logger.Error("falha ao gerar plano com Groq: %v", err)
		return nil, fmt.Errorf("falha ao gerar plano: %w", err)
	}
	return parsePlan(raw)
}

// Finalize gera resposta final em linguagem natural.
func (c *GroqClient) Finalize(history string) (string, error) {
	prompt := finalizerPrompt(history)
	raw, err := c.Generate(prompt)
	if err != nil {
		logger.Error("falha ao gerar resposta final com Groq: %v", err)
		return "", fmt.Errorf("falha ao gerar resposta final: %w", err)
	}
	return strings.TrimSpace(raw), nil
}
