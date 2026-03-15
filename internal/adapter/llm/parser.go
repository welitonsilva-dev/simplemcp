package llm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"humancli-server/internal/domain/plan"
)

// Plan monta o prompt com o histórico completo, envia ao LLM e parseia a resposta.
// É chamado em cada iteração do loop ReAct no AgentUseCase.
func (c *Client) Plan(history, tools string) (*plan.ExecutionPlan, error) {
	prompt := agentPrompt(history, tools)

	raw, err := c.generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar plano: %w", err)
	}

	return parse(raw)
}

// parse extrai e decodifica o ExecutionPlan do JSON bruto retornado pelo LLM.
// Trata tanto respostas de ação (steps) quanto de encerramento (final).
func parse(raw string) (*plan.ExecutionPlan, error) {
	cleaned := cleanJSON(raw)

	var result plan.ExecutionPlan
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("falha ao parsear JSON do plano: %w — raw: %s", err, cleaned)
	}

	// resposta de encerramento válida: final=true com mensagem
	if result.Final {
		if result.FinalMessage == "" {
			result.FinalMessage = "tarefa concluída"
		}
		return &result, nil
	}

	// resposta de ação: deve ter ao menos um step
	if len(result.Steps) == 0 {
		return nil, fmt.Errorf("plano sem steps e sem final: %s", cleaned)
	}

	return &result, nil
}

// cleanJSON extrai JSON válido de uma string com possíveis ruídos do LLM.
// Remove blocos de código markdown, comentários e texto fora do JSON.
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
