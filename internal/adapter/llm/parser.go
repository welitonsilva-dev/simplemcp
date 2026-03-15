package llm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"humancli-server/internal/domain/plan"
)

// Plan monta o prompt, envia ao Ollama e parseia a resposta para um ExecutionPlan.
// É o ponto de entrada público do pacote llm.
func (c *Client) Plan(userInput, tools string) (*plan.ExecutionPlan, error) {
	prompt := plannerPrompt(userInput, tools)

	raw, err := c.generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar plano: %w", err)
	}

	return parse(raw)
}

// parse extrai e decodifica o ExecutionPlan do JSON bruto retornado pela LLM.
func parse(raw string) (*plan.ExecutionPlan, error) {
	cleaned := cleanJSON(raw)

	var result plan.ExecutionPlan
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("falha ao parsear JSON do plano: %w — raw: %s", err, cleaned)
	}

	if len(result.Steps) == 0 {
		return nil, fmt.Errorf("plano sem steps: %s", cleaned)
	}

	return &result, nil
}

// cleanJSON extrai JSON válido de uma string com possíveis ruídos da LLM.
// Absorve a lógica de utils/json.go — responsabilidade do adaptador LLM.
func cleanJSON(input string) string {
	input = strings.TrimSpace(input)

	// remove blocos de código markdown
	input = strings.ReplaceAll(input, "```json", "")
	input = strings.ReplaceAll(input, "```", "")

	// remove comentários // e /* */
	reComments := regexp.MustCompile(`(?s)//.*?\n|/\*.*?\*/`)
	input = reComments.ReplaceAllString(input, "")

	// extrai apenas o bloco JSON entre { e }
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start >= 0 && end >= 0 && end > start {
		input = input[start : end+1]
	}

	// remove vírgulas trailing antes de } e ]
	reTrailing := regexp.MustCompile(`,\s*([}\]])`)
	input = reTrailing.ReplaceAllString(input, "$1")

	return strings.TrimSpace(input)
}
