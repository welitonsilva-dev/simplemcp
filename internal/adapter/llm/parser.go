package llm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"humancli-server/internal/domain/plan"
)

// Plan decide qual tool chamar nesta iteração.
// Retorna um plano com a tool escolhida, ou tool="none" se não há ação necessária.
func (c *Client) Plan(history, tools string) (*plan.ExecutionPlan, error) {
	prompt := plannerPrompt(history, tools)
	raw, err := c.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar plano: %w", err)
	}
	return parsePlan(raw)
}

// Finalize gera a resposta final em linguagem natural após o loop encerrar.
// Chamado uma única vez quando o agente decide que a tarefa foi concluída.
func (c *Client) Finalize(history string) (string, error) {
	prompt := finalizerPrompt(history)
	raw, err := c.Generate(prompt)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar resposta final: %w", err)
	}
	return strings.TrimSpace(raw), nil
}

// parsePlan decodifica a resposta do plannerPrompt.
// Espera: {"tool": "...", "params": {}, "confidence": 0.9}
func parsePlan(raw string) (*plan.ExecutionPlan, error) {
	cleaned := cleanJSON(raw)

	var result struct {
		Tool       string         `json:"tool"`
		Params     map[string]any `json:"params"`
		Confidence float64        `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("falha ao parsear plano: %w — raw: %s", err, cleaned)
	}

	if result.Tool == "" {
		return nil, fmt.Errorf("plano sem tool: %s", cleaned)
	}

	// "none" sinaliza que o LLM não precisa mais de tools — encerra o loop
	if result.Tool == "none" {
		return &plan.ExecutionPlan{Final: true}, nil
	}

	if result.Params == nil {
		result.Params = map[string]any{}
	}

	return &plan.ExecutionPlan{
		Steps:      []plan.ToolCall{{Tool: result.Tool, Params: result.Params}},
		Confidence: result.Confidence,
	}, nil
}

// cleanJSON extrai JSON válido de uma string com possíveis ruídos do LLM.
func cleanJSON(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "```json", "")
	input = strings.ReplaceAll(input, "```", "")

	reComments := regexp.MustCompile(`(?s)//.*?\n|/\*.*?\*/`)
	input = reComments.ReplaceAllString(input, "")

	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start >= 0 && end >= 0 && end > start {
		input = input[start : end+1]
	}

	reTrailing := regexp.MustCompile(`,\s*([}\]])`)
	input = reTrailing.ReplaceAllString(input, "$1")

	return strings.TrimSpace(input)
}
