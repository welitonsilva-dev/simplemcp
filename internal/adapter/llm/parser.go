package llm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"humancli-server/internal/domain/plan"
	"humancli-server/internal/infra/logger"
)

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
		logger.Error("falha ao parsear plano: %v — raw: %s", err, raw)
		return nil, fmt.Errorf("falha ao parsear plano: %w — raw: %s", err, cleaned)
	}

	if result.Tool == "" {
		logger.Error("plano sem tool: %s", cleaned)
		return nil, fmt.Errorf("plano sem tool: %s", cleaned)
	}

	// "none" sinaliza que o LLM não precisa mais de tools — encerra o loop
	if result.Tool == "none" {
		logger.Info("plano final recebido (tool=none)")
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
