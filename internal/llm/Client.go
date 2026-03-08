package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"simplemcp/internal/logger"
)

type Client struct {
	BaseURL string
	Model   string
}

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
}

func NewClient(url string, model string) *Client {
	return &Client{
		BaseURL: url,
		Model:   model,
	}
}

// Generate envia um prompt para o modelo e retorna a resposta gerada
func (c *Client) Generate(prompt string) (string, error) {

	reqBody := Request{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		c.BaseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		logger.Error("failed to send request: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to decode response: %v", err)
		return "", err
	}

	if result.Response == "" {
		logger.Error("empty response from ollama")
		return "", fmt.Errorf("empty response from ollama")
	}

	return result.Response, nil
}
