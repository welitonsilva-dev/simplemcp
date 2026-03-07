package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
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
		return "", err
	}

	defer resp.Body.Close()

	var result Response

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Response, nil
}
