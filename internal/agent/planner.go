package agent

import (
	"encoding/json"
	"errors"
	"fmt"

	"simplemcp/internal/llm"
	"simplemcp/internal/pipeline"
	"simplemcp/utils"
)

// Step representa um passo do plano
type Step struct {
	Tool   string                 `json:"tool"`
	Params map[string]interface{} `json:"params"`
}

// Plan representa um plano completo (multi-step)
type Plan struct {
	Steps      []Step  `json:"steps"`
	Confidence float64 `json:"confidence"`
}

// Planner é o responsável por gerar planos
type Planner struct {
	LLM llm.Client
}

// NewPlanner cria um planner com LLM
func NewPlanner(llmClient llm.Client) *Planner {
	return &Planner{
		LLM: llmClient,
	}
}

// Generate cria um plano a partir da entrada do usuário
func (p *Planner) Generate(userInput string, toolRegistry string) (*Plan, error) {

	input, err := pipeline.OptimizeInput(userInput)
	if err != nil {
		return nil, errors.New("Erro ao otimizar o input")
	}

	// montar prompt para o LLM
	prompt := llm.PlannerPrompt(input, toolRegistry)
	fmt.Println(prompt)

	// gerar plano usando o LLM
	llmResp, err := p.LLM.Generate(prompt)
	if err != nil {
		return nil, err
	}

	clean := utils.CleanJSON(llmResp)

	var plan Plan
	if err := json.Unmarshal([]byte(clean), &plan); err != nil {
		return nil, errors.New(clean + err.Error())
	}

	return &plan, nil
}
